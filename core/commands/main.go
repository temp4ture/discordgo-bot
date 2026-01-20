// Package providing slash & chat command functionality.
package commands

import (
	"discordgo-bot/globals"
	"discordgo-bot/utils/uembed"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	ENABLE_CHAT_COMMANDS  bool = true
	ENABLE_SLASH_COMMANDS bool = true
	// Unregister all slash commands out on exit.
	//
	// NOTE: Bots are rate limited to making 200 app commands per day, per guild.
	// Don't use on a big command list unless you need to remove them all, or for cache reasons.
	CLEAR_SLASH_ON_EXIT bool = false
	ChatPrefix               = []string{"+"}

	command_map = map[string]*CommandEntry{}
	regcmd_map  = map[string]*discordgo.ApplicationCommand{}
)

type CommandEntry struct {
	// 'discordgo.AppCommand' attached to this command.
	//
	// Commands will be automatically generate from it.
	AppCommand discordgo.ApplicationCommand
	Aliases    []string // Alternate names usable via chat commands.
	// Command function to call when invoked via chat message.
	FuncMessage func(data *DataMessage) error
	// Command function to call when invoked via slash command.
	FuncInteraction func(data *DataInteraction) error
}

// data struct for message calls containing special functions
type DataMessage struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
	Content string
}

// data struct for interaction calls containing special functions
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

func handle_command_chat(session *discordgo.Session, message *discordgo.MessageCreate) {
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
			command, ok, content = find_command_entry(content)
			break
		}
	}
	if !ok {
		return
	}

	err := command.FuncMessage(&DataMessage{session, message, content})
	if err != nil {
		log.Printf("error executin command: %s\n", err)
		uembed.GenerateErrorMessage(1)
	}
}

func find_command_entry(content string) (*CommandEntry, bool, string) {
	if len(content) < 1 {
		log.Println("can't run 'chcmd_match_command' with an empty string")
		return nil, false, ""
	}
	var trimmed string = content
	var to_trim string
	var command *CommandEntry

	// dissect and use the first word from our content string
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
		trimmed = strings.TrimSuffix(trimmed, to_trim)
		trimmed = strings.TrimSpace(trimmed)
		return command, true, trimmed
	}
	return nil, false, ""
}

func handle_command_slash(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// find a command and run it with our given parameters
	if cmd_entry, ok := command_map[interaction.ApplicationCommandData().Name]; ok {
		err := cmd_entry.FuncInteraction(&DataInteraction{session, interaction})
		if err != nil {
			log.Printf("error executin command: %s\n", err)
			uembed.GenerateErrorMessage(2)
		}
	}
}

// Register a new command.
func Register(cmd CommandEntry) {
	cmd_name := cmd.AppCommand.Name
	// only allow registering commands before the bot starts operations
	// it's overall a pretty bad practice to do otherwise
	if globals.IsRunning() {
		log.Println("can't register commands on runtime.")
		return
	}
	// prevent command name conflicts
	if _, ok := command_map[cmd_name]; ok {
		log.Panicf(
			"error: trying to register command with name \"%s\" twice.",
			cmd_name,
		)
	}
	command_map[cmd_name] = &cmd // register our CommandEntry type
}

// Return our list of registered commands.
//
// Useful for help-like commands that require the knowledge of external commands.
func GetCommandEntries() map[string]*CommandEntry {
	return command_map
}

// Start the creation and listening for commands.
func InitCommands() {
	s := globals.Session

	if !ENABLE_CHAT_COMMANDS && !ENABLE_SLASH_COMMANDS {
		return
	}
	if ENABLE_SLASH_COMMANDS {
		log.Println("Registering slash commands...")
		for i, ecmd := range command_map {
			cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", &ecmd.AppCommand)
			if err != nil {
				log.Panicf("Failed to register command \"%v\": %v", ecmd.AppCommand.Name, err)
			}
			regcmd_map[i] = cmd
		}
		s.AddHandler(handle_command_slash)
		log.Print("Slash command registration successful!")
	}
	if ENABLE_CHAT_COMMANDS {
		s.AddHandler(handle_command_chat)
		log.Println("Listening for commands in chat!")
	}
}

// Unregister all slash commands out of our bot.
//
// NOTE: Bots are rate limited to making 200 app commands per day, per guild.
// Don't use on a big command list unless you need to remove them all, or for cache reasons.
func ClearSlashCommands() {
	s := globals.Session

	if CLEAR_SLASH_ON_EXIT {
		log.Println("Removing commands...")
		for _, v := range regcmd_map {
			err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}
