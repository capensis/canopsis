#!/bin/sh
sed -i -e "s,{{ CPS_API_URL }},$CPS_API_URL,g" /etc/nginx/conf.d/default.conf
echo resolver $(awk 'BEGIN{ORS=" "} $1=="nameserver" {print $2}' /etc/resolv.conf) ";" > /etc/nginx/resolvers.conf
exec nginx -g 'daemon off;'
