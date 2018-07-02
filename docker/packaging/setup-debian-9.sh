#!/bin/bash
set -e
set -o pipefail

mkdir -p /root/canopsis-core-${CPS_PKG_TAG}/{DEBIAN,opt}
mkdir -p /root/canopsis-core-${CPS_PKG_TAG}/lib/systemd/system/
rsync -aKSH /opt/canopsis /root/canopsis-core-${CPS_PKG_TAG}/opt/
cp -ar /lib/systemd/system/canopsis-* /root/canopsis-core-${CPS_PKG_TAG}/lib/systemd/system/
