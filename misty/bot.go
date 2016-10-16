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

// Lang is language enum.
type Lang int

const (
	cht Lang = iota
	chs
	en
	jpn
)

// LangTypeCount is the type count of supported languages.
const LangTypeCount = 4

// LangName stores all language name.
var LangName = []string{"cht", "chs", "en", "jpn"}

// ExeParams store the parameters.
type ExeParams struct {
	email    string
	password string
	token    string
}

// CmdFunc is the function type for misty's commands.
type CmdFunc func(args []string) string

// Misty is the primary data used by misty. It's a cheap db repacement.
type Misty struct {
	params   ExeParams
	botID    string
	cmdFuncs map[string]CmdFunc
	lstrings map[string][]string
	// localizedStringGSheetData []gshelp.GSheetData
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
	helpMessage := ":secret:"
	helpMessage = helpMessage + " Try these keywords on me: "
	for k := range misty.cmdFuncs {
		if k != "help" {
			helpMessage = helpMessage + " [" + k + "]"
		}
	}
	helpMessage = helpMessage + " :secret:"
	return helpMessage
}

func (misty *Misty) cmdLString(args []string) string {
	if len(args) > 0 {
		// Check if the ID exist.
		content, exist := misty.lstrings[args[0]]
		if exist {
			result := "Result: " + fmt.Sprint("\n")
			for i := 0; i < LangTypeCount; i++ {
				// Add a language tag.
				if i < len(LangName) {
					result = result + LangName[i]
				}
				result = result + " [" + content[i] + "]" + fmt.Sprint("\n")
			}
			return result
		}
		return "There is no such string in game [" + args[0] + "]. :weary:"
	}
	// Show info message if there's no args.
	return "Use [misty lstring <StringID>] to query an in-game string."
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

// getLStrings fetches all data sheet from our Google Drive and return them.
func getLStrings() map[string][]string {
	// Sync LStrings.
	fmt.Print("Syncing LString Data...")
	localizedStringWorkSheetXMLContent, err := fetchFeedXML(localizedStringSheetFeedURL)

	// All tabs' GSeetData.
	sheetData := []gshelp.GSheetData{}

	if err != nil {
		//Oh carp!
		fmt.Println("[Error] " + err.Error())
	} else {
		fmt.Println("[Complete]")
		URLs := gshelp.WorkSheetFeedToCellFeedURLs(localizedStringWorkSheetXMLContent)

		// Get all cellfeeds.
		for i, URL := range URLs {
			fmt.Print("[Fetching Tab] : [" + strconv.Itoa(i) + "]...")
			loclizedStringCellXMLContent, err := fetchFeedXML(URL)
			if err != nil {
				fmt.Println("[Error] " + err.Error())
			} else {
				tabData := gshelp.CellFeedToGSheetData(loclizedStringCellXMLContent)

				// Store in the golbal var.
				sheetData = append(sheetData, tabData)
				fmt.Println("[Complete]")
			}
		}
	}

	// Form the returned map.
	lstrings := make(map[string][]string)

	// Iterate through tabs.
	for _, sheetTab := range sheetData {
		// Iterate through rows.
		for _, row := range sheetTab.StringTable {
			// Check if each row has an ID.
			if len(row) > 0 {
				if row[0] != "" {
					// Add this row.
					lstrings[row[0]] = row[1:5]
				}
			}
		}
	}

	return lstrings
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
	misty.lstrings = getLStrings()

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
