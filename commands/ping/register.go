// Ping command.
//
// Usage: /ping
package ping

// This file tells our bot to register this package as a command.

import (
	"discordgo-bot/core/commands"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands.Register(commands.CommandEntry{
		AppCommand: discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Pong! Responds with response latency.",
		},
		FuncMessage:     doCommandMessage,
		FuncInteraction: doCommandInteraction,
	})
}
