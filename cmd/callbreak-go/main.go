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
	// for now just check by playing 5 tricks
	for i := 0; i < 1; i++ {
		trick := callbreak.Trick{}
		for _, player := range bots {
			c, _ := player.Play(trick)
			fmt.Printf("%s plays [%c %s] to:   ", player.Name, c.Suit, c.Rank)
			for _, c := range trick.Cards {
				fmt.Printf("%c %s   ", c.Suit, c.Rank)
			}
			fmt.Println()
			trick.Cards = append(trick.Cards, c)
			trick.Lead = trick.Cards[0]
			err := game.Play(c)
			if err != nil {
				msg := fmt.Errorf("invalid move from a player")
				panic(msg)
			}
		}
	}

}
