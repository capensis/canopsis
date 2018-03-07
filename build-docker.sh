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

echo "Wheeldir: ${wheel_dir}"

./docker/build/bricks.sh "${bricks_tag}"

docker build ${opt_squash} --build-arg PROXY=$http_proxy --build-arg TAG=${tag} -f docker/Dockerfile.sysbase -t canopsis/canopsis-sysbase:${tag} .

if [ ! -d ${wheel_dir} ]; then
    mkdir -p ${wheel_dir}
fi

if [ ! -f ${wheel_req_control} ]; then
    touch ${wheel_req_control}
fi

current_requirements_control=$(md5sum sources/canopsis/requirements.txt | awk '{print $1}')
if [ "$(grep ${current_requirements_control} ${wheel_req_control})" = "" ]; then
    echo -n "${current_requirements_control}" > ${wheel_req_control}
    docker build ${opt_squash} --build-arg PROXY=$http_proxy -f docker/Dockerfile.wheel -t canopsis/debian-9-wheel:latest .
    mkdir -p ${wheel_dir}
    docker run -v ${wheel_dir}:/root/wheelrep/ canopsis/debian-9-wheel:latest "${fix_ownership}"
fi

rm -rf ${workdir}/docker/wheels/
cp -ar ${wheel_dir} ${workdir}/docker/wheels

docker build ${opt_squash} --build-arg PROXY=$http_proxy --build-arg TAG=${tag} -f docker/Dockerfile -t canopsis/canopsis-core:${tag} .

if [ "${3}" == "test" ]; then
    docker build ${opt_squash} --build-arg PROXY=$http_proxy --build-arg TAG=${tag} -f docker/Dockerfile.tests -t canopsis/canopsis-core:${tag}-test .
fi
