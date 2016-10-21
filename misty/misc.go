package misty

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

//Version number and program name define.
const programName = "Misty"
const version = "0.0.0.4"

// Color defines.
var Mag = color.New(color.FgHiMagenta).SprintFunc()
var Yellow = color.New(color.FgHiYellow).SprintFunc()
var Cyan = color.New(color.FgHiCyan).SprintFunc()
var Green = color.New(color.FgHiGreen).SprintFunc()

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
	Email    string
	Password string
	Token    string
}

// PrintWelcomeMessage does what is says.
func PrintWelcomeMessage() {
	//Plain welcome text.
	welcomeText := "*   " + programName + " Version [" + version + "]   *"

	//The colored welcome text.
	colorWelcomeText := fmt.Sprintf("%s   %s Version [%s]   %s", Cyan("*"), Mag(programName), Yellow(version), Cyan("*"))

	fmt.Println(Cyan(starStringLine(len(welcomeText))))
	fmt.Println(colorWelcomeText)
	fmt.Println(Cyan(starStringLine(len(welcomeText))))
}

// starStringLine print a line formed by "*" with specified length.
func starStringLine(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = '*'
	}
	return string(b)
}

// fetchFeed read feed content from given URL.
func fetchFeed(feedURL string) (content string, err error) {
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
