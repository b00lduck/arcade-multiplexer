package display

import (
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"arcade-multiplexer/internal/framebuffer"
)

type display struct {
	fb *framebuffer.Framebuffer
}

func NewDisplay(fb *framebuffer.Framebuffer) *display {
	return &display{
		fb: fb}
}

func (d *display) ShowImage(filename string) {
	img, err := LoadImage(filename)
	if err != nil {
		return
	}
	draw.Draw(d.fb, d.fb.Bounds(), img, image.ZP, draw.Src)
}

func LoadImage(filename string) (image.Image, error) {

	f, err := os.Open("images/" + filename)
	if err != nil {
		log.Error().Err(err).Str("filename", filename).Msg("Could not load image")
		return nil, err
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
		log.Error().Err(err).Str("filename", filename).Msg("Could not decode image")
	}

	return img, nil
}
