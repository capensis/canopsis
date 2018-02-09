#!/bin/sh
echo "deb http://ftp.fr.debian.org/debian/ jessie main contrib non-free" > /etc/apt/sources.list
echo "deb http://security.debian.org/ jessie/updates main" >> /etc/apt/sources.list

apt-get update || exit 1
apt-get -qqy install \
    build-essential \
    python2.7-dev \
    libcurl4-openssl-dev \
    libsasl2-dev \
    libxml2-dev \
    libxslt1-dev \
    libssl-dev \
    patch \
    libffi-dev \
    libxmlsec1-dev \
    libldap2-dev || exit 2

apt-get clean || exit 3