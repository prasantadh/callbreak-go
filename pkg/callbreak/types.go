package callbreak

import (
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

const (
	NCards   = 52
	NPlayers = 4
	NRounds  = 5
	NTricks  = NCards / NPlayers
)

// CallBreak implements the callbreak game
type CallBreak struct {
	// TODO: think about have a token that provides access to all data in game
	Players      [NPlayers]PlayerId `json:"players"`
	Rounds       [NRounds]Round     `json:"rounds"`
	Stage        Stage              `json:"stage"`
	TotalPlayers int                `json:"totalplayers"`
	RoundNumber  int                `json:"roundnumber"`
	// TotalPlayers and RoundNumber might be better as names
	workPermit chan struct{}
	debug      bool
}

type Assistant interface {
	SetToken(Token)
	SetStrategy(Strategy)
	Assist()
}

type PlayerId struct {
	Name  string `json:"name"`
	Token `json:"token"`
}

type Token string

type Strategy interface {
	Call(CallBreak) (Call, error)
	Break(CallBreak) (deck.Card, error)
}

type Round struct {
	Calls       [NPlayers]Score
	Scores      [NPlayers]Score
	Hands       [NPlayers]Hand
	Tricks      [NTricks]Trick
	TrickNumber int // current Trick number
}

type Score int
type Call Score

type Hand [NTricks]deck.Card

type Trick struct {
	Cards [NPlayers]deck.Card
	Lead  int // the card position that is 1st in this trick
	Size  int // number of cards played so far in this trick
}

// track the current state of the game
type Stage int

const (
	NOTFULL Stage = iota
	DEALT
	CALLED
	DONE
)
