#!/bin/bash

set -eu

dev_appserver.py --datastore_path /opt/local/datastore --host 0.0.0.0 --admin_host 0.0.0.0 --port 8080 appengine
