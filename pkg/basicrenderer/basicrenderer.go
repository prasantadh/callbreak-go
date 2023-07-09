package basicrenderer

import (
	"strings"

	"atomicgo.dev/cursor"
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
)

var (
	// TODO: this order will change when round is implemented
	left   = 1
	top    = 2
	right  = 3
	bottom = 0
)

type Renderer struct {
	area cursor.Area
}

func New() *Renderer {
	return &Renderer{
		area: cursor.NewArea(),
	}
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

	s := []byte{}
	addline := func(n int) {
		s = append(s, hands[left][n].String()...)
		s = append(s, strings.Repeat(" ", 7*4)...)
		if n == 6 {
			s = append(s, trick[left].String()...)
		} else {
			s = append(s, strings.Repeat(" ", 7)...)
		}
		s = append(s, strings.Repeat(" ", 7)...)
		if n == 4 {
			s = append(s, trick[top].String()...)
		} else if n == 8 {
			s = append(s, trick[bottom].String()...)
		} else {
			s = append(s, strings.Repeat(" ", 7)...)
		}
		s = append(s, strings.Repeat(" ", 7)...)
		if n == 6 {
			s = append(s, trick[right].String()...)
		} else {
			s = append(s, strings.Repeat(" ", 7)...)
		}
		s = append(s, strings.Repeat(" ", 7*4)...)
		s = append(s, hands[right][n].String()...)
		s = append(s, '\n')
	}

	s = append(s, strings.Repeat("-", 7)...)
	s = append(s, hands[top].String()...)
	s = append(s, strings.Repeat("-", 7)...)
	s = append(s, '\n')
	for i := 0; i < 13; i++ {
		addline(i)
	}
	s = append(s, strings.Repeat("-", 7)...)
	s = append(s, hands[bottom].String()...)
	s = append(s, strings.Repeat("-", 7)...)
	s = append(s, '\n')

	r.area.Update(string(s))
}
