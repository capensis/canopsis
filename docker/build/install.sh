#!/usr/bin/bash
set -e
set -o pipefail

virtualenv --no-download ${CPS_HOME}

find_links="file:///sources/python-libs/"

echo -e "[easy_install]\nallow_hosts = ''\nfind_links = ${find_links}" > /root/.pydistutils.cfg

source ${CPS_HOME}/bin/activate

pip install -b /tmp/pipbuild --no-index -f ${find_links} -U setuptools pip
pip install -b /tmp/pipbuild --no-index -f ${find_links} /sources/canopsis/

rm -rf /tmp/pipbuild
rm -rf ~/.cache
