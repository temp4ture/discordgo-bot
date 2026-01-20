// package in charge to load all commands.
package cmds

// make sure to append your own commands here for them to
// register when the bot boots up.
import (
	_ "discordgo-bot/commands/help"
	_ "discordgo-bot/commands/magic8"
	_ "discordgo-bot/commands/ping"
)
