package bot

import "gitlab.com/ahtram/misty/gshelp"
import "strconv"
import "strings"
import "errors"

const configKeyCmdPrefix = "commandPrefix"
const configKeyLineSheetID = "lineSheetID"
const configKeyLiteralCommandSheetID = "literalCommandSheetID"
const configKeyOnlineNotify = "onlineNotify"
const configKeyResidentDiscordChannelID = "residentDiscordChannelID"
const configKeyBroadcastDiscordChannelID = "broadcastDiscordChannelID"
const configKeyMixerWatchingChannelID = "mixerWatchingChannelID"
const configKeySmashcastWatchingChannelID = "smashcastWatchingChannelID"
const configKeyTwitchWatchingChannelID = "twitchWatchingChannelID"
const configKeyUCloudHookEndPoint = "uCloudHookEndPoint"
const configKeyGitLabHookEndPoint = "gitLabHookEndPoint"
const configKeyGitHubHookEndPoint = "gitHubHookEndPoint"

// botConfig stores the config values readed from the Google Sheet config file.
type botConfig struct {
	CommandPrefix            string
	lineSheetID              string
	literalCommandSheetID    string
	onlineNotify             bool
	ResidentDiscordChannelID string
	BroadcastDiscrdChannelID []string
	WatchingMixerChannel     string
	WatchingSmashcastChannel string
	WatchingTwitchChannel    string
	UCloudConfigs            []*uCloudConfig
	GitLabConfigs            []*gitLabConfig
	GitHubConfigs            []*gitHubConfig
}

type uCloudConfig struct {
	UCloudHookEndPoint string
	UCloudHookPort     string
	UCloudAccessToken  string
}

type gitLabConfig struct {
	GitLabHookEndPoint string
	GitLabHookPort     string
}

type gitHubConfig struct {
	GitHubHookEndPoint string
	GitHubHookPort     string
}

// ToString output the object's content and return as a formated string.
func (conf *botConfig) ToString() string {
	var returnString = "=================== [Config] ====================\n"
	returnString += "CommandPrefix: [" + conf.CommandPrefix + "]\n"
	returnString += "lineSheetID: [" + conf.lineSheetID + "]\n"
	returnString += "literalCommandSheetID: [" + conf.literalCommandSheetID + "]\n"
	returnString += "onlineNotify: [" + strconv.FormatBool(conf.onlineNotify) + "]\n"
	returnString += "ResidentDiscordChannelID: [" + conf.ResidentDiscordChannelID + "]\n"
	returnString += "WatchingMixerChannel: [" + conf.WatchingMixerChannel + "]\n"
	returnString += "WatchingSmashcastChannel: [" + conf.WatchingSmashcastChannel + "]\n"
	returnString += "WatchingTwitchChannel: [" + conf.WatchingTwitchChannel + "]\n"
	for _, value := range conf.UCloudConfigs {
		returnString += "UCloudConfig: [" + value.UCloudHookEndPoint + "] [" + value.UCloudHookPort + "] [" + value.UCloudAccessToken + "] \n"
	}
	for _, value := range conf.GitLabConfigs {
		returnString += "GitLabConfig: [" + value.GitLabHookEndPoint + "] [" + value.GitLabHookPort + "] \n"
	}
	for _, value := range conf.GitHubConfigs {
		returnString += "GitHubConfig: [" + value.GitHubHookEndPoint + "] [" + value.GitHubHookPort + "] \n"
	}

	//Due to REST API limitation. Watching multiple channels may not be a good idea...
	for _, v := range conf.BroadcastDiscrdChannelID {
		returnString += "BroadcastDiscrdChannelID: [" + v + "]\n"
	}

	returnString += "================================================="

	//[Dep]: Due to REST API limitation. Watching multiple channels may not be a good idea...
	// for _, v := range conf.WatchingBeamChannel {
	// 	returnString += "WatchingBeamChannel: [" + v + "]\n"
	// }

	// for _, v := range conf.WatchingHitboxChannel {
	// 	returnString += "WatchingHitboxChannel: [" + v + "]\n"
	// }

	return returnString
}

// LineSheetURL returns the Line Sheet feed URL.
func (conf *botConfig) LineSheetURL() string {
	return gshelp.SheetIDToFeedURL(conf.lineSheetID)
}

// LiteralCommandSheetURL returns the LiteralCommand Sheet feed URL.
func (conf *botConfig) LiteralCommandSheetURL() string {
	return gshelp.SheetIDToFeedURL(conf.literalCommandSheetID)
}

// Setup do what it says.
func (conf *botConfig) Setup(sheetData []gshelp.GSheetData) error {

	//Clear the previous setting.
	conf.BroadcastDiscrdChannelID = conf.BroadcastDiscrdChannelID[:0]

	// Iterate through tabs.
	for _, sheetTab := range sheetData {
		// Iterate through rows.
		for _, row := range sheetTab.StringTable {
			// Check if each row has an ID.
			if len(row) > 0 {
				// Scan and set all config value.
				if row[0] == configKeyCmdPrefix {
					conf.CommandPrefix = row[1]
				} else if row[0] == configKeyLineSheetID {
					conf.lineSheetID = row[1]
				} else if row[0] == configKeyLiteralCommandSheetID {
					conf.literalCommandSheetID = row[1]
				} else if row[0] == configKeyOnlineNotify {
					conf.onlineNotify, _ = strconv.ParseBool(row[1])
				} else if row[0] == configKeyResidentDiscordChannelID {
					conf.ResidentDiscordChannelID = row[1]
				} else if row[0] == configKeyMixerWatchingChannelID {
					conf.WatchingMixerChannel = row[1]
				} else if row[0] == configKeySmashcastWatchingChannelID {
					conf.WatchingSmashcastChannel = row[1]
				} else if row[0] == configKeyTwitchWatchingChannelID {
					conf.WatchingTwitchChannel = row[1]
				} else if row[0] == configKeyBroadcastDiscordChannelID {
					if row[1] != "" {
						conf.BroadcastDiscrdChannelID = append(conf.BroadcastDiscrdChannelID, row[1])
					}
				} else if row[0] == configKeyUCloudHookEndPoint {
					uCloudConfig := uCloudConfig{
						UCloudHookEndPoint: row[1],
						UCloudHookPort:     row[2],
						UCloudAccessToken:  row[3],
					}
					conf.UCloudConfigs = append(conf.UCloudConfigs, &uCloudConfig)
				} else if row[0] == configKeyGitLabHookEndPoint {
					gitLabConfig := gitLabConfig{
						GitLabHookEndPoint: row[1],
						GitLabHookPort:     row[2],
					}
					conf.GitLabConfigs = append(conf.GitLabConfigs, &gitLabConfig)
				} else if row[0] == configKeyGitHubHookEndPoint {
					gitHubConfig := gitHubConfig{
						GitHubHookEndPoint: row[1],
						GitHubHookPort:     row[2],
					}
					conf.GitHubConfigs = append(conf.GitHubConfigs, &gitHubConfig)
				}
			}
		}
	}

	// Check if everything is good.
	if isEmptyOrHasSpace(conf.CommandPrefix) {
		return errors.New("Oops! Illegal CommandPrefix in config file! Please fix this. ")
	}

	if isEmptyOrHasSpace(conf.lineSheetID) {
		return errors.New("Oops! Illegal lineSheetID in config file! Please fix this! ")
	}

	if isEmptyOrHasSpace(conf.literalCommandSheetID) {
		return errors.New("Oops! Illegal literalCommandSheetID in config file! Please fix this! ")
	}

	if isEmptyOrHasSpace(conf.ResidentDiscordChannelID) {
		return errors.New("Oops! Illegal ResidentDiscordChannelID in config file! Please fix this! ")
	}

	return nil
}

func isEmptyOrHasSpace(str string) bool {
	if str == "" {
		return true
	}

	if strings.Contains(str, " ") {
		return true
	}

	return false
}
