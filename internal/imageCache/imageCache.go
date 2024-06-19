package imageCache

import (
	"arcade-multiplexer/internal/config"
	"arcade-multiplexer/internal/framebuffer"
)

type ImageCache struct {
	Images map[string]*framebuffer.ResizedImage
}

func NewImageCache(c *config.Config) *ImageCache {
	return &ImageCache{
		Images: make(map[string]*framebuffer.ResizedImage),
	}
}

func (i *ImageCache) LoadAll(c *config.Config) {

	// Load HUD background
	i.GetImage(608, 259, "hud_1.jpg", true)
	i.GetImage(608, 259, "hud_2.jpg", true)
	i.GetImage(608, 259, "hud_3.jpg", true)
	i.GetImage(608, 259, "hud_4.jpg", true)

	// Load game covers
	for _, game := range c.Games {
		i.GetImage(555, 740, game.Image, false)
	}

}

func (i *ImageCache) GetImage(width, height int, filename string, flip bool) {
	i.Images[filename] = framebuffer.NewResizedImageFromImageFile(width, height, filename, flip)
}
