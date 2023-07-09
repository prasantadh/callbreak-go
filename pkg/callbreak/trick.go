package callbreak

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/deck"
)

// returns index of the winning card of a trick
// but only considers cards of a given suit
// -1 if no card of the suit exist
func (t Trick) SuitWinner(suit deck.Suit) int {
	winner := struct {
		deck.Card
		index int
	}{Card: deck.Card{Suit: suit, Rank: 0}, index: -1}

	for i := 0; i < NPlayers; i++ {
		this := t.Cards[i]
		if this.Suit == suit && this.Rank > winner.Card.Rank {
			winner.Card = this
			winner.index = i
		}
	}
	return winner.index
}

// returns the index of the winning card in a given trick
func (t Trick) Winner() int {
	// if Hukum is present, return index of highest ranking Hukum
	// else return index of highest ranking card of the opening suit
	for _, c := range t.Cards {
		if c.Suit == deck.Hukum {
			return t.SuitWinner(deck.Hukum)
		}
	}
	return t.SuitWinner(t.Cards[t.Lead].Suit)

}

func (t *Trick) Add(card deck.Card) error {
	if t.Size == NPlayers {
		return fmt.Errorf("trick already full")
	}
	card.Playable = false
	t.Cards[(t.Lead+t.Size)%NPlayers] = card
	t.Size += 1
	return nil
}
