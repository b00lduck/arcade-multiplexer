package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/b00lduck/arcade-multiplexer/internal/pcf8574"
	"github.com/tarent/logrus"
)

func main() {

	// capture exit signals
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		logrus.Info("Shutting down")
		//aux.SetPwr(false)
		os.Exit(0)
	}()

	mpx := pcf8574.NewPcf8574()
	go mpx.Run(func(ms data.MatrixState) {
		fmt.Println(ms.String())

		// TODO: transform input matrix

		//hc595.SetJoys(&ms.Player1Joystick, &ms.Player2Joystick,
		//	ms.Player1Keypad.Red, ms.Player1Keypad.Yellow,
		//	ms.Player2Keypad.Red, ms.Player2Keypad.Yellow)

		fmt.Printf("\033[7A")
	})

	/*err := gpio.Open()
	if err != nil {
		logrus.WithError(err).Fatal("Could not open GPIO")
	}
	defer gpio.Close()

	fb := framebuffer.NewFramebuffer("/dev/fb1")
	defer fb.Close()

	splash := LoadImage("turrican2.jpg")
	draw.Draw(*fb, (*fb).Bounds(), splash, image.ZP, draw.Src)

	aux := aux.NewAux(23, 20)
	aux.SetPwr(true)
	*/

	for {
		mpx.SetLeds(data.LedState{
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

		mpx.SetLeds(data.LedState{
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
	/*
		matrix := matrix.NewMatrix(func(row uint8) {
			hc595.SelectRow(row)
		}, 4, []uint8{14, 15, 18, 12, 16})

		rotary := rotary.NewRotary(0, 1, 21)
		go rotary.Run()

		posi := 0

		go func() {
			for {
				d := rotary.Delta()
				if d != 0 {
					posi -= d
					logrus.Info(posi / 4)
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()

	*/

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
