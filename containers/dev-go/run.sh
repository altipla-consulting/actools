#!/bin/bash

set -eu

APP=$(basename $WORKDIR)

cd /go/src/$PROJECT/$WORKDIR

echo """
**/*.go {
  prep: go install ./cmd/$APP
  daemon +sigterm: $APP
}
""" > /tmp/modd.conf

modd -f /tmp/modd.conf
