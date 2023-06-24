package card

import (
	"github.com/prasantadh/callbreak-go/pkg/rank"
	"github.com/prasantadh/callbreak-go/pkg/suit"
)

type Card struct {
	Suit suit.Suit
	Rank rank.Rank
}
