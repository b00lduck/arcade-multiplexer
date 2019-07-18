package data

type MatrixState struct {
	Player1Keypad   PlayerKeypad
	Player2Keypad   PlayerKeypad
	Player1Joystick Joystick
	Player2Joystick Joystick
	GlobalKeypad    GlobalKeypad
}

type LedState struct {
	Player1Keypad PlayerKeypad
	Player2Keypad PlayerKeypad
	GlobalKeypad  GlobalKeypad
}
