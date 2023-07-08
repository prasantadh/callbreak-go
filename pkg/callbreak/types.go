package callbreak

import "github.com/prasantadh/callbreak-go/pkg/deck"

const (
	NCards   = 52
	NPlayers = 4
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

type player struct {
	hand Hand
	PlayerInterface
}
