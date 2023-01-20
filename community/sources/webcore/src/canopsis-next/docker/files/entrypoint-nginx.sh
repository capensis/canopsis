#!/bin/sh
set -eu

NGINX_CONFIGURATION_DIRECTORY="/etc/nginx"

NGINX_OPENSHIFT="${NGINX_OPENSHIFT:-}"
CPS_ENABLE_HTTPS="${CPS_ENABLE_HTTPS:-}"
NGINX_DISABLE_IPV6="${NGINX_DISABLE_IPV6:-}"
ENABLE_RUNDECK="${ENABLE_RUNDECK:-}"

if [ "$NGINX_OPENSHIFT" = "true" ]
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
	sed -i -e "s,#include ${NGINX_CONFIGURATION_DIRECTORY}/https.inc;,include ${NGINX_CONFIGURATION_DIRECTORY}/https.inc;,g" "${NGINX_CONFIGURATION_DIRECTORY}"/conf.d/default.conf
	if [ ! -f /etc/nginx/ssl/cert.crt ] && [ ! -f /etc/nginx/ssl/key.key ]
	then
		openssl req -x509 -nodes -days 730 -newkey rsa:2048 -sha256 \
			-keyout /etc/nginx/ssl/key.key \
			-out /etc/nginx/ssl/cert.crt \
			-subj "/CN=${CPS_SERVER_NAME}"
	fi
fi

if [ "$NGINX_DISABLE_IPV6" = "true" ]
then
	sed -i -e '/listen \[\:\:\]\:.*/d' "${NGINX_CONFIGURATION_DIRECTORY}"/conf.d/default.conf
fi

if [ "$ENABLE_RUNDECK" = "true" ]
then
	sed -i -e "s,#include ${NGINX_CONFIGURATION_DIRECTORY}/rundeck.inc;,include ${NGINX_CONFIGURATION_DIRECTORY}/rundeck.inc;,g" "${NGINX_CONFIGURATION_DIRECTORY}"/conf.d/default.conf
        sed -i -e "s,{{ RUNDECK_GRAILS_URL }},$RUNDECK_GRAILS_URL,g" "${NGINX_CONFIGURATION_DIRECTORY}"/rundeck.inc
        sed -i -e "s,{{ NGINX_URL }},$NGINX_URL,g" "${NGINX_CONFIGURATION_DIRECTORY}"/rundeck.inc
fi

sed -i -e "s,http://127.0.0.1:8082,$CPS_API_URL,g" "${NGINX_CONFIGURATION_DIRECTORY}"/conf.d/default.conf
sed -i -e "s,\$canopsis_server_name localhost,\$canopsis_server_name $CPS_SERVER_NAME,g" "${NGINX_CONFIGURATION_DIRECTORY}"/conf.d/default.conf
echo "resolver $(awk 'BEGIN{ORS=" "} $1=="nameserver" {print $2}' /etc/resolv.conf) valid=20s;" > "${NGINX_CONFIGURATION_DIRECTORY}"/resolvers.inc

exec nginx -g 'daemon off;' -c "${NGINX_CONFIGURATION_DIRECTORY}/nginx.conf"
