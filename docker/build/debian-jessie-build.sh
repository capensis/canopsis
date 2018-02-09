#!/bin/sh
echo "deb http://ftp.fr.debian.org/debian/ jessie main contrib non-free" > /etc/apt/sources.list
echo "deb http://security.debian.org/ jessie/updates main" >> /etc/apt/sources.list

apt-get update || exit 1
apt-get -qqy install \
    python \
    apt-transport-https \
    python-pip \
    build-essential \
    base-files \
    bash \
    bash-completion \
    python2.7-dev \
    python-virtualenv \
    libcurl4-openssl-dev \
    libsasl2-dev \
    libxml2-dev \
    libxslt1-dev \
    lsb-core \
    lsb \
    libssl-dev \
    patch \
    ca-certificates \
    libffi6 \
    libffi-dev \
    libgmp10 \
    libgnutls-deb0-28 \
    libhogweed2 \
    libicu52 \
    libidn11 \
    libnettle4 \
    libp11-kit0 \
    libpsl0 \
    libssl1.0.0 \
    libtasn1-6 \
    libxmlsec1 \
    libxmlsec1-dev \
    libldap2-dev || exit 2

apt-get clean || exit 3