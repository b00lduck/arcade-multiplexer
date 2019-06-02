package oled

import (
	"image"
	"os"

	"github.com/goiot/devices/monochromeoled"
	"golang.org/x/exp/io/i2c"

	_ "image/png"
)

type oled struct {
	fs *monochromeoled.OLED
}

func NewOled(device string) *oled {

	fs, err := monochromeoled.Open(&i2c.Devfs{Dev: device})
	if err != nil {
		panic(err)
	}

	return &oled{
		fs: fs}

}

func (o *oled) Close() {
	o.fs.Close()
}

func (o *oled) ShowImage(filename string) error {
	rc, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer rc.Close()

	m, _, err := image.Decode(rc)
	if err != nil {
		return err
	}

	if err := o.fs.Clear(); err != nil {
		return err
	}
	if err := o.fs.SetImage(0, 0, m); err != nil {
		return err
	}
	if err := o.fs.Draw(); err != nil {
		return err
	}

	return nil
}
