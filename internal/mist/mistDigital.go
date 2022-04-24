package mist

import (
	"time"

	"arcade-multiplexer/internal/data"

	"periph.io/x/periph/conn/i2c"
)

type mistDigital struct {
	state         uint32
	chips         []i2c.Dev
	writtenStates []uint8
	readStates    []uint8
	buttonStates  []data.ButtonState
	afcount       uint32
	afstate       bool
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

func NewMistDigital(bus i2c.Bus) *mistDigital {

	// Address the devices on the I²C bus
	chip1 := i2c.Dev{Bus: bus, Addr: 0x23}
	chip2 := i2c.Dev{Bus: bus, Addr: 0x24}

	chip1.Write([]byte{0xff})
	chip2.Write([]byte{0xef})

	return &mistDigital{
		chips:         []i2c.Dev{chip1, chip2},
		writtenStates: []uint8{0xff, 0xef},
		readStates:    []uint8{0xff, 0xef},
		buttonStates:  make([]data.ButtonState, 4)}

}

func (o *mistDigital) AutoFired(bs *data.ButtonState) bool {
	if bs.Autofire {
		return bs.State && o.afstate
	}
	return bs.State
}

func (o *mistDigital) Run() {

	for {
		o.afcount++
		if o.afcount > 20 {
			o.afcount = 0
			o.afstate = !o.afstate
		}
		o.changeChipBit(1, 2, !o.AutoFired(&o.buttonStates[0]))
		o.changeChipBit(1, 3, !o.AutoFired(&o.buttonStates[1]))
		o.changeChipBit(0, 4, !o.AutoFired(&o.buttonStates[2]))
		o.changeChipBit(0, 5, !o.AutoFired(&o.buttonStates[3]))

		//logrus.Info(fmt.Sprintf("%08b %08bb", o.writtenStates[0], o.writtenStates[1]))
		for k, v := range o.chips {
			r := []byte{0}
			v.Tx([]byte{o.writtenStates[k]}, r)
			o.readStates[k] = r[0]
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func (o *mistDigital) SetJoystickButton(id uint8, bs data.ButtonState) {
	o.buttonStates[id] = bs
}

func (o *mistDigital) SetJoystick1(joy *data.Joystick) {
	o.changeChipBit(0, 7, !joy.Up)
	o.changeChipBit(0, 6, !joy.Down)
	o.changeChipBit(1, 0, !joy.Left)
	o.changeChipBit(1, 1, !joy.Right)
}

func (o *mistDigital) SetJoystick2(joy *data.Joystick) {
	o.changeChipBit(0, 0, !joy.Up)
	o.changeChipBit(0, 1, !joy.Down)
	o.changeChipBit(0, 2, !joy.Left)
	o.changeChipBit(0, 3, !joy.Right)
}

func (o *mistDigital) SetPower(state bool) {
	o.changeChipBit(1, 4, !state)
}

func (o *mistDigital) SetResetButton(state bool) {
	o.changeChipBit(1, 7, !state)
}

func (o *mistDigital) changeChipBit(chip uint8, bit uint8, state bool) {
	o.writtenStates[chip] = changeBit(o.writtenStates[chip], bit, state)
}

func changeBit(val uint8, shift uint8, b bool) uint8 {
	mask := uint8(1 << shift)
	if b {
		return val | mask
	}
	return val & (^mask)
}
