package tools

import (
	"fmt"
	"io/ioutil"
)

func Unexport(pin uint8) {
	d1 := []byte(fmt.Sprintf("%d\n", pin))
	ioutil.WriteFile("/sys/class/gpio/unexport", d1, 0644)
}
