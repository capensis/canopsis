#!/usr/bin/env bash
set -e
set -o pipefail
set -u

git submodule update --init
./build-docker.sh

# override .env default variable
export CANOPSIS_IMAGE_TAG=${CANOPSIS_DISTRIBUTION}-ci-test

docker exec -t ${COMPOSE_PROJECT_NAME}_webserver_1 /bin/bash -c "source ~/.bashrc && /opt/canopsis/bin/ut_runner /opt/canopsis/test"
docker cp ${COMPOSE_PROJECT_NAME}_webserver_1:/opt/canopsis/tmp/tests_report/ tests_report
