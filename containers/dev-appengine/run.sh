#!/bin/bash

set -eu

mkdir -p /opt/local/workdir/default

cd /workspace/$WORKDIR/appengine
dev_appserver.py . --datastore_path /opt/local/datastore --host 0.0.0.0 --admin_host 0.0.0.0 --port 8080 --go_work_dir /opt/local/workdir --enable_watching_go_path
