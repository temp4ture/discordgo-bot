// Ping command.
//
// Usage: /ping
package ping

// This file handles command interpretation.

import (
	"discordgo-bot/core/commands"

	"github.com/bwmarrin/discordgo"
)

// Handler when our command gets called via chat message.
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

// Handler when our command gets called via discord's slash command.
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
