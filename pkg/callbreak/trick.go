package callbreak

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/deck"
)

// returns index of the winning card of a trick
// but only considers cards of a given suit
// -1 if no card of the suit exist
func (t Trick) SuitWinner(suit deck.Suit) int {
	winner := deck.Card{}
	ans := -1

	for i := 0; i < NPlayers; i++ {
		this := t.Cards[i]
		if this.Suit == suit && this.Rank > winner.Rank {
			winner = this
			ans = i
		}
	}
	return ans
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
	return t.SuitWinner(t.LeadSuit())

}

func (t *Trick) Add(card deck.Card) {
	if t.Size == NPlayers {
		panic(fmt.Errorf("trick already full"))
	}
	card.Playable = false
	t.Cards[(t.Lead+t.Size)%NPlayers] = card
	t.Size += 1
}

func (t *Trick) LeadCard() deck.Card {
	return t.Cards[t.Lead]
}

func (t *Trick) LeadSuit() deck.Suit {
	return t.LeadCard().Suit
}

func (t *Trick) LeadRank() deck.Rank {
	return t.LeadCard().Rank
}
