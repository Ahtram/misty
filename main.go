package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

//Version number and program name define.
var programName = "Misty"
var version = "0.0.0.1"

//Color define.
var mag = color.New(color.FgHiMagenta).SprintFunc()
var yellow = color.New(color.FgHiYellow).SprintFunc()
var cyan = color.New(color.FgHiCyan).SprintFunc()

// red := color.New(color.FgHiRed).SprintFunc()
// green := color.New(color.FgHiGreen).SprintFunc()
// blue := color.New(color.FgHiBlue).SprintFunc()

func main() {
	printWelcomeMessage()
	discordgo.New("abx")
}

func printWelcomeMessage() {
	//Plain welcome text.
	welcomeText := "*   " + programName + " Version [" + version + "]   *"

	//The colored welcome text.
	colorWelcomeText := fmt.Sprintf("%s   %s Version [%s]   %s", cyan("*"), mag(programName), yellow(version), cyan("*"))

	fmt.Println(cyan(starStringLine(len(welcomeText))))
	fmt.Println(colorWelcomeText)
	fmt.Println(cyan(starStringLine(len(welcomeText))))
}

func starStringLine(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = '*'
	}
	return string(b)
}
