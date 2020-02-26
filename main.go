package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

var suggestions = []prompt.Suggest{
	// Commands
	{Text: "list", Description: "Display available input devices"},
	{Text: "grab", Description: "Gain exclusive control over input device"},
	{Text: "load", Description: "Load a specific YAML configuration file [default: $HOME/keyswitcher.yaml]"},
	{Text: "release", Description: "Release any activated keyboards"},
	{Text: "parse", Description: "Print the shortcuts loaded in keyswitcher"},
	{Text: "log", Description: "Toggle whether shortcuts, keypresses, or none are logged to stdout"},
	{Text: "exit", Description: "Exit"},
}

var detailedLogging int = 1 // 2 "Verbose", 1 "Normal", 0 "Quiet"
var releaseKeyboard = make(chan bool)

// Provides home directory and sets ownership on program execution
var normalUser *user.User
var uid, gid int

// Updates terminal title and CLI cursor
var cliTitle, cliPrefix string = "keyswitcher", "keyswitcher"

func executor(in string) {
	in = strings.TrimSpace(in)

	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "list":
		listKeyboards()
	case "grab":
		if len(blocks) != 2 {
			listKeyboards()
		} else {
			device, err := parseDevicePath(blocks[1])
			if err != nil {
				fmt.Printf("\nNOTE: Unable to access device. \n\tERROR: %v\n", err)
			} else {
				cliPrefix = "" + device.Name + " - [active]"
				livePrefix()

				go grabDevice(device)
			}
		}
	case "load":
		if len(blocks) > 1 {
			loadShortcutFile(blocks[1])
		} else {
			loadShortcutFile(shortcutFile)
		}
	case "release":
		fmt.Println("Press a key on the captured keyboard(s) to release.")
		cliPrefix = cliTitle
		livePrefix()
		releaseKeyboard <- true
	case "log":
		if detailedLogging == 0 {
			detailedLogging = 1
			fmt.Println("\nNormal logging enabled. Executed shortcuts will print to stdout.")
			fmt.Println()
		} else if detailedLogging == 1 {
			detailedLogging = 2
			fmt.Println("\nVerbose logging enabled. Both keypresses and shortcuts will print to stdout.")
			fmt.Println()
		} else {
			detailedLogging = 0
			fmt.Println("\nLogging has been disabled. No logging to stdout when using a keyboard.")
			fmt.Println()
		}
	case "parse":
		parseLoadedShortcuts()
	case "exit":
		fmt.Println("Bye!")
		clearTempDir()
		os.Exit(0)
	default:
		displayHelp()
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

func main() {

	initializeEnv()
	initializeShortcuts()
	loadShortcutFile(shortcutFile)
	displayHelp()

	p := prompt.New(
		executor,
		completer,
		prompt.OptionLivePrefix(livePrefix),
		prompt.OptionTitle(cliTitle),
	)

	p.Run()
	defer clearTempDir()
}
