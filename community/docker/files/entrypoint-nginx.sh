#!/bin/sh
set -eu

if [ "$CPS_ENABLE_HTTPS" = "true" ]
then
	sed -i -e "s,#include /etc/nginx/https.inc;,include /etc/nginx/https.inc;,g" /etc/nginx/conf.d/default.conf
	if [ ! -f /etc/nginx/ssl/cert.crt ] && [ ! -f /etc/nginx/ssl/key.key ]
	then
		openssl req -x509 -nodes -days 730 -newkey rsa:2048 -sha256 \
			-keyout /etc/nginx/ssl/key.key \
			-out /etc/nginx/ssl/cert.crt \
			-subj "/CN=${CPS_SERVER_NAME}"
	fi
fi

if [ "$ENABLE_RUNDECK" = "true" ]
then
	sed -i -e "s,#include /etc/nginx/rundeck.inc;,include /etc/nginx/rundeck.inc;,g" /etc/nginx/conf.d/default.conf
fi

sed -i -e "s,{{ CPS_API_URL }},$CPS_API_URL,g" /etc/nginx/conf.d/default.conf
sed -i -e "s,{{ CPS_OLD_API_URL }},$CPS_OLD_API_URL,g" /etc/nginx/conf.d/default.conf
sed -i -e "s,{{ CPS_SERVER_NAME }},$CPS_SERVER_NAME,g" /etc/nginx/conf.d/default.conf
sed -i -e "s,{{ RUNDECK_GRAILS_URL }},$RUNDECK_GRAILS_URL,g" /etc/nginx/rundeck.inc
sed -i -e "s,{{ NGINX_URL }},$NGINX_URL,g" /etc/nginx/rundeck.inc
echo "resolver $(awk 'BEGIN{ORS=" "} $1=="nameserver" {print $2}' /etc/resolv.conf) valid=20s;" > /etc/nginx/resolvers.inc
exec nginx -g 'daemon off;'
