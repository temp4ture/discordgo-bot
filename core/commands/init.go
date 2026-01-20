// Package providing slash & chat command functionality.
package commands

// This file has initialization related functions.

import (
	"discordgo-bot/globals"
	"log"
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
)

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
		s.AddHandler(handleCommandViaSlash)
		log.Print("Slash command registration successful!")
	}
	if ENABLE_CHAT_COMMANDS {
		s.AddHandler(handleCommandViaChat)
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
