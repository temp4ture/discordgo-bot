package uembed

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrorMessage = discordgo.MessageEmbed{
		Title:       "Error",
		Description: "Something wrong has happened!",
		Color:       0xff1313,
	}
)

// Generates an error message including non-explicit data of the probable cause.
func GenerateErrorMessage(id int) *discordgo.MessageEmbed {
	// todo: add more parameters
	error_code := fmt.Sprintf("0x0F%02dFFFF", id)
	error_embed := &ErrorMessage
	error_embed.Description =
		"We're sorry, something went wrong while processing your command...\n" +
			"-# code: " + error_code
	return error_embed
}
