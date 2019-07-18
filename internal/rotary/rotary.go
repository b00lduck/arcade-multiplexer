package rotary

import (
	"sync"
	"time"

	"github.com/warthog618/gpio"
)

type rotary struct {
	clk   *gpio.Pin
	data  *gpio.Pin
	btn   *gpio.Pin
	delta int
	mutex *sync.Mutex
}

func NewRotary(clkPin, dataPin, btnPin uint8) *rotary {

	clk := gpio.NewPin(clkPin)
	clk.Input()

	data := gpio.NewPin(dataPin)
	data.Input()

	btn := gpio.NewPin(btnPin)
	btn.Input()

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
