package data

type MatrixState struct {
	Player1Keypad   PlayerKeypad
	Player2Keypad   PlayerKeypad
	Player1Joystick Joystick
	Player2Joystick Joystick
	GlobalKeypad    GlobalKeypad
}

func (m MatrixState) Changed(old MatrixState) bool {
	return m.Player1Keypad.Changed(old.Player1Keypad) ||
		m.Player2Keypad.Changed(old.Player2Keypad) ||
		m.Player1Joystick.Changed(old.Player1Joystick) ||
		m.Player2Joystick.Changed(old.Player2Joystick) ||
		m.GlobalKeypad.Changed(old.GlobalKeypad)
}

func (m MatrixState) String() string {
	ret := "MatrixState:\n"
	ret += "Player 1 Keypad: " + m.Player1Keypad.String() + "\n"
	ret += "Player 2 Keypad: " + m.Player2Keypad.String() + "\n"
	ret += "Player 1 Joystick: " + m.Player1Joystick.String() + "\n"
	ret += "Player 2 Joystick: " + m.Player2Joystick.String() + "\n"
	ret += "Global Keypad: " + m.GlobalKeypad.String() + "\n"
	return ret
}
