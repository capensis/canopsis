#!/usr/bin/bash
set -e
set -o pipefail

apt-get update
apt-get -y --no-install-recommends install \
    build-essential \
    curl \
    libcurl4-openssl-dev \
    libsasl2-dev \
    libxml2-dev \
    libxslt1-dev \
    libssl-dev \
    libffi-dev \
    libxmlsec1-dev \
    libxmlsec1-openssl \
    libldap2-dev \
    python2.7-dev

apt-get clean
