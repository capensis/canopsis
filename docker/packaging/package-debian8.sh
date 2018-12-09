#!/bin/bash
set -e
set -o pipefail
# set -u

arch="amd64"

apt-get -y --no-install-recommends install virtualenv python-pip build-essential libssl-dev libffi-dev python-dev

echo "*** plop"
virtualenv --system-site-packages ${CPS_HOME}/venv-ansible
source ${CPS_HOME}/venv-ansible/bin/activate

deb_version="${CANOPSIS_PACKAGE_TAG}"
deb_release="${CANOPSIS_PACKAGE_REL}"
deb_name="canopsis-core"
deb_path="/root/${deb_name}-${CANOPSIS_PACKAGE_TAG}"

source /etc/os-release
repver="${ID}-${VERSION_ID}"

source /etc/os-release

pyopenssl="pyOpenSSL"
if [ "${VERSION_ID}" = "8" ]&&[ "${ID}" = "debian" ]; then
    pyopenssl="pyOpenSSL==17.5.0"
fi
echo "*** coucou"
pip install -U setuptools==18.5
pip install -U pip wheel jmespath "pymongo==3.5.0" "influxdb==5.1.0" "${pyopenssl}" "ansible==2.4.4"

mkdir ${deb_path}/DEBIAN -p

cat > ${deb_path}/DEBIAN/control << EOF
Package: ${deb_name}
Architecture: ${arch}
Maintainer: Capensis
Depends: base-files, bash, ca-certificates, libsasl2-2, libxml2, libxslt1.1, lsb-base, libffi6, libgmp10, libgnutlsxx28, libgnutls-openssl27, libhogweed2, libicu52, libidn11, libnettle4, libp11-kit0, libpsl0, libssl1.0.0, libtasn1-6, libxmlsec1, libxmlsec1-openssl, libldap-2.4-2, python, python2.7, rsync, snmp, smitools, sudo
Version: ${deb_version}-${deb_release}
Description: Canopsis with CAT package.
EOF

cat > ${deb_path}/DEBIAN/preinst << EOF
#!/bin/bash
grep canopsis /etc/passwd
if [ ! "\$?" = "0" ]; then
    useradd -d /opt/canopsis -M -s /bin/bash canopsis
fi
EOF

cat > ${deb_path}/DEBIAN/postinst << EOF
#!/bin/bash
chmod +x /opt/canopsis/bin/*
chown -R canopsis:canopsis /opt/canopsis/var/log
chown -R canopsis:canopsis /opt/canopsis/var/cache
chown -R canopsis:canopsis /opt/canopsis/tmp
EOF

find /opt/canopsis/{etc/,opt/} -type f > ${deb_path}/DEBIAN/conffile

chmod +x ${deb_path}/DEBIAN/*inst

cp -rp /opt/canopsis  ${deb_path}/opt/
mkdir -p ${deb_path}/etc/systemd/system/
mkdir -p ${deb_path}/usr/bin
ln -sf /opt/canopsis/bin/canoctl ${deb_path}/usr/bin/canoctl

cp -ar /packages/*.service ${deb_path}/etc/systemd/system

dpkg-deb -b ${deb_path}/ /packages/${deb_name}-${deb_version}-${deb_release}.${arch}.jessie.deb

chown -R ${FIX_OWNERSHIP} /packages/*
