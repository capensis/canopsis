#!/usr/bin/bash
set -e
set -o pipefail

virtualenv --system-site-packages ${CPS_HOME}/venv-ansible

source /etc/os-release
repver="${ID}-${VERSION_ID}"

source ${CPS_HOME}/venv-ansible/bin/activate
find_links="file:///sources/wheels/${repver}"

pip install -U setuptools pip wheel
pip install -b /tmp/pipbuild --no-index -f ${find_links} "ansible==2.4.4"

rm -rf /tmp/pipbuild
rm -rf ~/.cache
