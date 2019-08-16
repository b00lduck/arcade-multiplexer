package mist

import (
	"os"
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/config"
	"github.com/sirupsen/logrus"
)

type Hid interface {
	WriteSequence([]string, uint64, uint64) error
}

type MistDigital interface {
	SetResetButton(bool)
}

type mistControl struct {
	hid         Hid
	mistDigital MistDigital
}

func NewMistControl(h Hid, mistDigital MistDigital) *mistControl {
	return &mistControl{
		hid:         h,
		mistDigital: mistDigital}
}

func (m *mistControl) ChangeCore(newCore *config.Core) {

	if newCore == nil {
		return
	}

	m.mistDigital.SetResetButton(true)
	time.Sleep(50 * time.Millisecond)
	m.mistDigital.SetResetButton(false)

	time.Sleep(1000 * time.Millisecond)

	logrus.WithField("newCore", newCore).Info("Changing core")

	m.hid.WriteSequence(newCore.Enter, 1, 10)

	time.Sleep(time.Duration(newCore.BootSleep) * time.Millisecond)
}

func (m *mistControl) LoadGame(game *config.Game, core *config.Core) {

	logrus.WithField("name", game.Name).
		WithField("core", game.Core).
		Info("Loading game")

	file, err := os.OpenFile("/dev/hidg0", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		logrus.WithError(err).Fatal("Error opening /dev/hidg0")
	}
	defer file.Close()
	m.hid.WriteSequence(core.Load, core.Speed1, core.Speed2)

	for i := 0; i < game.Index; i++ {
		m.hid.WriteSequence([]string{"KEY_DOWN"}, core.Speed1, core.Speed2)
	}
	m.hid.WriteSequence([]string{"KEY_ENTER"}, core.Speed1, core.Speed2)

	m.hid.WriteSequence(core.Run, core.Speed1, core.Speed2)

}
