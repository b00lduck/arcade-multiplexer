package hc595

import (
	"fmt"
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/warthog618/gpio"
)

type hc595 struct {
	data  *gpio.Pin
	clk   *gpio.Pin
	latch *gpio.Pin
	State uint32
}

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
		latch: latch}

}

func (o *hc595) SendWord(b uint32) {

	var x uint32 = 1
	for i := 0; i < 24; i++ {

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

	o.State = b

	fmt.Printf("%x\n", o.State)

}

func SetJoystick(oldState uint32, index uint8, data *data.Joystick) uint32 {

	var value uint32
	if data.Up {
		value += 16
	}
	if data.Down {
		value += 8
	}
	if data.Left {
		value += 4
	}
	if data.Right {
		value += 2
	}

	switch index {
	case 0:
		return (oldState | 0x1e) - value
	case 1:
		return (oldState | 0x3c0) - value<<5
	}

	return oldState

}

func SetButton(oldState uint32, index uint8, button bool) uint32 {

	var value uint32
	if !button {
		value = 1
	}

	switch index {
	case 0:
		return oldState&0xfffffffe | value
	case 1:
		return oldState&0xffffffdf | value<<5
	}

	return oldState

}
