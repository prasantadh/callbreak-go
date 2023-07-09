package deck

type Suit int

const (
	Hukum Suit = iota
	Chidi
	Itta
	Paan
)

var Suits [4]Suit = [...]Suit{Hukum, Chidi, Itta, Paan}

func (s *Suit) String() string {
	switch *s {
	case Hukum:
		return "♠"
	case Chidi:
		return "♣"
	case Itta:
		return "♦"
	case Paan:
		return "♥"
	}
	return ""
}
