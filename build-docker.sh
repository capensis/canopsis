#!/usr/bin/env bash
set -e
set -o pipefail

opt_squash=""

workdir=$(dirname $(readlink -e $0))
cd $workdir

source ${workdir}/build-env.sh

fix_ownership="$(id -u):$(id -g)"
mode=${CANOPSIS_DOCKER_MODE}

function build_for_distribution() {
    if [ "${1}" = "" ]; then
        echo "wrong arguments"
        exit 2
    fi

    local distribution="${1}"
    local tag=${CANOPSIS_TAG}
    local docker_args="${opt_squash} --build-arg PROXY=$http_proxy --build-arg CANOPSIS_TAG=${tag} --build-arg CANOPSIS_DISTRIBUTION=${distribution}"
    local full_tag="${distribution}-${tag}"

    if [ ! "${mode}" == "test-ci" ]; then
        echo "BUILDING DISTRIBUTION ${distribution}"
        docker build ${docker_args} -f docker/Dockerfile.sysbase-${distribution} -t canopsis/canopsis-sysbase:${full_tag} .

        echo "BUILDING CORE ${distribution}"
        docker build ${docker_args} -f docker/Dockerfile -t canopsis/canopsis-core:${full_tag} .

		echo "Building provisionning image"
		docker build ${docker_args} -f docker/Dockerfile.prov -t canopsis/canopsis-prov:${full_tag} .

        if [ "${distribution}" = "debian-9" ]; then
            echo "TAGGING OFFICIAL CANOPSIS-CORE IMAGE"

            docker tag canopsis/canopsis-core:${full_tag} canopsis/canopsis-core:${tag}
            docker tag canopsis/canopsis-prov:${full_tag} canopsis/canopsis-prov:${tag}
        fi
    fi

    if [ "${CANOPSIS_DOCKER_MODE}" == "test" ]||[ "${CANOPSIS_DOCKER_MODE}" == "test-ci" ]; then
        echo "BUILDING TEST ${distribution}"
        docker tag canopsis/canopsis-prov:${full_tag} canopsis/canopsis-prov:${full_tag}-test
        docker build ${docker_args} -f docker/Dockerfile.tests -t canopsis/canopsis-core:${full_tag}-test .
    fi
}

function build() {
    cd ${workdir}

    ./docker/build/bricks.sh

    if [ "${CANOPSIS_DISTRIBUTION}" = "all" ]; then
        build_for_distribution "debian-9"
        build_for_distribution "debian-8"
        build_for_distribution "centos-7"
    else
        build_for_distribution "${CANOPSIS_DISTRIBUTION}"
    fi
}

build
