#!/bin/bash

set -eu

dev_appserver.py /workspace/$WORKDIR/app.yaml --datastore_path /home/container/datastore --host 0.0.0.0 --admin_host 0.0.0.0 --port 8080
