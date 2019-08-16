package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/config"
	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/b00lduck/arcade-multiplexer/internal/display"
	"github.com/b00lduck/arcade-multiplexer/internal/framebuffer"
	"github.com/b00lduck/arcade-multiplexer/internal/hid"
	"github.com/b00lduck/arcade-multiplexer/internal/inputProcessor"
	"github.com/b00lduck/arcade-multiplexer/internal/mist"
	"github.com/b00lduck/arcade-multiplexer/internal/panel"
	"github.com/b00lduck/arcade-multiplexer/internal/rotary"
	"github.com/b00lduck/arcade-multiplexer/internal/ui"
	"github.com/tarent/logrus"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

type Mist interface {
	SetJoystick1(joy *data.Joystick)
	SetJoystick2(joy *data.Joystick)
	SetJoystick1Button1(state bool)
	SetJoystick1Button2(state bool)
	SetJoystick2Button1(state bool)
	SetJoystick2Button2(state bool)
}

func main() {

	// capture exit signals
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Initialize periph.io library
	// see https://periph.io/project/library/
	_, err := host.Init()
	if err != nil {
		logrus.WithError(err).Fatal("Could not open initialize periph.io")
	}

	// Open the first available I²C bus
	bus, err := i2creg.Open("")
	if err != nil {
		logrus.WithError(err).Fatal("Could not open i2c bus")
	}

	// Initialize TFT framebuffer and display
	fb := framebuffer.NewFramebuffer("/dev/fb1")
	defer fb.Close()
	display := display.NewDisplay(fb)
	display.ShowImage("splash.jpg")

	// Load game config from yml file

	c := config.NewConfig()

	// Initialize connection to MiST-interface board. This
	// contains Joystick inputs, power and reset
	mistDigital := mist.NewMistDigital(bus)
	go mistDigital.Run()

	hid := hid.NewHid()
	defer hid.Close()

	mistControl := mist.NewMistControl(hid, mistDigital)

	// Initialize connection to panel board and set input processor that
	// translates the inputs from joysticks and buttons to the configured
	// outputs of the active game
	inputProcessor := inputProcessor.NewInputProcessor(mistDigital, hid)
	panel := panel.NewPanel(bus, inputProcessor)
	go panel.Run()

	// Exit handler routine, triggered by signal (see above)
	go func() {
		<-quit
		logrus.Info("Shutting down")
		//mist.SetPower(false)
		// give some time to shut down the power pin via i2c
		//time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()

	mistDigital.SetPower(true)

	ui := ui.NewUi(c, display, panel, inputProcessor, mistControl)

	rotary := rotary.NewRotary(4, 5, 6, len(c.Games), ui.StartGameById, ui.SelectGameById)
	go rotary.Run()

	for {
		time.Sleep(500 * time.Millisecond)
	}

}
