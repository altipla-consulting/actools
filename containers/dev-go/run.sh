#!/bin/bash

set -eu

if [[ -z $WORKDIR ]]; then
  APP=$SERVICE
else
  APP=$(basename $WORKDIR)
fi

cd /workspace/$WORKDIR

echo """
**/*.go /workspace/pkg/**/*.go {
  prep: go install ./cmd/$APP
}

/go/bin/$APP {
  daemon +sigterm: $APP
}
""" > /tmp/modd.conf

modd -f /tmp/modd.conf
