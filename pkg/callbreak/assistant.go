package callbreak

import (
	log "github.com/sirupsen/logrus"
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

		me := current.Turn(&p.token)

		// todo: before playing also have to check
		// it was for the same round and the same trick
		if current.Next() == me && p.last.Next() == me {
			if current.Stage == DEALT {
				call, err := p.strategy.Call(current)
				if err != nil {
					// TODO: log that autoplay failed
				}
				log.Infof("assistant took over: calling %d", call)
				p.game.Call(p.token, call)
			} else if current.Stage == CALLED {
				card, err := p.strategy.Break(current)
				if err != nil {
					// TODO: log that autoplay failed
				}
				log.Infof("assistant took over: breaking %s", card)
				p.game.Break(p.token, card)
			}
		}

		p.last = current
	}
}
