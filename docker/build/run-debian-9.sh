#!/usr/bin/bash
set -e
set -o pipefail
set -u

echo "deb http://ftp.fr.debian.org/debian/ stretch main contrib non-free" > /etc/apt/sources.list
echo "deb http://security.debian.org/ stretch/updates main" >> /etc/apt/sources.list

rm -f /etc/localtime
ln -s /usr/share/zoneinfo/UTC /etc/localtime

apt-get update

apt-get dist-upgrade -y

apt-get -y --no-install-recommends install locales

export LANG="en_US.UTF-8"
echo "LANG=${LANG}" > /etc/locale.conf
echo "${LANG} UTF-8" > /etc/locale.gen
locale-gen

apt-get -y --no-install-recommends install \
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
    libpsl5 \
    libssl1.1 \
    libtasn1-6 \
    libxmlsec1 \
    libxmlsec1-openssl \
    libldap-2.4-2 \
    python \
    python2.7 \
    rsync \
    snmp \
    snmp-mibs-downloader \
    smitools \
    sudo \
    tmux \
    vim

apt-get clean
