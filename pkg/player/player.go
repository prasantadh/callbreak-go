package player

import "github.com/prasantadh/callbreak-go/pkg/callbreak"

func New(nature string, name string, token callbreak.Token) (Player, error) {
	// todo: might have to implement length restriction for name
	var bot Player
	if nature == "human" {
		bot = &CliHuman{
			name:  name,
			token: token,
		}
	}
	bot = &CliBasicBot{
		name:  name,
		token: token,
	}
	return bot, nil
}
