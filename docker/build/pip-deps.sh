#!/usr/bin/bash
set -e
set -o pipefail

source /etc/os-release
repver="${ID}-${VERSION_ID}"

source ${CPS_HOME}/bin/activate

pip install -r /sources/canopsis/requirements.txt -b /tmp/pipbuild

rm -rf /tmp/pipbuild
rm -rf ~/.cache
