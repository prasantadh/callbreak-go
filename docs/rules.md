## Game Details
It is a game of call and break. You call how many tricks you can win. You
attempt to win your calls and break others' calls.


```txt
Glossary:
Hand: The cards that a player has.
Trick: A collection of 4 cards played face up from each player.
In Nepali, these would both be called *haat*.
```

### Objective of the Game
Be the player with the highest score at the end of the game.

### Rules
**Generics**: The game is played clockwise with a standard deck of 52 cards 
for 5 rounds among 4 players.

**Scoring**: At the beginning of a round, each player is given a hand of 13 cards and 
makes a "call". A call is the estimated number of tricks you'll win in that round.
  - If you win less than called tricks, you get -call points for that round.
  - If you win called number of tricks, you get +call points for that round.
  - If you win more than called number of tricks, you get +0.1 point for each
    extra trick won.

**Taking Turns**: The first player (me in single player mode) deals the card 
for first round. In each subsequent round, the player clockwise-next to the player 
who dealt the last round deals the card.
- The first player to receive a card for the round calls/plays first for that
round.
- Starting second trick of the round, the player to win the last trick, starts 
the trick.

**Valid Calls**: Minimum of 1 and maximum of 8. A player who wins 8 tricks
in a round wins the game without further play (currently not implemented). 

**Valid Play**: In a trick, any card in the hand is a valid first play. All
subsequent players MUST play the card in the following order of priority.:
1. A winning card of the same suit
2. A card of the same suit
3. A winning card of the Spade suit
4. Any card on hand

**Card Ranking**: For each suit, the descending order ranking of cards is
A, K, Q, X (10), 9, 8, 7, 6, 5, 4, 3, 2. Spade suit wins over all other suits.
Spade can only be played if the leading suit is either Spade or not available.
All other suits lose to the leading suit irrespective of their rank.
