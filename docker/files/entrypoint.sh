#!/bin/bash
set -e
set -o pipefail

export HOME=/opt/canopsis

if [ "x$1" == "xshell" ]; then
  exec /bin/bash
else
    sudo /opt/canopsis/bin/env2cfg
	echo $CPS_COMPONENT
    case $CPS_COMPONENT in
        webserver )
			echo "running webserver"
            /opt/canopsis/bin/webserver --access-logfile /opt/canopsis/var/log/webserver-access.log -k gevent -w 1 -b 0.0.0.0:8082
            ;;
        engine )
			echo "running engine $ENGINE_NAME"
            /opt/canopsis/bin/engine-launcher -e $ENGINE_MODULE -n $ENGINE_NAME -w 1 -l ${CPS_LOGGING_LEVEL:-info}
            ;;
        provisionning )
			echo "running provisionning"
            /opt/canopsis/bin/canopsinit
            ;;
        * )
            echo "unknow CPS_COMPONENT $CPS_COMPONENT"
            exit 5
    esac
	exit $?
fi
