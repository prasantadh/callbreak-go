package callbreak

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/deck"
	"golang.org/x/exp/slices"
)

func (g *CallBreak) GetState(token Token) Game {
	// if token matches game token return everything
	// if token matches a player token, blur out other players
	// if token isn't a match return error
	response := Game{}
	for i, p := range g.players {
		rounds := []Round{}
		for _, r := range g.rounds {
			rounds = append(rounds, Round{
				Call:   r.calls[i],
				Break:  r.breaks[i],
				Hand:   r.hands[i],
				Tricks: r.tricks,
			})
		}
		response.Players = append(response.Players, Player{
			Name:   p.name,
			Rounds: rounds,
		})
		response.Next = g.next
	}
	return response
}

func New() *CallBreak {
	return &CallBreak{}
}

func (g *CallBreak) newRound() error {

	if len(g.rounds) >= NRounds {
		return fmt.Errorf("no more rounds left to play in the game")
	}

	if r := g.CurrentRound(); r != nil && len(r.tricks) < NTricks {
		return fmt.Errorf("more tricks left to play in this round")
	}

	if t := g.CurrentTrick(); t != nil && t.Size < NPlayers {
		return fmt.Errorf("current trick not full, players remaining")
	}

	g.next = len(g.rounds) % NRounds
	g.rounds = append(g.rounds, round{tricks: make([]Trick, 1)})
	g.CurrentRound().deal()
	// TODO: implement CALLED
	g.state = CALLED
	return nil
}

// add a player to the game. returns an authentication token on success
// else return error on failure
func (g *CallBreak) AddPlayer(name string) (Token, error) {
	if len(g.players) == NPlayers {
		return "", fmt.Errorf("couldn't add more players: table already full")
	}

	g.players = append(g.players, player{
		name: name,
		// TODO implement more secure token mechanism
		token: Token(fmt.Sprint(len(g.players))),
	})

	if len(g.players) == NPlayers {
		g.newRound()
	}

	return g.players[len(g.players)-1].token, nil
}

// deal the cards to the players
// each player can now make a call to GetHand
// TODO: auto trigger this action when round starts
func (r *round) deal() {

	d := deck.New()
	// TODO: make sure each player is dealt at least one Hukum
	// and at least one of Q, K, A else shuffle again
	for i, c := range d {
		p := i % NPlayers
		r.hands[p] = append(r.hands[p], c)
	}

	for i := range r.hands {
		r.hands[i].Sort()
	}

}

// updates the data structure curr of type current
// returns error when called without a valid game state for play
func (g *CallBreak) updateCurrent(curr *current, token Token) error {
	round := g.CurrentRound()

	// current player
	curr.player = nil
	for i, p := range g.players {
		if p.token == token {
			curr.player = &g.players[i]
			curr.hand = &round.hands[i]
			curr.call = &round.calls[i]
			curr.score = &round.breaks[i]
			curr.trick = g.CurrentTrick()
			return nil
		}
	}
	if curr.player == nil {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// return sets of cards that are valid moves for current trick
// if two sets at i and j > i are non-empty,
// the played card must be in set at i
func (g *CallBreak) getValidMoves(curr *current) [][]deck.Card {

	leadSuitWinners := []deck.Card{}
	leadSuit := []deck.Card{}
	turupWinners := []deck.Card{}
	playables := []deck.Card{}

	if curr.trick.Size == 0 {
		for _, c := range *curr.hand {
			if c.Playable {
				playables = append(playables, c)
			}
		}
		return [][]deck.Card{leadSuitWinners, leadSuit, turupWinners, playables}
	}

	winner := curr.trick.Cards[curr.trick.Winner()]
	leader := curr.trick.Cards[curr.trick.Lead]

	for _, c := range *curr.hand {
		if !c.Playable {
			continue
		}

		if c.Suit == leader.Suit {
			if winner.Suit == leader.Suit && c.Rank > winner.Rank {
				leadSuitWinners = append(leadSuitWinners, c)
			} else {
				leadSuit = append(leadSuit, c)
			}
		} else if c.Suit == winner.Suit && c.Rank > winner.Rank {
			turupWinners = append(turupWinners, c)
		} else {
			playables = append(playables, c)
		}
	}
	return [][]deck.Card{leadSuitWinners, leadSuit, turupWinners, playables}

}

// the next player in line playes the card c
// TODO: authorize the player for this action
func (g *CallBreak) Play(token Token, card deck.Card) error {
	// assert players are dealt the card and have made the calls
	// assert only one player can be going through this path at a time
	// if operating asynchronously
	// assert Play is currently valid move
	//      players have been dealt the cards
	//      players have made the calls
	//      there is an active round and active trick

	if g.state != CALLED {
		return fmt.Errorf("invalid action at current game stage")
	}

	curr := &current{}
	err := g.updateCurrent(curr, token)
	if err != nil {
		return fmt.Errorf("cannot play: %v", err)
	}

	if g.state != CALLED {
		return fmt.Errorf("cannot play: not all players have called")
	}

	if *curr.player != g.players[g.next] {
		return fmt.Errorf("you are not up next")
	}

	if curr.trick.Size >= NPlayers {
		return fmt.Errorf("no active trick to play on")
	}

	validMoveSets := g.getValidMoves(curr)
	for i, validMoveSet := range validMoveSets {
		if len(validMoveSet) != 0 && !slices.Contains(validMoveSet, card) {
			switch i {
			//TODO update these ints with nicely named constants
			case 0:
				return fmt.Errorf("must play winning card of the leading suit")
			case 1:
				return fmt.Errorf("must play card of the leading suit")
			case 2:
				return fmt.Errorf("must play a Hukum card")
			case 3:
				return fmt.Errorf("must play a card in hand")
			}
		}
	}

	for i, c := range *curr.hand { // play the card
		if c.Suit == card.Suit && c.Rank == card.Rank {
			(*curr.hand)[i].Playable = false
			(*curr.trick).Add(card)
		}
	}

	if curr.trick.Size == NPlayers {
		if len(g.CurrentRound().tricks) < NTricks {
			g.CurrentRound().tricks = append(g.CurrentRound().tricks, Trick{
				Lead: curr.trick.Winner(),
			})
		} else {
			g.newRound()
		}
	}

	return nil
}

func (g *CallBreak) CurrentRound() *round {
	if len(g.rounds) == 0 {
		return nil
	}
	return &g.rounds[len(g.rounds)-1]
}

func (g *CallBreak) CurrentTrick() *Trick {
	round := g.CurrentRound()
	if round == nil {
		return nil
	}

	if len(round.tricks) == 0 {
		return nil
	}
	return &round.tricks[len(round.tricks)-1]
}
