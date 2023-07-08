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
			trick.Cards[len(trick.Cards)-1].Playable = false
			// TODO: render here after each player plays
		}

		// rendering.
		// TODO trick cards leans in the direction of player who played it
		// // upper half
		fmt.Printf("-------%s-------\n", bots[2].Hand)
		for i := 0; i < 4; i++ {
			fmt.Printf("%s", bots[1].Hand[i])
			fmt.Printf("%s", strings.Repeat(" ", 7*13))
			fmt.Printf("%s\n", bots[3].Hand[i])
		}

		// top player
		fmt.Printf("%s", bots[1].Hand[4])
		fmt.Printf("%s", strings.Repeat(" ", 7*6))
		fmt.Printf("%s", trick.Cards[2])
		fmt.Printf("%s", strings.Repeat(" ", 7*6))
		fmt.Printf("%s", bots[1].Hand[4])

		// gap
		fmt.Printf("%s", bots[1].Hand[5])
		fmt.Printf("%s", strings.Repeat(" ", 7*13))
		fmt.Printf("%s\n", bots[3].Hand[5])

		// side players
		fmt.Printf("%s", bots[1].Hand[6])
		fmt.Printf("%s", strings.Repeat(" ", 7*4))
		fmt.Printf("%s", trick.Cards[1])
		fmt.Printf("%s", strings.Repeat(" ", 7*3))
		fmt.Printf("%s", trick.Cards[3])
		fmt.Printf("%s", strings.Repeat(" ", 7*4))
		fmt.Printf("%s", bots[1].Hand[6])

		// gap
		fmt.Printf("%s", bots[1].Hand[7])
		fmt.Printf("%s", strings.Repeat(" ", 7*13))
		fmt.Printf("%s\n", bots[3].Hand[7])

		// bottom player
		fmt.Printf("%s", bots[1].Hand[8])
		fmt.Printf("%s", strings.Repeat(" ", 7*6))
		fmt.Printf("%s", trick.Cards[0])
		fmt.Printf("%s", strings.Repeat(" ", 7*6))
		fmt.Printf("%s", bots[1].Hand[8])

		// lower half
		for i := 9; i < 13; i++ {
			fmt.Printf("%s", bots[1].Hand[i])
			fmt.Printf("%s", strings.Repeat(" ", 7*13))
			fmt.Printf("%s\n", bots[3].Hand[i])
		}
		fmt.Printf("-------%s-------\n", bots[0].Hand)

		fmt.Println()
		fmt.Println()
	}

}
