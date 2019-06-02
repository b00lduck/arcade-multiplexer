package hc595

import (
	"fmt"
	"time"

	"github.com/warthog618/gpio"
)

type hc595 struct {
	data  *gpio.Pin
	clk   *gpio.Pin
	latch *gpio.Pin
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

func (o *hc595) SendByte(b uint16) {

	var x uint16 = 1
	for i := 0; i < 16; i++ {

		fmt.Printf("%d\n", b&x)

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

}
