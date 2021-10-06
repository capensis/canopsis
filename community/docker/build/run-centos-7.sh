#!/usr/bin/env bash
#
# This file contains runtime dependencies for CentOS docker images.
#
set -e
set -o pipefail
set -u

yum --color=never clean metadata
yum --color=never clean all
rm -rf /var/cache/yum/*

echo 'LANG=en_US.UTF-8' > /etc/locale.conf
rm -f /etc/localtime
ln -s /usr/share/zoneinfo/UTC /etc/localtime

# Make sure Let's Encrypt's previous root certificate won't be used by the buggy
# OpenSSL version on CentOS 7. Otherwise, Yum and EPEL will fail even MORE.
trust dump --filter "pkcs11:id=%c4%a7%b1%a4%7b%2c%71%fa%db%e1%4b%90%75%ff%c4%15%60%85%89%10" | openssl x509 > /etc/pki/ca-trust/source/blacklist/DST-Root-CA-X3.pem
update-ca-trust extract
yum --color=never makecache
yum --color=never -y update ca-certificates || true
yum --color=never update -y

yum --color=never install -y epel-release
yum --color=never makecache

yum --color=never install -y \
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
