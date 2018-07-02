#!/bin/bash
set -e
set -o pipefail

arch="amd64"

deb_version="${CPS_PKG_TAG}"
deb_release="${CPS_PKG_REL:-1}"
deb_name="canopsis-core"
deb_path="/root/${deb_name}-${CPS_PKG_TAG}"

cat > ${deb_path}/DEBIAN/control << EOF
Package: ${deb_name}
Architecture: ${arch}
Maintainer: Capensis
Depends: base-files, bash, ca-certificates, libsasl2-2, libxml2, libxslt1.1, lsb-base, libffi6, libgmp10, libgnutlsxx28, libgnutls-openssl27, libhogweed2, libicu52, libidn11, libnettle4, libp11-kit0, libpsl0, libssl1.0.0, libtasn1-6, libxmlsec1, libxmlsec1-openssl, libldap-2.4-2, python, python2.7, rsync, snmp, smitools
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

find /opt/canopsis/etc/ -type f > ${deb_path}/DEBIAN/conffiles
find /opt/canopsis/opt/ -type f >> ${deb_path}/DEBIAN/conffiles

chmod +x ${deb_path}/DEBIAN/*inst

dpkg-deb -b ${deb_path}/ /packages/${deb_name}-${deb_version}-${deb_release}.${arch}.jessie.deb

chown -R ${FIX_OWNERSHIP} /packages/*
