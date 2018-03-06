#!/usr/bin/bash
set -e
set -o pipefail

virtualenv --no-download /root/wheelbuild

source /root/wheelbuild/bin/activate
pip install -U wheel

source /etc/os-release
repver="${ID}-${VERSION_ID}"

mkdir -p /root/wheelrep/${repver}
cd /root/wheelrep/${repver}


export MAKEFLAGS="-j4"
pip wheel --no-index -f file:///sources/python-libs/ -r /sources/canopsis/requirements.txt