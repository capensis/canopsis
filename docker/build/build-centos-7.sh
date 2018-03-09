#!/bin/bash
set -e
set -o pipefail

yum makecache

yum install -y \
    bzip2-devel \
    cyrus-sasl-devel \
    gcc \
    gcc-c++ \
    git \
    libacl-devel \
    libcurl-devel \
    libcurl-devel \
    libevent-devel \
    libffi-devel \
    librsync-devel \
    libtool \
    libtool-ltdl-devel \
    libxml2-devel \
    libxslt-devel \
    libxslt-devel \
    make \
    ncurses-devel \
    openldap-devel \
    openldap-devel \
    openssl \
    openssl-devel \
    perl-ExtUtils-MakeMaker \
    python-devel \
    python-virtualenv
    rsync \
    wget \
    xmlsec1-devel \
    xmlsec1-openssl-devel \
    zlib-devel

rm -rf /var/cache/yum
