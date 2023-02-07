#!/usr/bin/env bash
set -e
set -o pipefail

# Make sure that virtualenv doesn't auto-upgrade pip, because
# recent versions have problems with our old Python 2 setup.
export PATH=/opt/python/3.9/bin:$PATH
virtualenv -p python3 --no-download --system-site-packages ${CPS_HOME}/venv-ansible

source ${CPS_HOME}/venv-ansible/bin/activate

pip3 install -U setuptools "pip==23.0"
pip3 --no-color install -U wheel

# pymongo is required by the Ansible mongodb/canopsis roles
# psycopg2 is required y the Ansible canopsis role (for PostgreSQL)
# pyopenssl is required to fix compatibility between systems
# crytography 2.9.2 to avoid Python 2 warning, for now
# NEVER, NEVER, NEVER ⚠️⚠️⚠️ UPGRADE ANSIBLE, OR BE PREPARED FOR PAIN.
cat > tmprequirements.txt << EOF
ansible==2.8.7
autopep8==1.7.0
cryptography==2.9.2
dnspython==2.2.1
pika==1.3.1
psycopg2==2.8.4
pybars3==0.9.7
pycodestyle==2.9.1
PyMeta3==0.5.1
pymongo==4.3.2
pyOpenSSL==19.1.0
pytz==2022.4
toml==0.10.2
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
PyYAML==5.3.1
requests==2.24.0
six==1.15.0
urllib3==1.25.11
EOF
pip3 --no-color install --no-use-pep517 -U -r ./tmprequirements.txt

rm -rf /tmp/pipbuild ~/.cache tmprequirements.txt
