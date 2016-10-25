package gshelp

// SheetIDToFeedURL returns the Google Sheet feed URL by sheetID.
func SheetIDToFeedURL(sheetID string) string {
	return "https://spreadsheets.google.com/feeds/worksheets/" + sheetID + "/public/full"
}
