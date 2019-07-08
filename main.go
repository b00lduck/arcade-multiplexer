package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/b00lduck/arcade-multiplexer/internal/hc595"
	hc595p "github.com/b00lduck/arcade-multiplexer/internal/hc595"
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

	hc595.SendWord(0x00ffe3ff)
	time.Sleep(1 * time.Second)
	hc595.SendWord(0x000003ff)

	oled := oled.NewOled("/dev/i2c-1")
	defer oled.Close()

	ui := ui.NewUi(oled)

	state := state.NewState(ui)

	rotary := rotary.NewRotary(5, 6, 19, state.Up, state.Down, state.Choose)
	defer rotary.Close()

	matrix := matrix.NewMatrix([]uint8{23, 24, 25, 26}, []uint8{14, 15, 18, 12, 16})
	go matrix.Run(func(ms *data.MatrixState) {
		fmt.Println("Jostick 1: ", ms.Player1Joystick.String())
		fmt.Println("Jostick 2: ", ms.Player2Joystick.String())
		fmt.Println("Keypad 1:  ", ms.Player1Keypad.String())
		fmt.Println("Keypad 2:  ", ms.Player2Keypad.String())
		fmt.Println("Global:    ", ms.GlobalKeypad.String())

		// TODO: transform input matrix

		state := hc595.State
		state = hc595p.SetJoystick(state, 0, &ms.Player1Joystick)
		state = hc595p.SetJoystick(state, 1, &ms.Player2Joystick)
		state = hc595p.SetButton(state, 0, ms.Player1Keypad.Red)
		state = hc595p.SetButton(state, 1, ms.Player2Keypad.Red)
		hc595.SendWord(state)

		fmt.Printf("\033[6A")
	})

	select {
	case <-time.After(5 * time.Hour):
	case <-quit:
	}

}
