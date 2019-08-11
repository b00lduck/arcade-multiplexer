package rotary

import (
	"fmt"
	"sync"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type rotary struct {
	clk   gpio.PinIO
	data  gpio.PinIO
	btn   gpio.PinIO
	delta int
	mutex *sync.Mutex
}

func NewRotary(clkPin, dataPin, btnPin uint8) *rotary {

	clk := gpioreg.ByName(fmt.Sprintf("%d", clkPin))
	clk.In(gpio.PullDown, gpio.BothEdges)

	data := gpioreg.ByName(fmt.Sprintf("%d", dataPin))
	clk.In(gpio.PullDown, gpio.BothEdges)

	btn := gpioreg.ByName(fmt.Sprintf("%d", btnPin))
	clk.In(gpio.PullDown, gpio.BothEdges)

	return &rotary{
		clk:   clk,
		data:  data,
		btn:   btn,
		mutex: &sync.Mutex{}}
}

func (o *rotary) Run() {

	encLast := 1

	for {

		o.mutex.Lock()
		i := 0

		if o.clk.Read() == gpio.Low {
			i = 1
		}
		if o.data.Read() == gpio.Low {
			i ^= 3
		}

		i -= encLast

		if i&1 == 1 {
			encLast += i
			o.delta += (i & 2) - 1
		}
		o.mutex.Unlock()

		time.Sleep(1 * time.Millisecond)

	}

}

func (o *rotary) Delta() int {

	o.mutex.Lock()
	ret := o.delta
	o.delta = 0
	o.mutex.Unlock()
	return ret

}
