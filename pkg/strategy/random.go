package strategy

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

func init() {
	callbreak.RegisterStrategy("random",
		func() callbreak.Strategy {
			return &Random{}
		})
}

type Random struct{}

func (s *Random) Call(game *callbreak.CallBreak) (callbreak.Score, error) {
	return callbreak.Score(1), nil
}

func (s *Random) Break(game *callbreak.CallBreak) (deck.Card, error) {
	moves, err := callbreak.GetValidMoves(game)
	if err != nil {
		return deck.Card{}, fmt.Errorf("could not get a valid move: %v", err)
	}
	rand.Seed(time.Now().Unix())
	r := rand.Intn(len(moves))
	return moves[r], nil
}
