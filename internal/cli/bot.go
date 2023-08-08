package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/prasantadh/callbreak-go/pkg/basicrenderer"
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
	renderer := basicrenderer.New()
	ticker := time.NewTicker(500 * time.Millisecond)

	go func() {
		for {
			<-ticker.C
			renderer.Render(game)
		}
	}()

	// add the players
	bots := [callbreak.NPlayers]player.Player{}
	for i := 0; i < callbreak.NPlayers; i++ {
		b, _ := player.New("bot", "bot"+fmt.Sprint(i))
		_, err := game.AddPlayer(b.Name())
		if err != nil {
			msg := fmt.Errorf("failed to setup game: %v", err)
			panic(msg)
		}
		bots[i] = b
	}

	r, err := json.Marshal(game.GetState("0"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(r))

	// play the cards
	for i := 0; i < callbreak.NCards; i++ {
		// game.Update()
		// trick := game.CurrentTrick()
		fmt.Printf("card %d\n", i)
		player := bots[game.GetState("0").Next]
		c, _ := player.Play(game)
		err := game.Play(player.Token(), c)
		if err != nil {
			msg := fmt.Errorf("invalid move from a player: %v", err)
			panic(msg)
		}
		time.Sleep(time.Millisecond * 500)
	}
	// game.Update()
}
