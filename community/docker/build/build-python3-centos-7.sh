#!/usr/bin/env bash
#
set -e
set -o pipefail
set -u

yum --color=never groupinstall "Development Tools" -y
yum --color=never install gcc open-ssl-devel bzip2-devel libffi-devel -y

py_ver="3.9.16"
curl -O https://www.python.org/ftp/python/$py_ver/Python-$py_ver.tgz
tar xzf Python-$py_ver.tgz

cd Python-$py_ver
./configure --enable-optimizations --prefix=/opt/python/3.9
make altinstall

alternatives --install /usr/bin/python3 python3 /opt/python/3.9/bin/python3.9 60

rm -rf /var/cache/yum
