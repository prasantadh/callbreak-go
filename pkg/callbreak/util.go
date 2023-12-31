package callbreak

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/deck"
)

// return a slice containing cards that can be played
// by the player with next turn in the ongoing trick
func GetValidMoves(game *CallBreak) ([]deck.Card, error) {

	round := game.Rounds[game.RoundNumber]
	trick := round.Tricks[round.TrickNumber]
	next := (trick.Lead + trick.Size) % NPlayers
	hand := round.Hands[next]

	if !hand.IsValid() || !hand.HasPlayable() {
		return nil, fmt.Errorf("invalid hand")
	}

	switch trick.Size {
	case 0:
		return hand.Playables(), nil
	case NPlayers:
		return nil, fmt.Errorf("trick full, no playable moves")
	}

	if trick.Winner() == -1 {
		return nil, fmt.Errorf("trick is in invalid state")
	}
	winner := trick.Cards[trick.Winner()]
	leader := trick.Cards[trick.Lead]

	const leadSuitWinners, leadSuit, turupWinners, playables = 0, 1, 2, 3
	var validmoves [4][]deck.Card

	for _, c := range hand {
		if !c.Playable {
			continue
		}
		group := playables
		if c.Suit == leader.Suit {
			if winner.Suit == leader.Suit && c.Rank > winner.Rank {
				group = leadSuitWinners
			} else {
				group = leadSuit
			}
		} else if c.Suit == winner.Suit && c.Rank > winner.Rank {
			group = turupWinners
		}
		validmoves[group] = append(validmoves[group], c)
	}

	for i := range validmoves {
		if len(validmoves[i]) != 0 {
			return validmoves[i], nil
		}
	}

	return nil, fmt.Errorf("no playable card from hand to trick")
}
