package player

import (
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
)

// a Player is expected to get the hand to Play from the game
type Player interface {
	Name() string
	Token() callbreak.Token
	Play(*callbreak.CallBreak)
}
