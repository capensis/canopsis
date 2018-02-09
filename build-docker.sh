#!/usr/bin/env bash
set -e
set -o pipefail

if [ "${1}" = "" ]; then
    echo "Usage: $0 <tag>"
    exit 1
fi

tag="${1}"

docker build -f docker/Dockerfile.system -t canopsis/canopsis-system:debian8 .
docker build -f docker/Dockerfile.core --rm=false -t canopsis/canopsis-core:${tag} .
