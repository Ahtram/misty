package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/ahtram/misty/bot"
)

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
	misty.GetVars()
	misty.Update()

	// Create a new Discord session using the provided login information.
	session, err := discordgo.New(misty.Params.Email, misty.Params.Password, misty.Params.Token)
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
	misty.BotID = user.ID

	fmt.Println("BotID: " + bot.Green(misty.BotID))

	// Register messageHandler as a callback for the messageHandler events.
	session.AddHandler(misty.MessageHandler)

	// Open the websocket and begin listening.
	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return err
	}

	session.ChannelMessageSend(bot.AsylumChannelID, "Misty is here! Hello world! :smile::smile::smile:")
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})

	return nil
}
