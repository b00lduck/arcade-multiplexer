package aux

import (
	"github.com/warthog618/gpio"
)

type aux struct {
	pwr *gpio.Pin
	rst *gpio.Pin
}

func NewAux(pwrPin, rstPin uint8) *aux {

	pwr := gpio.NewPin(pwrPin)
	pwr.Output()
	pwr.High()

	rst := gpio.NewPin(rstPin)
	rst.Output()
	rst.High()

	return &aux{
		pwr: pwr,
		rst: rst}

}

func (o *aux) SetPwr(state bool) {
	if state {
		o.pwr.Low()
	} else {
		o.pwr.High()
	}
}

func (o *aux) SetRst(state bool) {
	if state {
		o.pwr.Low()
	} else {
		o.pwr.High()
	}
}
