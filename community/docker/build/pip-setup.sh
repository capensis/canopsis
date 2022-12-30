#!/usr/bin/env bash
set -e
set -o pipefail

# Make sure that virtualenv doesn't auto-upgrade pip, because
# recent versions have problems with our old Python 2 setup.
export PATH=/opt/python/3.8/bin:$PATH
which python3
python3 --version
virtualenv -p python3 --no-download --system-site-packages ${CPS_HOME}

source /etc/os-release
repver="${ID}-${VERSION_ID}"

source ${CPS_HOME}/bin/activate

pip3 install -U setuptools "pip==22.3.1"
pip3 --no-color install -U wheel

rm -rf ~/.cache
