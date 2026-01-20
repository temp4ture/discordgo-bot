// Package contains relevant global variables
package globals

import "github.com/bwmarrin/discordgo"

var (
	Session *discordgo.Session
	Running bool
)

func IsRunning() bool {
	return Running
}
