package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/rewardian/keybd_event"
	"gopkg.in/yaml.v2"
)

var loadedShortcuts Shortcuts
var shortcutFile string
var tmpDir string = os.TempDir()

// Shortcut represents a Type of macro ("text" or "program")
// and the Value of what is to be output or executed.
type Shortcut struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

// Shortcuts is a map that uses the keybinding labels in keymaps.go
// to access a Shortcut Entry. e.g. Shortcuts.Entry['KEY_SPACE'].
type Shortcuts struct {
	Entry map[string]Shortcut `yaml:"shortcuts"`
}

func initializeShortcuts() {
	filename := "keyswitcher.yaml"
	shortcutFile = normalUser.HomeDir + "/" + filename
}

func loadShortcutFile(yamlFile string) {
	source, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		fmt.Printf("\nNOTE: Unable to access shortcut file.\n\tERROR: %v\n", err)
		return
	}

	err = yaml.Unmarshal(source, &loadedShortcuts)
	if err != nil {
		fmt.Printf("\nNOTE: %s appears to be malformed.\n\tERROR: %v\n", shortcutFile, err)
		return
	}

	fmt.Printf("\nShortcut file %s is loaded.\n", yamlFile)
}

func parseLoadedShortcuts() {
	for i, v := range loadedShortcuts.Entry {
		fmt.Printf("Key: %s\n\tType: %s \n\tValue: %s\n", i, v.Type, v.Value)
	}
}

func processShortcuts(shortcut Shortcut, keyPressed string) {
	if shortcut.Type == "text" {

		if detailedLogging > 0 {
			fmt.Printf("[%s]: Outputting text: %s\n", keyPressed, shortcut.Value)
		}
		err := KeyboardWrite(shortcut.Value)

		if err != nil {
			fmt.Printf("Unable to write text to keyboard for [%s]\n", keyPressed)
		}

	} else if shortcut.Type == "program" {
		var u32uid uint32 = uint32(uid)
		var u32gid uint32 = uint32(gid)
		var cmd *exec.Cmd

		if detailedLogging > 0 {
			fmt.Printf("[%s]: Executing program: %s\n", keyPressed, shortcut.Value)
		}

		args := strings.Split(shortcut.Value, " ")

		if len(args) == 1 {
			cmd = exec.Command(shortcut.Value)

		} else {
			tmpFile, err := ioutil.TempFile(tmpDir, "keyswitcher-")
			if err != nil {
				fmt.Println(err)
				return
			}

			tmpFile.Chown(uid, gid)
			tmpFile.Chmod(0700)
			tmpFile.WriteString(shortcut.Value)
			tmpFile.Close()

			// Linux-specific
			cmd = exec.Command("/bin/sh", "-c", tmpFile.Name())
		}

		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: u32uid, Gid: u32gid}
		cmd.Start()
	}
}

func defaultKeystroke(keyPressedCode int) {
	keybind, _ := keybd_event.NewKeyBonding()

	keybind.SetKeys(keyPressedCode)

	//Super shift!
	keybind.HasSuper(true)
	keybind.HasSHIFT(true)

	keybind.Launching()
}
