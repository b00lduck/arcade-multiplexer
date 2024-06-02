package hid

import (
	"os"
	"strconv"
	"strings"
	"time"
	"errors"

	"github.com/rs/zerolog/log"
)

/*
  "KEY_WAIT": 0xff

// Modifier

*/

var keyMap = map[string]byte{

	"KEY_MOD_LCTRL":  0x01,
	"KEY_MOD_LSHIFT": 0x02,
	"KEY_MOD_LALT":   0x04,
	"KEY_MOD_RCTRL":  0x10,
	"KEY_MOD_RSHIFT": 0x20,
	"KEY_MOD_RALT":   0x40,

	"KEY_UP":    0x52,
	"KEY_DOWN":  0x51,
	"KEY_RIGHT": 0x4f,
	"KEY_LEFT":  0x50,
	"KEY_F12":   0x45,
	"KEY_ENTER": 0x28,
	"KEY_ESC":   0x29,
	"KEY_SPACE": 0x2c,
	"KEY_F1":    0x3a,
	"KEY_F2":    0x3b,
	"KEY_F3":    0x3c,
	"KEY_F4":    0x3d,
	"KEY_HOME":  0x4a,

	"KEY_A": 0x04,
	"KEY_B": 0x05,
	"KEY_C": 0x06,
	"KEY_D": 0x07,
	"KEY_E": 0x08,
	"KEY_F": 0x09,
	"KEY_G": 0x0a,
	"KEY_H": 0x0b,
	"KEY_I": 0x0c,
	"KEY_J": 0x0d,
	"KEY_K": 0x0e,
	"KEY_L": 0x0f,
	"KEY_M": 0x10,
	"KEY_N": 0x11,
	"KEY_O": 0x12,
	"KEY_P": 0x13,
	"KEY_Q": 0x14,
	"KEY_R": 0x15,
	"KEY_S": 0x16,
	"KEY_T": 0x17,
	"KEY_U": 0x18,
	"KEY_V": 0x19,
	"KEY_W": 0x1a,
	"KEY_X": 0x1b,
	"KEY_Y": 0x1c,
	"KEY_Z": 0x1d}

type hid struct {
	file *os.File
}

func NewHid() *hid {

	file, err := os.OpenFile("/dev/hidg0", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening /dev/hidg0")
	}
	return &hid{
		file: file}
}

func (h *hid) Close() {
	h.file.Close()

}

func (h *hid) SetKeys(keys []string) {

	out := make([]byte, 8)

	modPtr := 0
	keyPtr := 2
	for _, v := range keys {
		if strings.HasPrefix(v, "KEY_MOD_") {
			out[modPtr] += keyMap[v]
		} else if strings.HasPrefix(v, "KEY_") {
			if keyPtr < 8 {
				out[keyPtr] = keyMap[v]
				keyPtr++
			}
		}
	}

	h.sendRaw(out)
}

func (h *hid) sendRaw(b []byte) error {

	// retry 100 times if error
	for i := 0; i < 100; i++ {
		err := h.sendRawOnce(b)
		if err == nil {
			time.Sleep(100 * time.Millisecond)
			return nil
		}
	}

	err := errors.New("Error writing to hid, giving up")
	log.Error().Err(err).Msg("Error writing to hid, giving up")

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

func (h *hid) WriteSequence(seq []string, speed1, speed2 uint64) error {
	for _, v := range seq {
		log.Info().Str("key", v).Msg("PRESS")

		if strings.HasPrefix(v, "KEY_WAIT_") {
			foo, err := strconv.Atoi(strings.Split(v, "KEY_WAIT_")[1])
			if err != nil {
				log.Info().Err(err).Msg("Wait failed")
			}
			log.Info().Int("ms", foo).Msg("sleeping")
			time.Sleep(time.Duration(foo) * time.Millisecond)
		} else {

			key := keyMap[v]

			h.sendRaw([]byte{0, 0, byte(key), 0, 0, 0, 0, 0})
			time.Sleep(time.Duration(speed1) * time.Millisecond)

			h.sendRaw([]byte{0, 0, 0, 0, 0, 0, 0, 0})
			time.Sleep(time.Duration(speed2) * time.Millisecond)
		}

	}
	return nil
}
