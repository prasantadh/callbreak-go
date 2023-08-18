package cli

import (
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

type Status string

type Response struct {
	Status `json:"status"`
	Data   string `json:"data"`
}

const (
	Success Status = "success"
	Failure        = "failure"
)

type callRequest struct {
	Token callbreak.Token `json:"token"`
	Call  callbreak.Score `json:"call"`
}

type queryRequest struct {
	Token callbreak.Token `json:"token"`
}

type breakRequest struct {
	Token callbreak.Token `json:"token"`
	Suit  deck.Suit       `json:"suit"`
	Rank  deck.Rank       `json:"rank"`
}

type registerRequest callbreak.PlayerConfig

/*
request->
   /new -> {
       -- "debug" : true,
       -- [{"name": "bot0", "strategy" : "basic", "timeout": "500"},{},{}]
   }
   /query -> {
       "token" : "deadbeef"
   }
   /register -> {
       {"name" : "bot0", "strategy" :"basic", "timeout" : "500"}
   }
   /call -> {
       {"token" : "deadbeef", "call" : "5"}
   }
   /break
       {"token" : "deadbeef", "suit" : "Hukum", "rank" : "Ekka"}
response->
 {"status" : "success/failure", "data" : game}
*/
