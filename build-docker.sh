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

docker build --build-arg PROXY=$http_proxy --build-arg TAG=${tag} -f docker/Dockerfile.sysbase -t canopsis/canopsis-sysbase:${tag} .
docker build --build-arg PROXY=$http_proxy --build-arg TAG=${tag} -f docker/Dockerfile -t canopsis/canopsis-core:${tag} .
