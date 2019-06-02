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
	//err := dt.Watch(gpio.EdgeBoth, func(pin *gpio.Pin) {
	//fmt.Printf("DT Pin is %v\n", pin.Read())
	//})
	//if err != nil {
	//	panic(err)
	//}

	clk := gpio.NewPin(clkPin)
	clk.Input()
	clk.PullUp()
	err := clk.Watch(gpio.EdgeBoth, func(pin *gpio.Pin) {
		if clk.Read() == dt.Read() {
			fmt.Println("right")
		} else {
			fmt.Println("left")
		}
	})
	if err != nil {
		panic(err)
	}

	btn := gpio.NewPin(btnPin)
	btn.Input()
	btn.PullUp()
	//	err = btn.Watch(gpio.EdgeFalling, func(pin *gpio.Pin) {
	//		fmt.Printf("BTN Pin is %v\n", pin.Read())
	//	})
	//	if err != nil {
	//		panic(err)
	//	}

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
