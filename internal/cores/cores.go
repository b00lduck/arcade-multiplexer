package cores

import (
	"os"
	"time"

	"github.com/b00lduck/arcade-multiplexer/internal/hid"
	"github.com/sirupsen/logrus"
)

type Core struct {
	name      string
	enter     []hid.Key
	exit      []hid.Key
	load      []hid.Key
	run       []hid.Key
	bootSleep time.Duration
	speed1    time.Duration
	speed2    time.Duration
}

var Menu = &Core{
	name:   "Menu",
	enter:  []hid.Key{},
	exit:   []hid.Key{},
	load:   []hid.Key{},
	run:    []hid.Key{},
	speed1: 1 * time.Millisecond,
	speed2: 10 * time.Millisecond,
}

var Amiga = &Core{
	name: "Amiga",
	enter: []hid.Key{
		hid.KEY_DOWN,
		hid.KEY_DOWN,
		hid.KEY_ENTER},
	exit: []hid.Key{
		hid.KEY_F12,
		hid.KEY_RIGHT,
		hid.KEY_RIGHT,
		hid.KEY_ENTER,
		hid.KEY_ENTER,
		hid.KEY_HOME},
	load: []hid.Key{
		hid.KEY_F12,
		hid.KEY_ENTER,
		hid.KEY_HOME,
		hid.KEY_DOWN,
		hid.KEY_DOWN,
		hid.KEY_DOWN,
		hid.KEY_DOWN,
		hid.KEY_ENTER},
	run: []hid.Key{
		hid.KEY_WAIT,
		hid.KEY_RIGHT,
		hid.KEY_RIGHT,
		hid.KEY_DOWN,
		hid.KEY_DOWN,
		hid.KEY_ENTER,
		hid.KEY_UP,
		hid.KEY_ENTER},
	bootSleep: 11000 * time.Millisecond,
	speed1:    15 * time.Millisecond,
	speed2:    25 * time.Millisecond,
}

var C64 = &Core{
	name: "C64",
	enter: []hid.Key{
		hid.KEY_DOWN,
		hid.KEY_ENTER},
	exit: []hid.Key{
		hid.KEY_F12,
		hid.KEY_RIGHT,
		hid.KEY_ENTER,
		hid.KEY_ENTER,
		hid.KEY_HOME},
	load: []hid.Key{
		hid.KEY_F12,
		hid.KEY_DOWN,
		hid.KEY_ENTER,
		hid.KEY_HOME},
	run: []hid.Key{
		hid.KEY_WAIT,
		hid.KEY_R,
		hid.KEY_U,
		hid.KEY_N,
		hid.KEY_ENTER},
	bootSleep: 3500 * time.Millisecond,
	speed1:    25 * time.Millisecond,
	speed2:    40 * time.Millisecond,
}

func ChangeCore(oldCore, newCore *Core) {

	if newCore == nil {
		return
	}

	if oldCore == nil {
		oldCore = Menu
	}

	logrus.WithField("oldCore", oldCore).
		WithField("newCore", newCore).
		Info("Changing core")

	file, err := os.OpenFile("/dev/hidg0", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		logrus.Error("Error opening /dev/hidg0")
		return
	}
	defer file.Close()
	if oldCore != nil {
		hid.WriteSequence(file, oldCore.exit, oldCore.speed1, oldCore.speed2)
	}
	hid.WriteSequence(file, newCore.enter, oldCore.speed1, oldCore.speed2)

	time.Sleep(newCore.bootSleep)
}

func LoadGame(core *Core, index int) {

	logrus.WithField("index", index).
		Info("Load game")

	file, err := os.OpenFile("/dev/hidg0", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		logrus.Error("Error opening /dev/hidg0")
		return
	}
	defer file.Close()
	hid.WriteSequence(file, core.load, core.speed1, core.speed2)

	for i := 0; i < index; i++ {
		hid.WriteSequence(file, []hid.Key{hid.KEY_DOWN}, core.speed1, core.speed2)
	}
	hid.WriteSequence(file, []hid.Key{hid.KEY_ENTER}, core.speed1, core.speed2)

	hid.WriteSequence(file, core.run, core.speed1, core.speed2)
}