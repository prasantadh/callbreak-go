package callbreak

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/prasantadh/callbreak-go/pkg/deck"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func New() *CallBreak {
	game := &CallBreak{
		workPermit: make(chan struct{}, 1),
	}
	return game
}

// Get the current state of the game as visible to a token
func (game *CallBreak) Query(token Token) (*CallBreak, error) {
	// todo: currently returns tricks for all players.
	// consider sending out only the tricks won by current player
	me := game.Turn(&token)
	if me == NPlayers {
		return nil, fmt.Errorf("unauthorized token")
	}
	response := *game
	for p := range game.Players {
		if p != me {
			response.Players[p].Token = Token("")
		}
	}
	for r := range response.Rounds {
		round := &response.Rounds[r]
		for p := range game.Players {
			if p != me {
				round.Hands[p] = Hand{}
			}
		}
	}
	return &response, nil
}

// add a player to the game. returns an authentication token on success
// else return error on failure
func (game *CallBreak) AddPlayer(name, strategy string, timeout time.Duration) (PlayerId, error) {
	// TODO eventually add a registry for assistant
	// this works for now

	game.workPermit <- struct{}{}
	defer func() { <-game.workPermit }()

	if game.TotalPlayers == NPlayers {
		return PlayerId{}, fmt.Errorf("could not add players: table full")
	}

	s, err := GetStrategy(strategy)
	if err != nil {
		return PlayerId{}, fmt.Errorf("could not set strategy: %v", err)
	}

	buffer := make([]byte, 32)
	_, err = rand.Read(buffer)
	if err != nil {
		return PlayerId{}, fmt.Errorf("could not generate a token")
	}
	token := Token(hex.EncodeToString(buffer))
	playerid := PlayerId{Name: name, Token: token}

	assistant := Assistant{strategy: s, game: game, last: game, token: token}
	assistant.ticker = time.NewTicker(timeout)
	go assistant.Assist()
	game.Players[game.TotalPlayers] = playerid

	log.Infof("add player %s with token %s", playerid.Name, playerid.Token)
	game.TotalPlayers += 1

	if game.TotalPlayers == NPlayers {
		game.Rounds[game.RoundNumber].deal()
		// TODO: implement called
		game.Stage = CALLED
	}

	return playerid, nil
}

// deal the cards to the players.
func (round *Round) deal() {

	d := deck.New()
	// TODO: make sure each player is dealt at least one Hukum
	// and at least one of Q, K, A else shuffle again
	for i, card := range d {
		player := i % NPlayers
		cardnumber := i / NPlayers
		round.Hands[player][cardnumber] = card
	}

	for i := range round.Hands {
		round.Hands[i].Sort()
	}

	log.Infof("server: dealt:")
	for i := 0; i < NPlayers; i++ {
		log.Infof("\t%s", round.Hands[i])
	}
}

// the next player in line playes the card c
func (game *CallBreak) Break(token Token, card deck.Card) error {

	game.workPermit <- struct{}{}
	defer func() { <-game.workPermit }()

	// assert Play is currently valid move
	//      players have been dealt the cards
	//      players have made the calls
	//      there is an active round and active trick
	log.Infof("server: player %s attempted play with %s", token[:4], card)
	player := game.Turn(&token)
	if player == NPlayers {
		return fmt.Errorf("invalid token")
	}

	if game.Stage != CALLED {
		return fmt.Errorf("not a valid stage for this move")
	}

	round := &game.Rounds[game.RoundNumber]
	// the TrickNumber should always be valid for current game
	// or the game is in inconsistent state
	trick := &round.Tricks[round.TrickNumber]
	next := (trick.Lead + trick.Size) % NPlayers

	if player != next {
		return fmt.Errorf("you are not up next")
	}

	// TODO: eventually move this "server:" as a logger field
	log.Infof("server: RoundNumber: %d\tTrickNumber: %d",
		game.RoundNumber, round.TrickNumber)
	log.Infof("server: trick: %s (size: %d lead: %d)",
		trick.Cards, trick.Size, trick.Lead)
	log.Infof("server: Hand: %s", round.Hands[player])

	validMoves, err := GetValidMoves(game)
	log.Infof("valid moves: %s", validMoves)
	if err != nil {
		return fmt.Errorf("could not get valid moves: %v", err)
	}
	if !slices.Contains(validMoves, card) {
		return fmt.Errorf("invalid move from player")
	}

	// TODO: all the following sections can likely be refactored out to a func
	hand := &round.Hands[player]
	for i, c := range hand { // play the card
		if c.Suit == card.Suit && c.Rank == card.Rank {
			hand[i].Playable = false
			round.Tricks[round.TrickNumber].Add(card)
			break
		}
	}
	log.Infof("server: move succeeded: trick: %v", *trick)

	if trick.Size == NPlayers { // update results
		winner := trick.Winner()
		round.Scores[winner] += 1
		round.TrickNumber += 1
		if round.TrickNumber < NTricks {
			round.Tricks[round.TrickNumber].Lead = winner
		}
	}

	// next round
	if round.TrickNumber == NTricks {
		game.RoundNumber += 1
		if game.RoundNumber < NRounds {
			round := &game.Rounds[game.RoundNumber]
			round.deal()
			round.Tricks[round.TrickNumber].Lead = game.RoundNumber % NPlayers
		} else {
			log.Infof("Round Scores: %v", round.Scores)
			game.Stage = DONE
		}
	}

	return nil
}

// returns the turn in the game of the given token.
// [0-3] on success, NPlayers on failure
func (game *CallBreak) Turn(token *Token) int {
	var i int
	for i = range game.Players {
		if game.Players[i].Token == *token {
			return i
		}
	}
	return i
}

func (game *CallBreak) Next() int {
	round := game.Rounds[game.RoundNumber]
	trick := round.Tricks[round.TrickNumber]
	return (trick.Lead + trick.Size) % NPlayers
}

func (game *CallBreak) Call(token Token, call Score) error {
	game.workPermit <- struct{}{}
	defer func() { <-game.workPermit }()

	turn := game.Turn(&token)
	if turn == NPlayers {
		return fmt.Errorf("invalid token")
	}

	if game.Stage != DEALT {
		return fmt.Errorf("not a valid move at current stage")
	}

	next := game.Next()
	if next != turn {
		return fmt.Errorf("not your turn")
	}

	game.Rounds[game.RoundNumber].Calls[turn] = call
	return nil
}
