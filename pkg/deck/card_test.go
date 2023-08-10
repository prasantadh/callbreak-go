package deck

import (
	"encoding/json"
	"testing"

)

func TestPrintCardJson(t *testing.T) {

    card := Card{Suit: Hukum, Rank: Ekka}
    b, _ := json.Marshal(card)
    t.Log(string(b))
}
