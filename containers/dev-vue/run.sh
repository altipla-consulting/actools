#!/bin/bash

set -eu

cd /workspace/$WORKDIR

exec npm run serve
