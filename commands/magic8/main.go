// 8ball command.
//
// Usage: /8ball (question)
package magic8

// This file handles command interpretation.

import (
	"discordgo-bot/core/commands"

	"github.com/bwmarrin/discordgo"
)

// Handler when our command gets called via chat message.
func do_command_message(data *commands.DataMessage) error {
	question := data.Content
	author := data.Message.Author.DisplayName()

	embed := create_embed(question, author)
	// send our fancy embed, responding to our user without pinging them
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
	question := data.GetOptions()["question"].StringValue()
	author := data.Interaction.Member.User.DisplayName()

	embed := create_embed(question, author)
	// send our fancy embed, responding to our user without pinging them
	err := data.Session.InteractionRespond(
		data.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		},
	)
	return err
}
