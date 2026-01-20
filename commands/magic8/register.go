// 8ball command.
//
// Usage: /8ball (question)
package magic8

// This file tells our bot to register this package as a command.

import (
	"discordgo-bot/core/commands"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands.Register(commands.CommandEntry{
		AppCommand: discordgo.ApplicationCommand{
			Name:        "8ball",
			Description: "Responds to a yes / no question.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "question",
					Description: "What is your question?",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		FuncMessage:     do_command_message,
		FuncInteraction: do_command_interaction,
	})
}
