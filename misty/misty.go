package misty

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/ahtram/misty/gshelp"

	"github.com/bwmarrin/discordgo"
)

// Hard coded URLs
const localizedStringSheetFeedURL = "https://spreadsheets.google.com/feeds/worksheets/1w0EKa3K7pNQHY5sAAlY6I-9wQgub9jAe2ozC_1_N7FU/public/full"
const literalCommandSheetFeedURL = "https://spreadsheets.google.com/feeds/worksheets/1haLbQuE7TtF79_J2XLbzFRYbAkfGRCmrXxwdbJ0d724/public/full"
const commandPrefix = "misty"
const guideReply = "Use [misty help] to get help info! :laughing:"

// CmdFunc is the function type for misty's commands.
type CmdFunc func(args []string) string

// Misty is the primary data used by misty. It's a cheap db repacement.
type Misty struct {
	Params ExeParams
	BotID  string
	// Bunit-in commands.
	cmdFuncs map[string]CmdFunc
	// This is the command name index. We need this to properly order the [help] command's output.
	cmdNames []string
	// User defined custom command return strings.
	literalCommands map[string]string
	// Localized string data from TET.
	lstrings map[string][]string
}

//=========== Define all build-in cmd process function here ===========

func (misty *Misty) cmdHelp(words []string) string {
	helpMessage := ":secret:"
	helpMessage = helpMessage + " Try these keywords on me: "
	for _, value := range misty.cmdNames {
		if value != "help" {
			helpMessage = helpMessage + " [" + value + "]"
		}
	}
	helpMessage = helpMessage + " :secret:"
	return helpMessage
}

func (misty *Misty) cmdLString(words []string) string {
	if len(words) > 0 {
		args := words[1:]
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
		return "Use [misty lstring <StringID>] to query an in-game string."
	}
	// Show info message if there's no args.
	return "Use [misty lstring <StringID>] to query an in-game string."
}

// cmdLiteral query the user define reply string and return it.
func (misty *Misty) cmdLiteral(words []string) string {
	if len(words) > 0 {
		cmd := words[0]
		return misty.literalCommands[cmd]
	}
	return "&^*(&^%$*&()*&$%#@))(*&^%$#@!!!!!!!)"
}

//=====================================================================

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func (misty *Misty) MessageHandler(session *discordgo.Session, messageCreate *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if messageCreate.Author.ID == misty.BotID {
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
			// args := words[1:]
			// Call the cmd func and input words.
			return misty.cmdFuncs[words[0]](words)
		}

		return "I don't know what you mean [" + words[0] + "]. " + guideReply
	}

	// Not response.
	return ""
}

// getVars will scan all vars with flag and return them.
func (misty *Misty) GetVars() {
	//Parse (read) parmeters.
	flag.StringVar(&misty.Params.Email, "e", "", "Account Email")
	flag.StringVar(&misty.Params.Password, "p", "", "Account Password")
	flag.StringVar(&misty.Params.Token, "t", "", "Account Token")
	flag.Parse()
}

// update do all data sync with sheet files on our Google Drive. And refresh anything needed.
func (misty *Misty) Update() {
	misty.syncLStrings()
	misty.syncLiteralCommands()
	misty.updateCommands()
}

// syncLStrings fetches lstrings from our Google Drive and return them.
func (misty *Misty) syncLStrings() {
	// Sync LStrings.
	fmt.Print("Syncing LString Data...")
	workSheetXMLContent, err := fetchFeed(localizedStringSheetFeedURL)

	// All tabs' GSeetData.
	sheetData := []gshelp.GSheetData{}

	if err != nil {
		//Oh carp!
		fmt.Println("[Error] " + err.Error())
	} else {
		fmt.Println("[Complete]")
		URLs := gshelp.WorkSheetFeedToCellFeedURLs(workSheetXMLContent)

		// Get all cellfeeds.
		for i, URL := range URLs {
			fmt.Print("[Fetching Tab] : [" + strconv.Itoa(i) + "]...")
			cellXMLContent, err := fetchFeed(URL)
			if err != nil {
				fmt.Println("[Error] " + err.Error())
			} else {
				tabData := gshelp.CellFeedToGSheetData(cellXMLContent)

				// Store in the golbal var.
				sheetData = append(sheetData, tabData)
				fmt.Println("[Complete]")
			}
		}
	}

	// This will empty this container.
	misty.lstrings = make(map[string][]string)

	// Iterate through tabs.
	for _, sheetTab := range sheetData {
		// Iterate through rows.
		for _, row := range sheetTab.StringTable {
			// Check if each row has an ID.
			if len(row) > 0 {
				if row[0] != "" {
					// Add this row.
					misty.lstrings[row[0]] = row[1:5]
				}
			}
		}
	}
}

func (misty *Misty) syncLiteralCommands() {
	// Sync LStrings.
	fmt.Print("Syncing LiteralCommands Data...")
	workSheetXMLContent, err := fetchFeed(literalCommandSheetFeedURL)

	// All tabs' GSeetData.
	sheetData := []gshelp.GSheetData{}

	if err != nil {
		//Oh carp!
		fmt.Println("[Error] " + err.Error())
	} else {
		fmt.Println("[Complete]")
		URLs := gshelp.WorkSheetFeedToCellFeedURLs(workSheetXMLContent)

		// Get all cellfeeds.
		for i, URL := range URLs {
			fmt.Print("[Fetching Tab] : [" + strconv.Itoa(i) + "]...")
			cellXMLContent, err := fetchFeed(URL)
			if err != nil {
				fmt.Println("[Error] " + err.Error())
			} else {
				tabData := gshelp.CellFeedToGSheetData(cellXMLContent)

				// Store in the golbal var.
				sheetData = append(sheetData, tabData)
				fmt.Println("[Complete]")
			}
		}
	}

	// This will empty this container.
	misty.literalCommands = make(map[string]string)

	// Iterate through tabs.
	for _, sheetTab := range sheetData {
		// Iterate through rows.
		for _, row := range sheetTab.StringTable {
			// Check if each row has an ID.
			if len(row) > 1 {
				if row[0] != "" && row[1] != "" {
					// Add this row.
					misty.literalCommands[row[0]] = row[1]
				}
			}
		}
	}
}

func (misty *Misty) updateCommands() {
	misty.cmdFuncs = make(map[string]CmdFunc)
	misty.cmdFuncs["help"] = misty.cmdHelp
	misty.cmdNames = append(misty.cmdNames, "help")
	misty.cmdFuncs["lstring"] = misty.cmdLString
	misty.cmdNames = append(misty.cmdNames, "lstring")
	// Add new built-in cmd func below.

	// Add all user define literal commands.
	for key := range misty.literalCommands {
		if _, exist := misty.cmdFuncs[key]; !exist {
			// cmdLiteral will query literalCommands for response.
			misty.cmdFuncs[key] = misty.cmdLiteral
			misty.cmdNames = append(misty.cmdNames, key)
		}
	}

	fmt.Println("updateCommands done. Command count: [" + strconv.Itoa(len(misty.cmdFuncs)) + "].")
}
