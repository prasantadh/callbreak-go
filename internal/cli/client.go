package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
	"github.com/prasantadh/callbreak-go/pkg/renderer"
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

var (
	flagAddress  string
	flagName     string
	flagStrategy string
	flagTimeout  time.Duration
)

func init() {
	clientCommand.PersistentFlags().StringVarP(&flagAddress, "address", "a",
		"http://127.0.0.1:8000/", "address/url of the callbreak server")
	clientCommand.PersistentFlags().StringVarP(&flagName, "name", "n",
		"me", "your name for the game")
	clientCommand.PersistentFlags().StringVarP(&flagStrategy, "strategy", "s",
		"basic", "your default strategy in case of timeout")
	clientCommand.PersistentFlags().DurationVarP(&flagTimeout, "timeout", "t",
		30*time.Second, "timeout for your play")
	rootCommand.AddCommand(clientCommand)
}

var baseurl string
var token *callbreak.Token
var me int

type response struct {
	Status `json:"status"`
	Data   string `json:"data"`
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

func postNew() error {
	client := http.Client{}
	resp, err := client.Post(baseurl+"new", "application/json", nil)
	if err != nil {
		panic(fmt.Errorf("could not post to /new: %v", err))
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

func postBreak(token *callbreak.Token, card deck.Card) string {
	data := breakRequest{Token: *token, Card: card}
	request, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Errorf("could not marshal card to break: %v", err))
	}
	response := post(baseurl+"break", request)
	if response.Status == Success {
		return "break succeeded"
	}
	return "break failed: " + response.Data
}

func registerPlayers() {
	for i := 1; i < callbreak.NPlayers; i++ {
		name := fmt.Sprintf("bot%d", i)
		config := callbreak.PlayerConfig{ // TODO make bots configurable
			Name:     name,
			Strategy: "basic",
			Timeout:  time.Second,
		}
		postRegister(&config)
	}
	postRegister(&callbreak.PlayerConfig{
		Name:     flagName,
		Strategy: flagStrategy,
		Timeout:  flagTimeout,
	})
	log.Infof("returning from registering players")
}

func runClient(cmd *cobra.Command, args []string) {
	// todo discard log for now. eventually put it in a file
	log.SetOutput(io.Discard)

	// TODO change this to use viper when viper is added
	baseurl = flagAddress

	err := postNew()
	if err != nil {
		panic(fmt.Errorf("failed to create a new game: %v", err))
	}
	registerPlayers()
	log.Info("registerPlayers returned")

	timeout := time.Second
	updateTicker := time.NewTicker(timeout)
	renderer := renderer.New()
	updateC := make(chan callbreak.CallBreak)
	go func() { renderer.Render(updateC, me) }()

	log.Info("before the loop")
	for {
		select {
		case <-updateTicker.C:
			game := postQuery(token)
			if game.Stage == callbreak.DONE {
				os.Exit(0)
			}
			go func() { updateC <- *game }()
		case call := <-renderer.Call():
			log.Info("received a call")
			status := postCall(token, &call)
			log.Info(status)
		case card := <-renderer.Break():
			log.Infof("received a break")
			status := postBreak(token, card)
			log.Info(status)
		}
	}
}
