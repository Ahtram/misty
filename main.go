package main

import (
	"flag"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

//Version number and program name define.
const programName = "Misty"
const version = "0.0.0.1"

//Color defines.
var mag = color.New(color.FgHiMagenta).SprintFunc()
var yellow = color.New(color.FgHiYellow).SprintFunc()
var cyan = color.New(color.FgHiCyan).SprintFunc()
var green = color.New(color.FgHiGreen).SprintFunc()

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
	printWelcomeMessage()

	// Create a new Discord session using the provided login information.
	dg, err := discordgo.New(Email, Password, Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	user, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = user.ID

	fmt.Println("BotID: " + green(BotID))

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

func printWelcomeMessage() {
	//Plain welcome text.
	welcomeText := "*   " + programName + " Version [" + version + "]   *"

	//The colored welcome text.
	colorWelcomeText := fmt.Sprintf("%s   %s Version [%s]   %s", cyan("*"), mag(programName), yellow(version), cyan("*"))

	fmt.Println(cyan(starStringLine(len(welcomeText))))
	fmt.Println(colorWelcomeText)
	fmt.Println(cyan(starStringLine(len(welcomeText))))
}

func starStringLine(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = '*'
	}
	return string(b)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(session *discordgo.Session, messageCreate *discordgo.MessageCreate) {

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
