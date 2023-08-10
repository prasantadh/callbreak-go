package deck

import "strings"

type Card struct {
	Suit
	Rank
	Playable bool
}

func (c Card) String() string {
	if c.Rank == 0 {
		return "(-----)"
	}
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(c.Suit.String())
	sb.WriteString(" ")
	sb.WriteString(c.Rank.String())
	sb.WriteString(" ")
	if c.Playable {
		sb.WriteString("✓")
	} else {
		sb.WriteString(" ")
	}
	sb.WriteString(")")
	return sb.String()
}

func (c *Card) IsValid() bool {
	for _, suit := range Suits {
		for _, rank := range Ranks {
			if c.Suit == suit && c.Rank == rank {
				return true
			}
		}
	}
	return false
}
