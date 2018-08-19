#!/bin/bash

set -eu

mkdir -p /go/src/$(dirname $PROJECT)
ln -s /workspace /go/src/$PROJECT

export GOPATH=/go

gcloud app deploy $*
