#!/bin/bash

set -eu

mkdir -p /go/src/$(dirname $PROJECT)
ln -s /go/src/$PROJECT /workspace

gcloud app deploy $*
