package hand

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/card"
)

type Hand struct {
	cards map[card.Card]bool // card: true if playable, false otherwise
}

func New() *Hand {
	return &Hand{}
}

func (h *Hand) Add(c card.Card) {
	h.cards[c] = true
}

func (h *Hand) Play(c card.Card) error {
	_, ok := h.cards[c]
	if !ok {
		return fmt.Errorf("the card is not in your hand")
	}
	if !h.cards[c] {
		return fmt.Errorf("the card has already been played")
	}
	h.cards[c] = false
	return nil
}
