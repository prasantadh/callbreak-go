package rank

import (
	"testing"
)

func TestRankValueIsCorrect(t *testing.T) {
	if Dua != 2 ||
		Tirka != 3 ||
		Chauka != 4 ||
		Panja != 5 ||
		Xakka != 6 ||
		Satta != 7 ||
		Attha != 8 ||
		Nahar != 9 ||
		Dahar != 10 ||
		Gulam != 11 ||
		Missi != 12 ||
		Bassa != 13 ||
		Ekka != 14 {
		t.Errorf("some rank name is wrong!")
	}
}
