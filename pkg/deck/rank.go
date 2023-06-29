package deck

import "fmt"

type Rank int8

const (
	Dua Rank = 2 + iota
	Tirka
	Chauka
	Panja
	Xakka
	Satta
	Attha
	Nahar
	Dahar
	Gulam
	Missi
	Bassa
	Ekka
)

var Ranks [13]Rank = [...]Rank{Dua,
	Tirka,
	Chauka,
	Panja,
	Xakka,
	Satta,
	Attha,
	Nahar,
	Dahar,
	Gulam,
	Missi,
	Bassa,
	Ekka}

func (r Rank) String() string {
	switch r {
	case Dahar:
		return "X" // TODO: revert this back to 10 once testing isn't "visual"
	case Gulam:
		return "J"
	case Missi:
		return "Q"
	case Bassa:
		return "K"
	case Ekka:
		return "A"
	}
	return fmt.Sprintf("%d", r)
}
