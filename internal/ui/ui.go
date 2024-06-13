package ui

import (
	"image"

	"github.com/rs/zerolog/log"

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
	framebuffer    *framebuffer.DisplayFramebuffer
	panel          Panel
	config         *config.Config
	oldGame        config.Game
	inputProcessor InputProcessor
	mistControl    MistControl
	imageCache     *imageCache.ImageCache
	selectedGame   config.Game
	hudPosition    image.Point
	panelUpdates   chan data.MatrixState
}

func NewUi(c *config.Config,
	framebuffer *framebuffer.DisplayFramebuffer,
	panel Panel,
	ip InputProcessor,
	mistControl MistControl,
	imageCache *imageCache.ImageCache,
	panelUpdates chan data.MatrixState) *ui {

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
		inputProcessor: ip,
		mistControl:    mistControl,
		imageCache:     imageCache,
		hudPosition:    image.Point{0, 765},
		panelUpdates:   panelUpdates,
	}
}

func (u *ui) Start() {
	u.framebuffer.Clear()
	u.framebuffer.Blit(u.hudPosition.X, u.hudPosition.Y, u.imageCache.Images["hud_1.jpg"])
	// wait for panel updates
	for {
		select {
		case state := <-u.panelUpdates:
			u.DrawHud(state)
		}
	}
}

func (u *ui) StartGameById(id uint32) {
	u.startGame(u.config.Games[id])
}

func (u *ui) SelectGameById(id uint32) {
	u.selectGame(u.config.Games[id])
}

func (u *ui) startGame(game config.Game) {

	log.Info().Interface("game", game).Msg("Starting game")

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

	log.Info().Interface("game", game).Msg("Selected game")
	if game.Image != "" {
		img := u.imageCache.Images[game.Image]
		if img != nil {
			u.framebuffer.Blit(25, 25, img)
		}
	}

	u.selectedGame = game

	u.DrawHud(data.MatrixState{})
}

func (u *ui) DrawHud(state data.MatrixState) {
	ledState := data.LedStateByMapping(u.selectedGame.Mappings)
	u.LedRegion(ledState.GlobalKeypad.WhiteLeft, "WHITE_LEFT")
	u.LedRegion(ledState.GlobalKeypad.WhiteRight, "WHITE_RIGHT")
	u.LedRegion(ledState.Player1Keypad.Red, "P1_RED")
	u.LedRegion(ledState.Player1Keypad.Yellow, "P1_YELLOW")
	u.LedRegion(ledState.Player1Keypad.Green, "P1_GREEN")
	u.LedRegion(ledState.Player1Keypad.Blue, "P1_BLUE")
	u.LedRegion(ledState.GlobalKeypad.FlipperLeft, "P1_FLIPPER")
	u.LedRegion(ledState.Player2Keypad.Red, "P2_RED")
	u.LedRegion(ledState.Player2Keypad.Yellow, "P2_YELLOW")
	u.LedRegion(ledState.Player2Keypad.Green, "P2_GREEN")
	u.LedRegion(ledState.Player2Keypad.Blue, "P2_BLUE")
	u.LedRegion(ledState.GlobalKeypad.FlipperRight, "P2_FLIPPER")
}

func (u *ui) LedRegion(state bool, key string) {
	reg := u.config.Ui.Regions[key]
	if state {
		u.ShowRegion(reg, "hud_2.jpg", u.hudPosition)
	} else {
		u.ShowRegion(reg, "hud_1.jpg", u.hudPosition)
	}
}

func (u *ui) ShowRegion(region config.Region, filename string, off image.Point) {
	src := u.imageCache.Images[filename]
	dstX := region.X + off.X
	dstY := region.Y + off.Y
	u.framebuffer.BlitSrcWindow(dstX, dstY, region.X, region.Y, region.Width, region.Height, src)
}
