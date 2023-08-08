package player

import (
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

type CliHuman struct {
	name  string
	token callbreak.Token
}

func (p *CliHuman) Name() string {
	return p.name
}

func (p *CliHuman) Token() callbreak.Token {
	return p.token
}

func (p *CliHuman) Play(game *callbreak.CallBreak) (deck.Card, error) {
	return deck.Card{}, nil
}
