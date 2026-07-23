#!/bin/bash

set -eu

. /opt/ci-toolset/functions.sh

GOOGLE_PROJECT=altipla-tools

configure-google-cloud

run "sed -i 's/dev/$(build-tag)/g' pkg/config/version.go"
run 'actools go build -o actools ./cmd/actools'
run "gcloud storage cp --cache-control=no-cache actools gs://tools.altipla.consulting/bin/actools"

run "echo $(build-tag) > version"
run "gcloud storage cp --cache-control=no-cache version gs://tools.altipla.consulting/version-manifest/actools"

for FILE in containers/*/Dockerfile; do
  APP=$(basename $(dirname $FILE))
  docker-build-autotag eu.gcr.io/$GOOGLE_PROJECT/$APP containers/$APP/Dockerfile containers/$APP
done

git-tag
