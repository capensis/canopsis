#!/usr/bin/env bash
#
# This file contains buildtime dependencies for Debian docker images.
#
set -e
set -o pipefail
set -u

apt-get -q -o=Dpkg::Use-Pty=0 update
apt-get -q -o=Dpkg::Use-Pty=0 -y --no-install-recommends install \
    build-essential \
    curl \
    libcurl4-openssl-dev \
    libsasl2-dev \
    libxml2-dev \
    libxslt1-dev \
    libssl-dev \
    libffi-dev \
    libpq-dev \
    libxmlsec1-dev \
    libxmlsec1-openssl \
    libldap2-dev \
    pkg-config \
    python2.7-dev \
    python-pip \
    python-pkg-resources \
    python-virtualenv \
    python-wheel \
    virtualenv \
    net-tools \
    procps

apt-get -q -o=Dpkg::Use-Pty=0 clean
