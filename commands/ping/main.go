// Ping command.
//
// Usage: /ping
package ping

import (
	"discordgo-bot/core/commands"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func do_command_message(data *commands.DataMessage) error {
	embed := create_embed(data.Session)

	_, err := data.Session.ChannelMessageSendComplex(
		data.Message.ChannelID,
		&discordgo.MessageSend{
			Embed: embed,
			Reference: &discordgo.MessageReference{
				MessageID: data.Message.ID,
				ChannelID: data.Message.ChannelID,
				GuildID:   data.Message.GuildID,
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{
				RepliedUser: false,
			},
		},
	)
	return err
}

func do_command_interaction(data *commands.DataInteraction) error {
	embed := create_embed(data.Session)

	err := data.Session.InteractionRespond(data.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	return err
}

// create an embed using the provided session
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

func init() {
	commands.Register(commands.CommandEntry{
		AppCommand: discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Pong! Responds with response latency.",
		},
		FuncMessage:     do_command_message,
		FuncInteraction: do_command_interaction,
	})
}
