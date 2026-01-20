package main

import (
	"discordgo-bot/core"
	"discordgo-bot/terminal"
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	ENV_FILE_PATH   = ".env"
	ENV_TOKEN_ENTRY = "BOT_TOKEN"
)

// load our '.env' file
func loadEnv() {
	err := godotenv.Load(ENV_FILE_PATH)
	if err != nil {
		log.Fatalf("failed to load .env file.\n%s", err)
	}
}

// get our discord token from our active env file using
// the token entry string we have.
func getToken() string {
	token, success := os.LookupEnv(ENV_TOKEN_ENTRY)
	if !success {
		log.Fatalf("token @ %s not found in .env file!", ENV_TOKEN_ENTRY)
	}
	return token
}

// run a discord bot session and return it or "nil".
func startBotSession(token string) any {
	log.Println("creating new discord bot session...")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("failed to create session.\n%s", err)
		return nil
	}
	session.Open()
	log.Println("session created successfully!")
	return session
}

// Exit if we quit manually from our terminal.
func doTerminalExit(manual_exit bool) {
	if manual_exit {
		os.Exit(0)
	}
}

func main() {
	var bot_session any = nil
	// var core_session any

	loadEnv()
	token := getToken()

	// run a loop to run our bot as long as we need to.
	for {
		bot_session = startBotSession(token)
		if session, ok := bot_session.(*discordgo.Session); ok {
			// on a successful launch, we'll get our core
			// services running to handle all requests.
			core.Start(session)
			// on a failed launch, we would immediately try again.

			terminal.Session = session
			var terminal_manual_exit bool
			// start a terminal and hold onto it
			if has_launch_arg("--no-terminal") {
				terminal_manual_exit = terminal.Start_No_Terminal()
			} else {
				terminal_manual_exit = terminal.Start()
			}
			// we stay up here until our terminal exits.

			// cleanup process when exiting
			fmt.Println()
			log.Println("shutting down session...")
			core.Stop()
			session.Close()
			time.Sleep(1.0)
			// quit our application if we did a manual exit
			doTerminalExit(terminal_manual_exit)
		}
	}
}

// return whether we launched with a specific argument.
func has_launch_arg(arg string) bool {
	return slices.Contains(os.Args, arg)
}
