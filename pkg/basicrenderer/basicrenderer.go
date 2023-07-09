package basicrenderer

import (
	"fmt"
	"strings"

	"atomicgo.dev/cursor"
	"github.com/fatih/color"
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

var (
	// TODO: this order will change when round is implemented
	left   = 1
	top    = 2
	right  = 3
	bottom = 0
)

var RedString func(a ...interface{}) string
var BlackString func(a ...interface{}) string
var BgWhiteString func(a ...interface{}) string
var UnderlinedString func(a ...interface{}) string

func init() {
	RedString = color.New(color.BgWhite).Add(color.FgRed).SprintFunc()
	BlackString = color.New(color.BgWhite).Add(color.FgBlack).SprintFunc()
	BgWhiteString = color.New(color.BgWhite).Add(color.FgBlack).SprintFunc()
	UnderlinedString = color.New(color.BgWhite).Add(color.Underline).SprintFunc()
}

type Renderer struct {
	area cursor.Area
}

func New() *Renderer {
	return &Renderer{
		area: cursor.NewArea(),
	}
}

func blank(repetition int) string {
	return BgWhiteString(strings.Repeat(" ", 7*repetition))
}

func ColoredCard(c deck.Card) string {
	if c.Rank == 0 {
		return blank(1)
	}
	s := fmt.Sprintf("[%s %s ", &c.Suit, &c.Rank)
	if c.Playable {
		s += "✓]"
	} else {
		s += " ]"
	}
	if c.Suit == deck.Hukum || c.Suit == deck.Chidi {
		return BlackString(s)
	}
	return RedString(s)
}

func ColoredHand(h callbreak.Hand) string {
	sb := strings.Builder{}
	for _, c := range h {
		sb.WriteString(ColoredCard(c))
	}
	return sb.String()
}

func (r *Renderer) Render(g *callbreak.CallBreak) {
	// TODO: eventually only get one hand and display that
	// for now display all hands
	hands := []callbreak.Hand{}
	for i := 0; i < callbreak.NPlayers; i++ {
		hands = append(hands, g.GetHand(i))
	}
	trick := g.CurrentTrick().Cards

	r.area.Clear()

	sb := strings.Builder{}
	addline := func(n int) {
		sb.WriteString(ColoredCard(hands[left][n]))
		sb.WriteString(blank(4))
		if n == 6 {
			sb.WriteString(ColoredCard(trick[left]))
		} else {
			sb.WriteString(blank(1))
		}
		sb.WriteString(blank(1))
		if n == 4 {
			sb.WriteString(ColoredCard(trick[top]))
		} else if n == 8 {
			sb.WriteString(ColoredCard(trick[bottom]))
		} else {
			sb.WriteString(blank(1))
		}
		sb.WriteString(blank(1))
		if n == 6 {
			sb.WriteString(ColoredCard(trick[right]))
		} else {
			sb.WriteString(blank(1))
		}
		sb.WriteString(blank(4))
		sb.WriteString(ColoredCard(hands[right][n]))
		sb.WriteString("\n")
	}

	// corner := BgWhiteString(strings.Repeat("-", 7))
	sb.WriteString(BgWhiteString("///////"))
	sb.WriteString(ColoredHand(hands[top]))
	sb.WriteString(BgWhiteString("\\\\\\\\\\\\\\"))
	sb.WriteString("\n")
	for i := 0; i < 13; i++ {
		addline(i)
	}
	sb.WriteString(BgWhiteString("\\\\\\\\\\\\\\"))
	sb.WriteString(ColoredHand(hands[bottom]))
	sb.WriteString(BgWhiteString("///////"))
	sb.WriteString("\n")
	r.area.Update(sb.String())
}
