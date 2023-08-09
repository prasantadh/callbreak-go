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
var FaintString func(a ...interface{}) string
var BgWhiteString func(a ...interface{}) string
var UnderlinedString func(a ...interface{}) string

func init() {
	RedString = color.New(color.BgWhite).Add(color.FgRed).SprintFunc()
	FaintString = color.New(color.Faint).Add(color.CrossedOut).SprintFunc()
	BlackString = color.New(color.BgWhite).Add(color.FgBlack).SprintFunc()
	BgWhiteString = color.New(color.BgWhite).Add(color.FgBlack).SprintFunc()
	UnderlinedString = color.New(color.Underline).SprintFunc()
}

type Renderer struct {
	area  cursor.Area
	token callbreak.Token
}

func New() *Renderer {
	return &Renderer{
		area: cursor.NewArea(),
	}
}

func (r *Renderer) SetToken(token callbreak.Token) {
	r.token = token
}

func blank(repetition int) string {
	return BgWhiteString(strings.Repeat(" ", 6*repetition))
}

func ColoredCard(c deck.Card) string {
	if c.Rank == 0 {
		return blank(1)
	}
	s := fmt.Sprintf("[%s %s ]", &c.Suit, &c.Rank)
	if c.Suit == deck.Hukum || c.Suit == deck.Chidi {
		if c.Playable {
			return BlackString(s)
		}
		return FaintString(BlackString(s))
	}
	if c.Playable {
		return RedString(s)
	}
	return FaintString(RedString(s))
}

func ColoredHand(h callbreak.Hand) string {
	sb := strings.Builder{}
	for _, c := range h {
		sb.WriteString(ColoredCard(c))
	}
	return sb.String()
}

func (r *Renderer) Render(g *callbreak.CallBreak) {
	if r.token == "" {
		fmt.Println("token not yet set!")
		return
	}

	// TODO: eventually only get one hand and display that
	// for now display all hands
	state, _ := g.Query(r.token)
	round := state.Rounds[state.RoundNumber]
	hands := round.Hands
	trick := round.Tricks[round.TrickNumber].Cards

	for i := range trick { // little hack to make trick not-CrossedOut
		trick[i].Playable = true
	}

	r.area.Clear()

	sb := strings.Builder{}
	addline := func(n int) {
		sb.WriteString(ColoredCard(hands[left][n]))
		sb.WriteString(blank(4))
		if n == 6 {
			// sb.WriteString(ColoredCard(trick[left]))
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

	// the display content
	sb.WriteString(BgWhiteString("ScoreBoard  "))
	sb.WriteString(blank(13))
	sb.WriteString("\n")
	sb.WriteString(BgWhiteString("Bots:  | bot0 | bot1 | bot2 | bot3 |"))
	sb.WriteString(blank(9))
	sb.WriteString("\n")
	sb.WriteString(blank(1))
	sb.WriteString(BgWhiteString(" |"))
	// printing the scoreboard values
	for _, score := range round.Breaks {
		sb.WriteString(BgWhiteString(fmt.Sprintf(" %2d/_ |", score)))
	}
	// the game board
	sb.WriteString(blank(9))
	sb.WriteString(BgWhiteString("\n"))
	sb.WriteString(BgWhiteString("bot2⮕ "))
	sb.WriteString(ColoredHand(hands[top]))
	sb.WriteString(BgWhiteString("⬇ bot3"))
	sb.WriteString("\n")
	for i := 0; i < 13; i++ {
		addline(i)
	}
	sb.WriteString(BgWhiteString("bot1⬆ "))
	sb.WriteString(ColoredHand(hands[bottom]))
	sb.WriteString(BgWhiteString("⬅ bot0"))
	sb.WriteString("\n")
	r.area.Update(sb.String())
}
