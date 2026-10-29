#!/bin/sh
set -eu

SRC_DIR="/etc/nginx"

CPS_ENABLE_HTTPS="${CPS_ENABLE_HTTPS:-}"
NGINX_DISABLE_IPV6="${NGINX_DISABLE_IPV6:-}"
ENABLE_RUNDECK="${ENABLE_RUNDECK:-}"

NGINX_CONFIGURATION_DIRECTORY="$(mktemp -d)"
for entry in "${SRC_DIR}"/*; do
	case "$entry" in */ssl) continue ;; esac
	cp -R "$entry" "${NGINX_CONFIGURATION_DIRECTORY}/"
done
chmod -R u+w "${NGINX_CONFIGURATION_DIRECTORY}"
mkdir -p "${NGINX_CONFIGURATION_DIRECTORY}/ssl"

sed -i -e "s,${SRC_DIR},${NGINX_CONFIGURATION_DIRECTORY},g" \
	"${NGINX_CONFIGURATION_DIRECTORY}/nginx.conf" \
	"${NGINX_CONFIGURATION_DIRECTORY}"/conf.d/*.conf \
	"${NGINX_CONFIGURATION_DIRECTORY}"/*.inc

if [ "$CPS_ENABLE_HTTPS" = "true" ]
then
	sed -i -e "s,#include ${NGINX_CONFIGURATION_DIRECTORY}/https.inc;,include ${NGINX_CONFIGURATION_DIRECTORY}/https.inc;,g" "${NGINX_CONFIGURATION_DIRECTORY}"/conf.d/canopsis.conf

	if [ -r "${SRC_DIR}/ssl/cert.crt" ] && [ -r "${SRC_DIR}/ssl/key.key" ]
	then
		echo "TLS: using the certificate mounted in ${SRC_DIR}/ssl"
		ln -sf "${SRC_DIR}/ssl/cert.crt" "${NGINX_CONFIGURATION_DIRECTORY}/ssl/cert.crt"
		ln -sf "${SRC_DIR}/ssl/key.key"  "${NGINX_CONFIGURATION_DIRECTORY}/ssl/key.key"
	elif [ -e "${SRC_DIR}/ssl/cert.crt" ] || [ -e "${SRC_DIR}/ssl/key.key" ]
	then
		echo "TLS: certificate files found in ${SRC_DIR}/ssl but cert.crt and key.key are not both readable by UID $(id -u), check names, ownership and permissions" >&2
		exit 1
	else
		echo "TLS: no certificate provided in ${SRC_DIR}/ssl, generating a self-signed one"
		openssl req -x509 -nodes -days 730 -newkey rsa:2048 -sha256 \
			-keyout "${NGINX_CONFIGURATION_DIRECTORY}/ssl/key.key" \
			-out "${NGINX_CONFIGURATION_DIRECTORY}/ssl/cert.crt" \
			-subj "/CN=${CPS_SERVER_NAME}"
	fi
fi

if [ "$NGINX_DISABLE_IPV6" = "true" ]
then
	sed -i -e '/listen \[\:\:\]\:.*/d' "${NGINX_CONFIGURATION_DIRECTORY}/conf.d/canopsis.conf" "${NGINX_CONFIGURATION_DIRECTORY}"/https.inc
fi

if [ "$ENABLE_RUNDECK" = "true" ]
then
	sed -i -e "s,#include ${NGINX_CONFIGURATION_DIRECTORY}/rundeck.inc;,include ${NGINX_CONFIGURATION_DIRECTORY}/rundeck.inc;,g" "${NGINX_CONFIGURATION_DIRECTORY}"/conf.d/canopsis.conf
	sed -i -e "s,{{ RUNDECK_GRAILS_URL }},$RUNDECK_GRAILS_URL,g" "${NGINX_CONFIGURATION_DIRECTORY}"/rundeck.inc
	sed -i -e "s,{{ NGINX_URL }},$NGINX_URL,g" "${NGINX_CONFIGURATION_DIRECTORY}"/rundeck.inc
fi

sed -i -e "s,http://127.0.0.1:8082,$CPS_INTERNAL_API_URL,g" "${NGINX_CONFIGURATION_DIRECTORY}"/conf.d/canopsis.conf
sed -i -e "s,http://127.0.0.1:8083,$CPS_PUBLIC_API_URL,g" "${NGINX_CONFIGURATION_DIRECTORY}"/conf.d/canopsis.conf
sed -i -e "s,\$canopsis_server_name localhost,\$canopsis_server_name $CPS_SERVER_NAME,g" "${NGINX_CONFIGURATION_DIRECTORY}"/conf.d/canopsis.conf
echo "resolver $(awk 'BEGIN{ORS=" "} $1=="nameserver" {print $2}' /etc/resolv.conf) valid=20s;" > "${NGINX_CONFIGURATION_DIRECTORY}"/resolvers.inc

exec nginx -g 'daemon off;' -c "${NGINX_CONFIGURATION_DIRECTORY}/nginx.conf"