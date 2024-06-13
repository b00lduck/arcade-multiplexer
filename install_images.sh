#!/bin/bash

set -e

ssh root@arcade "mkdir -p /root/arcade-multiplexer/images"

scp -r data/images root@arcade:/root/arcade-multiplexer
scp data/artwork/hud_1.jpg root@arcade:/root/arcade-multiplexer/images/
scp data/artwork/hud_2.jpg root@arcade:/root/arcade-multiplexer/images/
scp data/artwork/hud_3.jpg root@arcade:/root/arcade-multiplexer/images/
scp data/artwork/hud_4.jpg root@arcade:/root/arcade-multiplexer/images/

