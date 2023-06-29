package main

import (
	"fmt"

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

	// print the info
	for _, b := range bots {
		fmt.Printf("%s has cards: ", b.Name)
		for _, c := range b.Hand {
			fmt.Printf("[%s %s] ", string(c.Suit), c.Rank)
		}
		fmt.Println()
	}

	// play the cards
	fmt.Printf("The game:\n")
	for i := 0; i < callbreak.NTricks; i++ {
		for _, player := range bots {
			c, _ := player.Play(game.CurrentTrick)
			fmt.Printf("%s plays [%s %s]\n", player.Name, string(c.Suit), c.Rank)
			err := game.Play(c)
			if err != nil {
				msg := fmt.Errorf("invalid move from a player")
				panic(msg)
			}
		}
	}

}
