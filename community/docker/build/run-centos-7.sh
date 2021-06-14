#!/usr/bin/env bash
#
# This file contains runtime dependencies for CentOS docker images.
#
set -e
set -o pipefail
set -u

yum clean metadata
yum clean all
rm -rf /var/cache/yum/*

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
    postgresql-libs \
    python \
    redhat-lsb-core \
    rsync \
    sudo \
    tmux \
    xmlsec1 \
    xmlsec1-openssl \
    zlib

rm -rf /var/cache/yum
