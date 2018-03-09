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
workdir=$(dirname $(readlink -e $0))
cd $workdir

fix_ownership="$(id -u):$(id -g)"
wheel_dir=${WHEEL_DIR:-${workdir}/docker/wheelbuild/}
wheel_req_control="${wheel_dir}/requirements_control"

sysbase=${SYSBASE:-"debian-9"}

if [ ! "${3}" == "test-ci" ]; then
    echo "Wheeldir: ${wheel_dir}"
    ./docker/build/bricks.sh "${bricks_tag}"

    echo "BUILDING SYSBASE ${sysbase}"
    docker build ${opt_squash} --build-arg PROXY=$http_proxy --build-arg TAG=${tag} --build-arg SYSBASE=${sysbase} -f docker/Dockerfile.sysbase-${sysbase} -t canopsis/canopsis-sysbase:${sysbase}-${tag} .

    if [ ! -d ${wheel_dir} ]; then
        mkdir -p ${wheel_dir}
    fi

    if [ ! -f ${wheel_req_control} ]; then
        touch ${wheel_req_control}
    fi

    current_requirements_control=$(md5sum sources/canopsis/requirements.txt | awk '{print $1}')

    if [ "$(grep ${current_requirements_control} ${wheel_req_control})" = "" ]||[ ! -d ${workdir}/docker/wheels/${sysbase} ]; then
        echo -n "${current_requirements_control}" > ${wheel_req_control}
        echo "BUILDING WHEEL ${sysbase}"
        docker build ${opt_squash} --build-arg PROXY=$http_proxy --build-arg SYSBASE=${sysbase} --build-arg TAG=${tag} -f docker/Dockerfile.wheel -t canopsis/wheel-${sysbase}:latest .
        mkdir -p ${wheel_dir}
        echo "RUNNING WHEEL ${sysbase}"
        docker run -v ${wheel_dir}:/root/wheelrep/ canopsis/wheel-${sysbase}:latest "${fix_ownership}"
    fi

    rm -rf ${workdir}/docker/wheels/
    cp -ar ${wheel_dir} ${workdir}/docker/wheels

    echo "BUILDING CORE ${sysbase}"
    docker build ${opt_squash} --build-arg PROXY=$http_proxy --build-arg SYSBASE=${sysbase} --build-arg TAG=${tag} -f docker/Dockerfile -t canopsis/canopsis-core:${sysbase}-${tag} .

    if [ "${sysbase}" = "debian-9" ]; then
        echo "TAGGING OFFICIAL CANOPSIS-CORE IMAGE"
        docker tag canopsis/canopsis-core:${sysbase}-${tag} canopsis/canopsis-core:${tag}
    fi
fi

if [ "${3}" == "test" ]||[ "${3}" == "test-ci" ]; then
    echo "BUILDING TEST ${sysbase}"
    docker build ${opt_squash} --build-arg PROXY=$http_proxy --build-arg SYSBASE=${sysbase} --build-arg TAG=${tag} -f docker/Dockerfile.tests -t canopsis/canopsis-core:${sysbase}-${tag}-test .
fi
