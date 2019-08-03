package data

type PlayerKeypad struct {
	Red    bool
	Yellow bool
	Green  bool
	Blue   bool
}

func (k *PlayerKeypad) String() string {
	ret := ""
	if k.Red {
		ret += "R "
	} else {
		ret += "  "
	}
	if k.Yellow {
		ret += "Y "
	} else {
		ret += "  "
	}
	if k.Green {
		ret += "G "
	} else {
		ret += "  "
	}
	if k.Blue {
		ret += "B"
	} else {
		ret += " "
	}
	return ret
}

func (k *PlayerKeypad) Changed(old PlayerKeypad) bool {
	return old.Red != k.Red ||
		old.Yellow != k.Yellow ||
		old.Green != k.Green ||
		old.Blue != k.Blue
}
