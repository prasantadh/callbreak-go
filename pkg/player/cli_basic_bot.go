package player

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

type CliBasicBot struct {
	name  string
	token callbreak.Token
	Hand  callbreak.Hand
}

func (p *CliBasicBot) Name() string {
	return p.name
}

func (p *CliBasicBot) SetName(name string) error {
	return nil
}

func (p *CliBasicBot) Token() callbreak.Token {
	return p.token
}

func (p *CliBasicBot) SetToken(token callbreak.Token) {
	p.token = token
}

func (p *CliBasicBot) GetGameState() {
	return
}

func (p *CliBasicBot) Play(game *callbreak.CallBreak) (deck.Card, error) {
	fmt.Println(game.GetState(p.token))
	t := game.GetState(p.token).Players[0].Rounds[0].Tricks[0]
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
