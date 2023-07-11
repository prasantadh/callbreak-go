package cli

import (
	"fmt"
	"time"

	"github.com/prasantadh/callbreak-go/pkg/basicrenderer"
	"github.com/prasantadh/callbreak-go/pkg/bot"
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/spf13/cobra"
)

var singleCommand = &cobra.Command{
	Use:   "single",
	Short: "Play in single player mode against three bots",
	Long:  `This subcommand allows you to play against the implemented bots.`,
	Run: func(cmd *cobra.Command, args []string) {
		runSingle(cmd, args)
	},
}

func init() {
	rootCommand.AddCommand(singleCommand)
}

// TODO: take input from the user and play the card accordingly

func runSingle(cmd *cobra.Command, args []string) {

	game := callbreak.New()
	renderer := basicrenderer.New()

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
	for i := 0; i < callbreak.NCards; i++ {
		game.Update()
		trick := game.CurrentTrick()
		player := bots[game.NextPlayer]
		c, _ := player.Play(trick)
		err := game.Play(c)
		if err != nil {
			msg := fmt.Errorf("invalid move from a player: %v", err)
			panic(msg)
		}
		renderer.Render(game)
		time.Sleep(time.Millisecond * 500)
	}
	game.Update()
	renderer.Render(game)
}
