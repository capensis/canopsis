#!/usr/bin/env bash
set -e
set -o pipefail

if [ "${1}" = "" ]; then
    echo "Usage: $0 <tag>"
    exit 1
fi

opt_squash=""
tag="${1}"
release="${CPS_PKG_REL:-1}"
CPS_PKG_TAG=${CPS_PKG_TAG:="${tag}"}
workdir=$(dirname $(readlink -e $0))
export FIX_OWNERSHIP="$(id -u):$(id -g)"
cd $workdir

function build_package_for_sysbase() {
    if [ "${1}" = "" ]; then
        echo "wrong params: $0 sysbase"
        exit 2
    fi

    local sysbase="${1}"

    echo "BUILDING PACKAGE FOR ${sysbase}"
    docker_args="${opt_squash} --build-arg FIX_OWNERSHIP=${FIX_OWNERSHIP} --build-arg PROXY=$http_proxy --build-arg SYSBASE=${sysbase} --build-arg TAG=${tag} --build-arg CPS_PKG_TAG=${CPS_PKG_TAG} --build-arg CPS_PKG_REL=${release}"

    docker build ${docker_args} -f docker/Dockerfile.packaging --rm=false -t canopsis/canopsis-packaging:${sysbase}-${tag} .
    docker run -ti -v $(pwd)/packages:/packages canopsis/canopsis-packaging:${sysbase}-${tag}
}

if [ "${SYSBASE}" = "" ]; then
    build_package_for_sysbase "centos-7"
    build_package_for_sysbase "debian-8"
    build_package_for_sysbase "debian-9"
else
    build_package_for_sysbase "${SYSBASE}"
fi
