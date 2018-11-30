#!/usr/bin/env bash
set -e
set -o pipefail

COMPOSE_FILE="${1}"
PROJECT_NAME="${2}"

while [ ! "$(docker-compose -f ${COMPOSE_FILE} -p ${PROJECT_NAME} ps | grep provisionning | grep Up)" = "" ]; do
    echo waiting provisioning end
    sleep 1
done
docker-compose -f ${COMPOSE_FILE} -p ${PROJECT_NAME} ps
