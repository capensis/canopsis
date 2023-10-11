#!/usr/bin/env bash
set -e
set -o pipefail

DEFAULT_LIMIT_REQUEST_LINE=6800
DEFAULT_WORKERS_SIZE=1

if [ -z ${LIMIT_REQUEST_LINE} ]
then
    LIMIT_REQUEST_LINE=${DEFAULT_LIMIT_REQUEST_LINE}
fi
if [ -z ${WORKERS_SIZE} ]
then
    WORKERS_SIZE=${DEFAULT_WORKERS_SIZE}
fi 

if [ "x$1" == "xshell" ]; then
  exec /bin/bash
else
    if [ "${CPS_OLD_API}" = "1" ]; then
        /opt/canopsis/bin/env2cfg
        /opt/canopsis/bin/canopsis-oldapi --access-logfile /opt/canopsis/var/log/oldapi-access.log -k gevent --limit-request-line ${LIMIT_REQUEST_LINE} -w ${WORKERS_SIZE} -b 0.0.0.0:8081
    else
        # examples: ENGINE_MODULE=canopsis_pro.engines.snmp ENGINE_NAME=snmp CPS_LOGGING_LEVEL=debug
        . /opt/canopsis/bin/activate && python3 -m canopsis_pro.engine.launcher -m $ENGINE_MODULE -n $ENGINE_NAME -l ${CPS_LOGGING_LEVEL:-info}
        exit 1
    fi
fi
