#!/bin/bash
set -e
set -o pipefail

export HOME=/opt/canopsis

if [ "x$1" == "xshell" ]; then
  exec /bin/bash
else
    sudo /opt/canopsis/bin/env2cfg
    sudo /entrypoint-prov-sync.sh
    cmd="/opt/canopsis/bin/canopsinit"
    options=""
    [[ $CPS_STACK && $CPS_EDITION ]] && options="--canopsis-stack $CPS_STACK --canopsis-edition $CPS_EDITION"
    [[ $CPS_STACK && ! $CPS_EDITION ]] && options="--canopsis-stack $CPS_STACK"
    [[ $CPS_EDITION && ! $CPS_STACK ]] && options="--canopsis-edition $CPS_EDITION"
    eval "$cmd $options"
fi
