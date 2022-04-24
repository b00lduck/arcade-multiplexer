package converter

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/rs/zerolog/log"

	"os"
	"strings"
)

func LoadImage(filename string) (image.Image, error) {

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Str("filename", filename).Msg("Could not load image")
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
		log.Fatal().Err(err).Str("filename", filename).Msg("Could not decode image")
	}

	return img, nil
}
