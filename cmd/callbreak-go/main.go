package main

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/bot"
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
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
	game.CollectDeck()
	fmt.Printf("The Deck: ")
	for _, c := range game.Deck {
		fmt.Printf("[%s %d]", string(c.Suit), c.Rank)
	}
	fmt.Println()
	// fmt.Println(game.Deck, len(game.Deck), deck.Size)
	for i := 0; i < deck.Size; i++ {
		// fmt.Printf("dealing card %d\n", i)
		c, _ := game.Deal()
		bots[i%callbreak.NPlayers].Take(c)
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
		trick := callbreak.Trick{}
		for _, player := range bots {
			c := player.Play(trick)
			fmt.Printf("%s plays [%s %s]\n", player.Name, string(c.Suit), c.Rank)
			err := game.Play(c)
			if err != nil {
				msg := fmt.Errorf("invalid move from a player")
				panic(msg)
			}
		}
	}

}
