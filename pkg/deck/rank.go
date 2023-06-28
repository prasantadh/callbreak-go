package deck

const (
	Dua = 2 + iota
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

type Rank int8

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
