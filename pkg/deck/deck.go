package deck

type Deck struct {
	cards []Card
}

func New() Deck {
	d := Deck{}
	for _, r := range Ranks {
		for _, s := range Suits {
			d.cards = append(d.cards, Card{Rank: r, Suit: s, Playable: true})
		}
	}
	return d
}
