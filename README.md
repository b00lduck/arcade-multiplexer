# arcade-multiplexer

## Installation

### Install raspbian

```
sudo dd bs=4M if=2019-07-10-raspbian-buster-lite.img of=/dev/sdX conv=fsync
```

### SSH
Copy SSH key to raspbian 'root' user. The default user is 'pi' with password 'raspbian'.

### Install software

```
apt-get update
apt-get install i2c-tools
```

### Configuration

Enable I²C with
```
raspi-config
```

# Deprecated/Legacy information

```
# /boot/config.txt

disable_splash=1
dtoverlay=pi3-disable-bt
dtoverlay=sdtweak,overclock_50=100
boot_delay=0

gpu_mem=64
initramfs initramfs-linux.img followkernel
dtoverlay=dwc2
dtparam=i2c_arm=on,i2c_arm_baudrate=400000
dtoverlay=rpi-display,speed=32000000,rotate=270
```
