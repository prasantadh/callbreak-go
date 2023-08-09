package player

import (
	"fmt"
	"time"

	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
	// log "github.com/sirupsen/logrus"
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

func (p *CliBasicBot) Play(game *callbreak.CallBreak) {
	// TODO make the ticker configurable
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		<-ticker.C
		state, _ := game.Query(p.token)
		me := -1
		for i, player := range state.Players {
			// TODO find a more reliable way to indentify oneself
			// current approach requires unique names for players
			if player.Name == p.name {
				me = i
				break
			}
		}
		if me == -1 {
			panic(fmt.Errorf("can't find myself in game"))
		}

		if state.RoundNumber == callbreak.NRounds {
			break
		}

		// if it's time to call, call TODO implement calling mechanism
		// if it's time to play, play
		round := &state.Rounds[state.RoundNumber]
		trick := &round.Tricks[round.TrickNumber]
		hand := &round.Hands[me]
		next := (trick.Lead + trick.Size) % callbreak.NPlayers
		if state.Stage == callbreak.CALLED && next == me {
			card, _ := p.ChooseCard(hand, trick)
			game.Play(p.token, card)
			// log.Infof("client: %s data\n", p.name)
			// log.Infof("\thand: %s", *hand)
			// log.Infof("\ttrick: %s", trick.Cards)
			// log.Infof("\tcard: %s", card)
		}
	}
}

// TODO currently the chosen card is already marked as not-playable
// only do this if the game accepts the card
func (p *CliBasicBot) ChooseCard(
	hand *callbreak.Hand, trick *callbreak.Trick) (
	deck.Card, error) {
	// set what turn I am
	// then play card per my turn
	// basically the server seems correct, doubledown on trusting server

	// for an empty trick play the first playable card on hand
	if trick.Size == 0 {
		for i, c := range hand {
			if c.Playable {
				hand[i].Playable = false
				return c, nil
			}
		}
		msg := fmt.Errorf("no playable card left")
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
	for i, c := range hand {
		if c.Playable {
			hand[i].Playable = false
			return c, nil
		}
	}
	return deck.Card{}, fmt.Errorf("no playable card in hand")
}
