#!/usr/bin/bash
set -e
set -o pipefail

source /etc/os-release
repver="${ID}-${VERSION_ID}"

source ${CPS_HOME}/bin/activate
find_links="file:///sources/wheels/${repver}"

pip install -r /sources/canopsis/requirements.txt -b /tmp/pipbuild --no-index -f ${find_links}

rm -rf /tmp/pipbuild
rm -rf ~/.cache
