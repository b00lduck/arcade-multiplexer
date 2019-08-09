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

	Chip 1
	------
	0 O Joy P1 UP
	1
	2
	3
	4
	5
	6
	7

	Chip 2
	------
	0
	1
	2
	3
	4
	5
	6
	7


*/

func NewMist(bus i2c.Bus) *mist {

	// Address the devices on the I²C bus
	chip0 := i2c.Dev{Bus: bus, Addr: 0x23}
	chip1 := i2c.Dev{Bus: bus, Addr: 0x24}

	chip0.Write([]byte{0xff})
	chip1.Write([]byte{0xff})

	return &mist{
		chips:         []i2c.Dev{chip0, chip1},
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

func (o *mist) SetJoystick(joy1, joy2 *data.Joystick) {
	o.changeChipBit(0, 0, !joy1.Up)
	o.changeChipBit(0, 1, !joy1.Down)
	o.changeChipBit(0, 2, !joy1.Left)
	o.changeChipBit(0, 3, !joy1.Right)
	o.changeChipBit(0, 4, true)
	o.changeChipBit(0, 5, true)

	o.changeChipBit(0, 6, !joy2.Up)
	o.changeChipBit(0, 7, !joy2.Down)
	o.changeChipBit(1, 0, !joy2.Left)
	o.changeChipBit(1, 1, !joy2.Right)
	o.changeChipBit(1, 2, true)
	o.changeChipBit(1, 3, true)

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
