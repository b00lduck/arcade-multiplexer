package matrix

import (
	"fmt"
	"time"

	"github.com/warthog618/gpio"
)

type matrix struct {
	cols  []*gpio.Pin
	rows  []*gpio.Pin
	state [][]bool
}

func NewMatrix(rowPins []uint8, colPins []uint8) *matrix {

	rows := make([]*gpio.Pin, len(rowPins))
	cols := make([]*gpio.Pin, len(colPins))

	state := make([][]bool, len(rowPins))
	for i := range state {
		state[i] = make([]bool, len(colPins))
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
		rows:  rows,
		cols:  cols,
		state: state}

}

func (m *matrix) Run() {

	for {

		changed := false

		for colKey, col := range m.cols {
			col.High()
			time.Sleep(1 * time.Microsecond)

			for rowKey, row := range m.rows {
				oldValue := m.state[rowKey][colKey]
				m.state[rowKey][colKey] = row.Read() == gpio.High
				if m.state[rowKey][colKey] != oldValue {
					changed = true
				}
			}

			col.Low()
			time.Sleep(1 * time.Microsecond)
		}

		if changed {
			for rowKey := range m.rows {
				for colKey := range m.cols {
					fmt.Printf("%t ", m.state[rowKey][colKey])
				}
				fmt.Print("\n")
			}
			fmt.Printf("\033[4A")
		}
	}

}
