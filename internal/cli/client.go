package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"

	// "github.com/prasantadh/callbreak-go/pkg/basicrenderer"
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var clientCommand = &cobra.Command{
	Use:   "client",
	Short: "client a connect to callbreak server",
	Long: `This subcommand allows you start a callbreak server that will listen
        to incoming connections from callbreak clients`,
	Run: func(cmd *cobra.Command, args []string) {
		runClient(cmd, args)
	},
}

func init() {
	rootCommand.AddCommand(clientCommand)
}

var baseurl = "http://localhost:3333/"
var token *callbreak.Token
var me int

type choice struct {
	call callbreak.Score
	card deck.Card
}

type response struct {
	Status `json:"status"`
	Data   string `json:"data"`
}

func postNew() error {
	client := http.Client{}
	resp, err := client.Post(baseurl+"new", "application/json", nil)
	if err != nil {
		panic(fmt.Errorf("could not post to /new"))
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		panic(fmt.Errorf("could not read body"))
	}

	response := response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(fmt.Errorf("unmarshal failed: %v", err))
	}

	if response.Status == Failure {
		return fmt.Errorf("/new request failed: %s", response.Data)
	}

	return nil

}

func postQuery(token *callbreak.Token) *callbreak.CallBreak {
	data := []byte(fmt.Sprintf(`{"token": "%v"}`, token))
	response := post(baseurl+"query", data)

	if response.Status == Failure {
		panic(fmt.Errorf("query failed: %s", response.Data))
	}

	game := callbreak.CallBreak{}
	err := json.Unmarshal([]byte(response.Data), &game)
	if err != nil {
		panic(fmt.Errorf("server sent invalid response: %v", err))
	}

	return &game
}

func post(url string, data []byte) *Response {
	client := http.Client{}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		panic(fmt.Errorf("could not post to server: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("could not read body: %v", err))
	}

	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(fmt.Errorf("server sent invalid response: %v", err))
	}

	return &response
}

func postRegister(config *callbreak.PlayerConfig) string {
	data, _ := json.Marshal(*config)
	response := post(baseurl+"register", data)

	if response.Status == Failure {
		panic(fmt.Errorf("register failed: %s", response.Data))
	}

	err := json.Unmarshal([]byte(response.Data), &game)
	if err != nil {
		panic(fmt.Errorf("server sent invalid data: %v", game))
	}

	for _, player := range game.Players {
		if player.Token != callbreak.Token("") {
			token = &player.Token
			me = game.Turn(token)
		}
	}
	log.Infof("register set token to: %s", *token)

	return "register succeeded"

}

func postCall(token *callbreak.Token, call *callbreak.Score) string {
	data := []byte(fmt.Sprintf(`{"token": "%s", "call": "%v"}`, *token, *call))
	response := post(baseurl+"call", data)

	if response.Status == Success {
		return "call succeeded"
	}
	return "call failed: " + response.Data
}

func postBreak(token *callbreak.Token, card *deck.Card) string {
	data := []byte(fmt.Sprintf(`{"token": "%v", "suit": "%v", "rank": "%v"}`,
		*token, card.Suit, card.Rank))
	response := post(baseurl+"break", data)
	if response.Status == Success {
		return "break succeeded"
	}
	return "break failed: " + response.Data
}

func registerPlayers() {
	for i := 0; i < callbreak.NPlayers; i++ {
		timeout := time.Second
		name := fmt.Sprintf("bot%d", i)
		if i == callbreak.NPlayers-1 {
			timeout = 5 * time.Second // TODO change this to 30s
			name = "me"
		}
		config := callbreak.PlayerConfig{
			Name:     name,
			Strategy: "basic",
			Timeout:  timeout,
		}
		postRegister(&config)
	}
}

func runClient(cmd *cobra.Command, args []string) {
	// todo discard log for now. eventually put it in a file
	// log.SetOutput(io.Discard)

	err := postNew()
	if err != nil {
		panic(fmt.Errorf("failed to create a new game: %v", err))
	}
	registerPlayers()

	timeout := time.Second
	renderticker := time.NewTicker(timeout)
	// renderer := basicrenderer.New()

	for {
		// TODO: don't have to query that often, can reduce the query time
		// also might have to add timeout to the queries even though
		// there is also a default timeout on http.Client
		<-renderticker.C
		game = postQuery(token)
		if game.Stage == callbreak.DONE {
			break
		}
		// renderer.Render(game, me, "message area")

		next := game.Next()
		log.Infof("next player: %d, my turn: %d, stage: %d", game.Next(), me, game.Stage)
		if next == me {
			if game.Stage == callbreak.DEALT {
				log.Info("I am attempting a call")
				// TODO change this to a more reasonable number
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				keypressed := make(chan callbreak.Score)
				go listenKeyboardCall(keypressed)
				select {
				case <-ctx.Done():
					keyboard.SimulateKeyPress(keys.CtrlC)
					log.Infof("timed out, re-rendering")
				case call := <-keypressed:
					log.Infof("sending user call to server: %d", call)
					postCall(token, &call)
				}
				cancel()
			} else if game.Stage == callbreak.CALLED {
				choices, err := callbreak.GetValidMoves(game)
				if err != nil {
					panic(fmt.Errorf("could not get valid moves: %v", err))
				}
				log.Infof("valid moves: %s", choices)
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				keypressed := make(chan int)
				go listenKeyboardBreak(keypressed, len(choices))
				select {
				case <-ctx.Done():
					keyboard.SimulateKeyPress(keys.CtrlC)
				case card := <-keypressed:
					postBreak(token, &choices[card])
				}
				cancel()
			}
		}
	}
}

func listenKeyboardCall(C chan callbreak.Score) {
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.CtrlC {
			return true, nil
		}
		if key.Code == keys.RuneKey {
			for v := 1; v < 8; v++ {
				s := fmt.Sprint(v)
				if s == key.String() {
					C <- callbreak.Score(v)
				}
			}
		}
		return false, nil
	})
}

func listenKeyboardBreak(C chan int, limit int) {
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.RuneKey {
			for v := 1; v < limit; v++ {
				s := fmt.Sprint(v)
				if s == key.String() && v < limit {
					C <- v
				} else if s == key.String() {
					C <- limit - 1
				}
			}
		}
		return true, nil
	})
}
