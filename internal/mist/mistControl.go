package mist

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"arcade-multiplexer/internal/config"
)

type Hid interface {
	PressKey(string, uint64, uint64) error
}

type MistDigital interface {
	SetResetButton(bool)
	PressButton(uint8, uint8)
}

type mistControl struct {
	hid         Hid
	mistDigital MistDigital
	conf        config.Mist
}

func NewMistControl(h Hid, mistDigital MistDigital, conf config.Mist) *mistControl {

	mistDigital.SetResetButton(true)
	log.Info().Uint64("ms", conf.ResetDuration).Msg("RST MiSt")
	time.Sleep(time.Duration(conf.ResetDuration) * time.Millisecond)
	mistDigital.SetResetButton(false)

	return &mistControl{
		hid:         h,
		mistDigital: mistDigital,
		conf:        conf,
	}
}

func (m *mistControl) ExitCore(core *config.Core) {

	if core == nil {
		return
	}

	log.Info().Interface("oldCore", core).Msg("Exiting core")
	m.SendHybridSequence(core.Exit, 10, 10)
}

func (m *mistControl) ChangeCore(newCore *config.Core) {

	if newCore == nil {
		log.Warn().Msg("new core is nil")
		return
	}

	log.Info().Interface("newCore", newCore).Msg("Changing core")

	m.SendHybridSequence(newCore.Enter, 10, 10)

	log.Info().Uint64("ms", newCore.BootSleep).Msg("SLEEP")
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
		m.SendHybridSequence(core.LoadSameCore, core.Speed1, core.Speed2)
	} else {
		log.Info().Msg("other core")
		m.SendHybridSequence(core.Load, core.Speed1, core.Speed2)
	}

	for i := 0; i < game.Index; i++ {
		m.SendHybridSequence([]string{"KEY_DOWN"}, core.Speed1, core.Speed2)
	}
	m.SendHybridSequence([]string{"KEY_ENTER"}, core.Speed1, core.Speed2)

	if game.Disks == 2 {
		m.SendHybridSequence([]string{"KEY_ENTER"}, core.Speed1, core.Speed2)
		m.SendHybridSequence([]string{"KEY_DOWN"}, core.Speed1, core.Speed2)
		m.SendHybridSequence([]string{"KEY_ENTER"}, core.Speed1, core.Speed2)
	}

	m.SendHybridSequence(core.Run, core.Speed1, core.Speed2)

	m.SendHybridSequence(game.InitSequence, core.Speed1, core.Speed2)

}

// Send a sequence consisting of HID commands and Digital Input (Joystick, Buttons etc)
// to the MiST
func (m *mistControl) SendHybridSequence(seq []string, speed1, speed2 uint64) {
	for _, v := range seq {

		// check if string v starts with "KEY_"
		if strings.HasPrefix(v, "SLEEP_") {
			foo, err := strconv.Atoi(strings.Split(v, "SLEEP_")[1])
			if err != nil {
				log.Info().Err(err).Msg("Parsing of sleep command failed")
			}
			log.Info().Int("ms", foo).Msg("SLEEP")
			time.Sleep(time.Duration(foo) * time.Millisecond)
		} else if strings.HasPrefix(v, "KEY_") {
			m.hid.PressKey(v, speed1, speed2)
		} else {
			switch v {
			case "JOY1_BUTTON1":
				m.mistDigital.PressButton(1, 1)
			case "JOY1_BUTTON2":
				m.mistDigital.PressButton(1, 2)
			case "JOY2_BUTTON1":
				m.mistDigital.PressButton(2, 1)
			case "JOY2_BUTTON2":
				m.mistDigital.PressButton(2, 2)
			default:
				log.Warn().Str("key", v).Msg("Unknown key command")
			}
		}

	}

}
