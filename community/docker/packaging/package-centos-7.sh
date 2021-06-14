#!/usr/bin/env bash
set -e
set -o pipefail
set -u

echo "#!/bin/bash" > /usr/lib/rpm/brp-python-bytecompile

sed -i /packaging/centos-7-canopsis-core.spec -e "s/CANOPSIS_PACKAGE_TAG/${CANOPSIS_PACKAGE_TAG}/g"
sed -i /packaging/centos-7-canopsis-core.spec -e "s/CANOPSIS_PACKAGE_REL/${CANOPSIS_PACKAGE_REL}/g"

yum-builddep -y /packaging/centos-7-canopsis-core.spec

rpmbuild -bb /packaging/centos-7-canopsis-core.spec

rsync -vrc /root/rpmbuild/RPMS/* /packages/

chown -R ${FIX_OWNERSHIP} /packages/*
