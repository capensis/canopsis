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

deploy_catag
run_catag
