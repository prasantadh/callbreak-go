package bot

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

type Player struct {
	Name string
	Hand callbreak.Hand
}

func New(n string) *Player {
	return &Player{
		Name: n,
	}
}

func (p *Player) Play(t callbreak.Trick) (deck.Card, error) {

	// for an empty trick play the first playable card on hand
	if count := len(t.Cards); count == 0 {
		for i, c := range p.Hand {
			if c.Playable {
				p.Hand[i].Playable = false
				return c, nil
			}
		}
		msg := fmt.Errorf("No playable card left!")
		return deck.Card{}, msg
	}

	// if we have a card of the leading suit, play it
	for _, card := range p.Hand {
		if card.Suit == t.Lead.Suit &&
			card.Playable {
			// TODO: sort the hand by rank so that we play the highest
			// by default. IMP because if there is a winning card,
			// we must play the winning card
			return card, nil
		}
	}

	// if we have a Hukum card, play it
	for _, card := range p.Hand {
		// TODO: eventually check to see if the hukum we are playing
		// wins all cards in the trick, else play non-Hukum
		// for now, playing the largest hukum is a valid move.
		if card.Suit == deck.Hukum {
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
		}
	}
	return ans, nil
}

func (p *Player) Take(c deck.Card) error {
	p.Hand = append(p.Hand, c)
	return nil
}
