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

func (f ConverterFramebuffer) convertToFb(c color.Color) (msb, lsb uint8) {
	r, g, b, _ := c.RGBA()
	out := uint16(r>>11)<<11 + uint16(g>>10)<<5 + uint16(b>>11)
	lsb = uint8(out >> 8)
	msb = uint8(out & 0xff)
	return
}

func (f ConverterFramebuffer) Set(x, y int, c color.Color) {
	offset := x*2 + y*resx*2
	f.data[offset], f.data[offset+1] = f.convertToFb(c)
}

func (f ConverterFramebuffer) convertToRgba(lsb, msb uint8) color.Color {
	val := uint16(msb)<<8 + uint16(lsb)
	b := uint8((val & 0x001F) << 3)
	g := uint8((val & 0x07E0) >> 3)
	r := uint8((val & 0xF800) >> 8)
	a := uint8(0xFF)
	return color.RGBA{r, g, b, a}
}

func (f ConverterFramebuffer) At(x, y int) color.Color {
	offset := x*2 + y*resx*2
	return f.convertToRgba(f.data[offset], f.data[offset+1])
}

func (f ConverterFramebuffer) Bounds() image.Rectangle {
	return image.Rectangle{image.ZP, image.Point{resx, resy}}
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
