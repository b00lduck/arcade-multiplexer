package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/config"
	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/b00lduck/arcade-multiplexer/internal/display"
	"github.com/b00lduck/arcade-multiplexer/internal/framebuffer"
	"github.com/b00lduck/arcade-multiplexer/internal/inputProcessor"
	"github.com/b00lduck/arcade-multiplexer/internal/mist"
	"github.com/b00lduck/arcade-multiplexer/internal/panel"
	"github.com/b00lduck/arcade-multiplexer/internal/rotary"
	"github.com/b00lduck/arcade-multiplexer/internal/ui"
	"github.com/tarent/logrus"
	"gopkg.in/yaml.v2"
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
	c := config.Config{}
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		logrus.WithError(err).Fatal("Error reading yaml file")
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		logrus.WithError(err).Fatal("Error parsing yaml file")
	}

	// Initialize connection to MiST-interface board. This
	// contains Joystick inputs, power and reset
	mist := mist.NewMist(bus)
	go mist.Run()

	// Initialize connection to panel board and set input processor that
	// translates the inputs from joysticks and buttons to the configured
	// outputs of the active game
	inputProcessor := inputProcessor.NewInputProcessor(mist)
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

	mist.SetPower(true)

	ui := ui.NewUi(&c, display, panel, inputProcessor, mist)

	rotary := rotary.NewRotary(4, 5, 6, len(c.Games), ui.StartGameById, ui.SelectGameById)
	go rotary.Run()

	for {
		time.Sleep(500 * time.Millisecond)
	}

}
