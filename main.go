package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/config"
	"github.com/b00lduck/arcade-multiplexer/internal/cores"
	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/b00lduck/arcade-multiplexer/internal/framebuffer"
	"github.com/b00lduck/arcade-multiplexer/internal/mist"
	"github.com/b00lduck/arcade-multiplexer/internal/panel"
	"github.com/b00lduck/arcade-multiplexer/internal/rotary"
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

	// load config
	c := config.Config{}
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		logrus.WithError(err).Fatal("Error reading yaml file")
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		logrus.WithError(err).Fatal("Error parsing yaml file")
	}

	logrus.Info(c)

	// Initialize connection to MiST-interface board
	mist := mist.NewMist(bus)
	go mist.Run()

	mappings := []config.Mapping{}

	// Initialize connection to panel board
	panel := panel.NewPanel(bus)
	go panel.Run(func(ms data.MatrixState) {
		fmt.Println(ms.String())

		for _, v := range mappings {
			switch v.Input {
			case "P1_JOY":
				OutputJoystick(mist, &ms.Player1Joystick, v.Output)
			case "P1_RED":
				OutputButton(mist, ms.Player1Keypad.Red, v.Output)
			case "P1_YELLOW":
				OutputButton(mist, ms.Player1Keypad.Yellow, v.Output)
			case "P1_BLUE":
				OutputButton(mist, ms.Player1Keypad.Blue, v.Output)
			case "P1_GREEN":
				OutputButton(mist, ms.Player1Keypad.Green, v.Output)

			case "P2_JOY":
				OutputJoystick(mist, &ms.Player2Joystick, v.Output)
			case "P2_RED":
				OutputButton(mist, ms.Player2Keypad.Red, v.Output)
			case "P2_YELLOW":
				OutputButton(mist, ms.Player2Keypad.Yellow, v.Output)
			case "P2_BLUE":
				OutputButton(mist, ms.Player2Keypad.Blue, v.Output)
			case "P2_GREEN":
				OutputButton(mist, ms.Player2Keypad.Green, v.Output)

			case "WHITE_LEFT":
				OutputButton(mist, ms.GlobalKeypad.WhiteLeft, v.Output)
			case "WHITE_RIGHT":
				OutputButton(mist, ms.GlobalKeypad.WhiteRight, v.Output)
			}

		}

		fmt.Printf("\033[7A")
	})

	// Exit handler routine, triggered by signal (see above)
	go func() {
		<-quit
		logrus.Info("Shutting down")
		//mist.SetPower(false)
		// give some time to shut down the power pin via i2c
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()

	fb := framebuffer.NewFramebuffer("/dev/fb1")
	defer fb.Close()

	splash := LoadImage("turrican2.jpg")
	draw.Draw(*fb, (*fb).Bounds(), splash, image.ZP, draw.Src)

	mist.SetPower(true)

	rotary := rotary.NewRotary(4, 5, 6)
	go rotary.Run()

	posi := 0
	oldPosi := 0

	go func() {
		for {
			d := rotary.Delta()
			if d != 0 {
				posi -= d
				if oldPosi != posi/4 {
					oldPosi = posi / 4
					logrus.Info(oldPosi)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	panel.LedsOff()

	// LOAD GAME
	mist.SetResetButton(true)
	time.Sleep(200 * time.Millisecond)
	mist.SetResetButton(false)

	game := c.Games[1]
	mappings = game.Mappings
	panel.SetLeds(data.LedStateByMapping(game.Mappings))
	if game.Image != "" {
		gameImage := LoadImage(game.Image)
		draw.Draw(*fb, (*fb).Bounds(), gameImage, image.ZP, draw.Src)
	}

	time.Sleep(1 * time.Second)

	switch game.Core {
	case "C64":
		cores.ChangeCore(cores.Menu, cores.C64)
		cores.LoadGame(&game, cores.C64)
	case "Amiga":
		cores.ChangeCore(cores.Menu, cores.Amiga)
		cores.LoadGame(&game, cores.Amiga)
	}

	for {
		time.Sleep(500 * time.Millisecond)
	}

}

func LoadImage(filename string) image.Image {

	f, err := os.Open("images/" + filename)
	if err != nil {
		logrus.WithError(err).Fatal("Could not load image")
	}

	var img image.Image

	lowerFilename := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lowerFilename, ".jpg"):
		img, err = jpeg.Decode(f)
	case strings.HasSuffix(lowerFilename, ".png"):
		img, err = png.Decode(f)
	case strings.HasSuffix(lowerFilename, ".gif"):
		img, err = gif.Decode(f)
	}

	if err != nil {
		logrus.WithError(err).Fatal("Could not decode image")
	}

	return img
}

func OutputJoystick(mist Mist, in *data.Joystick, out string) {
	switch out {
	case "JOY1_AXES":
		mist.SetJoystick1(in)
	case "JOY2_AXES":
		mist.SetJoystick2(in)
	}
}

func OutputButton(mist Mist, in bool, out string) {
	switch out {
	case "JOY1_BUTTON1":
		mist.SetJoystick1Button1(in)
	case "JOY1_BUTTON2":
		mist.SetJoystick1Button2(in)
	case "JOY2_BUTTON1":
		mist.SetJoystick2Button1(in)
	case "JOY2_BUTTON2":
		mist.SetJoystick2Button2(in)

	case "JOY1_UP":
		// not implemented
	case "JOY1_DOWN":
		// not implemented
	case "JOY1_LEFT":
		// not implemented
	case "JOY1_RIGHT":
		// not implemented

	case "JOY2_UP":
		// not implemented
	case "JOY2_DOWN":
		// not implemented
	case "JOY2_LEFT":
		// not implemented
	case "JOY2_RIGHT":
		// not implemented

	case "KEY_SPACE":
		// not implemented
	case "KEY_ESCAPE":
		// not implemented
	case "KEY_LSHIFT":
		// not implemented
	case "KEY_RSHIFT":
		// not implemented
	case "KEY_1":
		// not implemented
	case "KEY_2":
		// not implemented
	case "KEY_3":
		// not implemented
	case "KEY_4":
		// not implemented
	}
}
