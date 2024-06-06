#!/bin/bash

docker build -t rpi-kernel-builder .
docker run --rm -v $(pwd)/build:/build rpi-kernel-builder