package gshelp

// A GSheetData represent one tab of sheet in your Google sheet doc.
type GSheetData struct {
	Title       string
	StringTable [][]string
}

// ToFormattedString returns an output string just for print and debuging.
func (sd GSheetData) ToFormattedString(columnSeperator string, rowSeperator string) string {

	var returnString = ""
	for i, row := range sd.StringTable {
		//Chech if the fist cell of this row is emety. If it is, ignore this row.
		if len(row) > 0 && row[0] == "" {
			continue
		}

		for j, cell := range row {
			returnString += cell
			if j < len(row)-1 {
				returnString = returnString + columnSeperator
			}
		}

		if i < len(sd.StringTable)-1 {
			returnString = returnString + rowSeperator
		}
	}

	return returnString
}

// ToDefaultString returns an output string just for print and debuging using default seperactors.
func (sd GSheetData) ToDefaultString() string {
	return sd.ToFormattedString(" ", "\n")
}
