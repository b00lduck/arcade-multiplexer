package ui

import (
	"image"
	"os"

	"github.com/rs/zerolog/log"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	"arcade-multiplexer/internal/config"
	"arcade-multiplexer/internal/data"
	"arcade-multiplexer/internal/framebuffer"
	"arcade-multiplexer/internal/imageCache"
)

type Panel interface {
	SetLeds(data.LedState)
}

type InputProcessor interface {
	SetMappings([]config.Mapping)
}

type MistControl interface {
	ExitCore(*config.Core)
	ChangeCore(*config.Core)
	LoadGame(*config.Game, *config.Core, bool)
}

type ui struct {
	framebuffer     *framebuffer.DisplayFramebuffer
	panel           Panel
	config          *config.Config
	oldGame         config.Game
	inputProcessor  InputProcessor
	mistControl     MistControl
	imageCache      *imageCache.ImageCache
	selectedGame    config.Game
	hudPosition     image.Point
	panelUpdates    *chan data.MatrixState
	lastRegionState map[string]int
	ftCtx           *freetype.Context
	font            *truetype.Font
	face            font.Face
}

const FONT = "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"
const FONT_DPI = 96
const FONT_SIZE = 10

func NewUi(c *config.Config,
	framebuffer *framebuffer.DisplayFramebuffer,
	panel Panel,
	ip InputProcessor,
	mistControl MistControl,
	imageCache *imageCache.ImageCache,
	panelUpdates *chan data.MatrixState) *ui {

	// create Font
	fontBytes, err := os.ReadFile(FONT)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not read font file")
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse font")
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: FONT_SIZE,
		DPI:  FONT_DPI,
	})

	ftCtx := freetype.NewContext()

	ftCtx.SetDPI(96)
	ftCtx.SetSrc(image.White)
	ftCtx.SetDst(&framebuffer.BGRA)
	ftCtx.SetFont(font)
	ftCtx.SetFontSize(10)
	ftCtx.SetClip(framebuffer.Bounds())

	return &ui{
		framebuffer: framebuffer,
		panel:       panel,
		config:      c,
		oldGame: config.Game{
			Core: "none",
		},
		selectedGame: config.Game{
			Core: "none",
		},
		inputProcessor:  ip,
		mistControl:     mistControl,
		imageCache:      imageCache,
		hudPosition:     image.Point{0, 765},
		panelUpdates:    panelUpdates,
		lastRegionState: make(map[string]int),
		ftCtx:           ftCtx,
		font:            font,
		face:            face,
	}
}

func (u *ui) Start() {
	u.framebuffer.Clear()

	// wait for panel updates
	for state := range *u.panelUpdates {
		u.DrawHud(state, false)
	}
}

func (u *ui) StartGameById(id uint32) {
	u.startGame(u.config.Games[id])
}

func (u *ui) SelectGameById(id uint32) {
	u.selectGame(u.config.Games[id])
}

func (u *ui) startGame(game config.Game) {

	//log.Info().Interface("game", game).Msg("Starting game")

	u.panel.SetLeds(data.LedStateByMapping(game.Mappings))

	u.inputProcessor.SetMappings(game.Mappings)

	newCore := u.config.GetCoreByName(game.Core)
	oldCore := u.config.GetCoreByName(u.oldGame.Core)

	log.Info().Str("oldCore", u.oldGame.Core).Str("newCore", game.Core).Msg("cores")

	if u.oldGame.Core != game.Core {
		u.mistControl.ExitCore(oldCore)
		u.mistControl.ChangeCore(newCore)
		u.mistControl.LoadGame(&game, newCore, false)
	} else {
		u.mistControl.LoadGame(&game, newCore, true)
	}
	u.oldGame = game
}

func (u *ui) selectGame(game config.Game) {

	//log.Info().Interface("game", game).Msg("Selected game")
	if game.Image != "" {
		img := u.imageCache.Images[game.Image]
		if img != nil {
			u.framebuffer.Blit(25, 25, img)
		}
	}

	u.framebuffer.Blit(
		u.hudPosition.X,
		u.hudPosition.Y,
		u.imageCache.Images["hud_1.jpg"],
	)

	u.selectedGame = game

	u.DrawHud(data.MatrixState{}, true)
	u.DrawLabels()
}

func (u *ui) DrawLabels() {

	for regName, reg := range u.config.Ui.Regions {

		text := ""
		for _, v := range u.selectedGame.Mappings {
			if v.Input == regName {
				if v.Text != "" {
					text = v.Text
				} else {
					text = u.config.Outputs[v.Output].Alias
				}
			}
		}

		// iterate over the string and sum the width of each character
		// to get the total width of the string
		var width fixed.Int26_6
		var height fixed.Int26_6

		for _, r := range text {
			width += u.font.HMetric(700, truetype.Index(r)).AdvanceWidth
			aheight := u.font.VMetric(700, truetype.Index(r))
			if aheight.AdvanceHeight > height {
				height = aheight.AdvanceHeight
			}
		}

		xPosInt := u.hudPosition.X + reg.X + reg.Textbox.X
		yPosInt := u.hudPosition.Y + reg.Y + reg.Textbox.Y

		switch reg.Textbox.Align {
		case "right":
			xPosInt -= u.measureStringWidth(text)
		case "center":
			xPosInt -= u.measureStringWidth(text) / 2
		}

		posPt := freetype.Pt(xPosInt, yPosInt)
		posPt.Y += height

		_, err := u.ftCtx.DrawString(text, posPt)
		if err != nil {
			log.Error().Err(err).Msg("Could not draw text")
		}

	}
}

func (u *ui) measureStringWidth(text string) int {
	var width int
	for _, ch := range text {
		advance, ok := u.face.GlyphAdvance(ch)
		if !ok {
			log.Fatal().Interface("rune", ch).Msg("failed to get advance for rune")
		}
		width += int(advance.Round())
	}
	return width
}

func (u *ui) DrawHud(state data.MatrixState, redraw bool) {
	ledState := data.LedStateByMapping(u.selectedGame.Mappings)

	//log.Info().Interface("state", state).Interface("ledState", ledState).Msg("Drawing HUD")

	u.LedRegion(state.GlobalKeypad.WhiteLeft, ledState.GlobalKeypad.WhiteLeft, "WHITE_LEFT", redraw)
	u.LedRegion(state.GlobalKeypad.WhiteRight, ledState.GlobalKeypad.WhiteRight, "WHITE_RIGHT", redraw)
	u.LedRegion(state.Player1Keypad.Red, ledState.Player1Keypad.Red, "P1_RED", redraw)
	u.LedRegion(state.Player1Keypad.Yellow, ledState.Player1Keypad.Yellow, "P1_YELLOW", redraw)
	u.LedRegion(state.Player1Keypad.Green, ledState.Player1Keypad.Green, "P1_GREEN", redraw)
	u.LedRegion(state.Player1Keypad.Blue, ledState.Player1Keypad.Blue, "P1_BLUE", redraw)
	u.LedRegion(state.GlobalKeypad.FlipperLeft, ledState.GlobalKeypad.FlipperLeft, "P1_FLIPPER", redraw)
	u.LedRegion(state.Player2Keypad.Red, ledState.Player2Keypad.Red, "P2_RED", redraw)
	u.LedRegion(state.Player2Keypad.Yellow, ledState.Player2Keypad.Yellow, "P2_YELLOW", redraw)
	u.LedRegion(state.Player2Keypad.Green, ledState.Player2Keypad.Green, "P2_GREEN", redraw)
	u.LedRegion(state.Player2Keypad.Blue, ledState.Player2Keypad.Blue, "P2_BLUE", redraw)
	u.LedRegion(state.GlobalKeypad.FlipperRight, ledState.GlobalKeypad.FlipperRight, "P2_FLIPPER", redraw)

	u.LedRegion(state.Player1Joystick.Up, true, "P1_JOY_UP", redraw)
	u.LedRegion(state.Player1Joystick.Down, true, "P1_JOY_DOWN", redraw)
	u.LedRegion(state.Player1Joystick.Left, true, "P1_JOY_LEFT", redraw)
	u.LedRegion(state.Player1Joystick.Right, true, "P1_JOY_RIGHT", redraw)
	u.LedRegion(state.Player2Joystick.Up, true, "P2_JOY_UP", redraw)
	u.LedRegion(state.Player2Joystick.Down, true, "P2_JOY_DOWN", redraw)
	u.LedRegion(state.Player2Joystick.Left, true, "P2_JOY_LEFT", redraw)
	u.LedRegion(state.Player2Joystick.Right, true, "P2_JOY_RIGHT", redraw)
}

func calcRegionState(led bool, pushed bool) int {
	if pushed {
		if led {
			return 4
		} else {
			return 3
		}
	} else {
		if led {
			return 2
		} else {
			return 1
		}
	}
}

func (u *ui) LedRegion(pushed bool, led bool, key string, redraw bool) {
	reg := u.config.Ui.Regions[key]

	newReqionState := calcRegionState(led, pushed)

	if u.lastRegionState[key] == newReqionState && !redraw {
		return
	}

	switch newReqionState {
	case 1:
		u.ShowRegion(reg, "hud_1.jpg", u.hudPosition)
	case 2:
		u.ShowRegion(reg, "hud_2.jpg", u.hudPosition)
	case 3:
		u.ShowRegion(reg, "hud_3.jpg", u.hudPosition)
	case 4:
		u.ShowRegion(reg, "hud_4.jpg", u.hudPosition)
	}

	// draw a border around the region
	/*
		for x := reg.X; x < reg.X+reg.Width; x++ {
			u.framebuffer.BGRA.Set(x+u.hudPosition.X, reg.Y+u.hudPosition.Y, image.White)
			u.framebuffer.BGRA.Set(x+u.hudPosition.X, reg.Y+reg.Height-1+u.hudPosition.Y, image.White)
		}

		for y := reg.Y; y < reg.Y+reg.Height; y++ {
			u.framebuffer.BGRA.Set(reg.X+u.hudPosition.X, y+u.hudPosition.Y, image.White)
			u.framebuffer.BGRA.Set(reg.X+reg.Width-1+u.hudPosition.X, y+u.hudPosition.Y, image.White)
		}
	*/

}

func (u *ui) ShowRegion(region config.Region, filename string, off image.Point) {
	src := u.imageCache.Images[filename]
	dstX := region.X + off.X
	dstY := region.Y + off.Y
	u.framebuffer.BlitSrcWindow(dstX, dstY, region.X, region.Y, region.Width, region.Height, src)
}
