package deck

import (
	"math/rand"
	"time"
)

const (
	Size = 52
)

type Deck []Card

func New() Deck {
	d := Deck{}
	for _, r := range Ranks {
		for _, s := range Suits {
			d = append(d, Card{Rank: r, Suit: s, Playable: true})
		}
	}
	d.Shuffle()
	return d
}

func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(*d), func(i, j int) { (*d)[i], (*d)[j] = (*d)[j], (*d)[i] })
	return
}
