// Help command.
//
// Usage: /help [page]
package help

// This file is in charge of creating an embed using a 'page' parameter.

import (
	"discordgo-bot/core/commands"
	"fmt"
	"math"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	// Amount of commands to be rendered per command call.
	//
	// Higher is more information at the cost of chat visibility.
	commands_per_page int = 11
	// Prefix used for commands on the help embed.
	embed_command_prefix string = "/"
)

// Create and return a pretty 'help' embed on a specific page.
func create_embed(page int) *discordgo.MessageEmbed {
	allcommands := commands.GetCommandEntries()
	var page_active int = 1
	page_max := int(
		math.Ceil(float64(len(allcommands)) /
			float64(commands_per_page)),
	)
	page_active = max(1, min(page_max, page))

	var embed_description strings.Builder

	i := 0
	for _, command := range allcommands {
		i++
		if i < commands_per_page*(page_active-1) {
			// offset our command list using help_page & commands_per_page
			continue
		} else if i+1 > commands_per_page*page_active {
			// stop generating if we go past our cmds limit
			break
		}

		var str_options string
		// generate usage line using command options
		command_options := command.AppCommand.Options
		for _, option := range command_options {
			str_options += fmt.Sprintf(" *[%s]*", option.Name)
		}
		fmt.Fprintf(
			&embed_description, "%s%s%s\n-# %s\n\n",
			embed_command_prefix, command.AppCommand.Name,
			str_options, command.AppCommand.Description,
		)
	}
	// footer showing our active page
	embed_footer := fmt.Sprintf(
		"Page %d of %d",
		page_active,
		page_max,
	)
	// finally, we generate an embed with
	// our embed description and footer
	embed := &discordgo.MessageEmbed{
		Title:       "Available Commands:",
		Description: embed_description.String(),

		Footer: &discordgo.MessageEmbedFooter{
			Text: embed_footer,
		},
		Color: 0x41aa0e,
	}
	return embed
}
