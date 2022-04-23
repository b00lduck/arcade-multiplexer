package data

import "arcade-multiplexer/internal/config"

type LedState struct {
	Player1Keypad PlayerKeypad
	Player2Keypad PlayerKeypad
	GlobalKeypad  GlobalKeypad
}

func LedStateByMapping(mapping []config.Mapping) LedState {
	ret := LedState{}

	for _, v := range mapping {
		switch v.Input {
		case "P1_RED":
			ret.Player1Keypad.Red = true
		case "P1_YELLOW":
			ret.Player1Keypad.Yellow = true
		case "P1_GREEN":
			ret.Player1Keypad.Green = true
		case "P1_BLUE":
			ret.Player1Keypad.Blue = true
		case "P2_RED":
			ret.Player2Keypad.Red = true
		case "P2_YELLOW":
			ret.Player2Keypad.Yellow = true
		case "P2_GREEN":
			ret.Player2Keypad.Green = true
		case "P2_BLUE":
			ret.Player2Keypad.Blue = true
		case "WHITE_LEFT":
			ret.GlobalKeypad.WhiteLeft = true
		case "WHITE_RIGHT":
			ret.GlobalKeypad.WhiteRight = true
		}
	}

	return ret
}
