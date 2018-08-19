#!/bin/bash

set -eu

. /opt/ci-toolset/functions.sh

GOOGLE_PROJECT=altipla-tools

for FILE in containers/*/Dockerfile; do
  APP=$(basename $(dirname $FILE))
  docker-build-autotag eu.gcr.io/$GOOGLE_PROJECT/$APP containers/$APP/Dockerfile containers/$APP
done

run "sed -i 's/dev/${build-tag}/g' pkg/config/version.go"
run 'actools go build -o actools ./cmd/actools'
run "gsutil -h 'Cache-Control: no-cache' cp actools gs://tools.altipla.consulting/bin/actools"

run "echo ${build-tag} > version"
run "gsutil -h 'Cache-Control: no-cache' cp version gs://tools.altipla.consulting/version-manifest/actools"

git-tag
