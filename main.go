package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/b00lduck/arcade-multiplexer/internal/framebuffer"
	"github.com/b00lduck/arcade-multiplexer/internal/mist"
	"github.com/b00lduck/arcade-multiplexer/internal/panel"
	"github.com/b00lduck/arcade-multiplexer/internal/rotary"
	"github.com/tarent/logrus"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

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

	// Initialize connection to MiST-interface board
	mist := mist.NewMist(bus)
	go mist.Run()

	// Initialize connection to panel board
	panel := panel.NewPanel(bus)
	go panel.Run(func(ms data.MatrixState) {
		fmt.Println(ms.String())

		mist.SetJoystick(&ms.Player1Joystick, &ms.Player2Joystick,
			ms.Player1Keypad.Red, ms.Player1Keypad.Yellow,
			ms.Player2Keypad.Red, ms.Player2Keypad.Yellow)

		fmt.Printf("\033[7A")
	})

	// Exit handler routine, triggered by signal (see above)
	go func() {
		<-quit
		logrus.Info("Shutting down")
		mist.SetPower(false)
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

	go func() {
		oldPosi := 0
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

	for {
		panel.SetLeds(data.LedState{
			Player1Keypad: data.PlayerKeypad{
				Red:    true,
				Yellow: true,
				Green:  true,
				Blue:   true},
			Player2Keypad: data.PlayerKeypad{
				Red:    true,
				Yellow: true,
				Green:  true,
				Blue:   true},
			GlobalKeypad: data.GlobalKeypad{
				WhiteLeft:  true,
				WhiteRight: true}})

		time.Sleep(500 * time.Millisecond)

		panel.SetLeds(data.LedState{
			Player1Keypad: data.PlayerKeypad{
				Red:    true,
				Yellow: false,
				Green:  false,
				Blue:   false},
			Player2Keypad: data.PlayerKeypad{
				Red:    true,
				Yellow: false,
				Green:  false,
				Blue:   false},
			GlobalKeypad: data.GlobalKeypad{
				WhiteLeft:  false,
				WhiteRight: false}})

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
