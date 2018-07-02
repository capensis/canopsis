#!/bin/bash
set -e
set -o pipefail

echo "#!/bin/bash" > /usr/lib/rpm/brp-python-bytecompile

sed -i /packaging/centos-7-canopsis-cat.spec -e "s/CPS_PKG_TAG/${CPS_PKG_TAG}/g"
sed -i /packaging/centos-7-canopsis-cat.spec -e "s/CPS_PKG_REL/${CPS_PKG_REL}/g"

yum-builddep -y /packaging/centos-7-canopsis-cat.spec

rpmbuild -bb /packaging/centos-7-canopsis-cat.spec

rsync -vrc /root/rpmbuild/RPMS/* /packages/

chown -R ${FIX_OWNERSHIP} /packages/*
