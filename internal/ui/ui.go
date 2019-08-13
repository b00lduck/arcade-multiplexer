package ui

import (
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/config"
	"github.com/b00lduck/arcade-multiplexer/internal/cores"
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

type Mist interface {
	SetResetButton(bool)
}

type ui struct {
	display        Display
	panel          Panel
	config         *config.Config
	oldGame        config.Game
	inputProcessor InputProcessor
	mist           Mist
}

func NewUi(c *config.Config, display Display, panel Panel, ip InputProcessor, mist Mist) *ui {
	return &ui{
		display:        display,
		panel:          panel,
		config:         c,
		oldGame:        config.Game{},
		inputProcessor: ip,
		mist:           mist}
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

	u.mist.SetResetButton(true)
	time.Sleep(50 * time.Millisecond)
	u.mist.SetResetButton(false)

	time.Sleep(1000 * time.Millisecond)

	newCore := cores.CoreFromString(game.Core)
	cores.ChangeCore(cores.Menu, newCore)
	cores.LoadGame(&game, newCore)

}

func (u *ui) selectGame(game config.Game) {
	logrus.WithField("game", game).Info("Selected game")
	if game.Image != "" {
		u.display.ShowImage(game.Image)
	}
}
