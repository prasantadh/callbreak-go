package player

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
	log "github.com/sirupsen/logrus"
)

type CliBasicBot struct {
	name  string
	token callbreak.Token
}

func (p *CliBasicBot) Name() string {
	return p.name
}

func (p *CliBasicBot) Token() callbreak.Token {
	return p.token
}

func (p *CliBasicBot) Play(game *callbreak.CallBreak) (deck.Card, error) {
	state := game.GetState(p.Token())
	var me *callbreak.Player
	for i, player := range state.Players {
		if player.Name == p.name {
			me = &state.Players[i]
		}
	}

	// current := me.Rounds[len(me.Rounds)-1]
	// t := current.Tricks[len(current.Tricks)-1]
	r := me.Rounds[len(me.Rounds)-1]
	hand := r.Hand
	trick := r.Tricks[len(r.Tricks)-1]
	log.Infof("current player hand: %s\n", hand)
	log.Infof("current trick: %s\n", trick.Cards)

	// for an empty trick play the first playable card on hand
	if trick.Size == 0 {
		for i, c := range hand {
			if c.Playable {
				hand[i].Playable = false
				log.Infof("playing: %s", hand[i])
				return c, nil
			}
		}
		msg := fmt.Errorf("no playable card left")
		log.Errorf("no playable card left!")
		return deck.Card{}, msg
	}

	// if we have a card of the leading suit, play it
	for i, card := range hand {
		if card.Suit == trick.Cards[trick.Lead].Suit &&
			card.Playable {
			// TODO: sort the hand by rank so that we play the highest
			// by default. IMP because if there is a winning card,
			// we must play the winning card
			hand[i].Playable = false
			return card, nil
		}
	}

	// if we have a Hukum card, play it
	for i, card := range hand {
		// TODO: eventually check to see if the hukum we are playing
		// wins all cards in the trick, else play non-Hukum
		// for now, playing the largest hukum is a valid move.
		if card.Suit == deck.Hukum &&
			card.Playable {
			hand[i].Playable = false
			return card, nil
		}
	}

	// we don't have the leading suit and we don't have Hukum
	// play any random playable card
	var ans deck.Card
	for i, c := range hand {
		if c.Playable {
			hand[i].Playable = false
			ans = c
			break
		}
	}
	return ans, nil
}
