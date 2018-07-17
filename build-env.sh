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

    if [ "${CANOPSIS_UIV2_BRICKS_TAG}" = "" ]; then
        echo "\$CANOPSIS_UIV2_BRICKS_TAG is not initialised: provide uiv2 bricks version"
        exit 4
    fi

    if [ "${CANOPSIS_DOCKER_MODE}" = "" ]; then
        # test or test-ci are used to build test containers
        export CANOPSIS_DOCKER_MODE=regular
    fi

    # just try to run these programs so with set -e if they are not found,
    # it will exit with command not found
    rsync --version 2>&1 > /dev/null
    git --version 2>&1 > /dev/null
    go version 2>&1 > /dev/null
}

ensure_env
