package rotary

import (
	"fmt"

	"github.com/warthog618/gpio"
)

type rotary struct {
	clk *gpio.Pin
	dt  *gpio.Pin
	btn *gpio.Pin
}

func NewRotary(clkPin, dtPin, btnPin uint8) *rotary {

	dt := gpio.NewPin(dtPin)
	dt.Input()
	dt.PullUp()

	clk := gpio.NewPin(clkPin)
	clk.Input()
	clk.PullUp()
	err := clk.Watch(gpio.EdgeFalling, func(pin *gpio.Pin) {
		if dt.Read() == gpio.Low {
			fmt.Println("down")
		} else {
			fmt.Println("up")
		}

	})
	if err != nil {
		panic(err)
	}

	btn := gpio.NewPin(btnPin)
	btn.Input()
	btn.PullUp()
	err = btn.Watch(gpio.EdgeFalling, func(pin *gpio.Pin) {
		fmt.Println("choose")
	})
	if err != nil {
		panic(err)
	}

	return &rotary{
		dt:  dt,
		clk: clk,
		btn: btn}

}

func (r *rotary) Close() {
	r.btn.Unwatch()
	r.dt.Unwatch()
	r.clk.Unwatch()
}
