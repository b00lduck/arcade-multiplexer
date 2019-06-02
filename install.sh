#!/bin/bash

set -e

GOARCH=arm GOOS=linux go build .

set +e
ssh root@arcade "killall -q /root/arcade-multiplexer"

set -e
scp arcade-multiplexer root@arcade:/root/
#scp test.png root@arcade:/root/

ssh root@arcade "/root/arcade-multiplexer"
