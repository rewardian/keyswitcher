package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

// Displays the primary menu options.
func displayHelp() {

	fmt.Println()
	fmt.Println("Options:\n---")
	for i := range suggestions {
		fmt.Println(suggestions[i].Text+"\t\t", suggestions[i].Description)
	}
	fmt.Println()
}

// Updates the CLI menu cursor.
func livePrefix() (string, bool) {
	return cliPrefix + "> ", true
}

func initializeEnv() {
	var err error
	if value := os.Getenv("USERNAME"); value != "" {
		normalUser, err = user.Lookup(os.Getenv("USERNAME"))
	} else {
		normalUser, err = user.Lookup(os.Getenv("SUDO_USER"))
	}

	if err != nil {
		fmt.Printf("\nNOTE: Unable to find environment variables.\nApplication may behave unexpectedly.\n\tERROR: %v\n", err)
	}

	uid, err = strconv.Atoi(normalUser.Uid)
	gid, err = strconv.Atoi(normalUser.Gid)

	if err != nil {
		fmt.Printf("\nNOTE: Unable to find normal user UID/GID.\n\tERROR: %v\n", err)
	}
}

// Output the list of keyboards we'll work with.
func listKeyboards() {
	var keyboards, err = pollKeyboards()

	if err != nil {
		fmt.Printf("\nNOTE: Unable to access keyboard devices. Are you running as root? \n\tERROR: %v\n", err)
	} else {

		fmt.Println()
		fmt.Println("Use the command 'grab ID' to take control of the device:")
		fmt.Println()
		for _, kb := range keyboards {
			fmt.Println(kb.ID, ":", kb.Name, "-", kb.Filepath)
		}
		fmt.Println()
	}
}

// Temporary executable shell scripts can be created by processShortcut() in shortcuts.go.
// Used to remove those temporary files, on program exit or release of the device.
func clearTempDir() error {
	tmpFiles, err := filepath.Glob(filepath.Join(tmpDir, "keyswitcher-*"))
	if err != nil {
		return err
	}

	for _, tmpFile := range tmpFiles {
		err = os.RemoveAll(tmpFile)
		if err != nil {
			return err
		}
	}

	return nil
}
