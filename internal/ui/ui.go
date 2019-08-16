package ui

import (
	"github.com/b00lduck/arcade-multiplexer/internal/config"
	"github.com/b00lduck/arcade-multiplexer/internal/data"
	"github.com/tarent/logrus"
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
	LoadGame(*config.Game, *config.Core)
}

type ui struct {
	display        Display
	panel          Panel
	config         *config.Config
	oldGame        config.Game
	inputProcessor InputProcessor
	mistControl    MistControl
}

func NewUi(c *config.Config, display Display, panel Panel, ip InputProcessor, mistControl MistControl) *ui {
	return &ui{
		display:        display,
		panel:          panel,
		config:         c,
		oldGame:        config.Game{},
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
	u.oldGame = game

	logrus.WithField("game", game).Info("Starting game")

	u.panel.SetLeds(data.LedStateByMapping(game.Mappings))
	if game.Image != "" {
		u.display.ShowImage(game.Image)
	}
	u.inputProcessor.SetMappings(game.Mappings)

	newCore := u.config.GetCoreByName(game.Core)
	u.mistControl.ChangeCore(newCore)
	u.mistControl.LoadGame(&game, newCore)

}

func (u *ui) selectGame(game config.Game) {
	logrus.WithField("game", game).Info("Selected game")
	if game.Image != "" {
		u.display.ShowImage(game.Image)
	}
}
