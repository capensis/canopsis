#!/usr/bin/env bash
#
# This file contains buildtime dependencies for Debian docker images.
#
set -e
set -o pipefail
set -u

apt-get -q -o=Dpkg::Use-Pty=0 update
apt-get -q -o=Dpkg::Use-Pty=0 -y --no-install-recommends install \
    checkinstall libreadline-gplv2-dev libncursesw5-dev libssl-dev \
    libsqlite3-dev libgdbm-dev libc6-dev libbz2-dev libffi-dev \
    tk-dev zlib1g-dev
  
py_ver="3.9.16"
curl -O https://www.python.org/ftp/python/$py_ver/Python-$py_ver.tgz
tar xzf Python-$py_ver.tgz

cd Python-$py_ver
./configure --enable-optimizations --prefix=/opt/python/3.9
make -j$(nproc)
make install

update-alternatives --install /usr/bin/python3 python3 /opt/python/3.9/bin/python3.9 60

apt-get -q -o=Dpkg::Use-Pty=0 clean
