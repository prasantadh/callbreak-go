package player

import (
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
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

func (p *CliHuman) Play(game *callbreak.CallBreak) {
	return
}
