package main

import (
"github.com/kidoman/embd"
_ "github.com/kidoman/embd/host/all"
"io/ioutil"
)


func main() {

    d1 := []byte("17\n")
    ioutil.WriteFile("/sys/class/gpio/unexport", d1, 0644)

    if err := embd.InitGPIO(); err != nil {
panic(err)
    }
    defer embd.CloseGPIO()

    led, err := embd.NewDigitalPin(17)
    if err != nil {
panic(err)
    }
    defer led.Close()
    if err := led.SetDirection(embd.In); err != nil {
panic(err)
    }

}

