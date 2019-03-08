#!/bin/bash

set -eu

if [[ -z $WORKDIR ]]; then
  APP=$SERVICE
else
  APP=$(basename $WORKDIR)
fi

cd /workspace/$WORKDIR

if [ ! -f containers/$APP/modd.conf ]; then
  echo """
  **/*.go /workspace/pkg/**/*.go {
    prep: go install ./cmd/$APP
    daemon +sigterm: $APP
  }
  """ > /tmp/modd.conf

  modd -f /tmp/modd.conf
else
  modd -f containers/$APP/modd.conf
fi
