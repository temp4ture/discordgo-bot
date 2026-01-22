// Help command.
//
// Usage: /help [page]
package help

// This file tells our bot to register this package as a command.

import (
	"discordgo-bot/core/commands"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands.Register(commands.CommandEntry{
		AppCommand: discordgo.ApplicationCommand{
			Name:        "help",
			Description: "Show a list and usage of available commands.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "page",
					Description: "Page to display.",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
				},
			},
		},
		FuncMessage:     doCommandMessage,
		FuncInteraction: doCommandInteraction,
	})
}
