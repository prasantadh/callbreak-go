package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
	_ "github.com/prasantadh/callbreak-go/pkg/strategy"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "callbreak server",
	Long: `This subcommand allows you start a callbreak server that will listen
        to incoming connections from callbreak clients`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer(cmd, args)
	},
}

func init() {
	rootCommand.AddCommand(serverCommand)
}

var game *callbreak.CallBreak

func failure(w http.ResponseWriter, data string) {
	response := Response{
		Status: Failure,
		Data:   data,
	}
	b, _ := json.Marshal(response)
	w.Write(b)
}

func success(w http.ResponseWriter, data any) {
	response := Response{Status: Success, Data: data}
	b, err := json.Marshal(response)
	if err != nil {
		panic("could not form sensible data to send")
	}
	w.Write(b)
}

func getNew(w http.ResponseWriter, r *http.Request) {
	game = callbreak.New()
	response := Response{
		Status: Success,
		Data:   "a new game was created. TODO: handle multiple games",
	}
	b, _ := json.Marshal(response)
	w.Write(b)
}

func getCall(w http.ResponseWriter, r *http.Request) {
	if game == nil {
		failure(w, "no existing game. create a game with /new endpoint")
		return
	}

	failure(w, "currently not implemented")
	return
}

func getBreak(w http.ResponseWriter, r *http.Request) {
	if game == nil {
		failure(w, "no existing game. create a game with /new endpoint")
		return
	}

	token := callbreak.Token(r.URL.Query().Get("token"))
	if len(token) == 0 {
		failure(w, "invalid request: missing `token` field")
		return
	}

	t, _ := strconv.Atoi(r.URL.Query().Get("suit"))
	suit := deck.Suit(t)
	if suit < 1 || suit > 4 {
		failure(w, "invalid suit")
		return
	}

	t, _ = strconv.Atoi(r.URL.Query().Get("rank"))
	rank := deck.Rank(t)
	if rank < 1 || rank > 14 {
		failure(w, "invalid rank")
		return
	}

	err := game.Play(token, deck.Card{Suit: suit, Rank: rank, Playable: true})
	if err != nil {
		failure(w, fmt.Sprintf("could not play card: %s", err))
		return
	}

	return
}

func getRegister(w http.ResponseWriter, r *http.Request) {
	if game == nil {
		failure(w, "no existing game. create a game with /new endpoint")
		return
	}

	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		failure(w, "invalid request: missing `name` field")
		return
	}

	strategy, err := callbreak.GetStrategy("basic")
	if err != nil {
		log.Errorf("error setting strategy: %v", err)
	}
	player, err := game.AddPlayer(name, strategy)
	if err != nil {
		failure(w, err.Error())
		return
	}

	response := Response{
		Status: Success,
		Data:   player,
	}
	b, _ := json.Marshal(response)
	w.Write(b)
	return
}

func getQuery(w http.ResponseWriter, r *http.Request) {
	if game == nil {
		failure(w, "no existing game. create a game with /new endpoint")
		return
	}
	success(w, game)
}

// TODO: add a template for an experimental bot
// then provide instructions on what can be changed to
// have your own bot playing in the arena
func runServer(cmd *cobra.Command, args []string) {

	// 4 people can connect and play
	http.HandleFunc("/new", getNew)
	http.HandleFunc("/register", getRegister)
	http.HandleFunc("/call", getCall)
	http.HandleFunc("/break", getBreak)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Infof("server closed.")
	} else {
		log.Infof("could not start server: %v", err)
	}
	// on timeout received, a random move is made

}
