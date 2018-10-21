#!/bin/sed -f
s|{{BASE_DOCKER_IMAGE}}|centos:centos7|g
s|{{BUILD_SETUP}}|yum makecache; \\ \
sed -i /etc/yum/pluginconf.d/fastestmirror.conf -e 's/enabled=.*/enabled=0/g'; \\ \
yum install -y \\ \
    bzip2-devel \\ \
    cyrus-sasl-devel \\ \
    gcc \\ \
    gcc-c++ \\ \
    libcurl-devel \\ \
    libevent-devel \\ \
    libffi-devel \\ \
    libtasn1 \\ \
    libtool \\ \
    libtool-ltdl-devel \\ \
    libxml2-devel \\ \
    libxslt-devel \\ \
    make \\ \
    nettle-devel \\ \
    openldap-devel \\ \
    openssl-devel \\ \
    python-devel \\ \
    python-virtualenv \\ \
    python-wheel \\ \
    sudo \\ \
    xmlsec1-devel \\ \
    xmlsec1-openssl-devel \\ \
    zlib-devel; \\ \
rm -rf /var/cache/yum|g
s|{{FINAL_IMAGE_SETUP}}|yum clean all; \\ \
# \
sed -i /etc/yum/pluginconf.d/fastestmirror.conf -e 's/enabled=.*/enabled=0/g'; \\ \
# \
echo 'LANG=en_US.UTF-8' > /etc/locale.conf; \\ \
rm -f /etc/localtime; \\ \
ln -s /usr/share/zoneinfo/UTC /etc/localtime; \\ \
# \
yum makecache; \\ \
yum install -y epel-release; \\ \
yum makecache; \\ \
# \
yum update -y; \\ \
# \
yum install -y \\ \
    bzip2 \\ \
    bind-utils \\  \
    curl \\ \
    cyrus-sasl \\ \
    htop \\ \
    iputils \\ \
    libcurl \\ \
    libevent \\ \
    libffi \\ \
    libicu \\ \
    libtasn1 \\ \
    libxml2 \\ \
    libxslt \\ \
    nettle \\ \
    net-snmp \\ \
    openldap \\ \
    openssl \\ \
    python \\ \
    redhat-lsb-core \\ \
    rsync \\ \
    sudo \\ \
    tmux \\ \
    xmlsec1 \\ \
    xmlsec1-openssl \\ \
    zlib; \\ \
rm -rf /var/cache/yum|g