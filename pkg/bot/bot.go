package bot

import (
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

func (p *Player) Play(t callbreak.Trick) deck.Card {

	for i, c := range p.Hand {
		if c.Playable {
			p.Hand[i].Playable = false
			return c
		}
	}

	// // for an empty trick, play the first playable card on hand
	// hand := p.Hand
	// if t.Size == 0 {
	// 	i := 0
	// 	for ; i < callbreak.NTricks && !hand[i].Playable; i++ {
	// 		fmt.Println(hand[i])
	// 	}
	// 	return hand[i]
	// }

	// // if we have a card of the same suit, play it
	// for _, card := range p.Hand {
	// 	if card.Suit == t.Lead.Suit &&
	// 		card.Playable {
	// 		// TODO: sort the hand by rank so that we play the highest
	// 		// by default. IMP because if there is a winning card,
	// 		// we must play the winning card
	// 		return card
	// 	}
	// }

	// // if we have a Hukum card, play it
	// for _, card := range p.Hand {
	// 	// TODO: eventually check to see if the hukum we are playing
	// 	// wins all cards in the trick, else play non-Hukum
	// 	// for now, playing the largest hukum is a valid move.
	// 	if card.Suit == deck.Hukum {
	// 		return card
	// 	}
	// }

	return deck.Card{}
}

func (p *Player) Take(c deck.Card) error {
	p.Hand = append(p.Hand, c)
	// fmt.Println(p.Name, p.Hand)
	return nil
}
