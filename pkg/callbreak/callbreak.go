package callbreak

import (
	"github.com/prasantadh/callbreak-go/pkg/deck"
	"github.com/prasantadh/callbreak-go/pkg/player"
)

type Game struct { // a game is
	Players   [NPlayers]*player.Player // between 4 players
	deck.Deck                          // with a standard deck of cards
}

func (g *Game) New(p string) *Game {
	game := &Game{}

	game.Deck = deck.New()

	game.Players[0] = player.New(p)
	for i := 1; i < NPlayers; i++ {
		game.Players[i] = player.New("bot" + string(i))
	}

	return game

}

func (g *Game) Play() {

	/*
		for each player,
			player.Play
			notify all other players of the played cards
	*/

}
