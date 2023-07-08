package callbreak

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/deck"
)

func (p *player) Take(c deck.Card) {
	p.hand = append(p.hand, c)
}

func (p *player) Play(c deck.Card) error {
	for i, card := range p.hand {
		if card.Suit == c.Suit && card.Rank == c.Rank {
			p.hand[i].Playable = false
			return nil
		}
	}
	return fmt.Errorf("card not in player's hand")
}
