package framebuffer

import (
	"image"
	"image/color"
	"os"

	"github.com/rs/zerolog/log"
)

type ConverterFramebuffer struct {
	data []byte
}

func NewConverterFramebuffer() *ConverterFramebuffer {
	return &ConverterFramebuffer{
		data: make([]byte, screensize),
	}
}

func (f ConverterFramebuffer) Close() {
}

func (f ConverterFramebuffer) Data() []byte {
	return f.data
}

func (f ConverterFramebuffer) Set(x, y int, c color.Color) {
	r, g, b, a := c.RGBA()
	offset := x*4 + y*resx*4
	f.data[offset] = uint8(b / 256)
	f.data[offset+1] = uint8(g / 256)
	f.data[offset+2] = uint8(r / 256)
	f.data[offset+3] = uint8(a / 256)
}

func (f ConverterFramebuffer) At(x, y int) color.Color {
	offset := x*4 + y*resx*4
	a := f.data[offset]
	r := f.data[offset+1]
	g := f.data[offset+2]
	b := f.data[offset+3]
	return color.RGBA{r, g, b, a}
}

func (f ConverterFramebuffer) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{}, image.Point{resx, resy}}
}

func (f ConverterFramebuffer) ColorModel() color.Model {
	return color.RGBAModel
}

func (f ConverterFramebuffer) Dump(target string) {
	err := os.WriteFile(target, f.data, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to write output file")
	}
	//os.Stdout.Write(f.data[:])
}
