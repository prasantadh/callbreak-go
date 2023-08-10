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
	// and if array makes sense here and in used types recursively
	Players      [NPlayers]Player `json:"players"`
	Rounds       [NRounds]Round   `json:"rounds"`
	Stage        Stage            `json:"stage"`
	TotalPlayers int              `json:"totalplayers"`
	RoundNumber  int              `json:"roundnumber"`
	// TotalPlayers and RoundNumber might be better as names
	workPermit chan struct{}
	debug      bool
	Input      chan any
	Update     chan struct{}
}

type Player struct {
	Id
	Strategy
	Client
	AutoPlay bool
}

type Id struct {
	Name  string `json:"name"`
	Token `json:"token"`
}

type Token string

type Strategy interface {
	Call(CallBreak) (Call, error)
	Break(CallBreak) (deck.Card, error)
}

type Client interface {
	Update()
	GetStrategy() Strategy
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
