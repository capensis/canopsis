#!/usr/bin/env bash
set -e
set -o pipefail

# Make sure that virtualenv doesn't auto-upgrade pip, because
# recent versions have problems with our old Python 2 setup.
virtualenv --no-download --system-site-packages ${CPS_HOME}

source /etc/os-release
repver="${ID}-${VERSION_ID}"

source ${CPS_HOME}/bin/activate

pip install -b /tmp/pipbuild -U setuptools "pip==20.1.1"
pip --no-color install -b /tmp/pipbuild -U wheel

rm -rf /tmp/pipbuild
rm -rf ~/.cache
