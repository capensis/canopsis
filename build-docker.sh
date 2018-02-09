#!/usr/bin/env bash

if [ "${1}" = "" ]; then
    echo "Usage: $0 <tag>"
    exit 1
fi

tag="${1}"

docker build -f docker/Dockerfile -t canopsis/canopsis-core:${tag} .
