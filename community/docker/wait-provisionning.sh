#!/usr/bin/env bash
set -e
set -o pipefail

COMPOSE_FILE="${1}"
PROJECT_NAME="${2}"

while [ ! "$(docker-compose --project-directory CI/ -f CI/${COMPOSE_FILE} -p ${PROJECT_NAME} ps | grep provisionning | grep Up)" = "" ]; do
    echo waiting provisioning end
    sleep 1
done
docker-compose --project-directory CI/ -f CI/${COMPOSE_FILE} -p ${PROJECT_NAME} ps
