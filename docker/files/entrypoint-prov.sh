#!/bin/bash
set -e
set -o pipefail

export HOME=/opt/canopsis

if [ "x$1" == "xshell" ]; then
  exec /bin/bash
else
    sudo /opt/canopsis/bin/env2cfg
    sudo /entrypoint-prov-sync.sh
    /opt/canopsis/bin/canopsinit
fi
