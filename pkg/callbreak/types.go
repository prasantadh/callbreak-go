package callbreak

import "github.com/prasantadh/callbreak-go/pkg/deck"

const (
	NCards   = 52
	NPlayers = 4
	NTricks  = NCards / NPlayers
)

type Hand []deck.Card

type Trick struct {
	Cards [NTricks]deck.Card
	Lead  deck.Card
	Size  int
}

// a player must be able to take a card dealt in the game
// and be able to play a card when it's their turn
type PlayerInterface interface {
	Take(deck.Card) error
	Play(Trick) deck.Card
}

type BotInterface PlayerInterface
