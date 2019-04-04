#!/bin/bash
set -e
set -o pipefail
set -u

source /etc/os-release

## TODO: facto debian 8 and 9
if [ "$PRETTY_NAME" = "Debian GNU/Linux 8 (stretch)" ]; then
  echo "build debian 8 packages ..."
  PKG_ROOT="/root/canopsis-core-${CANOPSIS_PACKAGE_TAG}"
  deb_path="/root/canopsis-core-${CANOPSIS_PACKAGE_TAG}"
  mkdir -p ${PKG_ROOT}/{DEBIAN,opt}
  mkdir -p ${PKG_ROOT}/lib/systemd/system/
  rsync -aKSH /opt/canopsis ${PKG_ROOT}/opt/
  mkdir -p ${PKG_ROOT}/usr/bin
  ln -sf /opt/canopsis/bin/canoctl ${PKG_ROOT}/usr/bin/canoctl
  cp -ar /lib/systemd/system/canopsis-* ${PKG_ROOT}/lib/systemd/system/
  apt-get -y --no-install-recommends install virtualenv python-pip build-essential libssl-dev libffi-dev python-dev
#  virtualenv --system-site-packages ${CPS_HOME}/venv-ansible
#  source ${CPS_HOME}/venv-ansible/bin/activate
  pyopenssl="pyOpenSSL==17.5.0"
#  pip install -U setuptools==18.5
#  pip install -U pip wheel jmespath "pymongo==3.5.0" "influxdb==5.1.0" "${pyopenssl}" "ansible==2.4.4"
  mkdir ${deb_path}/DEBIAN -p
  cp /packaging/debian8/* ${deb_path}/DEBIAN/
  sed -i ${deb_path}/control -e "s/DEB_VERSION/${CANOPSIS_PACKAGE_TAG}/g"
  sed -i ${deb_path}/control -e "s/DEB_RELEASE/${CANOPSIS_PACKAGE_REL}/g"
  find /opt/canopsis/{etc/,opt/} -type f > ${deb_path}/DEBIAN/conffile
  chmod +x ${deb_path}/DEBIAN/*inst
  cp -rp /opt/canopsis ${deb_path}/opt/
  mkdir -p ${deb_path}/etc/systemd/system/
  mkdir -p ${deb_path}/usr/bin
  ln -sf /opt/canopsis/bin/canoctl ${deb_path}/usr/bin/canoctl
  cp -ar /packaging/systemd-units/canopsis-* ${deb_path}/etc/systemd/system
  dpkg-deb -b ${deb_path}/ /packages/canopsis-core-${CANOPSIS_PACKAGE_TAG}-${CANOPSIS_PACKAGE_REL}.amd64.jessie.deb
  chown -R ${FIX_OWNERSHIP} /packages/*
elif [ "$PRETTY_NAME" = "Debian GNU/Linux 9 (stretch)" ]; then
  echo "build debian 9 packages ..."
  mkdir -p ${PKG_ROOT}/{DEBIAN,opt}
  mkdir -p ${PKG_ROOT}/lib/systemd/system/
  rsync -aKSH /opt/canopsis ${PKG_ROOT}/opt/
  mkdir -p ${PKG_ROOT}/usr/bin
  ln -sf /opt/canopsis/bin/canoctl ${PKG_ROOT}/usr/bin/canoctl
  cp -ar /lib/systemd/system/canopsis-* ${PKG_ROOT}/lib/systemd/system/
  apt-get -y --no-install-recommends install virtualenv python-pip
#  virtualenv --system-site-packages ${CPS_HOME}/venv-ansible
#  source ${CPS_HOME}/venv-ansible/bin/activate
  deb_path="/root/canopsis-core-${CANOPSIS_PACKAGE_TAG}"
#  pip install -U setuptools pip wheel jmespath "pymongo==3.5.0" "influxdb==5.1.0" "${pyopenssl}" "ansible==2.4.4"
  mkdir ${deb_path}/DEBIAN -p
  cp /packaging/debian9/* ${deb_path}/DEBIAN/
  sed -i ${deb_path}/control -e "s/DEB_VERSION/${CANOPSIS_PACKAGE_TAG}/g"
  sed -i ${deb_path}/control -e "s/DEB_RELEASE/${CANOPSIS_PACKAGE_REL}/g"
  find /opt/canopsis/{etc/,opt/} -type f > ${deb_path}/DEBIAN/conffile
  chmod +x ${deb_path}/DEBIAN/*inst
  cp -rp /opt/canopsis  ${deb_path}/opt/
  mkdir -p ${deb_path}/etc/systemd/system/
  mkdir -p ${deb_path}/usr/bin
  ln -sf /opt/canopsis/bin/canoctl ${deb_path}/usr/bin/canoctl
  cp -ar /packaging/systemd-units/canopsis-* ${deb_path}/etc/systemd/system
  dpkg-deb -b ${deb_path}/ /packages/canopsis-core-${CANOPSIS_PACKAGE_TAG}-${CANOPSIS_PACKAGE_REL}.amd64.stretch.deb
  chown -R ${FIX_OWNERSHIP} /packages/*
elif [ "$PRETTY_NAME" = "CentOS Linux 7 (Core)" ]; then
  echo "build centos 7 packages ..."
  head -n 1 /usr/lib/rpm/brp-python-bytecompile > /usr/lib/rpm/brp-python-bytecompile
  echo ${CANOPSIS_PACKAGE_TAG} > /opt/canopsis/VERSION.txt
  cp -r /root/deploy-ansible /opt/canopsis/
  yum makecache
  yum groupinstall -y "Development tools"
  yum install -y yum-utils
  cp /packaging/centos7/centos-7-canopsis-core.spec /root/centos-7-canopsis-core.spec
  sed -i /root/centos-7-canopsis-core.spec -e "s/CANOPSIS_PACKAGE_TAG/${CANOPSIS_PACKAGE_TAG}/g"
  sed -i /root/centos-7-canopsis-core.spec -e "s/CANOPSIS_PACKAGE_REL/${CANOPSIS_PACKAGE_REL}/g"
  yum-builddep -y /root/centos-7-canopsis-core.spec
  rpmbuild -bb /root/centos-7-canopsis-core.spec
  rsync -vrc /root/rpmbuild/RPMS/x86_64/* /packages/
  chown -R ${FIX_OWNERSHIP} /packages/*
else echo "Fail to OS detect !"
fi
