#!/usr/bin/bash
set -e
set -o pipefail

echo "deb http://ftp.fr.debian.org/debian/ stretch main contrib non-free" > /etc/apt/sources.list
echo "deb http://security.debian.org/ stretch/updates main" >> /etc/apt/sources.list

## Core and engines

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
    libldap2-dev \
    python2.7-dev

apt-get clean
