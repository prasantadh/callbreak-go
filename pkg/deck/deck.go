package deck

import (
	"github.com/prasantadh/callbreak-go/pkg/card"
	"github.com/prasantadh/callbreak-go/pkg/rank"
	"github.com/prasantadh/callbreak-go/pkg/suit"
)

type Deck struct {
	cards []card.Card
}

func New() *Deck {
	d := &Deck{}
	for _, r := range rank.All {
		for _, s := range suit.All {
			d.cards = append(d.cards, card.Card{Rank: r, Suit: s})
		}
	}
	return d
}
