package inputProcessor

import (
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
	hid      Hid
	mappings []config.Mapping
	mutex    *sync.Mutex

	buttonStates []data.ButtonState
	pressedKeys  []string
}

type Hid interface {
	SetKeys([]string)
}

func NewInputProcessor(mist Mist, hid Hid) *inputProcessor {
	return &inputProcessor{
		mist:     mist,
		hid:      hid,
		mappings: []config.Mapping{},
		mutex:    &sync.Mutex{}}
}

func (i *inputProcessor) SetMappings(mappings []config.Mapping) {
	i.mutex.Lock()
	i.mappings = mappings
	i.mutex.Unlock()
}

func (i *inputProcessor) ProcessMatrix(ms data.MatrixState) {
	//fmt.Println(ms.String())
	//fmt.Printf("\033[7A")

	i.mutex.Lock()

	i.buttonStates = make([]data.ButtonState, 4)
	i.pressedKeys = make([]string, 0)

	for _, v := range i.mappings {
		switch v.Input {

		case "P1_JOY":
			i.OutputJoystick(&ms.Player1Joystick, v.Output)
		case "P2_JOY":
			i.OutputJoystick(&ms.Player2Joystick, v.Output)

		case "P1_RED":
			i.RegisterButtonOrKey(ms.Player1Keypad.Red, v.Output, v.Autofire)
		case "P1_YELLOW":
			i.RegisterButtonOrKey(ms.Player1Keypad.Yellow, v.Output, v.Autofire)
		case "P1_BLUE":
			i.RegisterButtonOrKey(ms.Player1Keypad.Blue, v.Output, v.Autofire)
		case "P1_GREEN":
			i.RegisterButtonOrKey(ms.Player1Keypad.Green, v.Output, v.Autofire)

		case "P2_RED":
			i.RegisterButtonOrKey(ms.Player2Keypad.Red, v.Output, v.Autofire)
		case "P2_YELLOW":
			i.RegisterButtonOrKey(ms.Player2Keypad.Yellow, v.Output, v.Autofire)
		case "P2_BLUE":
			i.RegisterButtonOrKey(ms.Player2Keypad.Blue, v.Output, v.Autofire)
		case "P2_GREEN":
			i.RegisterButtonOrKey(ms.Player2Keypad.Green, v.Output, v.Autofire)

		case "WHITE_LEFT":
			i.RegisterButtonOrKey(ms.GlobalKeypad.WhiteLeft, v.Output, v.Autofire)
		case "WHITE_RIGHT":
			i.RegisterButtonOrKey(ms.GlobalKeypad.WhiteRight, v.Output, v.Autofire)
		}
	}

	i.mist.SetJoystickButton(0, i.buttonStates[0])
	i.mist.SetJoystickButton(1, i.buttonStates[1])
	i.mist.SetJoystickButton(2, i.buttonStates[2])
	i.mist.SetJoystickButton(3, i.buttonStates[3])

	i.hid.SetKeys(i.pressedKeys)

	i.mutex.Unlock()
}

func (i *inputProcessor) RegisterButtonOrKey(in bool, out string, autofire bool) {
	if !in {
		return
	}

	switch out {
	case "JOY1_BUTTON1":
		i.registerButton(0, autofire)
	case "JOY1_BUTTON2":
		i.registerButton(1, autofire)
	case "JOY2_BUTTON1":
		i.registerButton(2, autofire)
	case "JOY2_BUTTON2":
		i.registerButton(3, autofire)
	default:
		i.registerKey(out)
	}
}

func (i *inputProcessor) registerButton(index int, autofire bool) {
	i.buttonStates[index] = data.ButtonState{
		State:    true,
		Autofire: autofire}
}

func (i *inputProcessor) registerKey(out string) {
	i.pressedKeys = append(i.pressedKeys, out)
}

func (i *inputProcessor) OutputJoystick(in *data.Joystick, out string) {
	switch out {
	case "JOY1_AXES":
		i.mist.SetJoystick1(in)
	case "JOY2_AXES":
		i.mist.SetJoystick2(in)
	}
}
