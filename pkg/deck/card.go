package deck

import "strings"

type Card struct {
	Suit
	Rank
	Playable bool
}

func (c Card) String() string {
	var sb strings.Builder
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
	return sb.String()
}
