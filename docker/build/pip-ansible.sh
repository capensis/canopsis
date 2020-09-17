#!/usr/bin/env bash
set -e
set -o pipefail

virtualenv --system-site-packages ${CPS_HOME}/venv-ansible

source /etc/os-release
repver="${ID}-${VERSION_ID}"

source ${CPS_HOME}/venv-ansible/bin/activate
source /etc/os-release

# duplicate command for force version pip
pip install -U setuptools "pip==20.1.1"

# jmespath is required by ansible role repository adder maat
# pymongo 3.5.0 is required by ansible 2.4.x module mongodb
# influxdb is required for the modified ansible role influxdb
# pyopenssl is required to fix compatibility between systems
# psycopg2-binary is required for postgresql modules in ansible
# NEVER, NEVER, NEVER ⚠️⚠️⚠️ UPGRADE ANSIBLE, OR BE PREPARED FOR PAIN.
pip install -U wheel "jmespath==0.10.0" "pymongo==3.5.0" "influxdb==5.1.0" "pyOpenSSL==19.1.0" "ansible==2.8.7"

rm -rf /tmp/pipbuild
rm -rf ~/.cache
