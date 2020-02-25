#!/usr/bin/env bash
set -e
set -o pipefail

virtualenv --system-site-packages ${CPS_HOME}/venv-ansible

source /etc/os-release
repver="${ID}-${VERSION_ID}"

source ${CPS_HOME}/venv-ansible/bin/activate
source /etc/os-release

pip install -U setuptools pip

# jmespath is required by ansible role repository adder maat
# pymongo 3.5.0 is required by ansible 2.4.x module mongodb
# influxdb is required for the modified ansible role influxdb
# pyopenssl is required to fix compatibility between systems
pip install -U setuptools pip wheel jmespath "pymongo==3.5.0" "influxdb==5.1.0" "pyOpenSSL" "ansible==2.8.7"

rm -rf /tmp/pipbuild
rm -rf ~/.cache
