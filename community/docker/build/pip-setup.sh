#!/usr/bin/env bash
set -e
set -o pipefail

export PATH=/opt/python/3.9/bin:$PATH
python3 --version
type -P python3
pip3.9 install -U setuptools "pip==23.0"

python3 -m pip install --upgrade virtualenv
virtualenv -p python3 --no-download --system-site-packages ${CPS_HOME}

source /etc/os-release
repver="${ID}-${VERSION_ID}"

source ${CPS_HOME}/bin/activate

pip3 --no-color install -U wheel

rm -rf ~/.cache
