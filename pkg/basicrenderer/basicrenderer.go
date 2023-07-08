package basicrenderer

import (
	"fmt"
	"strings"

	"github.com/prasantadh/callbreak-go/pkg/callbreak"
)

var (
	// TODO: this order will change when round is implemented
	left   = 1
	top    = 2
	right  = 3
	bottom = 0
	me     = 0
)

type Renderer struct{}

func New() *Renderer {
	return &Renderer{}
}

func (r *Renderer) Render(g *callbreak.CallBreak) {
	// TODO: eventually only get one hand and display that
	// for now display all hands
	hands := []callbreak.Hand{}
	for i := 0; i < callbreak.NPlayers; i++ {
		hands = append(hands, g.GetHand(i))
	}
	trick := g.CurrentTrick().Cards

	// upper half
	fmt.Printf("-------%s-------\n", hands[top])
	for i := 0; i < 4; i++ {
		fmt.Printf("%s", hands[left][i])
		fmt.Printf("%s", strings.Repeat(" ", 7*13))
		fmt.Printf("%s\n", hands[right][i])
	}

	// top player card for the trick
	fmt.Printf("%s", hands[left][4])
	fmt.Printf("%s", strings.Repeat(" ", 7*6))
	fmt.Printf("%s", trick[top])
	fmt.Printf("%s", strings.Repeat(" ", 7*6))
	fmt.Printf("%s\n", hands[right][4])

	// gap
	fmt.Printf("%s", hands[left][5])
	fmt.Printf("%s", strings.Repeat(" ", 7*13))
	fmt.Printf("%s\n", hands[right][5])

	// left and right player card for the trick
	fmt.Printf("%s", hands[left][6])
	fmt.Printf("%s", strings.Repeat(" ", 7*4))
	fmt.Printf("%s", trick[left])
	fmt.Printf("%s", strings.Repeat(" ", 7*3))
	fmt.Printf("%s", trick[right])
	fmt.Printf("%s", strings.Repeat(" ", 7*4))
	fmt.Printf("%s\n", hands[right][6])

	// gap
	fmt.Printf("%s", hands[left][7])
	fmt.Printf("%s", strings.Repeat(" ", 7*13))
	fmt.Printf("%s\n", hands[right][7])

	// bottom player card for the trick
	fmt.Printf("%s", hands[left][8])
	fmt.Printf("%s", strings.Repeat(" ", 7*6))
	fmt.Printf("%s", trick[0])
	fmt.Printf("%s", strings.Repeat(" ", 7*6))
	fmt.Printf("%s\n", hands[left][8])

	// lower half
	for i := 9; i < 13; i++ {
		fmt.Printf("%s", hands[left][i])
		fmt.Printf("%s", strings.Repeat(" ", 7*13))
		fmt.Printf("%s\n", hands[right][i])
	}
	fmt.Printf("-------%s-------\n", hands[0])

	fmt.Println()
	fmt.Println()
}
