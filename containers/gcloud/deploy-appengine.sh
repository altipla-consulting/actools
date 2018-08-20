#!/bin/bash

set -eu

BUILD_DIR=$(mktemp -d)

cd $BUILD_DIR

echo " [*] Prepare build directory"
mkdir -p src/$(dirname $PROJECT)
rsync -r --exclude=.git --exclude=node_modules /workspace/ src/$PROJECT

echo " [*] Extract vendored files"
rsync -a src/$PROJECT/vendor/ src
rm -rf src/$PROJECT/vendor

export GOPATH=$BUILD_DIR

cd $BUILD_DIR/src/$PROJECT
gcloud app deploy $*
