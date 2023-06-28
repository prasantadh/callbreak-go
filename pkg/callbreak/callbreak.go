package callbreak

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/deck"
)

type CallBreak struct { // a game is
	players   [NPlayers]PlayerInterface // between 4 players
	deck.Deck                           // with a standard deck of cards
	nPlayers  int                       // number of players currently in the game
}

func New() *CallBreak {
	game := &CallBreak{}
	game.Deck = deck.New()
	return game
}

func (g *CallBreak) AddPlayer(p PlayerInterface) error {
	if g.nPlayers == NPlayers {
		return fmt.Errorf("couldn't add more players: table already full")
	}

	g.players[g.nPlayers] = p
	return nil

}

func (g *CallBreak) Play(c deck.Card) error {

	/*
		for each player,
			player.Play
			notify all other players of the played cards
	*/
	return nil

}
