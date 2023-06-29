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
	return &CallBreak{}
}

func (g *CallBreak) AddPlayer(p PlayerInterface) error {
	if g.nPlayers == NPlayers {
		return fmt.Errorf("couldn't add more players: table already full")
	}

	g.players[g.nPlayers] = p
	return nil

}

func (g *CallBreak) Deal() (deck.Card, error) {
	if len(g.Deck) == 1 {
		return g.Deck[0], nil
	}

	defer func() {
		g.Deck = g.Deck[1:]
	}()

	return g.Deck[0], nil
}

func (g *CallBreak) CollectDeck() {
	g.Deck = deck.New()
}

func (g *CallBreak) Play(c deck.Card) error {
	// update internal state to indicate that
	// a player has played a card in their turn
	return nil

}
