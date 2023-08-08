package cli

import (
	"fmt"
	"time"

	// "github.com/prasantadh/callbreak-go/pkg/basicrenderer"
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/player"
	"github.com/spf13/cobra"
)

var botCommand = &cobra.Command{
	Use:   "bot",
	Short: "Experiment with your custom bot",
	Long: `This subcommand allows you to implement your own calbreak bot
				and compete it against the existing callbreak bots.
				More details available here <TODO: link>`,
	Run: func(cmd *cobra.Command, args []string) {
		runBot(cmd, args)
	},
}

func init() {
	rootCommand.AddCommand(botCommand)
}

// TODO: add a template for an experimental bot
// then provide instructions on what can be changed to
// have your own bot playing in the arena

func runBot(cmd *cobra.Command, args []string) {

	game := callbreak.New()
	// renderer := basicrenderer.New()
	ticker := time.NewTicker(500 * time.Millisecond)

	go func() {
		for {
			<-ticker.C
			// todo looks like renderer isn't working
			// one possible issue is the logger takes stdout
			// renderer.Render(game)
		}
	}()

	// add the players
	bots := [callbreak.NPlayers]player.Player{}
	for i := 0; i < callbreak.NPlayers; i++ {
		name := "bot" + fmt.Sprint(i)
		token, err := game.AddPlayer(name)
		b, _ := player.New("bot", name, token)
		if err != nil {
			msg := fmt.Errorf("failed to setup game: %v", err)
			panic(msg)
		}
		bots[i] = b
	}

	// play the cards
	for i := 0; i < callbreak.NCards; i++ {
		player := bots[game.GetState("0").Next]
		c, _ := player.Play(game)
		err := game.Play(player.Token(), c)
		if err != nil {
			msg := fmt.Errorf("invalid move from a player: %v", err)
			panic(msg)
		}
		time.Sleep(time.Millisecond * 500)
	}
}
