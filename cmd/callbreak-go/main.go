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

	// play the cards
	fmt.Printf("The game:\n")
	// for now just check by playing 5 tricks
	for i := 0; i < callbreak.NTricks; i++ {
		trick := callbreak.Trick{}
		for _, player := range bots {
			fmt.Println(player.Hand)
			c, _ := player.Play(trick)
			trick.Cards = append(trick.Cards, c)
			trick.Lead = trick.Cards[0]
			err := game.Play(c)
			if err != nil {
				msg := fmt.Errorf("invalid move from a player")
				panic(msg)
			}
		}
		for _, c := range trick.Cards {
			fmt.Printf("%c %s   ", c.Suit, c.Rank)
		}
		fmt.Println()
		fmt.Println()
	}

}
