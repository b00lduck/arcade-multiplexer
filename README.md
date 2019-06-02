# arcade-multiplexer

```
# /etc/modules-load.d/
raspberrypi.conf

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
toverlay=sdtweak,overclock_50=100
boot_delay=0

gpu_mem=64
initramfs initramfs-linux.img followkernel
dtoverlay=dwc2
dtparam=i2c_arm=on
```

```
# installed packages
i2c-tools
```