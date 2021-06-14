#!/bin/sh
set -eu

sed -i -e "s,{{ CPS_API_URL }},$CPS_API_URL,g" /etc/nginx/conf.d/default.conf
echo "resolver $(awk 'BEGIN{ORS=" "} $1=="nameserver" {print $2}' /etc/resolv.conf);" > /etc/nginx/resolvers.inc
exec nginx -g 'daemon off;'
