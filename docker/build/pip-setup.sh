#!/usr/bin/bash
set -e
set -o pipefail

virtualenv --no-download ${CPS_HOME}

source /etc/os-release
repver="${ID}-${VERSION_ID}"
find_links="file:///sources/wheels/${repver}"

echo -e "[easy_install]\nallow_hosts = ''\nfind_links = ${find_links}" > /root/.pydistutils.cfg

source ${CPS_HOME}/bin/activate

pip install -b /tmp/pipbuild --no-index -f ${find_links} -U setuptools pip

rm -rf /tmp/pipbuild
rm -rf ~/.cache
