#!/usr/bin/env bash
set -e
set -o pipefail

if [ "${1}" = "" ]; then
    echo "Usage: $0 <tag>"
    exit 1
fi

tag="${1}"
workdir=$(dirname $(readlink -e $0))
cd $workdir

docker build --build-arg PROXY=$http_proxy -f docker/Dockerfile.system -t canopsis/canopsis-system:debian8 .
docker build --build-arg PROXY=$http_proxy --build-arg TAG=${tag} -f docker/Dockerfile.core --rm=false -t canopsis/canopsis-core:${tag} .
