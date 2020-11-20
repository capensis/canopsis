#!/bin/sh
set -x
function wait_for_backends {
        start_time=$(date +%s)
        timeout=${TIMEOUT:-30}
        until wget --spider ${@} -O /dev/null;
        do 
                echo "Waiting for backends ${@}";
                sleep 1
                
                # Check for timeout
                if [ $(($(date +%s) - start_time)) -gt ${timeout} ]
                then
                        echo "Backends not responding after $(($(date +%s) - start_time))s, exiting"
                        exit 1
                fi
        done
}

function fix_permissions {
        for f in ${@}
        do
                echo Fixing permissions for ${f}
                chmod 0755 -R ${f}
        done
}

sed -i -e "s,{{ TARGET }},$TARGET,g" /etc/nginx/conf.d/default.conf

wait_for_backends ${TARGET} ${ADDITIONAL_BACKENDS}
fix_permissions ${NGINX_ACCESSIBLE_FOLDERS}

exec nginx -g 'daemon off;'
