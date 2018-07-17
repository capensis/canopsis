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

function ensure_env() {
    cd ${workdir}

    if [ "${GOPATH}" = "" ]; then
        echo "\$GOPATH is not initialised: setup go environment properly."
        exit 1
    fi

    if [ ! -d "${GOPATH}/src" ]; then
        mkdir -p ${GOPATH}/src
    fi

    if [ "${CATAG_TOKEN}" = "" ]; then
        echo "\$CATAG_TOKEN is not initialised: provide gitlab api access token"
        exit 2
    fi

    if [ "${CANOPSIS_TAG}" = "" ]; then
        echo "\$CANOPSIS_TAG is not initialised: provide TAG for release"
        exit 3
    fi

    # just try to run these programs so with set -e if they are not found,
    # it will exit with command not found
    rsync -h 2>&1 > /dev/null
    git -h 2>&1 > /dev/null
    go -h 2>&1 > /dev/null
}

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

    ./tools/catag/catag -ini tools/catag/catag.ini -tag ${CANOPSIS_TAG} -token ${CATAG_TOKEN}
}

function bump_canopsis_next() {
    cd ${workdir}

    rm -rf tmp-canopsis-next
    git clone ssh://git@git.canopsis.net/canopsis/canopsis-next -b ${CANOPSIS_TAG} tmp-canopsis-next
    rm -rf tmp-canopsis-next/.git
    rsync -avKSH tmp-canopsis-next/ sources/webcore/src/canopsis-next/
    rm -rf tmp-canopsis-next

    git add sources/webcore/src/canopsis-next
    git commit --allow-empty -am "auto: bump canopsis-next ${CANOPSIS_TAG}"
    git push $(git remote) $(git rev-parse --abbrev-ref HEAD)
}

ensure_env
deploy_catag
run_catag
bump_canopsis_next
run_catag
