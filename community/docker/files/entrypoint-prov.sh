#!/usr/bin/env bash
set -e
set -o pipefail

export HOME=/opt/canopsis

if [ "x$1" == "xshell" ]; then
  exec /bin/bash
else
    sudo /opt/canopsis/bin/env2cfg
    sudo /entrypoint-prov-sync.sh

    case "$CPS_EDITION" in
    core|community)
        options="--canopsis-edition core"
	;;
    cat|pro)
        options="--canopsis-edition cat"
        ;;
    ""|*)
        options=""
        ;;
    esac
    /opt/canopsis/bin/canopsinit $options
fi
