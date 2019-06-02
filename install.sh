#!/bin/bash

set -e

GOARCH=arm go build .

scp arcade-multiplexer root@arcade:/root/
scp test.png root@arcade:/root/