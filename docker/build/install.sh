#!/usr/bin/bash
set -e
set -o pipefail

virtualenv ${CPS_HOME}

echo -e "[easy_install]\nallow_hosts = ''\nfind_links = file:///sources/python-libs/" > /root/.pydistutils.cfg

source ${CPS_HOME}/bin/activate

pip install -b /tmp/pipbuild --no-index --find-links=file:///sources/python-libs --upgrade setuptools distribute pip
pip install -b /tmp/pipbuild --no-index --find-links=file:///sources/python-libs /sources/canopsis/

rm -rf /tmp/pipbuild