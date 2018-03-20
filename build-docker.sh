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
wheel_dir=${WHEEL_DIR:-${workdir}/docker/wheelbuild/}
wheel_req_control="${wheel_dir}/requirements_control"

function build_for_sysbase() {
    if [ "${1}" = "" ]; then
        echo "wrong arguments"
        exit 2
    fi

    local sysbase="${1}"

    CPS_DOCKER_BUILD_ARGS="${opt_squash} --build-arg PROXY=$http_proxy --build-arg TAG=${tag} --build-arg SYSBASE=${sysbase}"

    if [ ! "${mode}" == "test-ci" ]; then
        echo "Wheeldir: ${wheel_dir}"

        echo "BUILDING SYSBASE ${sysbase}"
        docker build ${CPS_DOCKER_BUILD_ARGS} -f docker/Dockerfile.sysbase-${sysbase} -t canopsis/canopsis-sysbase:${sysbase}-${tag} .

        if [ ! -d ${wheel_dir} ]; then
            mkdir -p ${wheel_dir}
        fi

        if [ ! -f ${wheel_req_control} ]; then
            touch ${wheel_req_control}
        fi

        current_requirements_control=$(md5sum sources/canopsis/requirements.txt | awk '{print $1}')

        if [ "$(grep ${current_requirements_control} ${wheel_req_control})" = "" ]||[ ! -d ${wheel_dir}/${sysbase} ]; then
            echo -n "${current_requirements_control}" > ${wheel_req_control}
            echo "BUILDING WHEEL ${sysbase}"

            docker build ${CPS_DOCKER_BUILD_ARGS} -f docker/Dockerfile.wheel -t canopsis/wheel-${sysbase}:latest .
            mkdir -p ${wheel_dir}

            echo "RUNNING WHEEL ${sysbase}"

            docker run --rm -v ${wheel_dir}:/root/wheelrep/ canopsis/wheel-${sysbase}:latest "${fix_ownership}"
        fi

        rm -rf ${workdir}/docker/wheels/
        cp -ar ${wheel_dir} ${workdir}/docker/wheels

        echo "BUILDING CORE ${sysbase}"

        docker build ${CPS_DOCKER_BUILD_ARGS} -f docker/Dockerfile -t canopsis/canopsis-core:${sysbase}-${tag} .

		echo "Building provisonning image"
		docker build ${CPS_DOCKER_BUILD_ARGS} -f docker/Dockerfile.prov -t canopsis/canopsis-prov:${sysbase}-${tag} .

        if [ "${sysbase}" = "debian-9" ]; then
            echo "TAGGING OFFICIAL CANOPSIS-CORE IMAGE"

            docker tag canopsis/canopsis-core:${sysbase}-${tag} canopsis/canopsis-core:${tag}
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
