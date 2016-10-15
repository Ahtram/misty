package gshelp

import "encoding/xml"
import "fmt"
import "strings"
import "strconv"

type cellFeedEntryCell struct {
	XMLName xml.Name `xml:"cell"`
	Value   string   `xml:",chardata"`
	Row     string   `xml:"row,attr"`
	Col     string   `xml:"col,attr"`
}

type cellFeedEntry struct {
	Cell cellFeedEntryCell `xml:"cell"`
}

type cellFeed struct {
	XMLName  xml.Name        `xml:"feed"`
	Title    string          `xml:"title"`
	RowCount string          `xml:"rowCount"`
	ColCount string          `xml:"colCount"`
	Entries  []cellFeedEntry `xml:"entry"`
}

// CellFeedToGSheetData converts XML content to GSheetData object for us.
func CellFeedToGSheetData(cellFeedXMLContent string) GSheetData {
	// Parse the XML content.
	feed := cellFeed{}
	err := xml.Unmarshal([]byte(cellFeedXMLContent), &feed)
	returnGSheetData := GSheetData{"", make([][]string, 0)}

	if err != nil {
		fmt.Printf("error: %v", err)
	} else {
		tableHeight, _ := strconv.Atoi(feed.RowCount)
		tableWidth, _ := strconv.Atoi(feed.ColCount)
		returnGSheetData.Title = feed.Title
		//Allocate slices.
		for j := 0; j < tableHeight; j++ {
			returnGSheetData.StringTable = append(returnGSheetData.StringTable, make([]string, tableWidth))
		}

		//Assign all exist strings into table.
		for _, v := range feed.Entries {
			row, _ := strconv.Atoi(v.Cell.Row)
			col, _ := strconv.Atoi(v.Cell.Col)
			row = row - 1
			col = col - 1

			if row < len(returnGSheetData.StringTable) {
				if col < len(returnGSheetData.StringTable[row]) {
					returnGSheetData.StringTable[row][col] = v.Cell.Value
				}
			}
		}
	}

	return returnGSheetData
}

type workSheetFeedEntryLink struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type workSheetFeedEntry struct {
	Title string                   `xml:"title"`
	Links []workSheetFeedEntryLink `xml:"link"`
}

// workSheetFeed represent a simplified struct for a worksheet feed.
type workSheetFeed struct {
	XMLName     xml.Name             `xml:"feed"`
	ID          string               `xml:"id"`
	UpdatedTime string               `xml:"updated"`
	Entries     []workSheetFeedEntry `xml:"entry"`
}

// WorkSheetFeedToCellFeedURLs converts worksheet feed XML content to cell feed URLs for us.
func WorkSheetFeedToCellFeedURLs(worksheetFeedXMLContent string) []string {
	// Parse the XML content.
	feed := workSheetFeed{}
	err := xml.Unmarshal([]byte(worksheetFeedXMLContent), &feed)
	returnURLs := []string{}

	if err != nil {
		fmt.Printf("error: %v", err)
	} else {
		for _, entry := range feed.Entries {
			// fmt.Println("Got Tab: " + entry.Title)
			for _, link := range entry.Links {
				if strings.HasSuffix(link.Rel, "cellsfeed") {
					// fmt.Println("Link: " + link.Href)
					returnURLs = append(returnURLs, link.Href)
				}
			}
		}
	}

	return returnURLs
}
