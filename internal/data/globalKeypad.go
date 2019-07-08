package data

type GlobalKeypad struct {
	WhiteLeft    bool
	WhiteRight   bool
	FlipperLeft  bool
	FlipperRight bool
}

func (k *GlobalKeypad) String() string {
	ret := ""
	if k.WhiteLeft {
		ret += "WL "
	} else {
		ret += "   "
	}
	if k.WhiteRight {
		ret += "WR "
	} else {
		ret += "   "
	}
	if k.FlipperLeft {
		ret += "FL "
	} else {
		ret += "   "
	}
	if k.FlipperRight {
		ret += "FR "
	} else {
		ret += "   "
	}
	return ret
}
