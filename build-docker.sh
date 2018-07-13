#!/usr/bin/env bash
set -e
set -o pipefail

if [ "${2}" = "" ]; then
    echo "Usage: $0 <tag> <brick_branch> [test]"
    exit 1
fi

opt_squash=""
tag="${1}"
bricks_tag="${2}"
mode="${3}"
workdir=$(dirname $(readlink -e $0))

cd $workdir

fix_ownership="$(id -u):$(id -g)"

function build_for_sysbase() {
    if [ "${1}" = "" ]; then
        echo "wrong arguments"
        exit 2
    fi

    local sysbase="${1}"

    CPS_DOCKER_BUILD_ARGS="${opt_squash} --build-arg PROXY=$http_proxy --build-arg TAG=${tag} --build-arg SYSBASE=${sysbase}"

    if [ ! "${mode}" == "test-ci" ]; then
        echo "BUILDING SYSBASE ${sysbase}"
        docker build ${CPS_DOCKER_BUILD_ARGS} -f docker/Dockerfile.sysbase-${sysbase} -t canopsis/canopsis-sysbase:${sysbase}-${tag} .

        echo "BUILDING CORE ${sysbase}"
        docker build ${CPS_DOCKER_BUILD_ARGS} -f docker/Dockerfile -t canopsis/canopsis-core:${sysbase}-${tag} .

		echo "Building provisionning image"
		docker build ${CPS_DOCKER_BUILD_ARGS} -f docker/Dockerfile.prov -t canopsis/canopsis-prov:${sysbase}-${tag} .

        if [ "${sysbase}" = "debian-9" ]; then
            echo "TAGGING OFFICIAL CANOPSIS-CORE IMAGE"

            docker tag canopsis/canopsis-core:${sysbase}-${tag} canopsis/canopsis-core:${tag}
            docker tag canopsis/canopsis-prov:${sysbase}-${tag} canopsis/canopsis-prov:${tag}
        fi
    fi

    if [ "${mode}" == "test" ]||[ "${mode}" == "test-ci" ]; then
        echo "BUILDING TEST ${sysbase}"
        docker tag canopsis/canopsis-prov:${sysbase}-${tag} canopsis/canopsis-prov:${sysbase}-${tag}-test
        docker build ${CPS_DOCKER_BUILD_ARGS} -f docker/Dockerfile.tests -t canopsis/canopsis-core:${sysbase}-${tag}-test .
    fi
}

./docker/build/bricks.sh "${bricks_tag}"

if [ "${SYSBASE}" = "" ]; then
    build_for_sysbase "centos-7"
    build_for_sysbase "debian-8"
    build_for_sysbase "debian-9"
else
    build_for_sysbase "${SYSBASE}"
fi
