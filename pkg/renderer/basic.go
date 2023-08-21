package renderer

import (
	"fmt"
	"os"
	"strings"

	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/fatih/color"
	"github.com/prasantadh/callbreak-go/pkg/callbreak"
	cb "github.com/prasantadh/callbreak-go/pkg/callbreak"
	"github.com/prasantadh/callbreak-go/pkg/deck"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

var (
	RedString        func(a ...interface{}) string
	BlackString      func(a ...interface{}) string
	FaintString      func(a ...interface{}) string
	BgWhiteString    func(a ...interface{}) string
	UnderlinedString func(a ...interface{}) string
)

var bottom, left, top, right int

func init() {
	RedString = color.New(color.BgWhite).Add(color.FgRed).SprintFunc()
	FaintString = color.New(color.Faint).Add(color.CrossedOut).SprintFunc()
	BlackString = color.New(color.BgWhite).Add(color.FgBlack).SprintFunc()
	BgWhiteString = color.New(color.BgWhite).Add(color.FgBlack).SprintFunc()
	UnderlinedString = color.New(color.Underline).SprintFunc()
}

type Basic struct {
	area       cursor.Area
	token      cb.Token
	keypressed chan keys.KeyCode
	callC      chan callbreak.Score
	breakC     chan deck.Card
	validmoves []deck.Card
	choice     int
}

func (r *Basic) blank(repetition int) string {
	return BgWhiteString(strings.Repeat(" ", 6*repetition))
}

func (r *Basic) coloredCard(c deck.Card) string {
	if c.Playable == false {
		return BgWhiteString(fmt.Sprintf("  --  "))
	}
	s := fmt.Sprintf(" %s %s  ", &c.Suit, &c.Rank)
	if len(r.validmoves) > 0 && c == r.validmoves[r.choice] {
		s = fmt.Sprintf("[%s %s ]", &c.Suit, &c.Rank)
	} else if slices.Contains(r.validmoves, c) {
		s = fmt.Sprintf("(%s %s )", &c.Suit, &c.Rank)
	}
	if c.Suit == deck.Hukum || c.Suit == deck.Chidi {
		if c.Playable {
			return BlackString(s)
		}
		return FaintString(BlackString(s))
	}
	if c.Playable {
		return RedString(s)
	}
	return FaintString(RedString(s))
}

func (r *Basic) coloredHand(h cb.Hand) string {
	sb := strings.Builder{}
	for _, c := range h {
		sb.WriteString(r.coloredCard(c))
	}
	return sb.String()
}

func (r *Basic) activateKeypress() {
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		logrus.Info("got a keypress!")
		if key.Code == keys.CtrlC {
			os.Exit(0)
		}
		switch key.Code {
		case keys.Left, keys.Right, keys.Enter:
			logrus.Infof("sending key %v to the keypressed channel", key)
			go func() { r.keypressed <- key.Code }()
		}
		return false, nil
	})
}

func (r *Basic) renderScoreboard(game cb.CallBreak) string {
	var sb strings.Builder
	sb.WriteString(BgWhiteString("ScoreBoard  "))
	sb.WriteString(r.blank(13))
	sb.WriteString("\n")
	sb.WriteString(BgWhiteString("Bots:  |  me  | bot1 | bot2 | bot3 |"))
	sb.WriteString(r.blank(9))
	sb.WriteString("\n")
	total := [cb.NPlayers]float32{}
	for i := 0; i <= game.RoundNumber; i++ {
		sb.WriteString(r.blank(1))
		sb.WriteString(BgWhiteString(" |"))
		round := game.Rounds[i]
		for i := 0; i < cb.NPlayers; i++ {
			s := fmt.Sprintf(" %2d/%1d |", round.Scores[i], round.Calls[i])
			sb.WriteString(BgWhiteString(s))
			if round.Scores[i] < round.Calls[i] {
				total[i] -= float32(round.Calls[i])
			} else {
				total[i] += float32(round.Calls[i]) +
					float32(float32(round.Scores[i]-round.Calls[i])/10)
			}
		}
		sb.WriteString(r.blank(9))
		sb.WriteString(BgWhiteString("\n"))
	}
	s := fmt.Sprintf("| % 2.1f | % 2.1f | % 2.1f | % 2.1f |",
		total[0], total[1], total[2], total[3])
	sb.WriteString(BgWhiteString("Total: "))
	sb.WriteString(BgWhiteString(s))
	sb.WriteString(r.blank(9))
	sb.WriteString("\n")
	return sb.String()
}

func (r *Basic) renderTrickCard(card deck.Card) string {
	if card != (deck.Card{Playable: true}) {
		return r.coloredCard(card)
	}
	return BgWhiteString("______")
}

// renders the middle part of the table for left bot and right bot
// as well as the trick
func (r *Basic) renderTableMid(game cb.CallBreak, n int) string {
	round := &game.Rounds[game.RoundNumber]
	hands := &round.Hands
	trick := &round.Tricks[round.TrickNumber].Cards

	for i := range trick {
		trick[i].Playable = true
	}

	sb := strings.Builder{}
	sb.WriteString(r.coloredCard(hands[left][n]))
	sb.WriteString(r.blank(4))
	if n == 6 {
		sb.WriteString(r.renderTrickCard(trick[left]))
	} else {
		sb.WriteString(r.blank(1))
	}
	sb.WriteString(r.blank(1))
	if n == 4 {
		sb.WriteString(r.renderTrickCard(trick[top]))
	} else if n == 8 {
		sb.WriteString(r.renderTrickCard(trick[bottom]))
	} else {
		sb.WriteString(r.blank(1))
	}
	sb.WriteString(r.blank(1))
	if n == 6 {
		sb.WriteString(r.renderTrickCard(trick[right]))
	} else {
		sb.WriteString(r.blank(1))
	}
	sb.WriteString(r.blank(4))
	sb.WriteString(r.coloredCard(hands[right][n]))
	sb.WriteString("\n")
	return sb.String()
}

func (r *Basic) renderTable(game cb.CallBreak) string {
	round := game.Rounds[game.RoundNumber]
	hands := round.Hands
	var sb strings.Builder
	sb.WriteString(BgWhiteString("bot2⮕ "))
	sb.WriteString(r.coloredHand(hands[top]))
	sb.WriteString(BgWhiteString("⬇ bot3"))
	sb.WriteString("\n")
	for i := 0; i < cb.NTricks; i++ {
		s := r.renderTableMid(game, i)
		sb.WriteString(s)
	}
	sb.WriteString(BgWhiteString("bot1⬆ "))
	sb.WriteString(r.coloredHand(hands[bottom]))
	sb.WriteString(BgWhiteString("⬅ bot0"))
	sb.WriteString("\n")
	return sb.String()
}

func (r *Basic) Render(update <-chan cb.CallBreak, me int) {

	go r.activateKeypress()
	bottom, left, top, right =
		me, (me+1)%cb.NPlayers, (me+2)%cb.NPlayers, (me+3)%cb.NPlayers
	g := cb.CallBreak{}
	for {
		select {
		case key := <-r.keypressed:
			logrus.Info("re-rendering: received a keypress")
			logrus.Infof("choice is now: %d", r.choice)
			if g.Next() == me && g.Stage == cb.DEALT {
				if key == keys.Left {
					r.choice = (r.choice - 1 + 8) % 8
				} else if key == keys.Right {
					r.choice = (r.choice + 1) % 8
				} else if key == keys.Enter {
					go func() { r.callC <- cb.Score(r.choice + 1) }()
				}
			} else if g.Next() == me && g.Stage == cb.CALLED {
				validMoves, _ := cb.GetValidMoves(&g)
				limit := len(validMoves)
				if key == keys.Left {
					r.choice = (r.choice - 1 + limit) % limit
				} else if key == keys.Right {
					r.choice = (r.choice + 1) % limit
				} else if key == keys.Enter {
					go func() { r.breakC <- validMoves[r.choice] }()
				}
			}
		case current := <-update:
			logrus.Info("re-rendering: received game update")
			if current != g {
				r.choice = 0
				r.validmoves, _ = cb.GetValidMoves(&current)
				g = current
			}
		}

		sb := strings.Builder{}
		scoreboard := r.renderScoreboard(g)
		table := r.renderTable(g)
		sb.WriteString(scoreboard)
		sb.WriteString(table)
		if g.Next() == me && g.Stage == cb.DEALT {
			msg := fmt.Sprintf("Your Call: <- %d -> [Enter]    ", r.choice+1)
			sb.WriteString(BgWhiteString(msg))
			sb.WriteString(r.blank(10))
		} else if g.Next() == me && g.Stage == cb.CALLED {
			sb.WriteString(BgWhiteString("Your Play: <- "))
			sb.WriteString(r.coloredCard(r.validmoves[r.choice]))
			sb.WriteString(BgWhiteString(" -> [Enter]     "))
			sb.WriteString(r.blank(9))
			sb.WriteString("\n")
		}
		r.area.Clear()
		r.area.Update(sb.String())
	}
}

func (r *Basic) Call() <-chan callbreak.Score {
	return r.callC
}

func (r *Basic) Break() <-chan deck.Card {
	return r.breakC
}
