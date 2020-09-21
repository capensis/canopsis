#!/bin/sh
sed -i -e "s,{{ TARGET }},$TARGET,g" /etc/nginx/conf.d/default.conf
exec nginx -g 'daemon off;'
