package main

import (
	"fmt"
	"strconv"

	evdev "github.com/gvalkov/golang-evdev"
)

// Device represents an input device made available to the program.
type Device struct {
	ID             int
	Name, Filepath string
}

// pollKeyboards() generates a model of which input devices are deemed keyboards.
// If an input device exposes "KEY_ESC", we assume it's a valid keyboard.
// listKeyboards() is used for visual formatting.
//
// keyboardEventType is how we detect any key-like devices.
// keyboardEventCode serves to remove mice and other peripherals by looking for KEY_ESC.
//
// For more: https://www.kernel.org/doc/html/v4.15/input/event-codes.html &
// https://godoc.org/github.com/gvalkov/golang-evdev

func pollKeyboards() ([]Device, error) {
	var deviceID int = 1
	var keyboardEventType = evdev.CapabilityType{Type: 1, Name: "EV_KEY"}
	var keyboardEventCode = evdev.CapabilityCode{Code: 1, Name: "KEY_ESC"}
	var keyboards = []Device{}

	devices, err := evdev.ListInputDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {

		if deviceCapabilities, ok := device.Capabilities[keyboardEventType]; ok {
			for _, keyCode := range deviceCapabilities {
				if keyCode == keyboardEventCode {
					keyboards = append(keyboards, Device{deviceID, device.Name, device.Fn})
					deviceID++
				}
			}
		}

	}
	return keyboards, nil
}

func parseDevicePath(userInput string) (Device, error) {
	var keyboard Device

	deviceID, err := strconv.Atoi(userInput)
	keyboards, err := pollKeyboards()

	if err != nil {
		return keyboard, err
	}

	for _, keyboard := range keyboards {
		if deviceID == keyboard.ID {
			return keyboard, nil
		}
	}

	return keyboard, fmt.Errorf("could not find this device path or ID: %s", userInput)
}

func grabDevice(device Device) {
	keyboard, err := evdev.Open(device.Filepath)

	if err != nil {
		fmt.Printf("\nNOTE: Unable to take control of keyboard.\n\tERROR: %v\n", err)
		return
	}

	keyboard.Grab()
	defer keyboard.Release()
	defer clearTempDir()

	for {
		select {
		case <-releaseKeyboard:
			return
		default:
			keyEvent, _ := keyboard.ReadOne()

			// Value '01' is key depressed; Type 'EV_KEY' is keyboard event
			if keyEvent.Type == evdev.EV_KEY && keyEvent.Value == 01 {

				keyPressedCode := int(keyEvent.Code)
				keyPressed := keymaps[keyPressedCode]

				if detailedLogging == 2 {
					fmt.Printf("Key %s was pressed!\n", keyPressed)
				}

				if shortcut, ok := loadedShortcuts.Entry[keyPressed]; ok {
					processShortcuts(shortcut, keyPressed)
				} else {
					defaultKeystroke(keyPressedCode)

				}
			}
		}
	}
}
