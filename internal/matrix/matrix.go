package matrix

import (
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/warthog618/gpio"
)

type matrix struct {
	cols          []*gpio.Pin
	numRows       uint8
	state         [][]bool
	MatrixState   *data.MatrixState
	selectRowFunc func(uint8)
}

func NewMatrix(selectRowFunc func(uint8), numRows uint8, colPins []uint8) *matrix {

	cols := make([]*gpio.Pin, len(colPins))

	state := make([][]bool, len(colPins))
	for i := range state {
		state[i] = make([]bool, numRows)
	}

	for c := 0; c < len(colPins); c++ {
		cols[c] = gpio.NewPin(colPins[c])
		cols[c].Input()
		cols[c].PullUp()
	}

	return &matrix{
		numRows:       numRows,
		cols:          cols,
		state:         state,
		selectRowFunc: selectRowFunc,
		MatrixState:   &data.MatrixState{}}

}

func (m *matrix) Run(changedCallback func(*data.MatrixState)) {

	for {

		changed := false

		for row := uint8(0); row < m.numRows; row++ {
			m.selectRowFunc(row)
			time.Sleep(1 * time.Microsecond)
			for colKey, col := range m.cols {
				oldValue := m.state[colKey][row]
				m.state[colKey][row] = col.Read() == gpio.Low
				if m.state[colKey][row] != oldValue {
					changed = true
				}
			}
		}

		// C443FF 1100 0100 0100 0011 1111 1111

		// c447ff 1100 0100 0100 0111 1111 1111
		// c44bff 1100 0100 0100 1011 1111 1111
		// C453FF 1100 0100 0101 0011 1111 1111
		// c463ff 1100 0100 0110 0011 1111 1111

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
