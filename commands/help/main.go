// Help command.
//
// Usage: /help [page]
package help

import (
	"discordgo-bot/core/commands"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	// Amount of commands to be rendered per command call.
	//
	// Higher is more information at the cost of chat visibility.
	commands_per_page int = 11
	//
	em_cmdprefix string = "/"
)

func do_command_message(data *commands.DataMessage) error {
	parameters := strings.Split(data.Content, " ")
	page, err := strconv.Atoi(parameters[0])
	if err != nil {
		page = 1
	}

	embed := create_embed(page)

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

func do_command_interaction(data *commands.DataInteraction) error {
	// get our 'page' interaction option and default to 1 if none was provided.
	var page int = 1
	if option, ok := data.GetOptions()["page"]; ok {
		page = int(option.IntValue())
	}

	embed := create_embed(page)
	err := data.Session.InteractionRespond(data.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	},
	)
	return err
}

// Create and return a pretty 'help' embed on a specific page.
func create_embed(page int) *discordgo.MessageEmbed {
	allcommands := commands.GetCommandEntries()
	var page_active int = 1
	page_max := int(
		math.Ceil(float64(len(allcommands)) /
			float64(commands_per_page)),
	)
	page_active = max(1, min(page_max, page))

	// generate description from core
	em_description := ""
	i := 0 // ranging maps doesnt return a len variable so...
	for _, cmd := range allcommands {
		i++
		if i < commands_per_page*(page_active-1) { // offset our command list using help_page & commands_per_page
			continue
		} else if i+1 > commands_per_page*page_active { // stop generating if we go past our cmds limit
			break
		}

		var str_options string
		// generate usage line using command options
		command_options := cmd.AppCommand.Options
		for _, option := range command_options {
			str_options += fmt.Sprintf(" *[%s]*", option.Name)
		}
		em_description += fmt.Sprintf(
			"%s%s%s\n-# %s\n\n",
			em_cmdprefix, cmd.AppCommand.Name,
			str_options, cmd.AppCommand.Description,
		)
	}
	// footer showing our active page
	em_footer := fmt.Sprintf(
		"Page %d of %d",
		page_active,
		page_max,
	)

	embed := &discordgo.MessageEmbed{
		Title:       "Available Commands:",
		Description: em_description,
		Footer: &discordgo.MessageEmbedFooter{
			Text: em_footer,
		},
		Color: 0x41aa0e,
	}
	return embed
}

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
		FuncMessage:     do_command_message,
		FuncInteraction: do_command_interaction,
	})
}
