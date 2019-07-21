package hc595

import (
	"sync"
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/warthog618/gpio"
)

type hc595 struct {
	data  *gpio.Pin
	clk   *gpio.Pin
	latch *gpio.Pin
	state uint32
	mutex *sync.Mutex
}

/*
		Byte   1-(MSB)--------  2--------------  3--------------  4-(LSB)--------
		Bit    7 6 5 4 3 2 1 0  7 6 5 4 3 2 1 0  7 6 5 4 3 2 1 0  7 6 5 4 3 2 1 0

		Usage  unused  		    HC595 A	         HC595 B          HC595 C
		       L L L L L L L L  L L M M M M . .  B B B B B B . .  A A A A A A . .


	L = LED via ULN2003
	A = Atari Joystick Port A
	B = Atari Joystick Port B

*/

func NewHc595(dataPin, clkPin, latchPin uint8) *hc595 {

	// 595 data GPIO
	data := gpio.NewPin(dataPin)
	data.Output()
	data.Low()

	// 595 clk GPIO
	clk := gpio.NewPin(clkPin)
	clk.Output()
	clk.Low()

	// 595 latch GPIO
	latch := gpio.NewPin(latchPin)
	latch.Output()
	latch.Low()

	return &hc595{
		data:  data,
		clk:   clk,
		latch: latch,
		mutex: &sync.Mutex{}}

}

func (o *hc595) sendWord(b uint32) {

	var x uint32 = 1
	for i := 0; i < 32; i++ {

		if b&x > 0 {
			o.data.High()
		} else {
			o.data.Low()
		}
		x = x * 2

		time.Sleep(1 * time.Microsecond)
		o.clk.High()
		time.Sleep(1 * time.Microsecond)
		o.clk.Low()
	}

	o.latch.High()
	time.Sleep(1 * time.Microsecond)
	o.latch.Low()

	o.state = b
}

func (o *hc595) SelectRow(b uint8) {
	o.mutex.Lock()
	var value uint32 = (0xf - 1<<b) << 18

	o.sendWord(o.state&0xFFC3FFFF | value)
	o.mutex.Unlock()
}

func (o *hc595) SetJoys(joy1, joy2 *data.Joystick, butt1a, butt1b, butt2a, butt2b bool) {
	o.mutex.Lock()
	state := o.state
	state = setJoystick(state, 0, joy1, butt1a, butt1b)
	state = setJoystick(state, 1, joy2, butt2a, butt2b)
	o.sendWord(state)
	o.mutex.Unlock()
}

func (o *hc595) SetLeds(leds data.LedState) {
	o.mutex.Lock()
	ledState := 0

	ledState += B2i(leds.Player1Keypad.Red, 0)
	ledState += B2i(leds.Player1Keypad.Yellow, 1)
	ledState += B2i(leds.Player1Keypad.Blue, 2)
	ledState += B2i(leds.Player1Keypad.Green, 3)
	ledState += B2i(leds.GlobalKeypad.WhiteLeft, 4)

	ledState += B2i(leds.Player2Keypad.Red, 5)
	ledState += B2i(leds.Player2Keypad.Yellow, 6)
	ledState += B2i(leds.Player2Keypad.Blue, 7)
	ledState += B2i(leds.Player2Keypad.Green, 8)
	ledState += B2i(leds.GlobalKeypad.WhiteRight, 9)

	o.sendWord(o.state&0x001FFFFF | (uint32(ledState))<<22)
	o.mutex.Unlock()
}

func setJoystick(oldState uint32, index uint8, data *data.Joystick, button1, button2 bool) uint32 {

	var value = uint32(0xff)
	if data.Up {
		value -= 128
	}
	if data.Down {
		value -= 64
	}
	if data.Left {
		value -= 32
	}
	if data.Right {
		value -= 16
	}
	if button1 {
		value -= 8
	}
	if button2 {
		value -= 4
	}

	switch index {
	case 0:
		return oldState&0xFFFFFF03 | (uint32(value))<<0
	case 1:
		return oldState&0xFFFF03FF | (uint32(value))<<8
	}

	return oldState

}

func B2i(b bool, shift uint8) int {
	if b {
		return 1 << shift
	}
	return 0
}
