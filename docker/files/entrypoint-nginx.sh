#!/bin/sh
sed -i -e "s,___TARGET___,$TARGET,g" /etc/nginx/conf.d/default.conf
exec nginx -g 'daemon off;'
