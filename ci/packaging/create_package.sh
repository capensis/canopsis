#!/bin/bash
set -e
set -o pipefail
set -u

source /etc/os-release

if [ "$PRETTY_NAME" = "Debian GNU/Linux 9 (stretch)" ]; then
  echo "build debian 9 packages ..."
elif [ "$PRETTY_NAME" = "Debian GNU/Linux 8 (stretch)" ]; then
  echo "build debian 8 packages ..."
elif [ "$PRETTY_NAME" = "CentOS Linux 7 (Core)" ]; then
  echo "build centos 7 packages ..."
  yum makecache
  yum groupinstall -y "Development tools"
  yum install -y yum-utils
  sed -i /packaging/centos7/centos-7-canopsis-core.spec -e "s/CANOPSIS_PACKAGE_TAG/${CANOPSIS_PACKAGE_TAG}/g"
  sed -i /packaging/centos7/centos-7-canopsis-core.spec -e "s/CANOPSIS_PACKAGE_REL/${CANOPSIS_PACKAGE_REL}/g"
  yum-builddep -y /packaging/centos7/centos-7-canopsis-core.spec
  rpmbuild -bb /packaging/centos7/centos-7-canopsis-core.spec
  rsync -vrc /root/rpmbuild/RPMS/* /packages/
else echo "Fail to OS detect !"
fi
