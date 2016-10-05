package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

func main() {
	mag := color.New(color.FgHiMagenta).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()
	fmt.Printf("%s start to make some %s....", mag("Misty"), green("potions"))
	discordgo.New("abx")
}
