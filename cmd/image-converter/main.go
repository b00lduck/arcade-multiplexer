package main

import (
	"arcade-multiplexer/internal/converter"
	"arcade-multiplexer/internal/framebuffer"
	"image"
	"image/draw"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	source := os.Args[1]

	img, err := converter.LoadImage(source)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading image")
	}

	fb := framebuffer.NewConverterFramebuffer()

	draw.Draw(fb, fb.Bounds(), img, image.Point{}, draw.Src)

	fb.Dump(source + ".argb32.data")

}
