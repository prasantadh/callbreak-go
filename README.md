## Build and Usage (Installation)

```bash
go build ./cmd/callbreak-go
./callbreak-go
```

## API 

### Request Endpoints

Name                     | Description
-------------------------|--------------------------------------
new     | Restart a game. In doing so, your last game will no longer exist. TODO: choose mode: single/multi/debug
query   | Query the current game _optionally_ indexed by `roundnumber` and `tricknumber`
call    | Make your call (the call in callbreak) for the round.
break   | Play a card from your hand onto the table. The card is the json from of `deck.Card`
register| Register a player to the game. On success, response `data` has authentication token


### Response Values

Name        |  Description
------------|-----------------------------------------------------
status      | Specifies if the request was successfully executed. Returns `success` or `failure`. If it fails, `data` returns an explanation.
data        | When applicable, return data that was requested.
