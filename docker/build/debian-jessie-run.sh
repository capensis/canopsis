#!/usr/bin/bash
set -e
set -o pipefail

echo "deb http://ftp.fr.debian.org/debian/ jessie main contrib non-free" > /etc/apt/sources.list
echo "deb http://security.debian.org/ jessie/updates main" >> /etc/apt/sources.list

apt-get update
apt-get -y install \
    apt-transport-https \
    base-files \
    bash \
    bash-completion \
    ca-certificates \
    libsasl2-2 \
    libxml2 \
    libxslt1.1 \
    lsb-core \
    lsb \
    libffi6 \
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
    libldap-2.4-2 \
    python \
    python2.7 \
    python-pip \
    python-virtualenv \
    tmux \
    vim

apt-get clean