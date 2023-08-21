package callbreak

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Assistant struct {
	strategy Strategy
	game     *CallBreak
	ticker   *time.Ticker
	token    Token
}

func (p *Assistant) Assist() {
	type wait struct {
		round   int
		trick   int
		waiting bool
	}
	waitdata := wait{}
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
		if current.Next() == me {
			round := current.Rounds[current.RoundNumber]
			if waitdata.round != current.RoundNumber ||
				waitdata.trick != round.TrickNumber ||
				!waitdata.waiting {
				waitdata.round = current.RoundNumber
				waitdata.trick = round.TrickNumber
				waitdata.waiting = true
				continue
			}
			waitdata.waiting = false

			if current.Stage == DEALT {
				call, err := p.strategy.Call(current)
				if err != nil {
					log.Infof("call strategy failed for assistant %d", me)
				}
				log.Infof("assistant %d took over: calling %d", me, call)
				err = p.game.Call(p.token, call)
				if err != nil {
					log.Infof("auto call failed for assistant %d", me)
				}
			} else if current.Stage == CALLED {
				card, err := p.strategy.Break(current)
				if err != nil {
					log.Infof("break strategy failed for assistant %d", me)
				}
				log.Infof("assistant %d took over: breaking %s", me, card)
				err = p.game.Break(p.token, card)
				if err != nil {
					log.Infof("auto break failed for assistant %d", me)
				}
			}
		}
	}
}
