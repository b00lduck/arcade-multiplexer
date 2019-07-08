package matrix

import (
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/warthog618/gpio"
)

type matrix struct {
	cols        []*gpio.Pin
	rows        []*gpio.Pin
	state       [][]bool
	MatrixState *data.MatrixState
}

/*
		Byte   1-(MSB)--------  2--------------  3--------------  4-(LSB)--------
		Bit    7 6 5 4 3 2 1 0  7 6 5 4 3 2 1 0  7 6 5 4 3 2 1 0  7 6 5 4 3 2 1 0

		Usage  unused  		    HC595 A	         HC595 B          HC595 C
		       . . . . . . . .  L L L L L L L L  L L L . . . B B  B B B A A A A A

	L = LED via ULN2003
	A = Atari Joystick Port A
	B = Atari Joystick Port B

*/

func NewMatrix(rowPins []uint8, colPins []uint8) *matrix {

	rows := make([]*gpio.Pin, len(rowPins))
	cols := make([]*gpio.Pin, len(colPins))

	state := make([][]bool, len(colPins))
	for i := range state {
		state[i] = make([]bool, len(rowPins))
	}

	for r := 0; r < len(rowPins); r++ {
		rows[r] = gpio.NewPin(rowPins[r])
		rows[r].Input()
		rows[r].PullDown()
	}

	for c := 0; c < len(colPins); c++ {
		cols[c] = gpio.NewPin(colPins[c])
		cols[c].Output()
		cols[c].Low()
	}

	return &matrix{
		rows:        rows,
		cols:        cols,
		state:       state,
		MatrixState: &data.MatrixState{}}

}

func (m *matrix) Run(changedCallback func(*data.MatrixState)) {

	for {

		changed := false

		for colKey, col := range m.cols {
			col.High()
			time.Sleep(1 * time.Microsecond)

			for rowKey, row := range m.rows {
				oldValue := m.state[colKey][rowKey]
				m.state[colKey][rowKey] = row.Read() == gpio.High
				if m.state[colKey][rowKey] != oldValue {
					changed = true
				}
			}

			col.Low()
			time.Sleep(1 * time.Microsecond)
		}

		if changed {
			for colKey := range m.cols {

				col := m.state[colKey]

				switch colKey {
				case 0:
					m.MatrixState.Player1Keypad.Red = col[0]
					m.MatrixState.Player1Keypad.Yellow = col[1]
					m.MatrixState.Player1Keypad.Green = col[2]
					m.MatrixState.Player1Keypad.Blue = col[3]
				case 1:
					m.MatrixState.Player2Keypad.Red = col[0]
					m.MatrixState.Player2Keypad.Yellow = col[1]
					m.MatrixState.Player2Keypad.Green = col[2]
					m.MatrixState.Player2Keypad.Blue = col[3]
				case 2:
					m.MatrixState.GlobalKeypad.WhiteLeft = col[0]
					m.MatrixState.GlobalKeypad.WhiteRight = col[1]
					m.MatrixState.GlobalKeypad.FlipperRight = col[2]
					m.MatrixState.GlobalKeypad.FlipperLeft = col[3]
				case 3:
					m.MatrixState.Player1Joystick.Right = col[0]
					m.MatrixState.Player1Joystick.Left = col[1]
					m.MatrixState.Player1Joystick.Up = col[2]
					m.MatrixState.Player1Joystick.Down = col[3]
				case 4:
					m.MatrixState.Player2Joystick.Right = col[0]
					m.MatrixState.Player2Joystick.Left = col[1]
					m.MatrixState.Player2Joystick.Up = col[2]
					m.MatrixState.Player2Joystick.Down = col[3]
				}

			}

			changedCallback(m.MatrixState)

		}
	}

}
