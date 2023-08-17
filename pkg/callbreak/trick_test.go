package callbreak

import (
	"testing"

	"github.com/prasantadh/callbreak-go/pkg/deck"
)

func TestHukumWins(t *testing.T) {
	trick := Trick{}
	trick.Cards[0] = deck.Card{Suit: deck.Chidi, Rank: deck.Dua}
	trick.Cards[1] = deck.Card{Suit: deck.Chidi, Rank: deck.Nahar}
	trick.Cards[2] = deck.Card{Suit: deck.Chidi, Rank: deck.Ekka}
	trick.Cards[3] = deck.Card{Suit: deck.Hukum, Rank: deck.Dua}

	want := 3
	got := trick.Winner()
	if want != got {
		t.Errorf("hukum didn't win: %s won instead of %s",
			trick.Cards[got], trick.Cards[want])
	}
}

func TestHigherRankWins(t *testing.T) {
	trick := Trick{}
	trick.Cards[0] = deck.Card{Suit: deck.Chidi, Rank: deck.Dua}
	trick.Cards[1] = deck.Card{Suit: deck.Chidi, Rank: deck.Nahar}
	trick.Cards[2] = deck.Card{Suit: deck.Chidi, Rank: deck.Ekka}
	trick.Cards[3] = deck.Card{Suit: deck.Chidi, Rank: deck.Attha}

	want := 2
	got := trick.Winner()
	if want != got {
		t.Errorf("higher rank didn't win: %s won instead of %s",
			trick.Cards[got], trick.Cards[want])
	}
}

func TestNotFullTrickHasWinner(t *testing.T) {
	trick := Trick{Lead: 2}
	trick.Cards[2] = deck.Card{Suit: deck.Itta, Rank: deck.Bassa}
	want := 2
	got := trick.Winner()
	if want != got {
		t.Errorf("higher rank didn't win: %s won instead of %s",
			trick.Cards[got], trick.Cards[want])
	}
}
