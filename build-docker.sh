#!/usr/bin/env bash
set -e
set -o pipefail

if [ "${2}" = "" ]; then
    echo "Usage: $0 <tag> <brick_branch> [test]"
    exit 1
fi

opt_squash=""
tag="${1}"
brick_branch="${2}"
workdir=$(dirname $(readlink -e $0))
cd $workdir

./docker/build/bricks.sh "${brick_branch}"

docker build ${opt_squash} --build-arg PROXY=$http_proxy --build-arg TAG=${tag} -f docker/Dockerfile.sysbase -t canopsis/canopsis-sysbase:${tag} .
docker build ${opt_squash} --build-arg PROXY=$http_proxy --build-arg TAG=${tag} -f docker/Dockerfile -t canopsis/canopsis-core:${tag} .

if [ "${3}" == "test" ]; then
    docker build ${opt_squash} --build-arg PROXY=$http_proxy --build-arg TAG=${tag} -f docker/Dockerfile.tests -t canopsis/canopsis-core:${tag}-test .
fi
