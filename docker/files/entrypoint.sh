#!/bin/bash

if [ "x$1" == "xshell" ]; then
  exec /bin/bash
else
    if [ "${CPS_WEBSERVER}" = "1" ]; then
        GEVENT_RESOLVER=\"ares\" /opt/canopsis/bin/webserver -t 240 --access-logfile ~/var/log/webserver-access.log -k gevent -w 1 -b 0.0.0.0:8082
    else
        /opt/canopsis/bin/engine-launcher -e $ENGINE_MODULE -n $ENGINE_NAME -w 1 -l info
    fi
fi
