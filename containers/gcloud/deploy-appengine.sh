#!/bin/bash

set -eu

mkdir -p /go/src/$(dirname $PROJECT)
ln -s /workspace /go/src/$PROJECT

gcloud app deploy $*
