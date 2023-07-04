package main

import (
	"fmt"
	"strings"

	"github.com/prasantadh/callbreak-go/pkg/bot"
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
)

// as a start we are implementing 4 bots playing themselves
// with any one card that they have
func main() {

	game := callbreak.New()

	// add the players
	bots := [callbreak.NPlayers]*bot.Player{}
	for i := 0; i < callbreak.NPlayers; i++ {
		b := bot.New("bot" + fmt.Sprint(i))
		err := game.AddPlayer(b)
		if err != nil {
			msg := fmt.Errorf("failed to setup game: %v", err)
			panic(msg)
		}
		bots[i] = b
	}

	// deal a Deck of cards to players
	game.Deal()
	for i := range bots {
		bots[i].Hand = game.GetHand(i)
	}

	// play the cards
	for i := 0; i < callbreak.NTricks; i++ {

		// // play the trick
		trick := callbreak.Trick{}
		for _, player := range bots {
			c, _ := player.Play(trick)
			trick.Cards = append(trick.Cards, c)
			trick.Lead = trick.Cards[0]
			err := game.Play(c)
			if err != nil {
				msg := fmt.Errorf("invalid move from a player")
				panic(msg)
			}
		}

		// rendering.
		// TODO trick cards leans in the direction of player who played it
		// // upper half
		fmt.Printf("-------%s-------\n", bots[2].Hand)
		for i := 0; i < 6; i++ {
			fmt.Printf("%s", bots[1].Hand[i].String())
			fmt.Printf("%s", strings.Repeat(" ", 7*13))
			fmt.Printf("%s\n", bots[3].Hand[i].String())
		}
		// // the trick
		fmt.Printf("%s", bots[1].Hand[7].String())
		fmt.Printf("%s", strings.Repeat(" ", 7*3))
		for _, c := range trick.Cards {
			// TODO: cards in tricks are being rendered as playable, they shouldn't
			fmt.Printf("%s       ", c.String())
		}
		fmt.Printf("%s", strings.Repeat(" ", 7*2))
		fmt.Printf("%s\n", bots[3].Hand[7].String())
		// // lower half
		for i := 7; i < 13; i++ {
			fmt.Printf("%s", bots[1].Hand[i].String())
			fmt.Printf("%s", strings.Repeat(" ", 7*13))
			fmt.Printf("%s\n", bots[3].Hand[i].String())
		}
		fmt.Printf("-------%s-------\n", bots[0].Hand)
		fmt.Println()
		fmt.Println()
	}

}
