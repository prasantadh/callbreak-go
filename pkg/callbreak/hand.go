package callbreak

import (
	"sort"
	"strings"

	"github.com/prasantadh/callbreak-go/pkg/deck"
)

// TODO: see if it is better to provide a HandInterface
// so players can have their own Sort() and String() function

// sort by suit, Hukum -> itta -> chidi -> paan
// for the same suit, sort by rank descending
func (h *Hand) Sort() {
	sort.Slice(h[:], func(i, j int) bool {
		this := h[i]
		other := h[j]
		if this.Suit == other.Suit {
			return this.Rank > other.Rank
		}

		// TODO: check if the else ifs are necessary
		// or skipping them has the same results
		if this.Suit == deck.Hukum {
			return true
		} else if other.Suit == deck.Hukum {
			return false
		}

		if this.Suit == deck.Itta {
			return true
		} else if other.Suit == deck.Itta {
			return false
		}

		if this.Suit == deck.Chidi {
			return true
		} else if other.Suit == deck.Chidi {
			return false
		}

		return false
	})

}

func (h Hand) String() string {
	var sb strings.Builder
	for _, c := range h {
		sb.WriteString(c.String())
	}
	return sb.String()
}

func (h Hand) HasPlayable(card deck.Card) bool {
	for _, c := range h {
		if c.Suit == card.Suit && c.Rank == card.Rank {
			if c.Playable {
				return true
			}
			return false
		}
	}
	return false
}

func (h Hand) HasSuit(suit deck.Suit) bool {
	for _, c := range h {
		if c.Playable == true && c.Suit == suit {
			return true
		}
	}
	return false
}
