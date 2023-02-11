#!/bin/bash

set -e

/home/trichner/workspaces/go/tinygo-v0.26.0/bin/tinygo flash -target itsybitsy-m4 -scheduler tasks -gc conservative -size full

#arm-none-eabi-objcopy -O ihex trelligo.elf a.hex

#mv a.hex /run/media/trichner/NODE_L432KC/
