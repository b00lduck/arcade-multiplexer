package hid

import (
	"os"
	"time"

	"github.com/tarent/logrus"
)

const KEY_WAIT = 0xff

// Modifier
const KEY_MOD_LCTRL = 0x01
const KEY_MOD_LSHIFT = 0x02
const KEY_MOD_LALT = 0x04
const KEY_MOD_RCTRL = 0x10
const KEY_MOD_RSHIFT = 0x20
const KEY_MOD_RALT = 0x40

// Regular keys
const KEY_NONE = 0
const KEY_UP = 0x52
const KEY_DOWN = 0x51
const KEY_RIGHT = 0x4f
const KEY_LEFT = 0x50
const KEY_F12 = 0x45
const KEY_ENTER = 0x28
const KEY_ESC = 0x29
const KEY_SPACE = 0x2c
const KEY_F1 = 0x3a
const KEY_F2 = 0x3b
const KEY_HOME = 0x4a

const KEY_A = 0x04
const KEY_B = 0x05
const KEY_C = 0x06
const KEY_D = 0x07
const KEY_E = 0x08
const KEY_F = 0x09
const KEY_G = 0x0a
const KEY_H = 0x0b
const KEY_I = 0x0c
const KEY_J = 0x0d
const KEY_K = 0x0e
const KEY_L = 0x0f
const KEY_M = 0x10
const KEY_N = 0x11
const KEY_O = 0x12
const KEY_P = 0x13
const KEY_Q = 0x14
const KEY_R = 0x15
const KEY_S = 0x16
const KEY_T = 0x17
const KEY_U = 0x18
const KEY_V = 0x19
const KEY_W = 0x1a
const KEY_X = 0x1b
const KEY_Y = 0x1c
const KEY_Z = 0x1d

type Key byte

func WriteSequence(file *os.File, seq []Key, speed1, speed2 time.Duration) error {
	for _, k := range seq {
		logrus.WithField("key", k).Info("PRESS")

		if k == KEY_WAIT {
			time.Sleep(250 * time.Millisecond)
		} else {
			_, err := file.Write([]byte{0, 0, byte(k), 0, 0, 0, 0, 0})
			if err != nil {
				logrus.WithError(err).Error("Error writing keydown")
				return err
			}
			time.Sleep(speed1)
			_, err = file.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0})
			if err != nil {
				logrus.WithError(err).Error("Error writing keyup")
				return err
			}
			time.Sleep(speed2)
		}

	}
	return nil
}
