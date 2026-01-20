package terminal

import (
	"bufio"
	"discordgo-bot/utils"
	"discordgo-bot/utils/ucolor"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	Session           *discordgo.Session
	terminal_commands []TerminalCommand
)

////////////////////////////////////////////////////////////
// Commands ////////////////////////////////////////////////

func init() {
	terminal_commands := []TerminalCommand{
		{
			Name:        "help",
			Usage:       "[Page]",
			Description: "Show this list.",
			Handle: func(args []string) (bool, error) {
				command_show_limit := 12
				pages_curr := 1
				pages_max := int(math.Ceil(
					float64(len(terminal_commands)) /
						float64(command_show_limit),
				))

				// interpret first arg. as page
				spage := utils.GetSliceStr(args, 0, "1")
				page, err := strconv.Atoi(spage)
				if err == nil && page > 0 {
					pages_curr = min(page, pages_max)
				}

				// generate our string to print
				var str_to_print strings.Builder
				str_to_print.WriteString(ucolor.SUBTITLE)
				for i, command := range terminal_commands {
					// offset according to our page & visible commands per
					if i < command_show_limit*(pages_curr-1) {
						continue
					} else if i+1 > command_show_limit*pages_curr {
						break
					}
					tusage := ""
					// append usage if we have one
					if len(command.Usage) > 1 {
						tusage += " " + ucolor.ITALIC +
							command.Usage + ucolor.RESET + ucolor.SUBTITLE
					}
					fmt.Fprintf(
						&str_to_print, "%s%s - %s\n",
						command.Name, tusage, command.Description,
					)
				}

				fmt.Fprintf(
					&str_to_print, "%s%sPage %d of %d%s\n",
					ucolor.RESET, ucolor.BOLD, pages_curr,
					pages_max, ucolor.RESET,
				)
				fmt.Println(str_to_print.String())
				return true, nil
			},
		},
		{
			// fake quit command, actually handled by "interpret" func.
			// we keep this here so our help cmd shows what quit does.
			Name:        "quit",
			Usage:       "",
			Description: "Stop running the bot.",
			Handle:      func(args []string) (bool, error) { return true, nil },
		},
		{
			Name:        "speak",
			Usage:       "(ChannelID) (Message...)",
			Description: "Send a message to a channel.",
			Handle: func(args []string) (bool, error) {
				if len(args) < 2 {
					return false, nil
				}
				channel := args[0]
				message := strings.Join(args[1:], " ")
				_, err := Session.ChannelMessageSend(channel, message)
				return true, err
			},
		},
		{
			Name:        "clear",
			Description: "Clear the terminal.",

			Handle: func(args []string) (bool, error) {
				// different os' clear their terminals with different commands
				os_clearcmd := map[string]*exec.Cmd{
					"linux":   exec.Command("clear"),
					"windows": exec.Command("cmd", "/c", "cls"),
				}
				// get the proper command by fetching our os
				os_system := runtime.GOOS
				clear_cmd, success := os_clearcmd[os_system]
				if !success {
					// we fallback to linux as unix is the most common
					// os distribution out there, i believe.
					clear_cmd = os_clearcmd["linux"]
				}

				clear_cmd.Stdout = os.Stdout
				return true, clear_cmd.Run()
			},
		}}
	// register all commands listed
	for _, cmd := range terminal_commands {
		RegisterTerminalCommand(cmd)
	}
}

// Commands ////////////////////////////////////////////////
////////////////////////////////////////////////////////////

type TerminalCommand struct {
	Name        string
	Usage       string
	Description string // pref. 1st person
	// returns proper usage and error.
	// if false, handler will print the proper usage of the command.
	Handle func(args []string) (bool, error)
}

// Remove all nasty characters incoming from our console inputs.
func sanitizeInput(message string) string {
	return strings.TrimSuffix(strings.TrimSuffix(message, "\n"), "\r")
}

func getQuitShortcut() string {
	// show a different quit shortcut depending on the platform, as
	// it varies between unix and.. windows.
	os_quitsc := map[string]string{
		"linux":   "Ctrl + D",
		"windows": "Ctrl + C",
	}
	os_system := runtime.GOOS
	shortcut_str, success := os_quitsc[os_system]
	if !success {
		// fallback on linux because of unix likeliness.
		shortcut_str = os_quitsc["linux"]
	}
	return shortcut_str
}

// Start our terminal loop.
func Start() bool {
	// capture os.Interrupt to prevent hard quitting
	signal.Notify(make(chan os.Signal, 1), os.Interrupt)

	var help_string strings.Builder
	// build a pretty string using colors
	fmt.Fprintf(
		&help_string, "\nEnter %shelp%s for a list of available commands.\n",
		ucolor.OKBLUE, ucolor.RESET,
	)
	shortcut_str := getQuitShortcut()
	fmt.Fprintf(
		&help_string,
		"Quit the program by pressing '%s%s%s' or using %squit%s.\n",
		ucolor.OKCYAN, shortcut_str, ucolor.RESET, ucolor.OKBLUE, ucolor.RESET,
	)
	// finally, we print it!
	fmt.Println(help_string.String())

	run := true
	// read terminal inputs for as long as we're alive
	for {
		if !run {
			break
		}
		// create an input reader and wait for a command
		input_reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		input, err := input_reader.ReadString('\n')
		if err != nil {
			// os.Interrupt will land us here
			break
		}

		// input will contain a trailing line break as a side
		// effect to our reader: clean it before processing it
		input = sanitizeInput(input)
		code, err := interpret_terminal(input)
		if err != nil {
			// soft print errors
			log.Println(err)
		}

		// handle a quit command exit code
		switch code {
		case -1:
			run = false
		}
	}
	return true
}

// Start our terminal loop without an actual terminal.
// Returns manual exit as bool.
func Start_No_Terminal() bool {
	// capture 'os.Interrupt' to prevent hard quitting
	exit_signal := make(chan os.Signal, 1)
	signal.Notify(exit_signal, os.Interrupt)

	fmt.Printf(
		"Quit the program by pressing '%s%s%s'.\n",
		ucolor.OKCYAN,
		getQuitShortcut(),
		ucolor.RESET,
	)

	<-exit_signal
	return true
}

// Interpet commands sent via terminal.
// Returns (bool, error).
//
//	0: no command match | 1: command match | -1: quit executed
//	error: error returned by command execution
func interpret_terminal(message string) (int, error) {
	// ignore whitespace, return on nothing
	if parsed, _ := regexp.MatchString(`\w+`, message); !parsed {
		return 0, nil
	}

	message_spliced := strings.Split(message, " ")
	command_name := strings.ToLower(message_spliced[0])
	args := message_spliced[1:]

	if command_name == "quit" {
		return -1, nil
	}
	for _, terminal_command := range terminal_commands {
		if command_name == terminal_command.Name {
			ok, err := terminal_command.Handle(args)
			if !ok {
				// print command usage if formatted wrong
				fmt.Printf(
					"%sUsage: %s %s%s\n",
					ucolor.BOLD, terminal_command.Name,
					terminal_command.Usage, ucolor.RESET,
				)
			}
			return 1, err
		}
	}
	fmt.Printf(
		"%serror: \"%s\" not recognized as a terminal command.%s\n",
		ucolor.FAIL, command_name, ucolor.RESET,
	)
	return 0, nil

}

// Register terminal command.
func RegisterTerminalCommand(cmd TerminalCommand) {
	terminal_commands = append(terminal_commands, cmd)
}
