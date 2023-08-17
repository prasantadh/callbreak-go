package callbreak

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

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

func (game *CallBreak) Query(token Token) CallBreak {
	return *game
}

// add a player to the game. returns an authentication token on success
// else return error on failure
func (game *CallBreak) AddPlayer(name string, strategy string) (Id, error) {

	game.workPermit <- struct{}{}
	defer func() { <-game.workPermit }()

	if game.TotalPlayers == NPlayers {
		return Id{}, fmt.Errorf("could not add players: table already full")
	}

	player := &game.Players[game.TotalPlayers]
	player.Name = name

	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return Id{}, fmt.Errorf("could not generate a token")
	}
	player.Token = Token(hex.EncodeToString(token))

	s, err := GetStrategy(strategy)
	if err != nil {
		return Id{}, fmt.Errorf("could not set strategy: %v", err)
	}
	player.Strategy = s

	log.Infof("add player %s with token %s", game.Players[game.TotalPlayers].Name, player.Token)
	game.TotalPlayers += 1

	if game.TotalPlayers == NPlayers {
		game.Rounds[game.RoundNumber].deal()
		// TODO: implement called
		game.Stage = CALLED
	}

	return player.Id, nil
}

// deal the cards to the players
// each player can now make a call to GetHand
// TODO: auto trigger this action when round starts
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
	log.Infof("server: player %s attempted play with %s", token, card)

	player := -1
	for i, p := range game.Players {
		if p.Token == token {
			player = i
			break
		}
	}
	if player == -1 {
		return fmt.Errorf("cannot play: invalid token")
	}

	if game.RoundNumber == NRounds {
		return fmt.Errorf("game over")
	}

	if game.Stage != CALLED {
		return fmt.Errorf("cannot play: not all players have called")
	}

	round := &game.Rounds[game.RoundNumber]
	// the TrickNumber should always be valid for current game
	// or the game is in inconsistent state
	trick := &round.Tricks[round.TrickNumber]
	next := (trick.Lead + trick.Size) % NPlayers

	// TODO: eventually move this "server:" as a logger field
	log.Infof("server: RoundNumber: %d\tTrickNumber: %d",
		game.RoundNumber, round.TrickNumber)
	log.Infof("server: trick: %s (size: %d lead: %d)",
		trick.Cards, trick.Size, trick.Lead)
	log.Infof("server: Hand: %s", round.Hands[player])

	if player != next {
		return fmt.Errorf("you are not up next")
	}

	validMoves, err := GetValidMoves(*game)
	log.Infof("valid moves: %s", validMoves)
	if err != nil {
		return fmt.Errorf("could not get valid moves: %v", err)
	}
	if !slices.Contains(validMoves, card) {
		return fmt.Errorf("invalid move from player")
	}

	hand := &round.Hands[player]
	for i, c := range hand { // play the card
		if c.Suit == card.Suit && c.Rank == card.Rank {
			hand[i].Playable = false
			round.Tricks[round.TrickNumber].Add(card)
			break
		}
	}

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
		}
	}

	return nil
}

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

func (game *CallBreak) Call(token Token, call Call) {
	return
}
