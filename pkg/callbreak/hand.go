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

func (h *Hand) HasPlayable(cards ...deck.Card) bool {
	if len(cards) == 0 {
		for _, c := range h {
			if c.Playable {
				return true
			}
		}
		return false
	}

	for _, c := range cards {
		found := false
		for _, card := range h {
			if card.Playable && c.Suit == card.Suit && c.Rank == card.Rank {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (h *Hand) HasSuit(suit deck.Suit) bool {
	for _, c := range h {
		if c.Playable == true && c.Suit == suit {
			return true
		}
	}
	return false
}

func (h *Hand) IsValid() bool {
	if !h.HasPlayable() {
		return false
	}
	d := deck.New()
	for _, c := range *h {
		found := false
		for _, deckcard := range d {
			if c.Suit != deckcard.Suit && c.Rank != deckcard.Rank {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (h *Hand) Playables() []deck.Card {
	playables := []deck.Card{}
	for _, c := range h {
		if c.Playable {
			playables = append(playables, c)
		}
	}
	return playables
}
