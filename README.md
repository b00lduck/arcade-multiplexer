# arcade-multiplexer

## RPi configuration

```
# /etc/modules-load.d/raspberrypi.conf

snd-bcm2835
dwc2
libcomposite
i2c-dev
i2c-bcm2708
```

```
# /boot/config.txt

disable_splash=1
dtoverlay=pi3-disable-bt
dtoverlay=sdtweak,overclock_50=100
boot_delay=0

gpu_mem=64
initramfs initramfs-linux.img followkernel
dtoverlay=dwc2
dtparam=i2c_arm=on

dtoverlay=rpi-display,speed=32000000,rotate=270
```

```
# installed packages

i2c-tools
tslib
xorg-xinput



```