#!/bin/bash

set -eu

mkdir -p ~/bin
source ~/.bashrc

curl https://tools.altipla.consulting/bin/actools > ~/bin/actools
chmod +x ~/bin/actools

actools pull
