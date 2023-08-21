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

type PlayerConfig struct {
	Name     string        `json:"name"`
	Strategy string        `json:"strategy"`
	Timeout  time.Duration `json:"timeout"`
}

type Config struct {
	Debug   bool `json:"debug,omitempty"`
	Players []PlayerConfig
}

func New() *CallBreak {
	game := &CallBreak{
		workPermit: make(chan struct{}, 1),
	}
	go game.updateStage()
	return game
}

// Get the current state of the game as visible to a token
func (game *CallBreak) Query(token Token) (*CallBreak, error) {
	// todo: currently returns tricks for all players.
	// consider sending out only the tricks won by current player
	game.workPermit <- struct{}{}
	defer func() { <-game.workPermit }()

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
func (game *CallBreak) AddPlayer(config PlayerConfig) (PlayerId, error) {

	game.workPermit <- struct{}{}
	defer func() { <-game.workPermit }()

	if game.TotalPlayers == NPlayers {
		return PlayerId{}, fmt.Errorf("could not add players: table full")
	}

	s, err := GetStrategy(config.Strategy)
	if err != nil {
		return PlayerId{}, fmt.Errorf("could not set strategy: %v", err)
	}

	buffer := make([]byte, 32)
	_, err = rand.Read(buffer)
	if err != nil {
		return PlayerId{}, fmt.Errorf("could not generate a token")
	}
	token := Token(hex.EncodeToString(buffer))
	playerid := PlayerId{Name: config.Name, Token: token}

	// TODO eventually add a registry for assistant
	// this works for now
	assistant := Assistant{strategy: s, game: game, token: token}
	if config.Timeout < AssistMinTimeout {
		return PlayerId{}, fmt.Errorf("timeout must be more than %d", AssistMinTimeout)
	}
	if config.Timeout > AssistMaxTimeout {
		return PlayerId{}, fmt.Errorf("timeout must be less than %d", AssistMaxTimeout)
	}
	assistant.ticker = time.NewTicker(config.Timeout)
	go assistant.Assist()
	game.Players[game.TotalPlayers] = playerid

	log.Infof("add player: %s", playerid)
	log.Infof("\tassistant timeout: %s", config.Timeout)
	game.TotalPlayers += 1

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

// fires periodically to update the stage of the game
func (game *CallBreak) updateStage() {
	// TODO This might need to be configurable for the speed of the game
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		game.workPermit <- struct{}{}
		if game.Stage == DONE {
			break
		}

		round := &game.Rounds[game.RoundNumber]
		trick := &round.Tricks[round.TrickNumber]
		switch game.Stage {
		case NOTFULL:
			if game.TotalPlayers == NPlayers {
				round.deal()
				game.Stage = DEALT
			}

		case DEALT:
			if trick.Size == NPlayers {
				trick.Size = 0
				game.Stage = CALLED
			}

		case CALLED:
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
					game.Stage = DEALT
					round.Tricks[round.TrickNumber].Lead = game.RoundNumber % NPlayers
				} else {
					log.Infof("Round Scores: %v", round.Scores)
					game.Stage = DONE
				}
			}
		}
		<-game.workPermit
	}
}

// the next player in line playes the card c
func (game *CallBreak) Break(token Token, card deck.Card) error {

	game.workPermit <- struct{}{}
	defer func() { <-game.workPermit }()

	log.Infof("server: player %s attempted play with %s", token[:4], card)
	player := game.Turn(&token)
	if player == NPlayers {
		return fmt.Errorf("invalid token")
	}

	if game.Stage != CALLED {
		return fmt.Errorf("not a valid stage for this move")
	}

	round := &game.Rounds[game.RoundNumber]
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
	log.Infof("server: Calls: %v", round.Calls)

	validMoves, err := GetValidMoves(game)
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
	if trick.Size == NPlayers {
		return NPlayers
	}
	return (trick.Lead + trick.Size) % NPlayers
}

func (game *CallBreak) Call(token Token, call Score) error {
	game.workPermit <- struct{}{}
	defer func() { <-game.workPermit }()

	round := &game.Rounds[game.RoundNumber]
	trick := &round.Tricks[round.TrickNumber]

	if game.Stage != DEALT {
		return fmt.Errorf("not a valid move at current stage")
	}

	next := game.Next()
	turn := game.Turn(&token)
	if next != turn {
		return fmt.Errorf("not your turn")
	}

	if call < 1 || call > 8 {
		return fmt.Errorf("calls must be between 1 and 8")
	}

	round.Calls[next] = call
	trick.Size += 1

	return nil
}
