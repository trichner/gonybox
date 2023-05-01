#!/bin/bash

set -e

./flash_itsybitsy-m4.sh
rm screenlog.0
sleep 3
screen -L /dev/ttyACM0 9600
