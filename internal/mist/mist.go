package mist

import (
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/data"

	"periph.io/x/periph/conn/i2c"
)

type mist struct {
	state         uint32
	chips         []i2c.Dev
	writtenStates []uint8
	readStates    []uint8
}

/*
	Chip 1 (0x23)
	------
	0 O Joy P2 UP
	1 O Joy P2 DOWN
	2 O Joy P2 LEFT
	3 O Joy P2 RIGHT
	4 O Joy P2 BUT 1
	5 O Joy P2 BUT 2
	6 O Joy P1 DOWN
	7 O Joy P1 UP

	Chip 2 (0x24)
	------
	0 O Joy P1 LEFT
	1 O Joy P1 RIGHT
	2 O Joy P1 BUT 1
	3 O Joy P1 BUT 2
	4 O Power on
	5 O Reset Button
	6 I LED
	7 I LED
*/

func NewMist(bus i2c.Bus) *mist {

	// Address the devices on the I²C bus
	chip1 := i2c.Dev{Bus: bus, Addr: 0x23}
	chip2 := i2c.Dev{Bus: bus, Addr: 0x24}

	chip1.Write([]byte{0xff})
	chip2.Write([]byte{0xff})

	return &mist{
		chips:         []i2c.Dev{chip1, chip2},
		writtenStates: []uint8{0xff, 0xff},
		readStates:    []uint8{0xff, 0xff}}

}

func (o *mist) Run() {
	for {
		//logrus.Info(fmt.Sprintf("%08b %08bb", o.writtenStates[0], o.writtenStates[1]))
		for k, v := range o.chips {
			r := []byte{0}
			v.Tx([]byte{o.writtenStates[k]}, r)
			o.readStates[k] = r[0]
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func (o *mist) SetJoystick1(joy *data.Joystick) {
	o.changeChipBit(0, 7, !joy.Up)
	o.changeChipBit(0, 6, !joy.Down)
	o.changeChipBit(1, 0, !joy.Left)
	o.changeChipBit(1, 1, !joy.Right)
}

func (o *mist) SetJoystick1Button1(state bool) {
	o.changeChipBit(1, 2, !state)
}

func (o *mist) SetJoystick1Button2(state bool) {
	o.changeChipBit(1, 3, !state)
}

func (o *mist) SetJoystick2(joy *data.Joystick) {
	o.changeChipBit(0, 0, !joy.Up)
	o.changeChipBit(0, 1, !joy.Down)
	o.changeChipBit(0, 2, !joy.Left)
	o.changeChipBit(0, 3, !joy.Right)
}

func (o *mist) SetJoystick2Button1(state bool) {
	o.changeChipBit(0, 4, !state)
}

func (o *mist) SetJoystick2Button2(state bool) {
	o.changeChipBit(0, 5, !state)
}

func (o *mist) SetPower(state bool) {
	o.changeChipBit(1, 4, !state)
}

func (o *mist) SetResetButton(state bool) {
	o.changeChipBit(1, 5, !state)
}

func (o *mist) changeChipBit(chip uint8, bit uint8, state bool) {
	o.writtenStates[chip] = changeBit(o.writtenStates[chip], bit, state)
}

func changeBit(val uint8, shift uint8, b bool) uint8 {
	mask := uint8(1 << shift)
	if b {
		return val | mask
	}
	return val & (^mask)
}
