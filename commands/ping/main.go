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
func doCommandMessage(data *commands.DataMessage) error {
	embed := createEmbed(data.Session)

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
func doCommandInteraction(data *commands.DataInteraction) error {
	embed := createEmbed(data.Session)

	err := data.Session.InteractionRespond(data.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	return err
}
