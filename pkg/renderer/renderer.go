package renderer

import (
	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard/keys"
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

func New() Renderer {
	return &Basic{
		area:       cursor.NewArea(),
		keypressed: make(chan keys.KeyCode),
		callC:      make(chan callbreak.Score),
		breakC:     make(chan deck.Card),
	}
}
