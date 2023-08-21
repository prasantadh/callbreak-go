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
```json
{
    "players" : [
        {
            "name": "player0",
            "token": "token required to make all call/break in the game"
        },
        {
            "name": "player1",
            "token": "token is empty unless authorized or in debug mode<not implemented>"
        },
        {
            "name": "player2",
            "token": ""
        },
        {
            "name": "player3",
            "token": ""
        }
    ],
    "rounds" : [
    // round 1
    {
        "calls": ["0", "0", "0", "0"], // calls for the round indexed by players
        "scores": ["0", "0", "0", "0"], // score of each player in that round
        "hands": [
            [
                {"suit" : "0", "rank" : "0", "playable": true},
                ... all 13 cards in the hand, blurred out if not authorized...
            ],
            [],
            [],
            []
        ],
        "tricks" : [
        ]

    },
    // round 2, 3, 4, 5
    "roundnumber": "1" // current round number for the game
}
```



