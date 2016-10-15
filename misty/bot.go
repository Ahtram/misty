package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/ahtram/misty/gshelp"

	"io/ioutil"

	"github.com/bwmarrin/discordgo"
)

// Hard coded URLs
const localizedStringSheetFeedURL = "https://spreadsheets.google.com/feeds/worksheets/1w0EKa3K7pNQHY5sAAlY6I-9wQgub9jAe2ozC_1_N7FU/public/full"

// CommandParams store the parameters.
type CommandParams struct {
	email    string
	password string
	token    string
}

// MistyCauldron is the primary data used by the bot.
type MistyCauldron struct {
	params                    CommandParams
	botID                     string
	localizedStringGSheetData []gshelp.GSheetData
}

// getVars will scan all vars with flag and return them.
func getVars() CommandParams {
	returnParams := CommandParams{}
	//Parse (read) parmeters.
	flag.StringVar(&returnParams.email, "e", "", "Account Email")
	flag.StringVar(&returnParams.password, "p", "", "Account Password")
	flag.StringVar(&returnParams.token, "t", "", "Account Token")
	flag.Parse()
	return returnParams
}

// syncGSData fetches all data sheet from our Google Drive and return them.
func syncGSData() []gshelp.GSheetData {
	// Sync LStrings.
	fmt.Print("Syncing String Data...")
	localizedStringWorkSheetXMLContent, err := fetchFeedXML(localizedStringSheetFeedURL)

	// All tabs' GSeetData.
	data := []gshelp.GSheetData{}

	if err != nil {
		//Oh carp!
		fmt.Println("[Error] " + err.Error())
	} else {
		fmt.Println("[Complete]")
		// fmt.Println(loclizedStringXMLContent)
		URLs := gshelp.WorkSheetFeedToCellFeedURLs(localizedStringWorkSheetXMLContent)

		// Get all cellfeeds.
		for i, URL := range URLs {
			fmt.Print("[Fetching Tab] : [" + strconv.Itoa(i) + "]...")
			loclizedStringCellXMLContent, err := fetchFeedXML(URL)
			if err != nil {
				fmt.Println("[Error] " + err.Error())
			} else {
				tabData := gshelp.CellFeedToGSheetData(loclizedStringCellXMLContent)

				// fmt.Println("[Tab Result]: ")
				fmt.Println(tabData.ToDefaultString())

				// Store in the golbal var.
				data = append(data, tabData)
				fmt.Println("[Complete]")
			}
		}
	}

	// Sync Items.

	// Sync Recipe.

	return data
}

// StartBot gets the bot running.
func StartBot() {
	//The prime data object.
	mistyCauldron := MistyCauldron{}

	//Scan and store params.
	mistyCauldron.params = getVars()

	// Get all data from our Google sheets.
	mistyCauldron.localizedStringGSheetData = syncGSData()

	// Create a new Discord session using the provided login information.
	session, err := discordgo.New(mistyCauldron.params.email, mistyCauldron.params.password, mistyCauldron.params.token)
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
	mistyCauldron.botID = user.ID

	fmt.Println("BotID: " + green(mistyCauldron.botID))

	// Register messageHandler as a callback for the messageHandler events.
	session.AddHandler(mistyCauldron.messageHandler)

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

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func (cauldron MistyCauldron) messageHandler(session *discordgo.Session, messageCreate *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if messageCreate.Author.ID == cauldron.botID {
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

	if messageCreate.Content == "嚇死你!" {
		if len(cauldron.localizedStringGSheetData) > 0 {
			_, _ = session.ChannelMessageSend(messageCreate.ChannelID, ":pig:")
			_, _ = session.ChannelMessageSend(messageCreate.ChannelID, "多語系字串tab數量:"+strconv.Itoa(len(cauldron.localizedStringGSheetData)))
		} else {
			_, _ = session.ChannelMessageSend(messageCreate.ChannelID, "嚇不倒我的:yum:")
			_, _ = session.ChannelMessageSend(messageCreate.ChannelID, "len(localizedStringGSheetData): "+strconv.Itoa(len(cauldron.localizedStringGSheetData)))
		}
	}

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
