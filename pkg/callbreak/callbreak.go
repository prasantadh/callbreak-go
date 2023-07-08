package callbreak

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/deck"
)

type CallBreak struct { // a game is
	players    []player // between 4 players
	deck.Deck           // with a standard deck of cards
	tricks     []Trick  // and is played as a collection of Tricks
	NextPlayer int
}

func New() *CallBreak {
	return &CallBreak{
		tricks: make([]Trick, 1),
	}
}

func (g *CallBreak) AddPlayer(p PlayerInterface) error {
	if len(g.players) == NPlayers {
		return fmt.Errorf("couldn't add more players: table already full")
	}

	g.players = append(g.players, player{PlayerInterface: p})
	return nil

}

// deal the cards to the players
// each player can now make a call to GetHand
// TODO: auto trigger this action when round starts
func (g *CallBreak) Deal() error {
	if count := len(g.players); count != 4 {
		msg := fmt.Errorf("wrong number of players: wanted %d got %d", NPlayers, count)
		return msg
	}

	d := deck.New()
	// TODO: make sure each player is dealt at least one Hukum
	// and at least one of Q, K, A else shuffle again
	for i, c := range d {
		g.players[i%NPlayers].Take(c)
	}

	for i := range g.players {
		g.players[i].hand.Sort()
	}

	return nil
}

// get hand of the ith player of the game
// TODO: authenticate before returning
func (g *CallBreak) GetHand(i int) Hand {
	return g.players[i].hand
}

// necessary update at the end of a trick
func (g *CallBreak) Update() {
	trick := &g.tricks[len(g.tricks)-1]
	if trick.Size == NPlayers {
		if len(g.tricks) < NTricks {
			g.NextPlayer = trick.Winner()
			g.tricks = append(g.tricks, Trick{Lead: g.NextPlayer})
		}
	}
}

// the next player in line playes the card c
// TODO: authorize the player for this action
func (g *CallBreak) Play(c deck.Card) error {

	g.Update() // Critical to call this

	// the player plays the card
	player := &g.players[g.NextPlayer]
	err := player.Play(c)
	if err != nil {
		return fmt.Errorf("game could not play: %v", err)
	}
	g.NextPlayer = (g.NextPlayer + 1) % NPlayers

	// the card gets added to the trick
	// TODO: it's a problem if player.Play() succeeds
	// but trick.Add fails
	trick := &g.tricks[len(g.tricks)-1]
	err = trick.Add(c)
	if err != nil {
		return fmt.Errorf("game could not add to trick: %v", err)
	}

	return nil
}

func (g *CallBreak) CurrentTrick() Trick {
	return g.tricks[len(g.tricks)-1]
}
