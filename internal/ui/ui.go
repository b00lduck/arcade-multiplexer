package ui

import (
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type Display interface {
	ShowImage(image.Image) error
}

type ui struct {
	display      Display
	currentGame  string
	selectedGame string
}

func NewUi(display Display) *ui {
	return &ui{
		display: display}
}

func (u *ui) SelectedGame(s string) {
	u.selectedGame = s
	u.Draw()
}

func (u *ui) CurrentGame(s string) {
	u.currentGame = s
	u.Draw()
}

func (u *ui) Draw() {

	i := image.NewRGBA(image.Rect(0, 0, 128, 64))
	src := image.NewUniform(color.RGBA{255, 255, 255, 255})
	d := &font.Drawer{
		Dst:  i,
		Src:  src,
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{}}

	d.Dot = fixed.Point26_6{
		fixed.Int26_6(64),
		fixed.Int26_6(2 * 14 * 64)}
	d.DrawString(u.currentGame)

	d.Dot = fixed.Point26_6{
		fixed.Int26_6(64),
		fixed.Int26_6(4 * 14 * 64)}
	d.DrawString(u.selectedGame)

	u.display.ShowImage(i)

}
