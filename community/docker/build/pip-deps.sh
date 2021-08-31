#!/usr/bin/env bash
set -e
set -o pipefail

source /etc/os-release
repver="${ID}-${VERSION_ID}"

source ${CPS_HOME}/bin/activate

LC_ALL=en_US.UTF-8 LANG=en_US.UTF-8 pip --no-color install --no-use-pep517 -r /sources/canopsis/requirements.txt -b /tmp/pipbuild

rm -rf /tmp/pipbuild
rm -rf ~/.cache
