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

type Hand []deck.Card

type Trick struct {
	Cards [NPlayers]deck.Card
	Lead  int // the card position that is 1st in this trick
	Size  int // number of cards played so far in this trick
}

type Score int

type Token string

// the data that should be visible to anyone playing the game
// the players are in the order of who started the first round
// each trick also has cards in the same order
// TODO: check autorization before sending results
type Game struct {
	// in terms of data organization, there is a conflict on spatial locality.
	// [Rounds][Player] would offer more spatial locality when accessing
	// round-first information like a scoreboard
	// [Player][Round] would offer more spatial locality when accessing
	// player-first information like rendering tricks of current player
	Players []Player `json:"players"`
	Next    int      `json:"next"`
}

// information available to a player
type Player struct {
	Name   string  `json:"name"`
	Rounds []Round `json:"rounds"`
}

// information about a round available to a player
// Tricks only holds the tricks that are won by this player
type Round struct {
	Call   Score `json:"call"`
	Break  Score `json:"score"`
	Hand   `json:"hand"`
	Tricks []Trick `json:"tricks"`
}

type CallBreak struct {
	// TODO: think about have a token that provides access to all data in game
	players []player
	rounds  []round
	next    int // player that goes next
	state   State
}

type player struct {
	name  string
	token Token
}

type round struct {
	// the following fields are indexed using player as an index
	// ex. calls will have call for [player0, ... , player3]
	calls  [NPlayers]Score
	breaks [NPlayers]Score
	hands  [NTricks]Hand
	// tricks are indexed [0,13). EACH trick is indexed using player
	// ex. tricks[0] will have card played by [player0, ... , player3]
	tricks []Trick
}

type current struct {
	player *player
	hand   *Hand
	call   *Score
	score  *Score
	trick  *Trick
}

// track the current state of the game
type State int

const (
	NOTREADY State = iota
	DEALT
	CALLED
	DONE
)
