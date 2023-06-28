package callbreak

import "github.com/prasantadh/callbreak-go/pkg/deck"

const (
	NCards   = 52
	NPlayers = 4
	NTricks  = NCards / NPlayers
)

type Hand [NTricks]deck.Card
