package rank

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

var All [13]Rank = [...]Rank{Dua,
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
