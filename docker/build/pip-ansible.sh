#!/usr/bin/env bash
set -e
set -o pipefail

# Make sure that virtualenv doesn't auto-upgrade pip, because
# recent versions have problems with our old Python 2 setup.
virtualenv --no-download --system-site-packages ${CPS_HOME}/venv-ansible

source ${CPS_HOME}/venv-ansible/bin/activate

pip install -U setuptools "pip==20.1.1"
pip install -U wheel

# jmespath is required by ansible role repository adder maat
# pymongo 3.5.0 is required by ansible 2.4.x module mongodb
# influxdb is required for the modified ansible role influxdb
# pyopenssl is required to fix compatibility between systems
# crytography 2.9.2 to avoid Python 2 warning, for now
# NEVER, NEVER, NEVER ⚠️⚠️⚠️ UPGRADE ANSIBLE, OR BE PREPARED FOR PAIN.
cat > tmprequirements.txt << EOF
jmespath==0.10.0
pymongo==3.5.0
influxdb==5.1.0
pyOpenSSL==19.1.0
cryptography==2.9.2
ansible==2.8.7
# indirect dependencies follow
certifi==2020.6.20
cffi==1.14.3
chardet==3.0.4
enum34==1.1.10
idna==2.10
ipaddress==1.0.23
Jinja2==2.11.2
MarkupSafe==1.1.1
pycparser==2.20
python-dateutil==2.8.1
pytz==2020.1
PyYAML==5.3.1
requests==2.24.0
six==1.15.0
urllib3==1.25.11
EOF
pip install --no-use-pep517 -U -r ./tmprequirements.txt

rm -rf /tmp/pipbuild ~/.cache tmprequirements.txt
