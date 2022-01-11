#!/usr/bin/env bash
set -e
set -o pipefail
set -u

../build-docker.sh

# override .env default variable
export CANOPSIS_IMAGE_TAG=${CANOPSIS_DISTRIBUTION}-ci-test

docker-compose --project-directory CI/ -f CI/docker-compose.ci.yml -p ${COMPOSE_PROJECT_NAME} up -d webserver mongodb redis rabbitmq provisionning
./docker/wait-provisionning.sh CI/docker-compose.ci.yml ${COMPOSE_PROJECT_NAME}
docker-compose --project-directory CI/ -f CI/docker-compose.ci.yml -p ${COMPOSE_PROJECT_NAME} restart webserver mongodb redis rabbitmq provisionning

docker exec -t ${COMPOSE_PROJECT_NAME}_webserver_1 /bin/bash -c "source ~/.bashrc && /opt/canopsis/bin/ut_runner /opt/canopsis/test"
docker cp ${COMPOSE_PROJECT_NAME}_webserver_1:/opt/canopsis/tmp/tests_report/ tests_report
