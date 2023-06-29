package deck

type Suit rune

const (
	Hukum Suit = '♤'
	Chidi      = '♧'
	Itta       = '♢'
	Paan       = '♡'
)

var Suits [4]Suit = [...]Suit{Hukum, Chidi, Itta, Paan}

func (s *Suit) String() string {
	return string(*s)
}
