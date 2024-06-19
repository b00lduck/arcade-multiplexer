package hid

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"arcade-multiplexer/internal/config"
)

type hid struct {
	file   *os.File
	config *config.Config
}

func NewHid(config *config.Config) *hid {

	file, err := os.OpenFile("/dev/hidg0", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening /dev/hidg0")
	}
	return &hid{
		file:   file,
		config: config}
}

func (h *hid) Close() {
	h.file.Close()

}

func (h *hid) getScanCode(key string) byte {
	out := h.config.Outputs[key]
	if out.ScanCode == 0 {
		log.Fatal().Str("key", key).Msg("Unknown key")
	}
	return byte(h.config.Outputs[key].ScanCode)
}

func (h *hid) SetKeys(keys []string) {

	out := make([]byte, 8)

	modPtr := 0
	keyPtr := 2
	for _, v := range keys {
		if strings.HasPrefix(v, "KEY_MOD_") {
			out[modPtr] += h.getScanCode(v)
		} else if strings.HasPrefix(v, "KEY_") {
			if keyPtr < 8 {
				out[keyPtr] = h.getScanCode(v)
				keyPtr++
			}
		}
	}

	h.sendRaw(out)
}

func (h *hid) sendRaw(b []byte) error {

	// retry 1000 times if error
	for i := 0; i < 1000; i++ {
		err := h.sendRawOnce(b)
		if err == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}

	err := errors.New("error writing to hid, giving up")
	log.Error().Err(err)

	return err
}

func (h *hid) sendRawOnce(b []byte) error {

	_, err := h.file.Write(b)
	if err != nil {
		log.Error().Err(err).Msg("Error writing to hid")
		return err
	}
	return nil
}

func (h *hid) PressKey(v string, speed1, speed2 uint64) error {

	log.Info().Str("key", v).Uint64("holdMs", speed1).Uint64("pauseMs", speed2).Msg("PRESS")

	key := h.getScanCode(v)

	h.sendRaw([]byte{0, 0, byte(key), 0, 0, 0, 0, 0})
	time.Sleep(time.Duration(speed1) * time.Millisecond)

	h.sendRaw([]byte{0, 0, 0, 0, 0, 0, 0, 0})
	time.Sleep(time.Duration(speed2) * time.Millisecond)

	return nil
}
