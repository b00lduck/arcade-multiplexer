package framebuffer

import (
	"arcade-multiplexer/internal/converter"
	"image"
	"image/draw"
	"os"
	"syscall"

	"github.com/rs/zerolog/log"
)

const resx = 480
const resy = 640
const depth = 32
const stride = resx * depth / 8
const screensize = resy * stride

type DisplayFramebuffer struct {
	PixelFormat
	file *os.File
}

func NewDisplayFramebuffer(device string) *DisplayFramebuffer {

	log.Info().Str("device", device).Msg("initializing DisplayFramebuffer")

	file, err := os.OpenFile(device, os.O_RDWR, 0)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not open DisplayFramebuffer")
	}
	data, err := syscall.Mmap(int(file.Fd()), 0, screensize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not open DisplayFramebuffer")
	}
	return &DisplayFramebuffer{
		PixelFormat: PixelFormat{
			data: data,
		},
		file: file,
	}
}

func (f *DisplayFramebuffer) Close() {
	f.file.Close()
}

func (f *DisplayFramebuffer) ShowImage(i string) {

	filename := "images/" + i

	img, err := converter.LoadImage(filename)
	if err != nil {
		log.Error().Err(err).Str("filename", filename).Msg("Could not load image")
		return
	}

	draw.Draw(f, f.Bounds(), img, image.Point{}, draw.Src)
}
