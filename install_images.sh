#!/bin/bash

set -e

go build ./cmd/image-converter

ssh root@arcade "mkdir -p /root/arcade-multiplexer/images"

find data/images -type f -name *.jpg | xargs -n 1 ./image-converter

scp -r data/images root@arcade:/root/arcade-multiplexer

