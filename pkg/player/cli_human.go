package player

import "github.com/prasantadh/callbreak-go/pkg/deck"

type CliHuman struct {
	name string
}

func (p *CliHuman) GetGameState() {
	return
}

func (p *CliHuman) Play() (deck.Card, error) {
	return deck.Card{}, nil
}

func (p *CliHuman) SetName(name string) error {
	return nil
}
