#!/bin/bash

set -eu

cd /workspace/$WORKDIR

exec nodemon src/index.js
