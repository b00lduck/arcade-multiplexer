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
	0 O Joy P1 UP
	1 O Joy P1 DOWN
	2 O Joy P1 LEFT
	3 O Joy P1 RIGHT
	4 O Joy P1 BUT 1
	5 O Joy P1 BUT 2
	6 O Joy P2 UP
	7 O Joy P2 DOWN

	Chip 2 (0x24)
	------
	0 O Joy P2 LEFT
	1 O Joy P2 RIGHT
	2 O Joy P2 BUT 1
	3 O Joy P2 BUT 2
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
		//logrus.Info(fmt.Sprintf("%08b %08b %08b", o.writtenStates[0], o.writtenStates[1], o.writtenStates[2]))
		for k, v := range o.chips {
			r := []byte{0}
			v.Tx([]byte{o.writtenStates[k]}, r)
			o.readStates[k] = r[0]
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (o *mist) SetJoystick(joy1, joy2 *data.Joystick, j1b1, j1b2, j2b1, j2b2 bool) {
	o.changeChipBit(0, 0, !joy1.Up)
	o.changeChipBit(0, 1, !joy1.Down)
	o.changeChipBit(0, 2, !joy1.Left)
	o.changeChipBit(0, 3, !joy1.Right)
	o.changeChipBit(0, 4, !j1b1)
	o.changeChipBit(0, 5, !j1b2)

	o.changeChipBit(0, 6, !joy2.Up)
	o.changeChipBit(0, 7, !joy2.Down)
	o.changeChipBit(1, 0, !joy2.Left)
	o.changeChipBit(1, 1, !joy2.Right)
	o.changeChipBit(1, 2, !j2b1)
	o.changeChipBit(1, 3, !j2b2)
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
