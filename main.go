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

	"github.com/b00lduck/arcade-multiplexer/internal/aux"
	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/b00lduck/arcade-multiplexer/internal/framebuffer"
	"github.com/b00lduck/arcade-multiplexer/internal/hc595"
	"github.com/b00lduck/arcade-multiplexer/internal/matrix"
	"github.com/b00lduck/raspberry-datalogger-display/tools"
	"github.com/tarent/logrus"
	"github.com/warthog618/gpio"
)

func main() {

	// capture exit signals
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	err := gpio.Open()
	if err != nil {
		logrus.WithError(err).Fatal("Could not open GPIO")
	}
	defer gpio.Close()

	fb := framebuffer.NewFramebuffer("/dev/fb1")
	defer fb.Close()

	splash := LoadImage("turrican2.jpg")
	draw.Draw(*fb, (*fb).Bounds(), splash, image.ZP, draw.Src)

	hc595 := hc595.NewHc595(26, 27, 22)

	aux := aux.NewAux(23, 20)
	aux.SetPwr(true)

	hc595.SetLeds(data.LedState{
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

	time.Sleep(1 * time.Second)

	hc595.SetLeds(data.LedState{
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
			WhiteLeft:  true,
			WhiteRight: true}})

	matrix := matrix.NewMatrix(func(row uint8) {
		hc595.SelectRow(row)
	}, 4, []uint8{14, 15, 18, 12, 16})

	go func() {
		<-quit
		logrus.Info("Shuttong down")
		aux.SetPwr(false)
		os.Exit(0)
	}()

	matrix.Run(func(ms *data.MatrixState) {
		fmt.Println("Jostick 1: ", ms.Player1Joystick.String())
		fmt.Println("Jostick 2: ", ms.Player2Joystick.String())
		fmt.Println("Keypad 1:  ", ms.Player1Keypad.String())
		fmt.Println("Keypad 2:  ", ms.Player2Keypad.String())
		fmt.Println("Global:    ", ms.GlobalKeypad.String())

		// TODO: transform input matrix

		hc595.SetJoys(&ms.Player1Joystick, &ms.Player2Joystick, ms.Player1Keypad.Red, ms.Player2Keypad.Red)

		fmt.Printf("\033[5A")
	})

}

func LoadImage(filename string) image.Image {

	f, err := os.Open("images/" + filename)
	tools.ErrorCheck(err)

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

	tools.ErrorCheck(err)

	return img
}
