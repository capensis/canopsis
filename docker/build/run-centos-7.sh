#!/bin/bash
set -e
set -o pipefail

yum clean all

sed -i /etc/yum/pluginconf.d/fastestmirror.conf -e 's/enabled=.*/enabled=0/g'
find /etc/yum.repos.d/ -name "CentOS*.repo" -exec sed -e '/^mirrorlist=/d' -i {} \;
find /etc/yum.repos.d/ -name "CentOS*.repo" -exec sed -re 's/^#baseurl(.*)/baseurl\1/g' -i {} \;
echo 'LANG=en_US.UTF-8' > /etc/locale.conf
rm -f /etc/localtime
ln -s /usr/share/zoneinfo/UTC /etc/localtime

yum makecache
yum install -y epel-release
yum makecache

yum install -y \
    bzip2 \
    curl \
    cyrus-sasl \
    htop \
    libcurl \
    libevent \
    libffi \
    libicu \
    libtasn1 \
    libxml2 \
    libxslt \
    nettle \
    openldap \
    openssl \
    python \
    python-virtualenv \
    python-wheel \
    redhat-lsb-core \
    rsync \
    tmux \
    xmlsec1 \
    xmlsec1-openssl \
    zlib

rm -rf /var/cache/yum
