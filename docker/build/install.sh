#!/usr/bin/bash
set -e
set -o pipefail

echo "[easy_install]\nallow_hosts = ''\nfind_links = file:///sources/python-libs/" > ${CPS_HOME}/.pydistutils.cfg

source ${CPS_HOME}/bin/activate

pip install --no-index --find-links=file:///sources/python-libs --upgrade setuptools distribute pip
pip install --no-index --find-links=file:///sources/python-libs /sources/canopsis/