#!/bin/bash

set -e

GOARCH=arm GOARM=6 GOOS=linux go build ./cmd/arcade-multiplexer/

set +e
ssh root@arcade "pkill -f './arcade-multiplexer'"
ssh root@arcade "mkdir -p /root/arcade-multiplexer"
set -e


scp arcade-multiplexer root@arcade:/root/arcade-multiplexer
scp data/config.yml root@arcade:/root/arcade-multiplexer

ssh -t root@arcade "cd /root/arcade-multiplexer && ./arcade-multiplexer"

