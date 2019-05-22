#!/bin/bash
set -e
set -o pipefail

export HOME=/opt/canopsis

if [ "x$1" == "xshell" ]; then
  exec /bin/bash
else
    sudo /opt/canopsis/bin/env2cfg
    sudo /entrypoint-prov-sync.sh
    options=""
    [[ "$CPS_STACK" ]] && options="$options --canopsis-stack $CPS_STACK"
    [[ "$CPS_EDITION" ]] && options="$options --canopsis-edition $CPS_EDITION"
    /opt/canopsis/bin/canopsinit $options
fi
