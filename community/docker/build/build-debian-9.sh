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
    python3-dev \
    python3-pip \
    python-pkg-resources \
    python3-virtualenv \
    python-pip-whl \
    virtualenv \
    net-tools \
    procps

apt-get -q -o=Dpkg::Use-Pty=0 clean
