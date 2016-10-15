package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/ahtram/misty/gshelp"

	"io/ioutil"

	"github.com/bwmarrin/discordgo"
)

// Hard coded URLs
const localizedStringSheetFeedURL = "https://spreadsheets.google.com/feeds/worksheets/1w0EKa3K7pNQHY5sAAlY6I-9wQgub9jAe2ozC_1_N7FU/public/full"
const commandPrefix = "misty"
const guideReply = "Use [misty help] to get help info! :laughing:"

// ExeParams store the parameters.
type ExeParams struct {
	email    string
	password string
	token    string
}

// CmdFunc is the function type for misty's commands.
type CmdFunc func(args []string) string

// Misty is the primary data used by the bot.
type Misty struct {
	params                    ExeParams
	botID                     string
	localizedStringGSheetData []gshelp.GSheetData
	cmdFuncs                  map[string]CmdFunc
}

func (misty *Misty) init() {
	misty.cmdFuncs = make(map[string]CmdFunc)
	misty.cmdFuncs["help"] = misty.cmdHelp
	misty.cmdFuncs["lstring"] = misty.cmdLString

	fmt.Println("init done. Cmd count: [" + strconv.Itoa(len(misty.cmdFuncs)) + "].")

	// Add new built-in cmd func below.
}

//=========== Define all build-in cmd process function here ===========

func (misty *Misty) cmdHelp(args []string) string {
	return "This is help message"
}

func (misty *Misty) cmdLString(args []string) string {
	return "This is lstring message"
}

//=====================================================================

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func (misty *Misty) messageHandler(session *discordgo.Session, messageCreate *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if messageCreate.Author.ID == misty.botID {
		return
	}

	// Try response the message.
	reply := misty.responseMessage(messageCreate.Content)
	if reply != "" {
		_, _ = session.ChannelMessageSend(messageCreate.ChannelID, reply)
	}
}

// responseMessage return a suitable string as response message after decision.
// An empty string will be returned if not suitable reply found.
func (misty *Misty) responseMessage(message string) string {
	if strings.HasPrefix(message, commandPrefix+" ") {
		// Trim the prefix to get the message content.
		messageContent := strings.TrimPrefix(message, commandPrefix+" ")

		// get command and argument.(words) They should be devided by an empty character.
		words := strings.Split(messageContent, " ")

		return misty.responseCommand(words)
	} else if message == commandPrefix {
		return guideReply
	}

	// Not response.
	return ""
}

// responseCommand returns the command result by input wors.
// An empty string will be returned if this is not a legal command.
func (misty *Misty) responseCommand(words []string) string {
	if len(words) > 0 {
		// This maybe a command with arguments.
		//Check if misty actually has this command.
		if _, exist := misty.cmdFuncs[words[0]]; exist {
			args := words[1:]
			// Call the cmd func and input args.
			return misty.cmdFuncs[words[0]](args)
		}

		return "I don't know what you mean [" + words[0] + "]. " + guideReply
	}

	// Not response.
	return ""
}

// getVars will scan all vars with flag and return them.
func getVars() ExeParams {
	returnParams := ExeParams{}
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

				// fmt.Println(tabData.ToDefaultString())

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
	misty := Misty{}
	//This will initial the commands for misty.
	misty.init()

	//Scan and store params.
	misty.params = getVars()

	// Get all data from our Google sheets.
	misty.localizedStringGSheetData = syncGSData()

	// Create a new Discord session using the provided login information.
	session, err := discordgo.New(misty.params.email, misty.params.password, misty.params.token)
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
	misty.botID = user.ID

	fmt.Println("BotID: " + green(misty.botID))

	// Register messageHandler as a callback for the messageHandler events.
	session.AddHandler(misty.messageHandler)

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
