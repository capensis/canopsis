#!/usr/bin/bash
set -e
set -o pipefail

echo "deb http://ftp.fr.debian.org/debian/ stretch main contrib non-free" > /etc/apt/sources.list
echo "deb http://security.debian.org/ stretch/updates main" >> /etc/apt/sources.list

rm -f /etc/localtime
ln -s /usr/share/zoneinfo/UTC /etc/localtime

apt-get update

apt-get -y --no-install-recommends install locales

echo 'en_US.UTF-8 UTF-8' > /etc/locale.gen
locale-gen

export LANG="en_US.UTF-8"

apt-get -y --no-install-recommends install \
    apt-transport-https \
    base-files \
    bash \
    bash-completion \
    ca-certificates \
    curl \
    libsasl2-2 \
    libxml2 \
    libxslt1.1 \
    lsb-base \
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
    libxmlsec1-openssl \
    libldap-2.4-2 \
    python \
    python2.7 \
    python-pip \
    python-pkg-resources \
    python-virtualenv \
    python-wheel \
    rsync \
    sudo \
    tmux \
    vim \
    virtualenv \

apt-get clean
