#!/bin/bash
set -e
set -o pipefail
set -u

PKG_ROOT="/root/canopsis-core-${CANOPSIS_PACKAGE_TAG}"

mkdir -p ${PKG_ROOT}/{DEBIAN,opt}
mkdir -p ${PKG_ROOT}/lib/systemd/system/
rsync -aKSH /opt/canopsis ${PKG_ROOT}/opt/

mkdir -p ${PKG_ROOT}/usr/bin
ln -sf /opt/canopsis/bin/canoctl ${PKG_ROOT}/usr/bin/canoctl
cp -ar /lib/systemd/system/canopsis-* ${PKG_ROOT}/lib/systemd/system/
