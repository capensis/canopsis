#!/usr/bin/env bash
set -e
set -o pipefail

if [ ! -d /provisioning ]; then
    echo "directory /provisioning not found, exit 0"
    exit 0
fi

if [ -d /provisioning/mongo ]; then
    echo "directory /provisioning/mongo not found, exit 0"
    rsync --recursive --perms --omit-dir-times -vKSH /provisioning/mongo/ /opt/canopsis/opt/mongodb/load.d/
fi
