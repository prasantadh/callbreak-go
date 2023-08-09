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
	Players      [NPlayers]Player
	Rounds       [NRounds]Round
	Stage        Stage
	TotalPlayers int // number of players currently in the game
	RoundNumber  int // current round number
	// TotalPlayers and RoundNumber might be better as names
	workPermit chan struct{}
}

type Player struct {
	Name  string
	token Token
}

type Token string

type Round struct {
	Calls       [NPlayers]Score
	Breaks      [NPlayers]Score
	Hands       [NPlayers]Hand
	Tricks      [NTricks]Trick
	TrickNumber int // current Trick number
	// TrickNumber might be a better name here
}

type Score int

type Hand [NTricks]deck.Card

type Trick struct {
	Cards [NPlayers]deck.Card
	Lead  int // the card position that is 1st in this trick
	Size  int // number of cards played so far in this trick
}

// track the current state of the game
type Stage int

const (
	NOTREADY Stage = iota
	DEALT
	CALLED
	DONE
)
