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

# override package version/tag and release if it must be different from the
# regular CANOPSIS_TAG variable
export CANOPSIS_PACKAGE_TAG=${CANOPSIS_PACKAGE_TAG:="${CANOPSIS_TAG}"}
export CANOPSIS_PACKAGE_REL=${CANOPSIS_PACKAGE_REL:="1"}

# Only avoid undefined variable error
export http_proxy=${http_proxy:=""}
export https_proxy=${https_proxy:=""}

export FORCE_COLOR=0
export CLICOLOR=0
export NO_COLOR=1
export DOCKER_SCAN_SUGGEST=false
export NPM_CONFIG_PROGRESS=false
export NPM_CONFIG_SPIN=false
export NPM_CONFIG_COLOR=false
export BUILDKIT_PROGRESS=plain

function env_recap() {
    echo "CANOPSIS_TAG: ${CANOPSIS_TAG}"
    echo "CANOPSIS_DISTRIBUTION: ${CANOPSIS_DISTRIBUTION}"
    echo "CANOPSIS_PACKAGE_TAG: ${CANOPSIS_PACKAGE_TAG}"
    echo "CANOPSIS_PACKAGE_REL: ${CANOPSIS_PACKAGE_REL}"
    echo "CANOPSIS_DOCKER_MODE: ${CANOPSIS_DOCKER_MODE}"
}

function ensure_env() {
    if [ "${CANOPSIS_TAG}" = "" ]; then
        echo "\$CANOPSIS_TAG is not initialised: provide TAG for release"
        exit 3
    fi

    if [ "${CANOPSIS_DOCKER_MODE}" = "" ]; then
        # test or test-ci are used to build test containers
        export CANOPSIS_DOCKER_MODE=regular
    fi

    if [ "${CANOPSIS_DISTRIBUTION}" = "" ]; then
        echo "CANOPSIS_DISTRIBUTION is empty: give \"all\" to build on all systems, or:"
        echo "debian-9"
        echo "centos-7"
        exit 5
    fi


    # just try to run these programs so with set -e if they are not found,
    # it will exit with command not found
    rsync --version 2>&1 > /dev/null
    git --version 2>&1 > /dev/null

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
