package framebuffer

import (
	"image"
	"image/color"
)

type PixelFormat struct {
	data []byte
}

func (f PixelFormat) Data() []byte {
	return f.data
}

func (f PixelFormat) Set(x, y int, c color.Color) {
	offset := x*4 + y*stride
	r, g, b, a := c.RGBA()
	f.data[offset] = uint8(b >> 8)
	f.data[offset+1] = uint8(g >> 8)
	f.data[offset+2] = uint8(r >> 8)
	f.data[offset+3] = uint8(a >> 8)
}

func (f PixelFormat) At(x, y int) color.Color {
	offset := x*4 + y*stride
	a := f.data[offset]
	r := f.data[offset+1]
	g := f.data[offset+2]
	b := f.data[offset+3]
	return color.RGBA{r, g, b, a}
}

func (f PixelFormat) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{}, image.Point{resx, resy}}
}

func (f PixelFormat) ColorModel() color.Model {
	return color.RGBAModel
}
