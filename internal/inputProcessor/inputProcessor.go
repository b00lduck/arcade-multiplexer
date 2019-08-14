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
	SetJoystickButton(uint8, data.ButtonState)
}

type inputProcessor struct {
	mist     Mist
	mappings []config.Mapping
	mutex    *sync.Mutex

	buttonStates []data.ButtonState
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

	i.buttonStates = make([]data.ButtonState, 4)

	for _, v := range i.mappings {
		switch v.Input {

		case "P1_JOY":
			i.OutputJoystick(&ms.Player1Joystick, v.Output)
		case "P2_JOY":
			i.OutputJoystick(&ms.Player2Joystick, v.Output)

		case "P1_RED":
			i.RegisterButton(ms.Player1Keypad.Red, v.Output, v.Autofire)
		case "P1_YELLOW":
			i.RegisterButton(ms.Player1Keypad.Yellow, v.Output, v.Autofire)
		case "P1_BLUE":
			i.RegisterButton(ms.Player1Keypad.Blue, v.Output, v.Autofire)
		case "P1_GREEN":
			i.RegisterButton(ms.Player1Keypad.Green, v.Output, v.Autofire)

		case "P2_RED":
			i.RegisterButton(ms.Player2Keypad.Red, v.Output, v.Autofire)
		case "P2_YELLOW":
			i.RegisterButton(ms.Player2Keypad.Yellow, v.Output, v.Autofire)
		case "P2_BLUE":
			i.RegisterButton(ms.Player2Keypad.Blue, v.Output, v.Autofire)
		case "P2_GREEN":
			i.RegisterButton(ms.Player2Keypad.Green, v.Output, v.Autofire)

		case "WHITE_LEFT":
			i.RegisterButton(ms.GlobalKeypad.WhiteLeft, v.Output, v.Autofire)
		case "WHITE_RIGHT":
			i.RegisterButton(ms.GlobalKeypad.WhiteRight, v.Output, v.Autofire)
		}
	}

	i.mist.SetJoystickButton(0, i.buttonStates[0])
	i.mist.SetJoystickButton(1, i.buttonStates[1])
	i.mist.SetJoystickButton(2, i.buttonStates[2])
	i.mist.SetJoystickButton(3, i.buttonStates[3])

	i.mutex.Unlock()
}

func (i *inputProcessor) RegisterButton(in bool, out string, autofire bool) {

	if !in {
		return
	}

	index := -1
	switch out {
	case "JOY1_BUTTON1":
		index = 0
	case "JOY1_BUTTON2":
		index = 1
	case "JOY2_BUTTON1":
		index = 2
	case "JOY2_BUTTON2":
		index = 3
	default:
		return
	}

	i.buttonStates[index] = data.ButtonState{
		State:    in,
		Autofire: autofire}

}

func (i *inputProcessor) OutputJoystick(in *data.Joystick, out string) {
	switch out {
	case "JOY1_AXES":
		i.mist.SetJoystick1(in)
	case "JOY2_AXES":
		i.mist.SetJoystick2(in)
	}
}

/*
func (i *inputProcessor) OutputButton(in bool, out string, autofire bool) {

	switch out {
	case "JOY1_BUTTON1":
		i.mist.SetJoystick1Button1(in, autofire)
	case "JOY1_BUTTON2":
		i.mist.SetJoystick1Button2(in, autofire)
	case "JOY2_BUTTON1":
		i.mist.SetJoystick2Button1(in, autofire)
	case "JOY2_BUTTON2":
		i.mist.SetJoystick2Button2(in, autofire)

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
*/
