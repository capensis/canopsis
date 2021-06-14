#!/usr/bin/env bash
set -e
set -o pipefail
set -u

opt_squash=""
workdir=$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)
cd $workdir

source ${workdir}/build-env.sh

function build_package_for_distribution() {
    if [ "${1}" = "" ]; then
        echo "wrong params: $0 distribution"
        exit 2
    fi

    local distribution="${1}"
    local fix_ownership="$(id -u):$(id -g)"
    local tag=${CANOPSIS_TAG}
    local full_tag="${distribution}-${tag}"

    echo "BUILDING PACKAGE FOR ${distribution}"
    docker_args="${opt_squash} --build-arg FIX_OWNERSHIP=${fix_ownership} --build-arg PROXY=$http_proxy --build-arg CANOPSIS_DISTRIBUTION=${distribution} --build-arg CANOPSIS_TAG=${tag} --build-arg CANOPSIS_PACKAGE_TAG=${CANOPSIS_PACKAGE_TAG} --build-arg CANOPSIS_PACKAGE_REL=${CANOPSIS_PACKAGE_REL}"

    docker build ${docker_args} -f docker/Dockerfile.packaging --rm=false -t canopsis/canopsis-packaging:${full_tag} .
    docker run -v $(pwd)/packages:/packages canopsis/canopsis-packaging:${full_tag}
}

# Force only CentOS-7 here
build_package_for_distribution "centos-7"
