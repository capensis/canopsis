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
    /opt/canopsis/bin/env2cfg
    if [ "${CPS_OLD_API}" = "1" ]; then
        /opt/canopsis/bin/canopsis-oldapi --access-logfile /opt/canopsis/var/log/oldapi-access.log -k gevent --limit-request-line ${LIMIT_REQUEST_LINE} -w ${WORKERS_SIZE} -b 0.0.0.0:8081
    else
        echo "Launching Python engines is no longer supported" >&2
        exit 1
    fi
fi
