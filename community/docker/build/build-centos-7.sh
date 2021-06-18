#!/usr/bin/env bash
#
# This file contains buildtime dependencies for CentOS docker images.
#
set -e
set -o pipefail
set -u

yum --color=never makecache

yum --color=never install -y \
    bzip2-devel \
    cyrus-sasl-devel \
    gcc \
    gcc-c++ \
    libcurl-devel \
    libevent-devel \
    libffi-devel \
    libtasn1 \
    libtool \
    libtool-ltdl-devel \
    libxml2-devel \
    libxslt-devel \
    make \
    nettle-devel \
    openldap-devel \
    openssl-devel \
    postgresql \
    postgresql-devel \
    postgresql-libs \
    python-devel \
    python-virtualenv \
    python-wheel \
    sudo \
    xmlsec1-devel \
    xmlsec1-openssl-devel \
    zlib-devel

rm -rf /var/cache/yum
