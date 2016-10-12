package main

func init() {
	//Read bot vars.
	ScanVars()
}

func main() {
	// A welcome message with version number.
	PrintWelcomeMessage()

	// Sync Google Doc data.
	SyncGSData()

	// Start the bot.
	StartBot()

	return
}
