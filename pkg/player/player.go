package player

import (
	"github.com/prasantadh/callbreak-go/pkg/card"
	"github.com/prasantadh/callbreak-go/pkg/hand"
	"github.com/prasantadh/callbreak-go/pkg/trick"
)

type Player struct {
	Name   string
	hands  []hand.Hand
	tricks []trick.Trick
}

func New(n string) *Player {
	return &Player{
		Name: n,
	}
}

func (p *Player) Play() (card.Card, error) {
	// each player will be part of a game
	// query the game state, decide on the appropriate card to play
	// then play
	return card.Card{}, nil
}
