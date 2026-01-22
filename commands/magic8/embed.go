// 8ball command.
//
// Usage: /8ball (question)
package magic8

// This file is in charge of creating an embed by showing the user's question along
// with a randomly selected answer from a list of positive, neutral and negative answers.

import (
	"fmt"
	"math/rand/v2"

	"github.com/bwmarrin/discordgo"
)

var (
	response_positive = []string{
		"Yes!",
	}
	response_neutral = []string{
		"Maybe...",
	}
	response_negative = []string{
		"No.",
	}
)

// Create a pretty embed with a user and their question.
func createEmbed(question string, user string) *discordgo.MessageEmbed {
	var embed *discordgo.MessageEmbed

	// randomly pick a response
	response_pools := [][]string{
		response_positive,
		response_neutral,
		response_negative,
	}
	selected_pool := response_pools[rand.IntN(len(response_pools))]
	response := selected_pool[rand.IntN(len(selected_pool))]

	// put that response in our embed's description
	// and return that sucker!
	embed_description := fmt.Sprintf(
		"**%s's question:** %s\n**answer:** %s",
		user, question, response,
	)
	embed = &discordgo.MessageEmbed{
		Title:       "8 Ball",
		Description: embed_description,
	}

	return embed
}
