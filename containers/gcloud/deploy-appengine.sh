#!/bin/bash

set -eu

mkdir -p /go/src/$(dirname $PROJECT)
ln -s /workspace /go/src/$PROJECT

mkdir -p /vendor
ln -s /vendor/src /go/src/$PROJECT/vendor

export GOPATH=/vendor:/go

gcloud app deploy $*
