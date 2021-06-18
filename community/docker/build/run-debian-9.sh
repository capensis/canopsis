#!/usr/bin/env bash
#
# This file contains runtime dependencies for Debian docker images.
#
set -e
set -o pipefail
set -u

echo "deb http://deb.debian.org/debian/ stretch main contrib non-free" > /etc/apt/sources.list
echo "deb http://security.debian.org/ stretch/updates main" >> /etc/apt/sources.list

rm -f /etc/localtime
ln -s /usr/share/zoneinfo/UTC /etc/localtime

apt-get -q -o=Dpkg::Use-Pty=0 update

apt-get -q -o=Dpkg::Use-Pty=0 dist-upgrade -y

apt-get -q -o=Dpkg::Use-Pty=0 -y --no-install-recommends install locales

export LANG="en_US.UTF-8"
export LC_ALL="$LANG"
echo "LANG=${LANG}" > /etc/locale.conf
echo "${LANG} UTF-8" > /etc/locale.gen
locale-gen

apt-get -q -o=Dpkg::Use-Pty=0 -y --no-install-recommends install \
    apt-transport-https \
    base-files \
    bash \
    ca-certificates \
    curl \
    dnsutils \
    iputils-ping \
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
    libpq5 \
    libpsl5 \
    libssl1.1 \
    libtasn1-6 \
    libxmlsec1 \
    libxmlsec1-openssl \
    libldap-2.4-2 \
    python \
    python2.7 \
    procps \
    rsync \
    snmp \
    snmp-mibs-downloader \
    smitools \
    sudo \
    tmux \
    vim

apt-get -q -o=Dpkg::Use-Pty=0 clean
