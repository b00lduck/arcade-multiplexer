package ui

import (
	"github.com/rs/zerolog/log"

	"arcade-multiplexer/internal/config"
	"arcade-multiplexer/internal/data"
	"arcade-multiplexer/internal/framebuffer"
)

type Display interface {
	ShowImage(filename string)
}

type Panel interface {
	SetLeds(data.LedState)
}

type InputProcessor interface {
	SetMappings([]config.Mapping)
}

type MistControl interface {
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
}

func NewUi(c *config.Config, framebuffer *framebuffer.DisplayFramebuffer, panel Panel, ip InputProcessor, mistControl MistControl) *ui {
	return &ui{
		framebuffer: framebuffer,
		panel:       panel,
		config:      c,
		oldGame: config.Game{
			Core: "none",
		},
		inputProcessor: ip,
		mistControl:    mistControl}
}

func (u *ui) StartGameById(id uint32) {
	u.startGame(u.config.Games[id])
}

func (u *ui) SelectGameById(id uint32) {
	u.selectGame(u.config.Games[id])
}

func (u *ui) startGame(game config.Game) {

	//oldCore := cores.CoreFromString(u.oldGame.Core)

	log.Info().Interface("game", game).Msg("Starting game")

	u.panel.SetLeds(data.LedStateByMapping(game.Mappings))
	if game.Image != "" {
		u.framebuffer.ShowImage(game.Image)
	}
	u.inputProcessor.SetMappings(game.Mappings)

	newCore := u.config.GetCoreByName(game.Core)

	log.Info().Str("oldCore", u.oldGame.Core).Str("newCore", game.Core).Msg("cores")

	if u.oldGame.Core != game.Core {
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
		u.framebuffer.ShowImage(game.Image)
	}
}
