package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/hc595"
	"github.com/b00lduck/arcade-multiplexer/internal/matrix"
	"github.com/b00lduck/arcade-multiplexer/internal/oled"
	"github.com/b00lduck/arcade-multiplexer/internal/rotary"
	"github.com/b00lduck/arcade-multiplexer/internal/state"
	"github.com/b00lduck/arcade-multiplexer/internal/ui"
	"github.com/warthog618/gpio"
)

func main() {

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)

	err := gpio.Open()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()

	hc595 := hc595.NewHc595(17, 27, 22)
	hc595.SendByte(1023)
	time.Sleep(1 * time.Second)
	hc595.SendByte(0)

	oled := oled.NewOled("/dev/i2c-1")
	defer oled.Close()

	ui := ui.NewUi(oled)

	state := state.NewState(ui)

	rotary := rotary.NewRotary(5, 6, 19, state.Up, state.Down, state.Choose)
	defer rotary.Close()

	matrix := matrix.NewMatrix([]uint8{23, 24, 25, 26}, []uint8{14, 15, 18})
	go matrix.Run()

	select {
	case <-time.After(time.Minute):
	case <-quit:
	}

}
