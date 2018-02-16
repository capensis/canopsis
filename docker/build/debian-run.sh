#!/usr/bin/bash
set -e
set -o pipefail

echo "deb http://ftp.fr.debian.org/debian/ stretch main contrib non-free" > /etc/apt/sources.list
echo "deb http://security.debian.org/ stretch/updates main" >> /etc/apt/sources.list

apt-get update
apt-get -y --no-install-recommends install \
    apt-transport-https \
    base-files \
    bash \
    bash-completion \
    ca-certificates \
    libsasl2-2 \
    libxml2 \
    libxslt1.1 \
    lsb-base \
    lsb-release \
    libffi6 \
    libgmp10 \
    libgnutls30 \
    libgnutlsxx28 \
    libgnutls-openssl27 \
    libhogweed4 \
    libicu57 \
    libidn11 \
    libnettle6 \
    libp11-kit0 \
    libpsl5 \
    libssl1.1 \
    libtasn1-6 \
    libxmlsec1 \
    libldap-2.4-2 \
    python \
    python2.7 \
    python-pip \
    python-pkg-resources \
    python-virtualenv \
    python-wheel \
    rsync \
    tmux \
    vim \
    virtualenv \

apt-get clean
