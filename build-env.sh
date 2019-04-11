set -e
set -o pipefail
set -u

# docker mode is only a hint for build scripts to operate.
# regular doesn't mean anything in particular, only that it is not
# test or test-ci for example, which is set by .gitlab-ci.yml
export CANOPSIS_DOCKER_MODE=${CANOPSIS_DOCKER_MODE:=regular}

# print env vars and ask for confirmation
export CANOPSIS_ENV_RECAP=${CANOPSIS_ENV_RECAP:=1}
export CANOPSIS_ENV_CONFIRM=${CANOPSIS_ENV_CONFIRM:=1}

# which version of uiv2 bricks to pull
export CANOPSIS_UIV2_BRICKS_TAG=${CANOPSIS_UIV2_BRICKS_TAG:="${CANOPSIS_TAG}"}

# override package version/tag and release if it must be different from the
# regular CANOPSIS_TAG variable
export CANOPSIS_PACKAGE_TAG=${CANOPSIS_PACKAGE_TAG:="${CANOPSIS_TAG}"}
export CANOPSIS_PACKAGE_REL=${CANOPSIS_PACKAGE_REL:="1"}

# GitLab Access Token
export CANOPSIS_CATAG_TOKEN=${CANOPSIS_CATAG_TOKEN:=""}

# Do not check GOPATH and Go installation.
# Can be used by any other script to avoid doing anything with Go.
export CANOPSIS_SKIP_GO=${CANOPSIS_SKIP_GO:="0"}

# By default, pull, commit and push sources/webcore/src/canopsis-next from
# the canopsis/canopsis-next repository.
export CANOPSIS_BUILD_NEXT=${CANOPSIS_BUILD_NEXT:="0"}

# Only avoid undefined variable error
export GOPATH=${GOPATH:=""}
export http_proxy=${http_proxy:=""}
export https_proxy=${https_proxy:=""}

function env_recap() {
    echo "CANOPSIS_TAG: ${CANOPSIS_TAG}"
    echo "CANOPSIS_UIV2_BRICKS_TAG: ${CANOPSIS_UIV2_BRICKS_TAG}"
    echo "CANOPSIS_DISTRIBUTION: ${CANOPSIS_DISTRIBUTION}"
    echo "CANOPSIS_PACKAGE_TAG: ${CANOPSIS_PACKAGE_TAG}"
    echo "CANOPSIS_PACKAGE_REL: ${CANOPSIS_PACKAGE_REL}"
    echo "CANOPSIS_DOCKER_MODE: ${CANOPSIS_DOCKER_MODE}"
    echo "CANOPSIS_BUILD_NEXT: ${CANOPSIS_BUILD_NEXT}"
    echo "CANOPSIS_CATAG_TOKEN: set, hidden."
    echo "GOPATH: ${GOPATH}"
}

function ensure_env() {
    if [ ! "${CANOPSIS_SKIP_GO}" = "1" ]; then
        if [ "${GOPATH}" = "" ]; then
            echo "\$GOPATH is not initialised: setup go environment properly."
            exit 1
        fi

        if [ ! -d "${GOPATH}/src" ]; then
            mkdir -p ${GOPATH}/src
        fi
    fi

    if [ "${CANOPSIS_CATAG_TOKEN}" = "" ]; then
        echo "\$CANOPSIS_CATAG_TOKEN is not initialised: provide gitlab api access token"
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

    if [ "${CANOPSIS_DISTRIBUTION}" = "" ]; then
        echo "CANOPSIS_DISTRIBUTION is empty: give \"all\" to build on all systems, or:"
        echo "debian-8"
        echo "debian-9"
        echo "centos-7"
        exit 5
    fi


    # just try to run these programs so with set -e if they are not found,
    # it will exit with command not found
    rsync --version 2>&1 > /dev/null
    git --version 2>&1 > /dev/null
    if [ ! "${CANOPSIS_SKIP_GO}" = "1" ]; then
        go version 2>&1 > /dev/null
    fi

    if [ ! "${CANOPSIS_ENV_RECAP}" = "0" ]; then
        env_recap

        if [ "${CANOPSIS_ENV_CONFIRM}" = "1" ]; then
            echo -n "Confirm environment [y/N]: "
            local confirm="N"
            read confirm
            if [ ! "${confirm}" = "Y" ]&&[ ! "${confirm}" = "y" ]; then
                echo "Environment not confirmed, exit."
                exit 6
            else
                export CANOPSIS_ENV_CONFIRM=0
            fi
        fi
    fi
}

ensure_env
