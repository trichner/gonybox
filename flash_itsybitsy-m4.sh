#!/bin/bash

set -e

/home/trichner/workspaces/tinygo/build/tinygo \
    flash -target itsybitsy-m4 -scheduler tasks -gc conservative -size full

#arm-none-eabi-objcopy -O ihex trelligo.elf a.hex

#mv a.hex /run/media/trichner/NODE_L432KC/
