package renderer

import (
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

type Renderer interface {
	Render(<-chan callbreak.CallBreak, int)
	Call() <-chan callbreak.Score
	Break() <-chan deck.Card
}
