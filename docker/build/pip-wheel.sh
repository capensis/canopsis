#!/usr/bin/bash
set -u
set -e
set -o pipefail

virtualenv --no-download /root/wheelbuild

export PS1="$ "
source /root/wheelbuild/bin/activate
pip install -U wheel

source /etc/os-release
repver="${ID}-${VERSION_ID}"

mkdir -p /root/wheelrep/${repver}
cd /root/wheelrep/${repver}


export MAKEFLAGS="-j4"
pip wheel --no-index -f file:///sources/python-libs/ -r /sources/canopsis/requirements.txt

chown -R ${1} /root/wheelrep
