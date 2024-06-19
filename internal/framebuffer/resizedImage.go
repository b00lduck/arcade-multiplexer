package framebuffer

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/jmigpin/editor/util/imageutil"

	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
)

type ResizedImage struct {
	imageutil.BGRA
	filename      string
	cacheFilename string
}

func NewResizedImage(width, height int, filename string) *ResizedImage {
	image := imageutil.NewBGRA(&image.Rectangle{image.Point{}, image.Point{width, height}})
	return &ResizedImage{
		BGRA:          *image,
		filename:      filename,
		cacheFilename: "images/" + filename + ".cache",
	}
}

func NewResizedImageFromImageFile(width, height int, filename string, flip bool) *ResizedImage {

	ret := NewResizedImage(width, height, filename)

	if _, err := os.Stat(ret.cacheFilename); err == nil {

		log.Info().Str("filename", filename).Msg("loading image from cache")

		// read whole file at once into byte array
		data, err := os.ReadFile(ret.cacheFilename)
		if err != nil {
			log.Fatal().Err(err).Str("filename", ret.cacheFilename).Msg("Could not read cache file")
		}
		ret.Pix = data

	} else {
		ret.load(filename, flip)
		ret.CacheStore()
	}

	return ret
}

func (f ResizedImage) CacheStore() {
	err := os.WriteFile(f.cacheFilename, f.Pix, 0644)
	if err != nil {
		log.Fatal().Err(err).Str("filename", f.cacheFilename).Msg("Could not write cache file")
	}
}

func (f ResizedImage) load(filename string, flip bool) {
	log.Info().Str("filename", filename).Msg("loading image")

	img, err := loadImage("images/" + filename)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading image")
	}

	// scale image to fit the framebuffer limits
	// keep the aspect ratio
	fbWidth := f.BGRA.Bounds().Dx()
	fbHeight := f.BGRA.Bounds().Dy()

	srcRatio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())
	dstRatio := float64(fbWidth) / float64(fbHeight)

	var dstWidth, dstHeight int
	if srcRatio > dstRatio {
		dstWidth = fbWidth
		dstHeight = int(float64(fbWidth) / srcRatio)
	} else {
		dstHeight = fbHeight
		dstWidth = int(float64(fbHeight) * srcRatio)
	}

	x := (fbWidth - dstWidth) / 2
	y := (fbHeight - dstHeight) / 2

	dstSize := image.Rectangle{
		image.Point{x, y},
		image.Point{x + dstWidth, y + dstHeight},
	}

	draw.NearestNeighbor.Scale(&f.BGRA, dstSize, img, img.Bounds(), draw.Src, nil)

	if flip {
		f.flipChannels()
	}
}

func (f *ResizedImage) flipChannels() {
	for i := 0; i < len(f.Pix); i += 4 {
		f.Pix[i], f.Pix[i+2] = f.Pix[i+2], f.Pix[i]
	}
}

func loadImage(filename string) (image.Image, error) {

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
