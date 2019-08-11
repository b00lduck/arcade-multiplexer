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
dtparam=i2c_arm=on,i2c_arm_baudrate=400000
dtoverlay=tft35a:rotate=90,speed=62000000
gpu_mem=16
```

# Notes
https://github.com/goodtft/LCD-show.git 
scp tft35a-overlay.dtb root@arcade:/boot/overlays/tft35a.dtbo