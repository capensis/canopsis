#!/usr/bin/env bash
set -e
set -o pipefail
set -u

workdir=$(dirname $(readlink -e $0))
cd ${workdir}

./build-docker.sh
./build-packages.sh
