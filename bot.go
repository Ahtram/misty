package main

import (
	"flag"
	"fmt"
	"net/http"

	"gitlab.com/ahtram/misty/gshelp"

	"io/ioutil"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters.
var (
	email    string
	password string
	token    string
	botID    string
)

// Hard coded URLs
const localizedStringSheetFeedURL = "https://spreadsheets.google.com/feeds/worksheets/1w0EKa3K7pNQHY5sAAlY6I-9wQgub9jAe2ozC_1_N7FU/public/full"

// ScanVars will scan all vars with flag.
func ScanVars() {
	//Parse (read) parmeters.
	flag.StringVar(&email, "e", "", "Account Email")
	flag.StringVar(&password, "p", "", "Account Password")
	flag.StringVar(&token, "t", "", "Account Token")
	flag.Parse()
}

// SyncGSData will fetch all data sheet from our Google Drive.
func SyncGSData() {
	// Sync LStrings.
	fmt.Print("Syncing String Data...")
	localizedStringWorkSheetXMLContent, err := fetchFeedXML(localizedStringSheetFeedURL)

	// All tabs' GSeetData.
	localizedStringGSheetData := []gshelp.GSheetData{}

	if err != nil {
		//Oh carp!
		fmt.Println("[Error] " + err.Error())
	} else {
		fmt.Println("[Complete]")
		// fmt.Println(loclizedStringXMLContent)
		URLs := gshelp.WorkSheetFeedToCellFeedURLs(localizedStringWorkSheetXMLContent)

		// Get all cellfeeds.
		for _, URL := range URLs {
			loclizedStringCellXMLContent, err := fetchFeedXML(URL)
			if err != nil {
				fmt.Println("[Error] " + err.Error())
			} else {
				fmt.Println("[Tab Result]: ")
				// fmt.Println(loclizedStringCellXMLContent)
				tabData := gshelp.CellFeedToGSheetData(loclizedStringCellXMLContent)
				fmt.Println(tabData.ToDefaultString())
				localizedStringGSheetData = append(localizedStringGSheetData, tabData)

				// gSheetData := gshelp.CellFeedToGSheetData(loclizedStringCellXMLContent)
				// if gSheetData != nil {
				// 	append(localizedStringGSheetData, gSheetData)
				// }
			}
		}
	}

	// Sync Items.

	// Sync Recipe.

}

// StartBot will get the bot running.
func StartBot() {
	// Create a new Discord session using the provided login information.
	session, err := discordgo.New(email, password, token)
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
	botID = user.ID

	fmt.Println("BotID: " + green(botID))

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
}

// fetchFeedXML read XML content from given URL.
func fetchFeedXML(feedURL string) (XMLContent string, err error) {
	response, httpErr := http.Get(feedURL)
	if httpErr != nil {
		return "", httpErr
	}

	defer response.Body.Close()

	// Read the body with ioutil.
	htmlData, ioErr := ioutil.ReadAll(response.Body)
	if ioErr != nil {
		return "", ioErr
	}

	return string(htmlData), nil
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageHandler(session *discordgo.Session, messageCreate *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if messageCreate.Author.ID == botID {
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
