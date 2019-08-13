package inputProcessor

import (
	"fmt"
	"sync"

	"github.com/b00lduck/arcade-multiplexer/internal/config"
	"github.com/b00lduck/arcade-multiplexer/internal/data"
)

type Mist interface {
	SetJoystick1(in *data.Joystick)
	SetJoystick2(in *data.Joystick)
	SetJoystick1Button1(bool)
	SetJoystick1Button2(bool)
	SetJoystick2Button1(bool)
	SetJoystick2Button2(bool)
}

type inputProcessor struct {
	mist     Mist
	mappings []config.Mapping
	mutex    *sync.Mutex
}

func NewInputProcessor(mist Mist) *inputProcessor {
	return &inputProcessor{
		mist:     mist,
		mappings: []config.Mapping{},
		mutex:    &sync.Mutex{}}
}

func (i *inputProcessor) SetMappings(mappings []config.Mapping) {
	i.mutex.Lock()
	i.mappings = mappings
	i.mutex.Unlock()
}

func (i *inputProcessor) ProcessMatrix(ms data.MatrixState) {
	fmt.Println(ms.String())
	fmt.Printf("\033[7A")

	i.mutex.Lock()
	for _, v := range i.mappings {
		switch v.Input {
		case "P1_JOY":
			i.OutputJoystick(&ms.Player1Joystick, v.Output)
		case "P1_RED":
			i.OutputButton(ms.Player1Keypad.Red, v.Output)
		case "P1_YELLOW":
			i.OutputButton(ms.Player1Keypad.Yellow, v.Output)
		case "P1_BLUE":
			i.OutputButton(ms.Player1Keypad.Blue, v.Output)
		case "P1_GREEN":
			i.OutputButton(ms.Player1Keypad.Green, v.Output)

		case "P2_JOY":
			i.OutputJoystick(&ms.Player2Joystick, v.Output)
		case "P2_RED":
			i.OutputButton(ms.Player2Keypad.Red, v.Output)
		case "P2_YELLOW":
			i.OutputButton(ms.Player2Keypad.Yellow, v.Output)
		case "P2_BLUE":
			i.OutputButton(ms.Player2Keypad.Blue, v.Output)
		case "P2_GREEN":
			i.OutputButton(ms.Player2Keypad.Green, v.Output)

		case "WHITE_LEFT":
			i.OutputButton(ms.GlobalKeypad.WhiteLeft, v.Output)
		case "WHITE_RIGHT":
			i.OutputButton(ms.GlobalKeypad.WhiteRight, v.Output)
		}
	}
	i.mutex.Unlock()
}

func (i *inputProcessor) OutputJoystick(in *data.Joystick, out string) {
	switch out {
	case "JOY1_AXES":
		i.mist.SetJoystick1(in)
	case "JOY2_AXES":
		i.mist.SetJoystick2(in)
	}
}

func (i *inputProcessor) OutputButton(in bool, out string) {
	switch out {
	case "JOY1_BUTTON1":
		i.mist.SetJoystick1Button1(in)
	case "JOY1_BUTTON2":
		i.mist.SetJoystick1Button2(in)
	case "JOY2_BUTTON1":
		i.mist.SetJoystick2Button1(in)
	case "JOY2_BUTTON2":
		i.mist.SetJoystick2Button2(in)

	case "JOY1_UP":
		// not implemented
	case "JOY1_DOWN":
		// not implemented
	case "JOY1_LEFT":
		// not implemented
	case "JOY1_RIGHT":
		// not implemented

	case "JOY2_UP":
		// not implemented
	case "JOY2_DOWN":
		// not implemented
	case "JOY2_LEFT":
		// not implemented
	case "JOY2_RIGHT":
		// not implemented

	case "KEY_SPACE":
		// not implemented
	case "KEY_ESCAPE":
		// not implemented
	case "KEY_LSHIFT":
		// not implemented
	case "KEY_RSHIFT":
		// not implemented
	case "KEY_1":
		// not implemented
	case "KEY_2":
		// not implemented
	case "KEY_3":
		// not implemented
	case "KEY_4":
		// not implemented
	}
}
