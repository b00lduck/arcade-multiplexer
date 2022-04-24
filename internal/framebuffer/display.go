package framebuffer

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/rs/zerolog/log"
)

const resx = 320
const resy = 480
const depth = 16
const screensize = resx * resy * depth / 8

type DisplayFramebuffer struct {
	file *os.File
	data []byte
}

func NewDisplayFramebuffer(device string) *DisplayFramebuffer {

	fmt.Printf("initializing DisplayFramebuffer on device %s\n", device)

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
		file: file,
		data: data}

}

func (f *DisplayFramebuffer) Close() {
	f.file.Close()
}

func (f *DisplayFramebuffer) ShowImage(image string) {

	filename := "images/" + image + ".565.data"

	s, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Str("filename", filename).Msg("Could not load image")
		return
	}
	defer s.Close()

	_, err = io.ReadFull(s, f.data)

	if err != nil {
		if err != io.EOF {
			log.Error().Err(err).Str("filename", filename).Msg("error reading image data")
			return
		}
	}

}
