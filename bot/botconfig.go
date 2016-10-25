package bot

import "gitlab.com/ahtram/misty/gshelp"
import "strings"
import "errors"

const configKeyCmdPrefix = "commandPrefix"
const configKeyLineSheetID = "lineSheetID"
const configKeyLiteralCommandSheetID = "literalCommandSheetID"
const configKeyResidentDiscordChannelID = "residentDiscordChannelID"
const configKeyBeamWatchingChannelID = "beamWatchingChannelID"
const configKeyHitboxWatchingChannelID = "hitboxWatchingChannelID"

// botConfig stores the config values readed from our Google Sheet config file.
type botConfig struct {
	CommandPrefix            string
	lineSheetID              string
	literalCommandSheetID    string
	ResidentDiscordChannelID string
	WatchingBeamChannel      []string
	WatchingHitboxChannel    []string
}

// ToString output the object's content and return as a formated string.
func (conf *botConfig) ToString() string {
	var returnString = ""
	returnString += "CommandPrefix: [" + conf.CommandPrefix + "]\n"
	returnString += "lineSheetID: [" + conf.lineSheetID + "]\n"
	returnString += "literalCommandSheetID: [" + conf.literalCommandSheetID + "]\n"
	returnString += "ResidentDiscordChannelID: [" + conf.ResidentDiscordChannelID + "]\n"

	for _, v := range conf.WatchingBeamChannel {
		returnString += "WatchingBeamChannel: [" + v + "]\n"
	}

	for _, v := range conf.WatchingHitboxChannel {
		returnString += "WatchingHitboxChannel: [" + v + "]\n"
	}
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
				} else if row[0] == configKeyResidentDiscordChannelID {
					conf.ResidentDiscordChannelID = row[1]
				} else if row[0] == configKeyBeamWatchingChannelID {
					conf.WatchingBeamChannel = append(conf.WatchingBeamChannel, row[1])
				} else if row[0] == configKeyHitboxWatchingChannelID {
					conf.WatchingHitboxChannel = append(conf.WatchingHitboxChannel, row[1])
				}
			}
		}
	}

	// Check if everything is good.
	if isEmptyOrHasSpace(conf.CommandPrefix) {
		return errors.New("Oops! Illegal CommandPrefix in config file! Please fix this!")
	}

	if isEmptyOrHasSpace(conf.lineSheetID) {
		return errors.New("Oops! Illegal lineSheetID in config file! Please fix this!")
	}

	if isEmptyOrHasSpace(conf.literalCommandSheetID) {
		return errors.New("Oops! Illegal literalCommandSheetID in config file! Please fix this!")
	}

	if isEmptyOrHasSpace(conf.ResidentDiscordChannelID) {
		return errors.New("Oops! Illegal ResidentDiscordChannelID in config file! Please fix this!")
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
