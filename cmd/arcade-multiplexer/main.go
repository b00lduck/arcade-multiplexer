package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"arcade-multiplexer/internal/config"
	"arcade-multiplexer/internal/data"
	"arcade-multiplexer/internal/framebuffer"
	"arcade-multiplexer/internal/hid"
	"arcade-multiplexer/internal/inputProcessor"
	"arcade-multiplexer/internal/mist"
	"arcade-multiplexer/internal/panel"
	"arcade-multiplexer/internal/rotary"
	"arcade-multiplexer/internal/ui"
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

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Initialize periph.io library
	// see https://periph.io/project/library/
	_, err := host.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not open initialize periph.io")
	}

	// Open the first available I²C bus
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not open i2c bus")
	}

	// Initialize TFT framebuffer and display
	fb := framebuffer.NewDisplayFramebuffer("/dev/fb0")
	defer fb.Close()

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
		log.Info().Msg("Shutting down")
		//mist.SetPower(false)
		// give some time to shut down the power pin via i2c
		//time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()

	mistDigital.SetPower(true)

	ui := ui.NewUi(c, fb, panel, inputProcessor, mistControl)

	rotary := rotary.NewRotary(4, 5, 6, len(c.Games), ui.StartGameById, ui.SelectGameById)
	go rotary.Run()

	fb.ShowImage("splash.jpg")

	for {
		time.Sleep(500 * time.Millisecond)
	}

}
