package callbreak

import (
	"time"
)

type Assistant struct {
	strategy Strategy
	game     *CallBreak
	last     *CallBreak
	ticker   *time.Ticker
	token    Token
}

func (p *Assistant) Assist() {
	for {
		<-p.ticker.C
		current, _ := p.game.Query(p.token)

		// TODO implement callbreak.DONE
		if current.Stage == DONE {
			p.ticker.Stop()
			break
		}

		if current.Stage < DEALT {
			continue
		}

		turn := current.Turn(&p.token)
		next := current.Next()

		// todo: before playing also have to check
		// it was for the same round and the same trick
		if turn == next && p.last.Turn(&p.token) == turn {
			if current.Stage == DEALT {
				call, err := p.strategy.Call(current)
				if err != nil {
					// TODO: log that autoplay failed
				}
				p.game.Call(p.token, call)
			} else if current.Stage == CALLED {
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
