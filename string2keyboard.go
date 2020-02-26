// Sourced from https://github.com/taglme/string2keyboard
package main

import (
	"runtime"
	"time"

	"github.com/rewardian/keybd_event"
)

type keySet struct {
	code  int
	shift bool
}

var startUp int = 0

//KeyboardWrite emulate keyboard input from string
func KeyboardWrite(textInput string) error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}

	if runtime.GOOS == "linux" && startUp == 0 {
		time.Sleep(1500 * time.Millisecond)
		startUp++
	}

	//Should we skip next character in string
	//Used if we found some escape sequence
	skip := false
	for i, c := range textInput {
		if !skip {
			if c != '\\' {
				kb.SetKeys(s2kNames[string(c)].code)
				kb.HasSHIFT(s2kNames[string(c)].shift)
			} else {
				//Found backslash escape character
				//Check next character
				switch textInput[i+1] {
				case 'n':
					//Found newline character sequence
					kb.SetKeys(s2kNames["ENTER"].code)
					skip = true
				case '\\':
					//Found backslash character sequence
					kb.SetKeys(s2kNames["\\"].code)
					kb.HasSHIFT(s2kNames["\\"].shift)
					skip = true
				case 'b':
					//Found backspace character sequence
					kb.SetKeys(s2kNames["BACKSPACE"].code)
					skip = true
				case 't':
					//Found tab character sequence
					kb.SetKeys(s2kNames["TAB"].code)
					skip = true
				case '"':
					//Found double quote character sequence
					kb.SetKeys(s2kNames["\""].code)
					kb.HasSHIFT(s2kNames["\""].shift)
					skip = true
				default:
					//Nothing special, jsut backslash output
					kb.SetKeys(s2kNames["\\"].code)
					kb.HasSHIFT(s2kNames["\\"].shift)
				}

			}
			err = kb.Launching()
			if err != nil {
				return err
			}
		} else {
			skip = false
		}

	}
	return nil

}

var (
	s2kNames = map[string]keySet{
		"a": keySet{keybd_event.VK_A, false},
		"b": keySet{keybd_event.VK_B, false},
		"c": keySet{keybd_event.VK_C, false},
		"d": keySet{keybd_event.VK_D, false},
		"e": keySet{keybd_event.VK_E, false},
		"f": keySet{keybd_event.VK_F, false},
		"g": keySet{keybd_event.VK_G, false},
		"h": keySet{keybd_event.VK_H, false},
		"i": keySet{keybd_event.VK_I, false},
		"j": keySet{keybd_event.VK_J, false},
		"k": keySet{keybd_event.VK_K, false},
		"l": keySet{keybd_event.VK_L, false},
		"m": keySet{keybd_event.VK_M, false},
		"n": keySet{keybd_event.VK_N, false},
		"o": keySet{keybd_event.VK_O, false},
		"p": keySet{keybd_event.VK_P, false},
		"q": keySet{keybd_event.VK_Q, false},
		"r": keySet{keybd_event.VK_R, false},
		"s": keySet{keybd_event.VK_S, false},
		"t": keySet{keybd_event.VK_T, false},
		"u": keySet{keybd_event.VK_U, false},
		"v": keySet{keybd_event.VK_V, false},
		"w": keySet{keybd_event.VK_W, false},
		"x": keySet{keybd_event.VK_X, false},
		"y": keySet{keybd_event.VK_Y, false},
		"z": keySet{keybd_event.VK_Z, false},

		"A": keySet{keybd_event.VK_A, true},
		"B": keySet{keybd_event.VK_B, true},
		"C": keySet{keybd_event.VK_C, true},
		"D": keySet{keybd_event.VK_D, true},
		"E": keySet{keybd_event.VK_E, true},
		"F": keySet{keybd_event.VK_F, true},
		"G": keySet{keybd_event.VK_G, true},
		"H": keySet{keybd_event.VK_H, true},
		"I": keySet{keybd_event.VK_I, true},
		"J": keySet{keybd_event.VK_J, true},
		"K": keySet{keybd_event.VK_K, true},
		"L": keySet{keybd_event.VK_L, true},
		"M": keySet{keybd_event.VK_M, true},
		"N": keySet{keybd_event.VK_N, true},
		"O": keySet{keybd_event.VK_O, true},
		"P": keySet{keybd_event.VK_P, true},
		"Q": keySet{keybd_event.VK_Q, true},
		"R": keySet{keybd_event.VK_R, true},
		"S": keySet{keybd_event.VK_S, true},
		"T": keySet{keybd_event.VK_T, true},
		"U": keySet{keybd_event.VK_U, true},
		"V": keySet{keybd_event.VK_V, true},
		"W": keySet{keybd_event.VK_W, true},
		"X": keySet{keybd_event.VK_X, true},
		"Y": keySet{keybd_event.VK_Y, true},
		"Z": keySet{keybd_event.VK_Z, true},

		"0": keySet{keybd_event.VK_0, false},
		"1": keySet{keybd_event.VK_1, false},
		"2": keySet{keybd_event.VK_2, false},
		"3": keySet{keybd_event.VK_3, false},
		"4": keySet{keybd_event.VK_4, false},
		"5": keySet{keybd_event.VK_5, false},
		"6": keySet{keybd_event.VK_6, false},
		"7": keySet{keybd_event.VK_7, false},
		"8": keySet{keybd_event.VK_8, false},
		"9": keySet{keybd_event.VK_9, false},
		" ": keySet{keybd_event.VK_SPACE, false},

		")": keySet{keybd_event.VK_0, true},
		"!": keySet{keybd_event.VK_1, true},
		"@": keySet{keybd_event.VK_2, true},
		"#": keySet{keybd_event.VK_3, true},
		"$": keySet{keybd_event.VK_4, true},
		"%": keySet{keybd_event.VK_5, true},
		"^": keySet{keybd_event.VK_6, true},
		"&": keySet{keybd_event.VK_7, true},
		"*": keySet{keybd_event.VK_8, true},
		"(": keySet{keybd_event.VK_9, true},

		"-":  keySet{keybd_event.VK_MINUS, false},
		"=":  keySet{keybd_event.VK_EQUAL, false},
		"[":  keySet{keybd_event.VK_LEFTBRACE, false},
		"]":  keySet{keybd_event.VK_RIGHTBRACE, false},
		"\\": keySet{keybd_event.VK_BACKSLASH, false},
		";":  keySet{keybd_event.VK_SEMICOLON, false},
		"'":  keySet{keybd_event.VK_APOSTROPHE, false},
		",":  keySet{keybd_event.VK_COMMA, false},
		".":  keySet{keybd_event.VK_DOT, false},
		"/":  keySet{keybd_event.VK_SLASH, false},
		"`":  keySet{keybd_event.VK_GRAVE, false},

		"_":  keySet{keybd_event.VK_MINUS, true},
		"+":  keySet{keybd_event.VK_EQUAL, true},
		"{":  keySet{keybd_event.VK_LEFTBRACE, true},
		"}":  keySet{keybd_event.VK_RIGHTBRACE, true},
		"|":  keySet{keybd_event.VK_BACKSLASH, true},
		":":  keySet{keybd_event.VK_SEMICOLON, true},
		"\"": keySet{keybd_event.VK_APOSTROPHE, true},
		"<":  keySet{keybd_event.VK_COMMA, true},
		">":  keySet{keybd_event.VK_DOT, true},
		"?":  keySet{keybd_event.VK_SLASH, true},
		"~":  keySet{keybd_event.VK_GRAVE, true},

		"ENTER":     keySet{keybd_event.VK_ENTER, false},
		"TAB":       keySet{keybd_event.VK_TAB, false},
		"BACKSPACE": keySet{keybd_event.VK_BACKSPACE, false},
	}
)
