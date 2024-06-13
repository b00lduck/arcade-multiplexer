#!/bin/bash

scp data/config.txt root@arcade:/boot/firmware/config.txt
#scp data/cmdline.txt root@arcade:/boot/firmware/cmdline.txt
scp data/modules root@arcade:/etc/modules
scp data/install_hid.sh root@arcade:/root/install_hid.sh
ssh root@arcade "chmod +x /root/install_hid.sh"
scp data/rc.local root@arcade:/etc/rc.local
ssh root@arcade "chmod +x /etc/rc.local"
ssh root@arcade "systemctl disable --now getty@tty0"
ssh root@arcade "reboot"
