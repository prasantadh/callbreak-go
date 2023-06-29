package callbreak

import (
	"sort"
	"strings"

	"github.com/prasantadh/callbreak-go/pkg/deck"
)

// sort by suit, Hukum -> itta -> chidi -> paan
// for the same suit, sort by rank descending
func (h *Hand) Sort() {
	sort.Slice(*h, func(i, j int) bool {
		this := (*h)[i]
		other := (*h)[j]
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
		sb.WriteString("[")
		sb.WriteString(c.Suit.String())
		sb.WriteString(" ")
		sb.WriteString(c.Rank.String())
		sb.WriteString(" ")
		if c.Playable {
			sb.WriteString("✓")
		} else {
			sb.WriteString(" ")
		}
		sb.WriteString("]")
	}
	return sb.String()
}
