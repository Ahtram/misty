package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/ahtram/misty/misty"
)

func main() {
	// A welcome message with version number.
	misty.PrintWelcomeMessage()

	// Start the bot.
	StartBot()

	return
}

// StartBot gets the bot running.
func StartBot() error {
	//The prime data object.
	bot := misty.Misty{Updating: false}
	bot.GetVars()
	bot.Update()

	// Create a new Discord session using the provided login information.
	session, err := discordgo.New(bot.Params.Email, bot.Params.Password, bot.Params.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return err
	}

	// Get the account information.
	user, err := session.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
		return err
	}

	// Store the account ID for later use.
	bot.BotID = user.ID

	fmt.Println("BotID: " + misty.Green(bot.BotID))

	// Register messageHandler as a callback for the messageHandler events.
	session.AddHandler(bot.MessageHandler)

	// Open the websocket and begin listening.
	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return err
	}

	session.ChannelMessageSend(misty.AsylumChannelID, "Misty is here! Hello world! :smile::smile::smile:")
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})

	return nil
}
