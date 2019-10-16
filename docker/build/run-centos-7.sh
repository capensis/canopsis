#!/usr/bin/env bash
set -e
set -o pipefail
set -u

repo_baseurl="http://centos.mirrors.ovh.net/ftp.centos.org/"

yum clean all

sed -i /etc/yum/pluginconf.d/fastestmirror.conf -e 's/enabled=.*/enabled=0/g'

find /etc/yum.repos.d/ -name "CentOS*.repo" -exec sed -e '/^mirrorlist=/d' -i {} \;
find /etc/yum.repos.d/ -name "*.repo" -exec sed -re "s@^#baseurl=http://mirror\.centos\.org/centos/(.*)\$@baseurl=${repo_baseurl}\1@g" -i {} \;

echo 'LANG=en_US.UTF-8' > /etc/locale.conf
rm -f /etc/localtime
ln -s /usr/share/zoneinfo/UTC /etc/localtime

yum makecache
yum install -y epel-release
yum makecache

yum update -y

yum install -y \
    bzip2 \
    bind-utils \
    curl \
    cyrus-sasl \
    htop \
    iputils \
    libcurl \
    libevent \
    libffi \
    libicu \
    libsmi \
    libtasn1 \
    libxml2 \
    libxslt \
    nettle \
    net-snmp \
    net-snmp-utils \
    openldap \
    openssl \
    python \
    redhat-lsb-core \
    rsync \
    sudo \
    tmux \
    xmlsec1 \
    xmlsec1-openssl \
    zlib

rm -rf /var/cache/yum
