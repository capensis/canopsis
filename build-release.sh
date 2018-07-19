#!/usr/bin/env bash
set -e
set -o pipefail

# Helper script to be ran before any other build-*.sh script. It will:
# * build catag
# * run catag to apply tags on all projects configured
# * import canopsis-next (pull, rsync, add, commit, push)
# * run catag again to sync with any added files

workdir=$(dirname $(readlink -e $0))
cd ${workdir}

source ${workdir}/build-env.sh

function deploy_catag() {
    cd ${workdir}

    rm -rf ${GOPATH}/src/canopsis-catag
    ln -s ${workdir}/tools/catag ${GOPATH}/src/canopsis-catag
    cd ${GOPATH}/src/canopsis-catag
    glide install
    go build .
    mv canopsis-catag catag
}

function run_catag() {
    cd ${workdir}

    ./tools/catag/catag -ini tools/catag/catag.ini -tag ${CANOPSIS_TAG} -token ${CANOPSIS_CATAG_TOKEN}
}

function bump_canopsis_next() {
    cd ${workdir}

    rm -rf tmp-canopsis-next
    git clone ssh://git@git.canopsis.net/canopsis/canopsis-next -b ${CANOPSIS_TAG} tmp-canopsis-next

    rm -rf tmp-canopsis-next/.git

    rsync -avKSH tmp-canopsis-next/ sources/webcore/src/canopsis-next/

    rm -rf tmp-canopsis-next

    git add sources/webcore/src/canopsis-next
    git commit --allow-empty --author="build-release.sh <canopsis@canopsis.fr>" -m "auto: bump canopsis-next ${CANOPSIS_TAG}"
    git push $(git remote) $(git rev-parse --abbrev-ref HEAD)
}

deploy_catag
run_catag
bump_canopsis_next
run_catag
