package player

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

type CliBasicBot struct {
	Name string
	Hand callbreak.Hand
}

func (p *CliBasicBot) GetGameState() {
	return
}

func (p *CliBasicBot) Play() (deck.Card, error) {

	// for an empty trick play the first playable card on hand
	if count := len(t.Cards); count == 0 {
		for i, c := range p.Hand {
			if c.Playable {
				p.Hand[i].Playable = false
				return c, nil
			}
		}
		msg := fmt.Errorf("no playable card left")
		return deck.Card{}, msg
	}

	// if we have a card of the leading suit, play it
	for i, card := range p.Hand {
		if card.Suit == t.Cards[t.Lead].Suit &&
			card.Playable {
			// TODO: sort the hand by rank so that we play the highest
			// by default. IMP because if there is a winning card,
			// we must play the winning card
			p.Hand[i].Playable = false
			return card, nil
		}
	}

	// if we have a Hukum card, play it
	for i, card := range p.Hand {
		// TODO: eventually check to see if the hukum we are playing
		// wins all cards in the trick, else play non-Hukum
		// for now, playing the largest hukum is a valid move.
		if card.Suit == deck.Hukum &&
			card.Playable {
			p.Hand[i].Playable = false
			return card, nil
		}
	}

	// we don't have the leading suit and we don't have Hukum
	// play any random playable card
	var ans deck.Card
	for i, c := range p.Hand {
		if c.Playable {
			p.Hand[i].Playable = false
			ans = c
			break
		}
	}
	return ans, nil
}

func (p *CliBasicBot) Take(c deck.Card) error {
	p.Hand = append(p.Hand, c)
	return nil
}

func (p *CliBasicBot) SetName(name string) error {
	return nil
}
