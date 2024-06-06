#!/bin/bash

make -j 16 ARCH=arm CROSS_COMPILE=arm-linux-gnueabihf- bcmrpi_defconfig
sed -i 's/# CONFIG_USB_OTG is not set/CONFIG_USB_OTG=y/' .config
make -j 16 ARCH=arm CROSS_COMPILE=arm-linux-gnueabihf- zImage modules dtbs

mkdir -p /build/root/
make ARCH=arm CROSS_COMPILE=arm-linux-gnueabihf- INSTALL_MOD_PATH=/build/root modules_install

mkdir -p /build/boot/overlays
cp arch/arm/boot/zImage /build/boot/$KERNEL.img
cp arch/arm/boot/dts/broadcom/*.dtb /build/boot/
cp arch/arm/boot/dts/overlays/*.dtb* /build/boot/overlays/
cp arch/arm/boot/dts/overlays/README /build/boot/overlays/
