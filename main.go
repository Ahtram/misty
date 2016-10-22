package main

import "gitlab.com/ahtram/misty/bot"

func main() {
	// A welcome message with version number.
	bot.PrintWelcomeMessage()

	// Start the bot.
	StartBot()

	return
}

// StartBot gets the bot running.
func StartBot() error {
	//The prime data object.
	misty := bot.Misty{Updating: false}
	err := misty.Start()
	if err != nil {
		return err
	}

	return nil
}
