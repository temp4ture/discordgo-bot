// Package providing slash & chat command functionality.
package commands

// This file runs handling and logic functions.

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	command_map = map[string]*CommandEntry{}
	regcmd_map  = map[string]*discordgo.ApplicationCommand{}
)

type CommandEntry struct {
	// 'discordgo.AppCommand' attached to this command.
	//
	// Commands will automatically generate from it.
	AppCommand discordgo.ApplicationCommand
	Aliases    []string // Alternate names usable via chat commands.
	// Command function to call when invoked via chat message.
	FuncMessage func(data *DataMessage) error
	// Command function to call when invoked via slash command.
	FuncInteraction func(data *DataInteraction) error
}

// Data struct. for message interactions.
// Generated when a user calls a command via chat message.
type DataMessage struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
	Content string
}

// Data struct. for slash interactions.
// Generated when a user calls a command via Discord's slash commands.
type DataInteraction struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
}

// Return a map containing interaction options provided by the user.
func (o DataInteraction) GetOptions() map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options_list := o.Interaction.ApplicationCommandData().Options
	options_map := make(
		map[string]*discordgo.ApplicationCommandInteractionDataOption,
		len(options_list),
	)
	for _, option := range options_list {
		options_map[option.Name] = option
	}
	return options_map
}

func findCommandEntry(content string) (*CommandEntry, bool, string) {
	// return nothing if we don't got any content
	if len(content) < 1 {
		return nil, false, ""
	}
	var trimmed string = content
	var to_trim string
	var command *CommandEntry

	// split and use the first word from our content string
	cmdstr := strings.Split(content, " ")[0]
	for _, command_entry := range command_map {
		command_name := command_entry.AppCommand.Name
		if cmdstr == command_name {
			command = command_entry
			to_trim = command_name
		}
		// check for aliases too
		for _, alias := range command_entry.Aliases {
			if cmdstr == alias {
				command = command_entry
				to_trim = alias
			}
		}
	}
	if command != nil {
		// trim the command & any trailing spaces from our message output
		trimmed = strings.TrimLeft(trimmed, to_trim)
		trimmed = strings.TrimSpace(trimmed)
		return command, true, trimmed
	}
	return nil, false, ""
}

func handleCommandViaChat(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot { // ignore bots
		return
	}

	// message content minus prefix and command ran
	var content string
	// command to execute
	var command *CommandEntry
	var ok bool

	// confirm we're running a command by checking for prefixes
	// sanitize "content" if match is found
	for _, pref := range ChatPrefix {
		if content, ok = strings.CutPrefix(message.Content, pref); ok {
			// fetch for a command, replace content for trimmed
			command, ok, content = findCommandEntry(content)
			break
		}
	}
	if !ok {
		return
	}

	go func() {
		err := command.FuncMessage(&DataMessage{session, message, content})
		if err != nil {
			log.Printf("error executing command: %s\n", err)
			// uembed.GenerateErrorMessage(1)
		}
	}()
}

func handleCommandViaSlash(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
) {
	// find a command and run it with our given parameters
	if cmd_entry, ok := command_map[interaction.ApplicationCommandData().Name]; ok {
		go func() {
			err := cmd_entry.FuncInteraction(&DataInteraction{session, interaction})
			if err != nil {
				log.Printf("error executing command: %s\n", err)
				// uembed.GenerateErrorMessage(2)
			}
		}()
	}
}
