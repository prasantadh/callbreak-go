package deck

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
	return d
}
