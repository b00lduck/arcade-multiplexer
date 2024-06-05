package framebuffer

import (
	"os"

	"github.com/rs/zerolog/log"
)

type ConverterFramebuffer struct {
	PixelFormat
}

func NewConverterFramebuffer() *ConverterFramebuffer {
	return &ConverterFramebuffer{
		PixelFormat: PixelFormat{
			data: make([]byte, screensize),
		},
	}
}

func (f ConverterFramebuffer) Close() {
}

func (f ConverterFramebuffer) Dump(target string) {
	err := os.WriteFile(target, f.data, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to write output file")
	}
}
