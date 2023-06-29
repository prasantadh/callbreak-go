package callbreak

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/deck"
)

type CallBreak struct { // a game is
	players      []player // between 4 players
	deck.Deck             // with a standard deck of cards
	CurrentTrick Trick
}

func New() *CallBreak {
	return &CallBreak{}
}

func (g *CallBreak) AddPlayer(p PlayerInterface) error {
	if len(g.players) == NPlayers {
		return fmt.Errorf("couldn't add more players: table already full")
	}

	g.players = append(g.players, player{PlayerInterface: p})
	return nil

}

// deal the next card to the next player who should get it
func (g *CallBreak) Deal() error {
	if count := len(g.players); count != 4 {
		msg := fmt.Errorf("wrong number of players: wanted %d got %d", NPlayers, count)
		return msg
	}

	d := deck.New()
	for i, c := range d {
		g.players[i%NPlayers].Take(c)
	}

	return nil
}

// get hand of the ith player of the game
// TODO: authenticate before returning
func (g *CallBreak) GetHand(i int) Hand {
	return g.players[i].hand
}

// the next player in line playes the card c
func (g *CallBreak) Play(c deck.Card) error {
	// update internal state to indicate that
	// a player has played a card in their turn
	return nil

}
