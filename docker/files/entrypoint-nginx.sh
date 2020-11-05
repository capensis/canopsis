#!/bin/sh
sed -i -e "s,{{ CPS_OLD_API_URL }},$CPS_OLD_API_URL,g" /etc/nginx/conf.d/default.conf
exec nginx -g 'daemon off;'
