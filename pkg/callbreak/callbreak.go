package callbreak

import (
	"fmt"

	"github.com/prasantadh/callbreak-go/pkg/deck"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func (g *CallBreak) Query(token Token) (*CallBreak, error) {

	g.workPermit <- struct{}{}
	defer func() { <-g.workPermit }()

	// TODO if token matches game token return everything
	// if token matches a player token, blur out other players
	// if token isn't a match return error

	// find the current player requesting data
	current := -1
	for current = range g.Players {
		if g.Players[current].token == token {
			break
		}
	}
	if current == -1 {
		return nil, fmt.Errorf("invalid token")
	}

	// construct the response for this player
	response := *g
	for i := range g.Players {
		if i != current {
			response.Players[i].token = ""
		}
	}

	for r, round := range g.Rounds {
		for hand := range round.Hands {
			if hand != current {
				response.Rounds[r].Hands[hand] = Hand{}
			}
		}
	}

	return &response, nil

}

func New() *CallBreak {
	return &CallBreak{
		workPermit: make(chan struct{}, 1),
	}
}

// add a player to the game. returns an authentication token on success
// else return error on failure
func (game *CallBreak) AddPlayer(name string) (Token, error) {

	game.workPermit <- struct{}{}
	defer func() { <-game.workPermit }()

	if game.TotalPlayers == NPlayers {
		return "", fmt.Errorf("couldn't add more players: table already full")
	}

	player := &game.Players[game.TotalPlayers]
	player.Name = name
	//TODO implement more secure token mechanism
	player.token = Token(fmt.Sprint(game.TotalPlayers))
	game.TotalPlayers += 1

	if game.TotalPlayers == NPlayers {
		game.Rounds[game.RoundNumber].deal()
		// TODO: implement called
		game.Stage = CALLED
	}

	log.Infof("add player %s with token %s", player.Name, player.token)
	return player.token, nil
}

// deal the cards to the players
// each player can now make a call to GetHand
// TODO: auto trigger this action when round starts
func (round *Round) deal() {

	d := deck.New()
	// TODO: make sure each player is dealt at least one Hukum
	// and at least one of Q, K, A else shuffle again
	for i, card := range d {
		player := i % NPlayers
		cardnumber := i / NPlayers
		round.Hands[player][cardnumber] = card
	}

	for i := range round.Hands {
		round.Hands[i].Sort()
	}
	log.Infof("server: dealt: %s", round.Hands)
}

// return sets of cards that are valid moves for current trick
// if two sets at i and j > i are non-empty,
// the played card must be in set at i
// TODO: invalid card in hands will result in valid moves all empty
func (round *Round) getValidMoves(player int) [][]deck.Card {

	leadSuitWinners := []deck.Card{}
	leadSuit := []deck.Card{}
	turupWinners := []deck.Card{}
	playables := []deck.Card{}

	trick := round.Tricks[round.TrickNumber]
	hand := round.Hands[player]

	if trick.Size == 0 {
		for _, c := range hand {
			if c.Playable {
				playables = append(playables, c)
			}
		}
		return [][]deck.Card{leadSuitWinners, leadSuit, turupWinners, playables}
	}

	winner := trick.Cards[trick.Winner()]
	leader := trick.Cards[trick.Lead]

	for _, c := range hand {
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
func (game *CallBreak) Play(token Token, card deck.Card) error {

	game.workPermit <- struct{}{}
	defer func() { <-game.workPermit }()

	// assert Play is currently valid move
	//      players have been dealt the cards
	//      players have made the calls
	//      there is an active round and active trick
	log.Infof("server: player %s attempted play with %s", token, card)

	player := -1
	for i, p := range game.Players {
		if p.token == token {
			player = i
			break
		}
	}
	if player == -1 {
		return fmt.Errorf("cannot play: invalid token")
	}

	if game.Stage != CALLED {
		return fmt.Errorf("cannot play: not all players have called")
	}

	// TODO check that the RoundNumber and TrickNumber are valid
	// before accessing the array with those values
	round := &game.Rounds[game.RoundNumber]
	trick := &round.Tricks[round.TrickNumber]
	next := (trick.Lead + trick.Size) % NPlayers

	// TODO: eventually move this "server:" as a logger field
	log.Infof("server: RoundNumber: %d\tTrickNumber: %d",
		game.RoundNumber, round.TrickNumber)
	log.Infof("server: trick: %s (size: %d lead: %d)",
		trick.Cards, trick.Size, trick.Lead)
	log.Infof("server: Hand: %s", round.Hands[player])

	if player != next {
		return fmt.Errorf("you are not up next")
	}

	if game.RoundNumber == NRounds {
		return fmt.Errorf("no active trick to play on")
	}

	validMoveSets := round.getValidMoves(player)
	log.Infof("server: valid move sets: \n")
	log.Infof("\tleading suit winners: %s", validMoveSets[0])
	log.Infof("\tleading suit: %s", validMoveSets[1])
	log.Infof("\tturup winners: %s", validMoveSets[2])
	log.Infof("\tremaining playables: %s", validMoveSets[3])
	// TODO should be able to refactor this into a function
	for i, validMoveSet := range validMoveSets {
		if len(validMoveSet) == 0 {
			continue
		}

		if slices.Contains(validMoveSet, card) {
			break
		}

		// temporary log
		log.Infof("invalid move from the player")
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
		return fmt.Errorf("invalid card")
	}

	hand := &round.Hands[player]
	for i, c := range hand { // play the card
		if c.Suit == card.Suit && c.Rank == card.Rank {
			hand[i].Playable = false
			err := round.Tricks[round.TrickNumber].Add(card)
			if err != nil {
				panic(fmt.Errorf("cannot play: %v", err))
			}
			break
		}
	}

	if trick.Size == NPlayers { // update results
		winner := trick.Winner()
		round.Breaks[winner] += 1
		round.TrickNumber += 1
		if round.TrickNumber < NTricks {
			round.Tricks[round.TrickNumber].Lead = winner
		}
	}

	// next round
	if round.TrickNumber == NTricks {
		game.RoundNumber += 1
		if game.RoundNumber < NRounds {
			round := &game.Rounds[game.RoundNumber]
			round.deal()
			round.Tricks[round.TrickNumber].Lead = game.RoundNumber % NPlayers
		}
	}

	return nil
}
