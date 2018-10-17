#!/bin/bash
set -e
set -o pipefail

export HOME=/opt/canopsis

if [ "x$1" == "xshell" ]; then
  exec /bin/bash
else
    sudo /opt/canopsis/bin/env2cfg
    if [ "${CPS_WEBSERVER}" = "1" ]; then
        /opt/canopsis/bin/webserver --access-logfile /opt/canopsis/var/log/webserver-access.log -k gevent -w 1 -b 0.0.0.0:8082
    else
        /opt/canopsis/bin/engine-launcher -e $ENGINE_MODULE -n $ENGINE_NAME -w 1 -l ${CPS_LOGGING_LEVEL:-info}
    fi
fi
