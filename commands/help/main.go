// Help command.
//
// Usage: /help [page]
package help

// This file handles command interpretation.

import (
	"discordgo-bot/core/commands"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Handler when our command gets called via chat message.
func doCommandMessage(data *commands.DataMessage) error {
	// get our first parameter and use it as our page number
	// todo: need a better way of parsing chat parameters
	parameters := strings.Split(data.Content, " ")
	page, err := strconv.Atoi(parameters[0])
	if err != nil {
		page = 1
	}

	embed := createEmbed(page)
	// send our fancy embed, responding to our user without pinging them
	_, err = data.Session.ChannelMessageSendComplex(
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
	// get our 'page' interaction option and default to 1 if none was provided.
	var page int = 1
	if option, ok := data.GetOptions()["page"]; ok {
		page = int(option.IntValue())
	}

	embed := createEmbed(page)
	// send our fancy embed, responding to our user without pinging them
	err := data.Session.InteractionRespond(data.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	},
	)
	return err
}
