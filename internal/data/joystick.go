package data

type Joystick struct {
	Up    bool
	Down  bool
	Left  bool
	Right bool
}

func (k *Joystick) String() string {
	ret := ""
	if k.Up {
		ret += "U "
	} else {
		ret += "  "
	}
	if k.Down {
		ret += "D "
	} else {
		ret += "  "
	}
	if k.Left {
		ret += "L "
	} else {
		ret += "  "
	}
	if k.Right {
		ret += "R "
	} else {
		ret += "  "
	}
	return ret
}
