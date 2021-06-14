#!/usr/bin/env bash
set -e
set -o pipefail
set -u

workdir=$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)
cd ${workdir}

# just avoid double confirmations
./build-env.sh
export CANOPSIS_ENV_CONFIRM=0

# launch all builds
./build-docker.sh
./build-packages.sh
