package player

import (
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

// a Player is expected to get the hand to Play from the game
type Player interface {
	Name() string
	SetName(string) error
	Token() callbreak.Token
	SetToken(callbreak.Token)
	Play(*callbreak.CallBreak) (deck.Card, error)
}
