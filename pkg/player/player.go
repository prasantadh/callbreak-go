package autoplay

import (
	"time"

	"github.com/prasantadh/callbreak-go/pkg/callbreak"
)

type Config struct {
	Auto     bool
	Strategy callbreak.Strategy
	Game     *callbreak.CallBreak
	Token    callbreak.Token
}

type Player struct {
	strategy callbreak.Strategy
	game     *callbreak.CallBreak
	last     callbreak.CallBreak
	ticker   *time.Ticker
	token    callbreak.Token
}

func New(config Config) *Player {
	player := &Player{
		strategy: config.Strategy,
		game:     config.Game,
		ticker:   time.NewTicker(10 * time.Millisecond),
		token:    config.Token,
	}
	go player.Clock()
	return player
}

func (p *Player) Clock() {
	for {
		<-p.ticker.C
		current := p.game.Query(p.token)

		if current.Stage == callbreak.DONE || current.RoundNumber >= callbreak.NRounds {
			p.ticker.Stop()
			break
		}

		if current.Stage < callbreak.DEALT {
			continue
		}

		turn := current.Turn(&p.token)
		next := current.Next()

		// todo: before playing also have to check
		// it was for the same round and the same trick
		if turn == next && p.last.Turn(&p.token) == turn {
			if current.Stage == callbreak.DEALT {
				call, err := p.strategy.Call(current)
				if err != nil {
					// TODO: log that autoplay failed
				}
				p.game.Call(p.token, call)
			} else if current.Stage == callbreak.CALLED {
				card, err := p.strategy.Break(current)
				if err != nil {
					// TODO: log that autoplay failed
				}
				p.game.Break(p.token, card)
			}
		}

		p.last = current
	}
}
