package callbreak

import "github.com/prasantadh/callbreak-go/pkg/deck"

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

// a player must be able to take a card dealt in the game
// and be able to play a card when it's their turn
type PlayerInterface interface {
	Play(Trick) (deck.Card, error)
	Take(deck.Card) error
}

type BotInterface PlayerInterface

type Score int

type Token string

// type CallBreak struct { // a game is
// players    []player // between 4 players
// deck.Deck           // with a standard deck of cards
// tricks     []Trick  // and is played as a collection of Tricks
// NextPlayer int
// }

// Game holds the data for :TODO an authoized player.
// in terms of data organization, there is a conflict on spatial locality.
// [Rounds][Player] would offer more spatial locality when accessing
// round-first information like a scoreboard
// [Player][Round] would offer more spatial locality when accessing
// player-first information like rendering tricks of current player

// the data that should be visible to anyone playing the game
// the players are in the order of who started the first round
// each trick also has cards in the same order
// TODO: check autorization before sending results
type Game struct {
	Players []struct {
		Name   string `json:"name"`
		Rounds []struct {
			Call   Score   // Call made by the player at the beginning of round
			Break  Score   // Number of tricks won by the player this round
			Hand           // Hand of the player TODO: empty if not authoized
			Tricks []Trick // tricks won by the player requesting data TODO authorize
		}
	}
}

type CallBreak struct {
	// TODO: think about have a token that provides access to all data in game
	players []player
	rounds  []round
	next    int // player that goes next
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
	// there are 13 tricks in a round. each trick is indexed using player
	// ex. tricks[0] will have card played by [player0, ... , player3]
	tricks []Trick
}

type State int

const (
	NOT_DEALT = iota
)

var States [1]State = [...]State{NOT_DEALT}
