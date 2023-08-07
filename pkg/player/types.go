package player

import (
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

// a Player is expected to get the hand to Play from the game
type Player interface {
	SetName(string) error
	GetGameState()
	Play() (deck.Card, error)
}
