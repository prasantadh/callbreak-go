package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	_ "github.com/prasantadh/callbreak-go/pkg/strategy"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "start a callbreak server for incoming players",
	Long: `This subcommand allows you start a callbreak server that will listen
        to incoming connections from callbreak clients`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer(cmd, args)
	},
}

var (
	flagHost string
	flagPort int
)

func init() {
	serverCommand.PersistentFlags().StringVarP(&flagHost, "address", "a", "0.0.0.0",
		"address to listen on ex. 127.0.0.1")
	serverCommand.PersistentFlags().IntVarP(&flagPort, "port", "p", 8000,
		"port to listen on")
	rootCommand.AddCommand(serverCommand)
}

var game *callbreak.CallBreak

var unsupported_method_response, _ = json.Marshal(Response{
	Status: Failure,
	Data:   "unsupported method: please use POST",
})

var game_nil_response, _ = json.Marshal(Response{
	Status: Failure,
	Data:   "no active game: create a new game with /new",
})

func handleNew(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		game = callbreak.New()
		gamejson, _ := json.Marshal(game)
		response := Response{
			Status: Success,
			Data:   string(gamejson),
		}
		data, _ := json.Marshal(response)
		w.Write(data)
	default:
		log.Infof("/new request received with method %s", r.Method)
		w.Write(unsupported_method_response)
	}
}

func error_message(msg string, err error) []byte {
	response := Response{
		Status: Failure,
		Data:   msg + err.Error(),
	}
	data, _ := json.Marshal(response)
	return data
}

func handleCall(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if game == nil {
			w.Write(game_nil_response)
			return
		}

		request := callRequest{}
		body, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.Write(error_message("invalid POST data: ", err))
			return
		}

		err = game.Call(request.Token, request.Call)
		if err != nil {
			w.Write(error_message("call failed: ", err))
			return
		}

		q, _ := game.Query(request.Token)
		gamejson, _ := json.Marshal(q)
		response := Response{
			Status: Success,
			Data:   string(gamejson),
		}
		data, _ := json.Marshal(response)
		w.Write(data)
	default:
		w.Write(unsupported_method_response)
	}
}

func handleBreak(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if game == nil {
			w.Write(game_nil_response)
			return
		}

		request := breakRequest{}
		body, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(body, &request)
		if err != nil {
			log.Infof("break failed: %v", err)
			w.Write(error_message("invalid POST data: ", err))
			return
		}

		request.Card.Playable = true // make it Playable
		err = game.Break(request.Token, request.Card)
		if err != nil {
			w.Write(error_message("break failed: ", err))
			return
		}

		q, _ := game.Query(request.Token)
		gamejson, _ := json.Marshal(q)
		response := Response{
			Status: Success,
			Data:   string(gamejson),
		}
		data, _ := json.Marshal(response)
		w.Write(data)
	default:
		w.Write(unsupported_method_response)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if game == nil {
			w.Write(game_nil_response)
			return
		}

		request := registerRequest{}
		body, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.Write(error_message("invalid POST data: ", err))
			return
		}

		config := callbreak.PlayerConfig{
			Name:     request.Name,
			Strategy: request.Strategy,
			Timeout:  request.Timeout,
		}
		playerid, err := game.AddPlayer(config)
		if err != nil {
			w.Write(error_message("register failed: ", err))
			return
		}

		q, _ := game.Query(playerid.Token)
		gamejson, _ := json.Marshal(q)
		response := Response{
			Status: Success,
			Data:   string(gamejson),
		}
		data, _ := json.Marshal(response)
		w.Write(data)
	default:
		w.Write(unsupported_method_response)
	}
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if game == nil {
			w.Write(game_nil_response)
			return
		}

		request := queryRequest{}
		body, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.Write(error_message("invalid POST data: ", err))
			return
		}

		q, err := game.Query(request.Token)
		if err != nil {
			w.Write(error_message("query failed: ", err))
			return
		}

		gamejson, _ := json.Marshal(q)
		response := Response{
			Status: Success,
			Data:   string(gamejson),
		}
		data, _ := json.Marshal(response)
		w.Write(data)
	default:
		w.Write(unsupported_method_response)
	}
}

// TODO: add a template for an experimental bot
// then provide instructions on what can be changed to
// have your own bot playing in the arena
func runServer(cmd *cobra.Command, args []string) {

	// 4 people can connect and play
	http.HandleFunc("/new", handleNew)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/call", handleCall)
	http.HandleFunc("/break", handleBreak)
	http.HandleFunc("/query", handleQuery)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", flagHost, flagPort), nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Infof("server closed.")
	} else {
		log.Infof("could not start server: %v", err)
	}
	// on timeout received, a random move is made

}
