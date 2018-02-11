#!/bin/bash
set -e
set -o pipefail

echo "deb http://ftp.fr.debian.org/debian/ jessie main contrib non-free" > /etc/apt/sources.list
echo "deb http://security.debian.org/ jessie/updates main" >> /etc/apt/sources.list

## Core and engines

apt-get update || exit 1
apt-get -qqy install \
    build-essential \
    curl \
    git-core \
    libcurl4-openssl-dev \
    libsasl2-dev \
    libxml2-dev \
    libxslt1-dev \
    libssl-dev \
    libffi-dev \
    libxmlsec1-dev \
    libldap2-dev \
    patch \
    python2.7-dev  || exit 2


## Webserver
#curl -sL https://deb.nodesource.com/setup_6.x | bash -

#apt-get install -y nodejs wget

apt-get clean || exit 4