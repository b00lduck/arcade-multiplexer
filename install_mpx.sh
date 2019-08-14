#!/bin/bash

set -e

GOARCH=arm GOOS=linux go build .

set +e
ssh root@arcade "killall -q ./arcade-multiplexer"
set -e

ssh root@arcade "mkdir -p /root/arcade-multiplexer/images"

scp -r data/images root@arcade:/root/arcade-multiplexer
scp arcade-multiplexer root@arcade:/root/arcade-multiplexer
scp data/config.yml root@arcade:/root/arcade-multiplexer

ssh -t root@arcade "cd /root/arcade-multiplexer && ./arcade-multiplexer"
