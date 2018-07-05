#!/usr/bin/bash
set -e
set -o pipefail

virtualenv /root/wheelbuild

export PS1="$ "
source /root/wheelbuild/bin/activate

pip install -U setuptools pip wheel

source /etc/os-release
repver="${ID}-${VERSION_ID}"

mkdir -p /root/wheelrep/${repver}
cd /root/wheelrep/${repver}

pip wheel -r /sources/canopsis/requirements.txt

chown -R ${1} /root/wheelrep
