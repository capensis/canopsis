#!/bin/sh
set -eu
NGINX_CONFIGURATION_DIRECTORY="/etc/nginx"

if [ "$NGINX_OPENSHIFT" = "1" ]
then
	echo "NGINX runnning in OpenShift-compatible mode"
	NGINX_CONFIGURATION_DIRECTORY="$(mktemp -d)"

	mkdir -p "${NGINX_CONFIGURATION_DIRECTORY}"
	cp -TR /etc/nginx/ "${NGINX_CONFIGURATION_DIRECTORY}"
	rm -rf "${NGINX_CONFIGURATION_DIRECTORY}/ssl"

	sed -i -e "s,/etc/nginx,${NGINX_CONFIGURATION_DIRECTORY},g" "${NGINX_CONFIGURATION_DIRECTORY}/nginx.conf"
	sed -i -e "s,/etc/nginx,${NGINX_CONFIGURATION_DIRECTORY},g" "${NGINX_CONFIGURATION_DIRECTORY}/conf.d/default.conf"
fi

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

sed -i -e "s,{{ CPS_API_URL }},$CPS_API_URL,g" /etc/nginx/conf.d/default.conf
sed -i -e "s,{{ CPS_OLD_API_URL }},$CPS_OLD_API_URL,g" /etc/nginx/conf.d/default.conf
sed -i -e "s,{{ CPS_SERVER_NAME }},$CPS_SERVER_NAME,g" /etc/nginx/conf.d/default.conf
echo "resolver $(awk 'BEGIN{ORS=" "} $1=="nameserver" {print $2}' /etc/resolv.conf) valid=20s;" > /etc/nginx/resolvers.inc
exec nginx -g 'daemon off;'
