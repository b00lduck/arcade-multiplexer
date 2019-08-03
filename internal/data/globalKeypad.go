package data

type GlobalKeypad struct {
	WhiteLeft    bool
	WhiteRight   bool
	FlipperLeft  bool
	FlipperRight bool
}

func (k *GlobalKeypad) Changed(old GlobalKeypad) bool {
	return old.WhiteLeft != k.WhiteLeft ||
		old.WhiteRight != k.WhiteRight ||
		old.FlipperLeft != k.FlipperLeft ||
		old.FlipperRight != k.FlipperRight
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
