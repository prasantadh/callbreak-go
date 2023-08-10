package strategy

import (
	"fmt"
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
)

func init() {
	callbreak.RegisterStrategy("basic",
		func() callbreak.Strategy {
			return &Basic{}
		})
}

type Basic struct{}

func (s *Basic) Call(game callbreak.CallBreak) (callbreak.Call, error) {
	return callbreak.Call(1), nil
}

func (s *Basic) Break(game callbreak.CallBreak) (deck.Card, error) {
	moves, err := callbreak.GetValidMoves(game)
	if err != nil {
		fmt.Errorf("could not get a valid move: %v", err)
	}
	return moves[0], nil
}
