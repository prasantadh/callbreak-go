package callbreak

import "github.com/prasantadh/callbreak-go/pkg/deck"

func (p *player) Take(c deck.Card) {
	p.hand = append(p.hand, c)
}
