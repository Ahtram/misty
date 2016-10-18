package main

import "fmt"
import "github.com/fatih/color"

//Version number and program name define.
const programName = "Misty"
const version = "0.0.0.3"

//Color defines.
var mag = color.New(color.FgHiMagenta).SprintFunc()
var yellow = color.New(color.FgHiYellow).SprintFunc()
var cyan = color.New(color.FgHiCyan).SprintFunc()
var green = color.New(color.FgHiGreen).SprintFunc()

// PrintWelcomeMessage does what is says.
func PrintWelcomeMessage() {
	//Plain welcome text.
	welcomeText := "*   " + programName + " Version [" + version + "]   *"

	//The colored welcome text.
	colorWelcomeText := fmt.Sprintf("%s   %s Version [%s]   %s", cyan("*"), mag(programName), yellow(version), cyan("*"))

	fmt.Println(cyan(starStringLine(len(welcomeText))))
	fmt.Println(colorWelcomeText)
	fmt.Println(cyan(starStringLine(len(welcomeText))))
}

// starStringLine print a line formed by "*" with specified length.
func starStringLine(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = '*'
	}
	return string(b)
}
