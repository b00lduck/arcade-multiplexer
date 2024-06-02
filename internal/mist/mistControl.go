package mist

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"arcade-multiplexer/internal/config"
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

	mistDigital.SetResetButton(true)
	time.Sleep(50 * time.Millisecond)
	mistDigital.SetResetButton(false)

	return &mistControl{
		hid:         h,
		mistDigital: mistDigital}
}

func (m *mistControl) ChangeCore(newCore *config.Core) {

	if newCore == nil {
		log.Warn().Msg("new core is nil")
		return
	}

	m.mistDigital.SetResetButton(true)
	time.Sleep(50 * time.Millisecond)
	m.mistDigital.SetResetButton(false)
	time.Sleep(2000 * time.Millisecond)

	log.Info().Interface("newCore", newCore).Msg("Changing core")

	m.hid.WriteSequence(newCore.Enter, 1, 10)

	time.Sleep(time.Duration(newCore.BootSleep) * time.Millisecond)
}

func (m *mistControl) LoadGame(game *config.Game, core *config.Core, sameCore bool) {

	log.Info().Str("name", game.Name).
		Str("core", game.Core).
		Msg("Loading game")

	file, err := os.OpenFile("/dev/hidg0", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening /dev/hidg0")
	}
	defer file.Close()

	if sameCore {
		log.Info().Msg("same core")
		m.hid.WriteSequence(core.LoadSameCore, core.Speed1, core.Speed2)
	} else {
		log.Info().Msg("other core")
		m.hid.WriteSequence(core.Load, core.Speed1, core.Speed2)
	}

	for i := 0; i < game.Index; i++ {
		m.hid.WriteSequence([]string{"KEY_DOWN"}, core.Speed1, core.Speed2)
	}
	m.hid.WriteSequence([]string{"KEY_ENTER"}, core.Speed1, core.Speed2)

	if game.Disks == 2 {
		m.hid.WriteSequence([]string{"KEY_ENTER"}, core.Speed1, core.Speed2)
		m.hid.WriteSequence([]string{"KEY_DOWN"}, core.Speed1, core.Speed2)
		m.hid.WriteSequence([]string{"KEY_ENTER"}, core.Speed1, core.Speed2)
	}

	m.hid.WriteSequence(core.Run, core.Speed1, core.Speed2)

}
