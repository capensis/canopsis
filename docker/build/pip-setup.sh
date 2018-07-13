#!/usr/bin/bash
set -e
set -o pipefail

virtualenv ${CPS_HOME}

source /etc/os-release
repver="${ID}-${VERSION_ID}"

source ${CPS_HOME}/bin/activate

pip install -b /tmp/pipbuild -U setuptools pip wheel

rm -rf /tmp/pipbuild
rm -rf ~/.cache
