#!/bin/bash
set -e
set -o pipefail

yum clean all

sed -i /etc/yum/pluginconf.d/fastestmirror.conf -e 's/enabled=.*/enabled=0/g'
find /etc/yum.repos.d/ -name "CentOS*.repo" -exec sed -e '/^mirrorlist=/d' -i {} \;
find /etc/yum.repos.d/ -name "CentOS*.repo" -exec sed -re 's/^#baseurl(.*)/baseurl\1/g' -i {} \;
echo 'LANG=en_US.UTF-8' > /etc/locale.conf

yum makecache
yum install -y epel-release
yum makecache

yum install -y xz zlib-devel libevent-devel libevent ncurses-devel libcurl-devel curl tmux htop libtool openssl openssl-devel bzip2-devel cyrus-sasl-devel openldap-devel libcurl-devel python-devel openldap-devel libxml2-devel libxslt-devel perl-ExtUtils-MakeMaker git rsync librsync-devel uthash-devel.noarch libacl-devel libxslt-devel libffi-devel xmlsec1-devel xmlsec1-openssl-devel libtool-ltdl-devel python-virtualenv

rm -rf /var/cache/yum
