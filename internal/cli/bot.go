package cli

import (
	"fmt"
	"sync"
	"time"

	"github.com/prasantadh/callbreak-go/pkg/basicrenderer"
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/player"
	// log "github.com/sirupsen/logrus"
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
	// TODO: make this configurable
	ticker := time.NewTicker(500 * time.Millisecond)

	go func() {
		for {
			<-ticker.C
			// TODO fix the renderer call with the new game architecture
			// renderer.Render(game)
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < callbreak.NPlayers; i++ {
		name := "bot" + fmt.Sprint(i)
		token, _ := game.AddPlayer(name)
		b, _ := player.New("bot", name, token)
		if i == 0 {
			// TODO currently rendering the game from the
			// first player to register's perspective
			// eventually change this or allow game to change the order
			renderer.SetToken(token)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			// TODO: when there are 4 players playing
			// this will eventually have to time out for each game
			b.Play(game)
		}()
	}
	wg.Wait()

}
