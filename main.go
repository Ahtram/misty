package main

import (
	"flag"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters.
var (
	Email    string
	Password string
	Token    string
	BotID    string
)

func init() {
	//Parse (read) parmeters.
	flag.StringVar(&Email, "e", "", "Account Email")
	flag.StringVar(&Password, "p", "", "Account Password")
	flag.StringVar(&Token, "t", "", "Account Token")
	flag.Parse()
}

func main() {
	// A welcome message with version number.
	PrintWelcomeMessage()

	// Create a new Discord session using the provided login information.
	session, err := discordgo.New(Email, Password, Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	user, err := session.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = user.ID

	fmt.Println("BotID: " + green(BotID))

	// Register messageHandler as a callback for the messageHandler events.
	session.AddHandler(messageHandler)

	// Open the websocket and begin listening.
	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageHandler(session *discordgo.Session, messageCreate *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if messageCreate.Author.ID == BotID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if messageCreate.Content == "ping" {
		_, _ = session.ChannelMessageSend(messageCreate.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if messageCreate.Content == "pong" {
		_, _ = session.ChannelMessageSend(messageCreate.ChannelID, "Ping!")
	}

	if messageCreate.Content == "misty" {
		_, _ = session.ChannelMessageSend(messageCreate.ChannelID, "蝦? 叫我喔?")
	}

	if messageCreate.Content == "今天幾點開會?" {
		_, _ = session.ChannelMessageSend(messageCreate.ChannelID, "9點啦!")
	}

	if messageCreate.Content == "吃飽沒?" {
		_, _ = session.ChannelMessageSend(messageCreate.ChannelID, "還沒啦!")
	}

	if messageCreate.Content == "開會了" {
		_, _ = session.ChannelMessageSend(messageCreate.ChannelID, "https://hangouts.google.com/call/wpi5vlbz6bcc5bm7nromkueyrae")
	}

	if messageCreate.Content == "?" {
		_, _ = session.ChannelMessageSend(messageCreate.ChannelID, "?")
	}

}
