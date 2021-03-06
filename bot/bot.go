package bot

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ahtram/misty/gshelp"

	"github.com/bwmarrin/discordgo"
)

// CmdFunc is the function type for misty's commands.
type CmdFunc func(words []string, channelID string) string

// Misty is the primary data used by misty. It's a cheap db repacement.
type Misty struct {
	Params ExeParams
	conf   botConfig
	// Store the least watching channel's online status.
	streamingStatus streamingStatusCache
	session         *discordgo.Session
	BotID           string
	// Command functions.
	cmdFuncs map[string]CmdFunc
	// User defined custom command return strings. (map[Name]([Content Str ID][Desc Str ID]))
	literalCommands map[string][2]string
	// This is the command index. We need this to properly order the [help] command's output. ([]([Name][Desc Str ID]))
	cmdIndex [][2]string
	// Localized lines for the bot. [key][localized string array]
	lines map[string][]string
	// Is executing an updating.
	Updating bool
	// The Unity Cloud Hooks we are listening.
	uCloudHooks []*UCloudHook
	// The GitLab Hooks we are listening.
	gitLabHooks []*GitLabHook
	// The GitHub Hooks we are listening.
	gitHubHooks []*GitHubHook
}

// Start the bot.
func (misty *Misty) Start() error {
	// Get all commandline vars.
	misty.GetVars()

	// Check vars.
	if misty.Params.ConfigSheetID == "" || (misty.Params.Token == "" && (misty.Params.Email == "" || misty.Params.Password == "")) {
		// If the user does not behave as we think...
		fmt.Println(Red("Input vars not legal!"))
		flag.Usage()
		return errors.New("Program exit gracefully. ")
	}

	// Update data.
	misty.Update()

	// Create a new Discord session using the provided login information.
	var err error
	misty.session, err = discordgo.New(misty.Params.Email, misty.Params.Password, misty.Params.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return err
	}

	// Get the account information.
	user, err := misty.session.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
		return err
	}

	// Start observe the watching streaming channel.
	misty.startObserveStreamingStatus()

	//Start listen to the Unity Cloud hooks. (check if we have a uCloud End Point and Port setting)
	for _, value := range misty.conf.UCloudConfigs {
		if value.UCloudHookEndPoint != "" && value.UCloudHookPort != "" && value.UCloudAccessToken != "" {
			fmt.Println("Start listen to Unity Cloud hook: " + Yellow("["+value.UCloudHookEndPoint+"] ["+value.UCloudHookPort+"]"))
			uCloudHook := UCloudHook{
				UCloudHookEndPoint: value.UCloudHookEndPoint,
				UCloudHookPort:     value.UCloudHookPort,
				UCloudAccessToken:  value.UCloudAccessToken,
				MistyRef:           misty,
			}
			//Store the hook reference for good.
			misty.uCloudHooks = append(misty.uCloudHooks, &uCloudHook)
			go uCloudHook.StartUCloudHook()
		}
	}

	//Start listen to the GitLab hooks.
	for _, value := range misty.conf.GitLabConfigs {
		if value.GitLabHookEndPoint != "" && value.GitLabHookPort != "" {
			fmt.Println("Start listen to GitLab hook: " + Yellow("["+value.GitLabHookEndPoint+"] ["+value.GitLabHookPort+"]"))
			gitLabHook := GitLabHook{
				GitLabHookEndPoint: value.GitLabHookEndPoint,
				GitLabHookPort:     value.GitLabHookPort,
				MistyRef:           misty,
			}
			//Store the hook reference for good.
			misty.gitLabHooks = append(misty.gitLabHooks, &gitLabHook)
			go gitLabHook.StartGitLabHook()
		}
	}

	//Start listen to the GitHub hooks.
	for _, value := range misty.conf.GitHubConfigs {
		if value.GitHubHookEndPoint != "" && value.GitHubHookPort != "" {
			fmt.Println("Start listen to GitHub hook: " + Yellow("["+value.GitHubHookEndPoint+"] ["+value.GitHubHookPort+"]"))
			gitHubHook := GitHubHook{
				GitHubHookEndPoint: value.GitHubHookEndPoint,
				GitHubHookPort:     value.GitHubHookPort,
				MistyRef:           misty,
			}
			//Store the hook reference for good.
			misty.gitHubHooks = append(misty.gitHubHooks, &gitHubHook)
			go gitHubHook.StartGitHubHook()
		}
	}

	// Store the account ID for later use.
	misty.BotID = user.ID

	fmt.Println("BotID: " + Yellow(misty.BotID))

	// Register messageHandler as a callback for the messageHandler events.
	misty.session.AddHandler(misty.MessageHandler)

	// Open the websocket and begin listening.
	err = misty.session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return err
	}

	//Send online notify message?
	if misty.conf.onlineNotify {
		misty.session.ChannelMessageSend(misty.conf.ResidentDiscordChannelID, misty.Line("onlineNotify", 0))
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})

	return nil
}

// Line returns the line string by ID and language.
func (misty *Misty) Line(lineID string, lang int) string {
	value, exist := misty.lines[lineID]
	if exist {
		if lang >= 0 && len(value) > lang {
			return value[lang]
		}
		return ""
	}
	return ""
}

//=========== Define all build-in cmd process function here ===========

func (misty *Misty) cmdHelp(words []string, channelID string) string {
	helpMessage := ":secret::secret::secret:\n"
	helpMessage += "```Markdown\n"
	for _, value := range misty.cmdIndex {
		if value[0] != "help" && value[0] != "update" && value[0] != "cid" {
			helpMessage += "#[" + value[0] + "]\n"
			helpMessage += "    " + misty.Line(value[1], 0) + "\n"
		}
	}
	helpMessage += "```"
	return helpMessage
}

func (misty *Misty) cmdUpdate(words []string, channelID string) string {
	go misty.Update(channelID)
	return misty.Line("updateStart", 0)
}

func (misty *Misty) cmdChannelID(words []string, channelID string) string {
	returnMessage := "```Markdown\n"
	returnMessage += "#ChannelID:\n"
	returnMessage += channelID + "\n"
	returnMessage += "```"
	return returnMessage
}

// cmdLiteral query the user define reply string and return it.
func (misty *Misty) cmdLiteral(words []string, channelID string) string {
	if len(words) > 0 {
		//Read the first index: the content of this literal command.
		return misty.Line(misty.literalCommands[words[0]][0], 0)
	}
	return "&^*(&^%$*&()*&$%#@))(*&^%$#@!!!!!!!)"
}

//=====================================================================

// MessageHandler be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func (misty *Misty) MessageHandler(session *discordgo.Session, messageCreate *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if messageCreate.Author.ID == misty.BotID {
		return
	}

	// Try response the message.
	reply := misty.responseMessage(messageCreate.Content, messageCreate.ChannelID)

	if reply != "" {
		// fmt.Println("ChannelMessageSend ChannelID: " + messageCreate.ChannelID)
		_, _ = session.ChannelMessageSend(messageCreate.ChannelID, reply)
	}
}

// responseMessage return a suitable string as response message after decision.
// An empty string will be returned if not suitable reply found.
func (misty *Misty) responseMessage(message string, channelID string) string {
	if strings.HasPrefix(message, misty.conf.CommandPrefix+" ") {
		//Check if misty is updating anything
		if !misty.Updating {
			//Could response commands now.
			// Trim the prefix to get the message content.
			messageContent := strings.TrimPrefix(message, misty.conf.CommandPrefix+" ")

			// get command and argument.(words) They should be devided by an empty character.
			words := strings.Split(messageContent, " ")

			return misty.responseCommand(words, channelID)
		}

		return "I'm a little busy right now. Talk to me later. :smile: (Misty is updating data)"
	} else if message == misty.conf.CommandPrefix {
		return misty.Line("guideReply", 0)
	}
	// Not response.
	return ""
}

// responseCommand returns the command result by input wors.
// An empty string will be returned if this is not a legal command.
func (misty *Misty) responseCommand(words []string, channelID string) string {
	if len(words) > 0 {
		// This maybe a command with arguments.
		//Check if misty actually has this command.
		if _, exist := misty.cmdFuncs[words[0]]; exist {
			// args := words[1:]
			// Call the cmd func and input words.
			return misty.cmdFuncs[words[0]](words, channelID)
		}

		return "I don't know what you mean [" + words[0] + "]. " + misty.Line("guideReply", 0)
	}

	// Not response.
	return ""
}

// broadcastMessage send a message to all broadcast channel in config.
func (misty *Misty) broadcastMessage(message string) {
	if misty.session != nil {
		for _, v := range misty.conf.BroadcastDiscrdChannelID {
			misty.session.ChannelMessageSend(v, message)
		}
	}
}

// attempt to remove the previous messages with prefix.
func (misty *Misty) deletePreviousBroadcastMessage(messagePrefix string) {
	for _, v := range misty.conf.BroadcastDiscrdChannelID {
		//Get previous messages first.
		previousMessages, err := misty.session.ChannelMessages(v, 100, "", "", "")
		if err == nil {
			for _, msg := range previousMessages {
				if strings.HasPrefix(msg.Content, messagePrefix) || strings.Compare(msg.Content, messagePrefix) == 0 {
					err = misty.session.ChannelMessageDelete(v, msg.ID)
					if err != nil {
						fmt.Println(Red("[Error] ") + err.Error())
					}
				}
			}
		} else {
			fmt.Println(Red("[Error] ") + err.Error())
		}
	}
}

// GetVars will scan all vars with flag and return them.
func (misty *Misty) GetVars() {
	//Parse (read) parmeters.
	flag.StringVar(&misty.Params.Email, "e", "", "Account Email")
	flag.StringVar(&misty.Params.Password, "p", "", "Account Password")
	flag.StringVar(&misty.Params.Token, "t", "", "Bot Token")
	flag.StringVar(&misty.Params.ConfigSheetID, "c", "", "Config Sheet")
	flag.Parse()
}

// Update do all data sync with sheet files on our Google Drive. And refresh anything needed.
func (misty *Misty) Update(channelID ...string) {
	// Check if we are already updating.
	if !misty.Updating {
		// Not updating. So we do update.
		misty.Updating = true
		misty.syncConfig()
		misty.syncLines()

		misty.cmdFuncs = make(map[string]CmdFunc)
		misty.cmdIndex = [][2]string{}

		misty.updateBuiltInCommands()
		misty.syncLiteralCommands()

		if misty.session != nil {
			if len(channelID) > 0 {
				misty.session.ChannelMessageSend(channelID[0], misty.Line("updateComplete", 0))
			} else {
				misty.session.ChannelMessageSend(misty.conf.ResidentDiscordChannelID, misty.Line("updateComplete", 0))
			}
		}
		misty.Updating = false
	}
}

func (misty *Misty) syncConfig() {
	// Sync LStrings.
	fmt.Print("Syncing Config Data...")
	workSheetXMLContent, err := fetchFeed(misty.Params.ConfigSheetURL())

	// All tabs' GSeetData.
	sheetData := []gshelp.GSheetData{}

	if err != nil {
		//Oh carp!
		fmt.Println(Red("[Error] ") + err.Error())
	} else {
		fmt.Println(Green("[Complete]"))
		URLs := gshelp.WorkSheetFeedToCellFeedURLs(workSheetXMLContent)

		// Get all cellfeeds.
		for i, URL := range URLs {
			fmt.Print("[Fetching Tab] : [" + strconv.Itoa(i) + "]...")
			cellXMLContent, err := fetchFeed(URL)
			if err != nil {
				fmt.Println(Red("[Error] ") + err.Error())
			} else {
				tabData := gshelp.CellFeedToGSheetData(cellXMLContent)

				// Store in the golbal var.
				sheetData = append(sheetData, tabData)
				fmt.Println(Green("[Complete]"))
			}
		}
	}

	misty.conf.Setup(sheetData)

	//Print the config of this bot.
	fmt.Println(misty.conf.ToString())
}

// syncLStrings fetches lstrings from our Google Drive and return them.
func (misty *Misty) syncLines() {
	// Sync LStrings.
	fmt.Print("Syncing Line Data...")
	workSheetXMLContent, err := fetchFeed(misty.conf.LineSheetURL())

	// All tabs' GSeetData.
	sheetData := []gshelp.GSheetData{}

	if err != nil {
		//Oh carp!
		fmt.Println(Red("[Error] ") + err.Error())
	} else {
		fmt.Println(Green("[Complete]"))
		URLs := gshelp.WorkSheetFeedToCellFeedURLs(workSheetXMLContent)

		// Get all cellfeeds.
		for i, URL := range URLs {
			fmt.Print("[Fetching Tab] : [" + strconv.Itoa(i) + "]...")
			cellXMLContent, err := fetchFeed(URL)
			if err != nil {
				fmt.Println(Red("[Error] ") + err.Error())
			} else {
				tabData := gshelp.CellFeedToGSheetData(cellXMLContent)

				// Store in the golbal var.
				sheetData = append(sheetData, tabData)
				fmt.Println(Green("[Complete]"))
			}
		}
	}

	// This will empty this container.
	misty.lines = make(map[string][]string)

	// Iterate through tabs.
	for _, sheetTab := range sheetData {
		// Iterate through rows.
		for _, row := range sheetTab.StringTable {
			// Check if each row has an ID.
			if len(row) > 0 {
				if row[0] != "" {
					// Add this row.
					misty.lines[row[0]] = row[1:5]
				}
			}
		}
	}
}

func (misty *Misty) updateBuiltInCommands() {
	//build-in commands
	misty.cmdFuncs["help"] = misty.cmdHelp
	misty.cmdIndex = append(misty.cmdIndex, [2]string{"help", ""})
	misty.cmdFuncs["update"] = misty.cmdUpdate
	misty.cmdIndex = append(misty.cmdIndex, [2]string{"update", ""})
	misty.cmdFuncs["cid"] = misty.cmdChannelID
	misty.cmdIndex = append(misty.cmdIndex, [2]string{"cid", ""})
	// Add new built-in cmd func here...

	// // Add all user define literal commands.
	// for key, value := range misty.literalCommands {
	// 	if _, exist := misty.cmdFuncs[key]; !exist {
	// 		// cmdLiteral will query literalCommands for response.
	// 		misty.cmdFuncs[key] = misty.cmdLiteral
	// 		misty.cmdIndex = append(misty.cmdIndex, [2]string{key, value[1]})
	// 	}
	// }

	// fmt.Println("updateCommands done. Command count: [" + strconv.Itoa(len(misty.cmdFuncs)) + "].")
}

func (misty *Misty) syncLiteralCommands() {
	// Sync LStrings.
	fmt.Print("Syncing LiteralCommands Data...")
	workSheetXMLContent, err := fetchFeed(misty.conf.LiteralCommandSheetURL())

	// All tabs' GSeetData.
	sheetData := []gshelp.GSheetData{}

	if err != nil {
		//Oh carp!
		fmt.Println(Red("[Error] ") + err.Error())
	} else {
		fmt.Println(Green("[Complete]"))
		URLs := gshelp.WorkSheetFeedToCellFeedURLs(workSheetXMLContent)

		// Get all cellfeeds.
		for i, URL := range URLs {
			fmt.Print("[Fetching Tab] : [" + strconv.Itoa(i) + "]...")
			cellXMLContent, err := fetchFeed(URL)
			if err != nil {
				fmt.Println(Red("[Error] ") + err.Error())
			} else {
				tabData := gshelp.CellFeedToGSheetData(cellXMLContent)

				// Store in the golbal var.
				sheetData = append(sheetData, tabData)
				fmt.Println(Green("[Complete]"))
			}
		}
	}

	// This will empty this container.
	misty.literalCommands = make(map[string][2]string)

	// Iterate through tabs.
	for _, sheetTab := range sheetData {
		// Iterate through rows.
		for _, row := range sheetTab.StringTable {
			// Check if each row has an ID.
			if len(row) > 1 {
				if row[0] != "" && row[1] != "" && row[2] != "" {
					// Add this row. {Content Str ID, Desc Str ID}
					misty.literalCommands[row[0]] = [2]string{row[1], row[2]}
					// {Name, Desc Str ID}
					misty.cmdIndex = append(misty.cmdIndex, [2]string{row[0], row[2]})
					misty.cmdFuncs[row[0]] = misty.cmdLiteral
				}
			}
		}
	}
}

func (misty *Misty) startObserveStreamingStatus() {
	//Observe the watching mixer channel.
	mixerTicker := time.NewTicker(time.Second * 3)
	go func() {
		for _ = range mixerTicker.C {
			//Prevent observing when the bot is updating or do not have a Mixer channel name.
			if !misty.Updating && misty.conf.WatchingMixerChannel != "" {
				isOnline, err := isMixerChannelOnline(misty.conf.WatchingMixerChannel)
				if err != nil {
					fmt.Println(err)
				}
				//Compare to the cache status vars.
				if isOnline {
					if !misty.streamingStatus.MixerOnline {
						//Watching channel become online. Inform this in the resident channel.
						misty.streamingStatus.MixerOnline = true

						informMessage := misty.Line("mixerStreamingOnline", 0) + "\n"
						informMessage += mixerChannelURLPrefix + misty.conf.WatchingMixerChannel
						misty.deletePreviousBroadcastMessage(misty.Line("mixerStreamingOnline", 0))
						misty.broadcastMessage(informMessage)
					} //Okey. Do nothing.
				} else {
					if misty.streamingStatus.MixerOnline {
						//Watching channel become online. Inform this in the resident channel.
						misty.streamingStatus.MixerOnline = false
						informMessage := misty.Line("mixerStreamingOffline", 0)
						misty.deletePreviousBroadcastMessage(misty.Line("mixerStreamingOffline", 0))
						misty.broadcastMessage(informMessage)
					}
				}
			}
		}
	}()

	//Observe the watching Smashcast channel.
	smashcastTicker := time.NewTicker(time.Second * 40)
	go func() {
		for _ = range smashcastTicker.C {
			//Prevent observing when the bot is updating or do not have a Smashcast channel name.
			if !misty.Updating && misty.conf.WatchingSmashcastChannel != "" {
				isOnline, err := isSmashcastChannelOnline(misty.conf.WatchingSmashcastChannel)
				if err != nil {
					fmt.Println(err)
				}
				//Compare to the cache status vars.
				if isOnline {
					if !misty.streamingStatus.SmashcastOnline {
						//Watching channel become online. Inform this in the resident channel.
						misty.streamingStatus.SmashcastOnline = true

						informMessage := misty.Line("smashcastStreamingOnline", 0) + "\n"
						informMessage += smashcastChannelURLPrefix + misty.conf.WatchingSmashcastChannel
						misty.deletePreviousBroadcastMessage(misty.Line("smashcastStreamingOnline", 0))
						misty.broadcastMessage(informMessage)
					} //Okey. Do nothing.
				} else {
					if misty.streamingStatus.SmashcastOnline {
						//Watching channel become online. Inform this in the resident channel.
						misty.streamingStatus.SmashcastOnline = false
						informMessage := misty.Line("smashcastStreamingOffline", 0)
						misty.deletePreviousBroadcastMessage(misty.Line("smashcastStreamingOffline", 0))
						misty.broadcastMessage(informMessage)
					}
				}
			}
		}
	}()

	//Observe the watching Twitch channel.
	twitchTicker := time.NewTicker(time.Second * 50)
	go func() {
		for _ = range twitchTicker.C {
			//Prevent observing when the bot is updating or do not have a Twitch channel name.
			if !misty.Updating && misty.conf.WatchingTwitchChannel != "" {
				isOnline, err := isTwitchChannelOnline(misty.conf.WatchingTwitchChannel)
				if err != nil {
					fmt.Println(err)
				}
				//Compare to the cache status vars.
				if isOnline {
					if !misty.streamingStatus.TwitchOnline {
						//Watching channel become online. Inform this in the resident channel.
						misty.streamingStatus.TwitchOnline = true
						informMessage := misty.Line("twitchStreamingOnline", 0) + "\n"
						informMessage += twitchChannelURLPrefix + misty.conf.WatchingTwitchChannel
						misty.deletePreviousBroadcastMessage(misty.Line("twitchStreamingOnline", 0))
						misty.broadcastMessage(informMessage)
					} //Okey. Do nothing.
				} else {
					if misty.streamingStatus.TwitchOnline {
						//Watching channel become online. Inform this in the resident channel.
						misty.streamingStatus.TwitchOnline = false
						informMessage := misty.Line("twitchStreamingOffline", 0)
						misty.deletePreviousBroadcastMessage(misty.Line("twitchStreamingOffline", 0))
						misty.broadcastMessage(informMessage)
					}
				}
			}
		}
	}()

}
