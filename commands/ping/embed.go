// Ping command.
//
// Usage: /ping
package ping

// This file is in charge of creating an embed using the ping from our bot's session.

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Create an embed using the provided session.
func create_embed(s *discordgo.Session) *discordgo.MessageEmbed {
	var embed *discordgo.MessageEmbed
	var em_description string = fmt.Sprintf(
		"# Pong! 🏓\n-# %dms response time",
		s.HeartbeatLatency().Milliseconds(),
	)

	embed = &discordgo.MessageEmbed{
		Title:       " ",
		Description: em_description,
		Color:       0x41aa0e,
	}
	return embed
}
