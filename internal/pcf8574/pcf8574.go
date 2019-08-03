package pcf8574

import (
	"sync"
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/tarent/logrus"

	"github.com/jinzhu/copier"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

type pcf8574 struct {
	state         uint32
	mutex         *sync.Mutex
	chips         []i2c.Dev
	writtenStates []uint8
	readStates    []uint8

	selectedRow uint8

	matrixState data.MatrixState
}

/*

	Chip 1
	------
	0 O LED P1 red
	1 O LED P1 yellow
	2 O LED P1 green
	3 O LED P1 blue
	4 O LED P1 select
	5 O LED P2 red
	6 O LED P2 yellow
	7 O LED P2 green

	Chip 2
	------
	0 O LED P2 blue
	1
	2
	3
	4 I Column 4
	5
	6
	7 O LED P2 select


	Chip 3
	------
	0 O Row 0
	1 O Row 1
	2 O Row 2
	3 O Row 3
	4 I Column 0
	5 I Column 1
	6 I Column 2
	7 I Column 3

*/

const NUM_ROWS = 4

func NewPcf8574() *pcf8574 {

	_, err := host.Init()
	if err != nil {
		logrus.WithError(err).Fatal("Could not open initialize periph.io")
	}

	// Open the first available I²C bus
	bus, err := i2creg.Open("")
	if err != nil {
		logrus.WithError(err).Fatal("Could not open i2c bus")
	}

	// Address the devices on the I²C bus
	chip0 := i2c.Dev{Bus: bus, Addr: 0x20}
	chip1 := i2c.Dev{Bus: bus, Addr: 0x21}
	chip2 := i2c.Dev{Bus: bus, Addr: 0x22}

	chip0.Write([]byte{0xff})
	chip1.Write([]byte{0xff})
	chip2.Write([]byte{0xff})

	return &pcf8574{
		chips:         []i2c.Dev{chip0, chip1, chip2},
		writtenStates: []uint8{0xff, 0xff, 0xff},
		readStates:    []uint8{0xff, 0xff, 0xff},
		mutex:         &sync.Mutex{}}

}

func (o *pcf8574) Run(changeEvent func(data.MatrixState)) {

	for {

		o.mutex.Lock()

		o.selectNextRow()

		//logrus.Info(fmt.Sprintf("%08b %08b %08b", o.writtenStates[0], o.writtenStates[1], o.writtenStates[2]))
		for k, v := range o.chips {
			r := []byte{0}
			v.Tx([]byte{o.writtenStates[k]}, r)
			o.readStates[k] = r[0]
		}

		time.Sleep(10 * time.Millisecond)

		newMatrix := o.decodeMatrix()

		if newMatrix.Changed(o.matrixState) {
			o.matrixState = newMatrix
			changeEvent(o.matrixState)
		}

		o.mutex.Unlock()

		time.Sleep(1000 * time.Microsecond)

	}

}

func (o *pcf8574) decodeMatrix() data.MatrixState {

	var newMatrix data.MatrixState
	copier.Copy(&newMatrix, &o.matrixState)

	col0 := (o.readStates[2] & 0x10) == 0
	col1 := (o.readStates[2] & 0x20) == 0
	col2 := (o.readStates[2] & 0x40) == 0
	col3 := (o.readStates[2] & 0x80) == 0
	col4 := (o.readStates[1] & 0x10) == 0

	switch o.selectedRow {
	case 0:
		newMatrix.Player1Keypad.Red = col0
		newMatrix.Player2Keypad.Red = col1
		newMatrix.GlobalKeypad.WhiteLeft = col2
		newMatrix.Player1Joystick.Right = col3
		newMatrix.Player2Joystick.Down = col4

	case 1:
		newMatrix.Player1Keypad.Yellow = col0
		newMatrix.Player2Keypad.Yellow = col1
		newMatrix.GlobalKeypad.WhiteRight = col2
		newMatrix.Player1Joystick.Left = col3
		newMatrix.Player2Joystick.Right = col4

	case 2:
		newMatrix.Player1Keypad.Green = col0
		newMatrix.Player2Keypad.Green = col1
		newMatrix.GlobalKeypad.FlipperRight = col2
		newMatrix.Player1Joystick.Up = col3
		newMatrix.Player2Joystick.Left = col4

	case 3:
		newMatrix.Player1Keypad.Blue = col0
		newMatrix.Player2Keypad.Blue = col1
		newMatrix.GlobalKeypad.FlipperLeft = col2
		newMatrix.Player1Joystick.Down = col3
		newMatrix.Player2Joystick.Up = col4
	}

	return newMatrix
}

func (o *pcf8574) selectNextRow() {

	for i := uint8(0); i < NUM_ROWS; i++ {
		state := i != o.selectedRow
		o.changeChipBit(2, i, state)
	}

	// Advance to next row
	o.selectedRow++
	if o.selectedRow >= NUM_ROWS {
		o.selectedRow = 0
	}

}

func (o *pcf8574) SetLeds(leds data.LedState) {
	o.mutex.Lock()

	o.changeChipBit(0, 0, leds.Player1Keypad.Red)
	o.changeChipBit(0, 1, leds.Player1Keypad.Yellow)
	o.changeChipBit(0, 2, leds.Player1Keypad.Green)
	o.changeChipBit(0, 3, leds.Player1Keypad.Blue)
	o.changeChipBit(0, 4, leds.GlobalKeypad.WhiteLeft)

	o.changeChipBit(0, 5, leds.Player1Keypad.Red)
	o.changeChipBit(0, 6, leds.Player1Keypad.Yellow)
	o.changeChipBit(0, 7, leds.Player1Keypad.Green)
	o.changeChipBit(1, 0, leds.Player1Keypad.Blue)
	o.changeChipBit(1, 7, leds.GlobalKeypad.WhiteLeft)

	o.mutex.Unlock()
}

func (o *pcf8574) changeChipBit(chip uint8, bit uint8, state bool) {
	o.writtenStates[chip] = changeBit(o.writtenStates[chip], bit, state)
}

func changeBit(val uint8, shift uint8, b bool) uint8 {
	mask := uint8(1 << shift)
	if b {
		return val | mask
	}
	return val & (^mask)
}
