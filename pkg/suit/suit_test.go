package suit

import "testing"

func TestSuitValueIsCorrect(t *testing.T) {
	if Hukum != '♤' ||
		Chidi != '♧' ||
		Itta != '♢' ||
		Paan != '♡' {
		t.Errorf("some suit name is wrong")
	}
}
