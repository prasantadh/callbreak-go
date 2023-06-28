package callbreak

import "github.com/prasantadh/callbreak-go/pkg/deck"

const (
	NCards   = 52
	NPlayers = 4
	NTricks  = NCards / NPlayers
)

type Hand [NTricks]deck.Card

type Trick struct {
	Cards [NTricks]deck.Card
	Lead  deck.Card
	Size  int
}

// a player bot needs to provide a Play function
// that will take the current trick and return the card to play
type PlayerInterface interface {
	Play(Trick) deck.Card
}

type BotInterface PlayerInterface
