// 8ball command.
//
// Usage: /8ball (question)
package magic8

import (
	"discordgo-bot/core/commands"
	"fmt"
	"math/rand/v2"

	"github.com/bwmarrin/discordgo"
)

var (
	resp_positive = []string{
		"Yes!",
	}
	resp_neutral = []string{
		"Maybe...",
	}
	resp_negative = []string{
		"No.",
	}
)

func do_command_message(data *commands.DataMessage) error {
	question := data.Content
	author := data.Message.Author.DisplayName()

	embed := create_embed(question, author)
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
	question := data.GetOptions()["question"].StringValue()
	author := data.Interaction.Member.User.DisplayName()

	embed := create_embed(question, author)
	err := data.Session.InteractionRespond(data.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	},
	)
	return err
}

func create_embed(question string, user string) *discordgo.MessageEmbed {
	var embed *discordgo.MessageEmbed

	// pick a response!
	r_ipool := [][]string{resp_positive, resp_neutral, resp_negative}
	r_spool := r_ipool[rand.IntN(len(r_ipool))]
	response := r_spool[rand.IntN(len(r_spool))]

	em_content := fmt.Sprintf(
		"**%s's question:** %s\n**answer:** %s",
		user, question, response,
	)
	embed = &discordgo.MessageEmbed{
		Title:       "8 Ball",
		Description: em_content,
	}

	return embed
}

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
