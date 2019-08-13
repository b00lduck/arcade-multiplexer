# arcade-multiplexer

## Installation

### Install raspbian

```
sudo dd bs=4M if=2019-07-10-raspbian-buster-lite.img of=/dev/sdX conv=fsync
```

### SSH
Copy SSH key to raspbian 'root' user. The default user is 'pi' with password 'raspbian'.

### Install optional software

```
apt-get update
apt-get install i2c-tools
```


# Notes
https://github.com/goodtft/LCD-show.git 
scp tft35a-overlay.dtb root@arcade:/boot/overlays/tft35a.dtbo