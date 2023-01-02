#!/bin/bash

set -e

tinygo-dev build -target=feather-m0
arm-none-eabi-objcopy -O ihex trelligo.elf trelligo.hex
openocd

