package callbreak

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
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
		log.Infof("assistant %d ticker fired: ", me)

		// todo: before playing also have to check
		// it was for the same round and the same trick
		if current.Next() == me && p.last.Next() == me {
			if me == 3 {
				panic(fmt.Errorf("playing now"))
			}
			if current.Stage == DEALT {
				call, err := p.strategy.Call(current)
				if err != nil {
					// TODO: log that autoplay failed
				}
				log.Infof("assistant %d took over: calling %d", me, call)
				p.game.Call(p.token, call)
			} else if current.Stage == CALLED {
				card, err := p.strategy.Break(current)
				if err != nil {
					// TODO: log that autoplay failed
				}
				log.Infof("assistant %d took over: breaking %s", me, card)
				p.game.Break(p.token, card)
			}
		}

		p.last = current
	}
}
