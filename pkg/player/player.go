package player

import "fmt"

func New(n string, name string) (Player, error) {
	var bot Player
	if n == "human" {
		bot = &CliHuman{}
	}
	bot = &CliBasicBot{}
	if err := bot.SetName(name); err != nil {
		return nil, fmt.Errorf("error initializing bot: %v", err)
	}
	return bot, nil
}
