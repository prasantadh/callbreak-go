## API 
### Response Values
All response are json strings of the format
```json
{
"status" : "success/failure"
"data" : "<json string representing game data on success>/<error message on failure>"
}
```

Name        |  Description
------------|-----------------------------------------------------
status      | `success` on success and `failure` on failure
data        | When applicable, return data that was requested.

### Request Endpoints

Name                     | Description
-------------------------|--------------------------------------
new     | Restart a game. In doing so, your last game will no longer exist. TODO: choose mode: single/multi/debug
query   | Query the current game
call    | Make your call for the round.
break   | Play a card from your hand onto the table. The card is the json from of `deck.Card`
register| Register a player to the game. On success, response `data` has authentication token

### Game Data
The `data` field returned by the server is a string. However, it is a string that is valid json object.
The json object is in the following format.
```json
{
  "players": [
    { "name": "bot0", "token": "" },
    { "name": "bot1", "token": "" },
    { "name": "bot2", "token": "" },
    { "name": "me", "token": "1a4bb00eb4b2dccf3e47029628fb42555a7bf94f5965fc4e9dc87de63b34727f" }
  ],
  "rounds": [
    {
      "Calls": [ 1, 1, 1, 1],
      "Scores": [ 0, 3, 9, 1],
      "Hands": [
        [ // 13 cards of player 0's hand
          { "Suit": 0, "Rank": 0, "Playable": false },
            ...13 cards of player 0's hand...
          { "Suit": 0, "Rank": 0, "Playable": false }
        ],
        [ // 13 cards of player1's hand ],
        [ // 13 cards of player2's hand ],
        [ // 13 cards of player3's hand ],
      ],
      "Tricks": [ // 13 tricks of the round, each trick is 4 cards, one from each player
        { // trick 0
          "Cards": [
            { "Suit": 3, "Rank": 8,  "Playable": false },
            { "Suit": 3, "Rank": 14, "Playable": false },
            { "Suit": 3, "Rank": 13, "Playable": false },
            { "Suit": 3, "Rank": 10, "Playable": false }
          ],
          "Lead": 0,
          "Size": 4
        },
        { // trick 1
          "Cards": [ {}, {}, {}, {} ], "Lead": 1, "Size": 4
        },
        { // trick 2 },
        { // trick 3 },
        { // trick 4 },
        { // trick 5 },
        { // trick 6 },
        { // trick 7 },
        { // trick 8 },
        { // trick 9 },
        { // trick 10 },
        { // trick 11 },
        { // trick 12 }
      ],
      "TrickNumber": 13
    },
    { // round 1 },
    { // round 2 },
    { // round 3 },
    { // round 4 }
  ],
  "stage": 2,
  "totalplayers": 4,
  "roundnumber": 4
}
```
