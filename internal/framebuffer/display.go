package framebuffer

import (
	"image"
	"os"
	"syscall"

	"github.com/jmigpin/editor/util/imageutil"
	"github.com/rs/zerolog/log"
)

type DisplayFramebuffer struct {
	imageutil.BGRA
	file *os.File
}

const BPP = 4 // Bytes per pixel

func NewDisplayFramebuffer(device string) *DisplayFramebuffer {

	log.Info().Str("device", device).Msg("initializing DisplayFramebuffer")

	file, err := os.OpenFile(device, os.O_RDWR, 0)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not open DisplayFramebuffer")
	}

	width := 608
	height := 1024
	bufferSize := width * height * BPP

	data, err := syscall.Mmap(int(file.Fd()), 0, bufferSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not open DisplayFramebuffer")
	}

	image := imageutil.NewBGRA(&image.Rectangle{image.Point{}, image.Point{width, height}})
	image.Pix = data

	return &DisplayFramebuffer{
		BGRA: *image,
		file: file,
	}
}

func (f *DisplayFramebuffer) Close() {
	f.file.Close()
}

func (f *DisplayFramebuffer) Clear() {
	for i := range f.Pix {
		f.Pix[i] = 0
	}
}

/**
 * Blit (Block image ttransfer) copies the contents of the src image to the framebuffer
 * at the specified x and y coordinates.
 */
func (f *DisplayFramebuffer) Blit(x, y int, src *ResizedImage) {
	dst := &f.BGRA
	dstPtr := x*BPP + y*dst.Stride
	srcEnd := src.Stride * src.Rect.Dy()
	sourceWidthBytesMinusOne := src.Stride - 1
	for srcPtr := 0; srcPtr < srcEnd; srcPtr += src.Stride {
		copy(dst.Pix[dstPtr:], src.Pix[srcPtr:srcPtr+sourceWidthBytesMinusOne])
		dstPtr += dst.Stride
	}
}

/**
 * BlitSrcWindow copies a window of the src image to the framebuffer
 * at the specified x and y coordinates.
 */
func (f *DisplayFramebuffer) BlitSrcWindow(x, y, sx, sy, sw, sh int, src *ResizedImage) {
	dst := &f.BGRA
	dstPtr := x*BPP + y*dst.Stride
	srcStart := sx*BPP + sy*src.Stride
	srcEnd := srcStart + src.Stride*sh
	sourceWidthBytesMinusOne := sw*BPP - 1
	for srcPtr := srcStart; srcPtr < srcEnd; srcPtr += src.Stride {
		copy(dst.Pix[dstPtr:], src.Pix[srcPtr:srcPtr+sourceWidthBytesMinusOne])
		dstPtr += dst.Stride
	}
}
