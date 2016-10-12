package gshelp

import "encoding/xml"

// ToFormattedString returns an output string just for print and debuging.
func (sd GSheetData) ToFormattedString(columnSeperator string, rowSeperator string) string {

	var returnString = ""
	for i, row := range sd.stringTable {
		for j, cell := range row {
			returnString += cell
			if j < len(row)-1 {
				returnString = returnString + columnSeperator
			}
		}

		if i < len(sd.stringTable)-1 {
			returnString = returnString + rowSeperator
		}
	}

	return returnString
}

// ToDefaultString returns an output string just for print and debuging using default seperactors.
func (sd GSheetData) ToDefaultString() string {
	return sd.ToFormattedString(" ", "\n")
}

// CellFeedToGSheetData converts XML content to GSheetData object for us.
func CellFeedToGSheetData(cellFeedXMLContent string) GSheetData {
	returnGSheetData := GSheetData{"", make([][]string, 0)}

	return returnGSheetData
}

// WorkSheetFeedToCellFeedURLs converts worksheet feed XML content to cell feed URLs for us.
func WorkSheetFeedToCellFeedURLs(worksheetFeedXMLContent string) []string {
	return make([]string, 0)
}
