package rotary

import (
	"fmt"
	"sync"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"

	"github.com/rs/zerolog/log"
)

type rotary struct {
	clk            gpio.PinIO
	data           gpio.PinIO
	btn            gpio.PinIO
	delta          int
	oldButtonState bool
	mutex          *sync.Mutex
	clickCallback  func(uint32)
	chooseCallback func(uint32)
	count          int
	currentPosi    int
}

func NewRotary(clkPin, dataPin, btnPin uint8, count int, clickCallback, chooseCallback func(uint32)) *rotary {

	clk := gpioreg.ByName(fmt.Sprintf("%d", clkPin))
	clk.In(gpio.Float, gpio.NoEdge)

	data := gpioreg.ByName(fmt.Sprintf("%d", dataPin))
	clk.In(gpio.Float, gpio.NoEdge)

	btn := gpioreg.ByName(fmt.Sprintf("%d", btnPin))
	clk.In(gpio.Float, gpio.NoEdge)

	return &rotary{
		clk:            clk,
		data:           data,
		btn:            btn,
		mutex:          &sync.Mutex{},
		clickCallback:  clickCallback,
		count:          count,
		chooseCallback: chooseCallback}
}

func (o *rotary) Run() {
	go o.runAcquisition()
	o.runEvaluation()
}

func (o *rotary) runAcquisition() {
	encLast := 1
	o.oldButtonState = bool(o.btn.Read())
	for {

		btnRead := bool(o.btn.Read())
		if btnRead != o.oldButtonState {
			o.oldButtonState = btnRead
			if !btnRead {
				o.clickCallback(uint32(o.currentPosi))
			}
		}

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

func (o *rotary) runEvaluation() {
	posi := 0

	for {
		d := o.fetchDelta()
		if d != 0 {
			posi -= d

			if posi > (o.count-1)*4 {
				posi = 2
			} else if posi < 0 {
				posi = o.count*4 - 2
			}

			if o.currentPosi != posi/4 {
				o.currentPosi = posi / 4
				log.Info().Int("pos", o.currentPosi).Msg("Rotary pos")
				o.chooseCallback(uint32(o.currentPosi))
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func (o *rotary) fetchDelta() int {
	o.mutex.Lock()
	ret := o.delta
	o.delta = 0
	o.mutex.Unlock()
	return ret
}
