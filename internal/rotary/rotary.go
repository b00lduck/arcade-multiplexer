package rotary

import (
	"github.com/b00lduck/arcade-multiplexer/internal/tools"
	"github.com/warthog618/gpio"
)

type rotary struct {
	clkPin         *gpio.Pin
	btnPin         *gpio.Pin
	dtPin          *gpio.Pin
	upCallback     func()
	downCallback   func()
	chooseCallback func()
}

func NewRotary(clkPinNo, dtPinNo, btnPinNo uint8, up, down, choose func()) *rotary {

	tools.Unexport(dtPinNo)
	dtPin := gpio.NewPin(dtPinNo)
	dtPin.Input()
	dtPin.PullUp()

	tools.Unexport(clkPinNo)
	clkPin := gpio.NewPin(clkPinNo)
	clkPin.Input()
	clkPin.PullUp()
	err := clkPin.Watch(gpio.EdgeFalling, func(pin *gpio.Pin) {
		if dtPin.Read() == gpio.Low {
			down()
		} else {
			up()
		}

	})
	if err != nil {
		panic(err)
	}

	tools.Unexport(btnPinNo)
	btnPin := gpio.NewPin(btnPinNo)
	btnPin.Input()
	btnPin.PullUp()
	err = btnPin.Watch(gpio.EdgeFalling, func(pin *gpio.Pin) {
		choose()
	})
	if err != nil {
		panic(err)
	}

	return &rotary{
		clkPin:         clkPin,
		dtPin:          dtPin,
		btnPin:         btnPin,
		upCallback:     up,
		downCallback:   down,
		chooseCallback: choose}

}

func (r *rotary) Close() {
	r.btnPin.Unwatch()
	r.dtPin.Unwatch()
	r.clkPin.Unwatch()
}
