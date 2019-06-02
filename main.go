package main

import (
	"io/ioutil"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"

	"image"
	"os"

	"github.com/goiot/devices/monochromeoled"
	"golang.org/x/exp/io/i2c"

	_ "image/png"
)

func main() {

	d1 := []byte("17\n")
	ioutil.WriteFile("/sys/class/gpio/unexport", d1, 0644)

	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	led, err := embd.NewDigitalPin(17)
	if err != nil {
		panic(err)
	}
	defer led.Close()
	if err := led.SetDirection(embd.In); err != nil {
		panic(err)
	}

	rc, err := os.Open("./test.png")
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	m, _, err := image.Decode(rc)
	if err != nil {
		panic(err)
	}

	d, err := monochromeoled.Open(&i2c.Devfs{Dev: "/dev/i2c-1"})
	if err != nil {
		panic(err)
	}
	defer d.Close()

	// clear the display before putting on anything
	if err := d.Clear(); err != nil {
		panic(err)
	}
	if err := d.SetImage(0, 0, m); err != nil {
		panic(err)
	}
	if err := d.Draw(); err != nil {
		panic(err)
	}
}
