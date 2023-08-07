package callbreak

import (
	"errors"
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/deck"
	"golang.org/x/exp/slices"
)

var current struct {
	player_index int
	round_index  int
	tricks_index int
	hand         *Hand
	player       *player
	round        *round
	trick        *Trick
}

func (p *CallBreak) GetState(token Token) Game {
	// if token matches game token return everything
	// if token matches a player token, blur out other players
	// if token isn't a match return error
	return Game{}
}

func New() *CallBreak {
	return &CallBreak{}
}

func (g *CallBreak) newRound() error {

	if len(g.rounds) >= NRounds {
		return fmt.Errorf("no more rounds left to play in the game")
	}

	r := g.rounds[len(g.rounds)-1] // current round
	if len(r.tricks) < NTricks {
		return fmt.Errorf("more tricks left to play in this round")
	}

	t := r.tricks[len(r.tricks)-1] // current trick
	if t.Size != NPlayers {
		return fmt.Errorf("current trick not full, players remaining")
	}

	r = round{}
	g.rounds = append(g.rounds, r)
	g.next = len(g.rounds) % NRounds
	r.deal()
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
		token: Token(len(g.players)),
	})

	if len(g.players) == NPlayers {
		// the game is on
		g.rounds = append(g.rounds, round{})
		g.next = 0
	}

	return Token(len(g.players)), nil
}

// deal the cards to the players
// each player can now make a call to GetHand
// TODO: auto trigger this action when round starts
func (r *round) deal() error {

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

	return nil
}

// // get hand of the ith player of the game
// // TODO: authenticate before returning
// func (g *CallBreak) GetHand(i int) Hand {
// 	return g.players[i].hand
// }

// necessary update at the end of a trick
// func (g *CallBreak) Update() {
// 	trick := &g.tricks[len(g.tricks)-1]
// 	if trick.Size == NPlayers {
// 		winner := trick.Winner()
// 		g.NextPlayer = winner
// 		g.players[winner].Score += 1
// 		g.tricks = append(g.tricks, Trick{Lead: winner})
// 	}
// }

func (g *CallBreak) updateCurrent() error {
	current.round_index = g.CurrentRound()
	if current.round_index >= NRounds {
		return fmt.Errorf("game over")
	}
	current.round = &g.rounds[current.round_index]

	current.tricks_index = g.CurrentTrick()
	if current.tricks_index >= NTricks {
		return fmt.Errorf("round over, consider starting a new round")
	}
	current.trick = &current.round.tricks[current.tricks_index]
	if current.trick.Size >= NPlayers {
		return fmt.Errorf("trick full, consider starting a new trick")
	}

	return nil
}

func (g *CallBreak) verifyValidMove(card deck.Card) error {
	current.hand = &current.round.hands[current.player_index]
	if current.hand.HasPlayable(c) {
		return fmt.Errorf("failed to verify card: %v")
	}

    candidates := []deck.Card
    for c := range *current.hand {
        if c.Playable && c.Suit == card.Suit && c.Rank > card.Rank {
            candidates = append(candidates, c)
        }
    }
    if len(candidates) > 0 {
        if slices.Contains(candidates, card) {
            return nil
        } 
        return fmt.Errorf("must play leading suit winning card when available")
    }

    for c := range *current.hand {
        if c.Playable && c.Suit == card.Suit {
            candidates = append(candidates, c)
        }
    }
    if len(candidates) > 0 && !slices.Contains(candidates, card) {
        return fmt.Errorf("must play leading suit card when available")
    }

    // must play hukum when leading suit isn't available and hukum is
    // must play higher hukum when leading suit isn't available and higher hukum is
	if current.hand.HasSuit(current.trick.LeadSuit()) &&
		c.Suit != current.trick.LeadSuit() {
		return fmt.Errorf("must play same suit when available")
	}

    if current.hand.HasHigherRankForSuit()
	// played same suit lower card, has higher
	// played different suit card has same

	// has a same suit higher card as trick and didn't play
	// has a same suit card as trick and didn't play
	// has a hukum higher card as trick and didn't play
	// has a hukum card and didn't play
	return nil
}

// the next player in line playes the card c
// TODO: authorize the player for this action
func (g *CallBreak) Play(token Token, c deck.Card) error {
	// assert players are dealt the card
	// assert < NRounds
	// assert < NTricks on the last round
	// assert correct player token
	// assert the card is in player's hand
	// assert a valid play
	//      - same suit
	//      - winning card of the same suit
	//      - winning card of the hukum suit

	// TODO: only one player can be going through this path at a time
	// if operating asynchronously
	// and this user has to be authorized to make the next move

	err := g.updateCurrent()
	if err != nil {
		fmt.Errorf("play is currently not a valid move %v", err)
	}

	err = g.verifyValidMove(c)
	if err != nil {
		fmt.Errorf("invalid move: %v", err)
	}

	// if this is the first card on the trick, anything goes
	if current.trick.Lead == current.player_index {
		// play the card
	}

	// if this is the last card on the trick, update trick winner
	// setup new round if necessary

	// g.Update() // Critical to call this

	// TODO: check for all possible errors first
	// current trick must not be full
	// nextPlayer must have the card
	// nextPlayer must play current suit if in hand
	// nextPlayer must play spade if current suit not in hand
	// nextPlayer can play whatever if no space and no current suit in hand
	// we should be able to eat up all subsequent errors after these checks

	// the player plays the card
	// 	player := &g.players[g.NextPlayer]
	// 	err := player.Play(c)
	// 	if err != nil {
	// 		return fmt.Errorf("game could not play: %v", err)
	// 	}
	// 	g.NextPlayer = (g.NextPlayer + 1) % NPlayers
	//
	// 	// the card gets added to the trick
	// 	// TODO: it's a problem if player.Play() succeeds
	// 	// but trick.Add fails
	// 	trick := &g.tricks[len(g.tricks)-1]
	// 	err = trick.Add(c)
	// 	if err != nil {
	// 		return fmt.Errorf("game could not add to trick: %v", err)
	// 	}
	//
	// 	// if this was the last card of the trick, compute winners

	return nil
}

func (g *CallBreak) CurrentRound() int {
	return len(g.rounds)
}

func (g *CallBreak) CurrentTrick() int {
	tricks := g.rounds[g.CurrentRound()].tricks
	return len(tricks) - 1
}
