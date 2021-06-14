#!/usr/bin/env bash
set -e
set -o pipefail

cmd="docker-compose exec event_filter"
if [ ! "${1}" = "" ]; then
    cmd="docker-compose -p ${1} exec event_filter"
fi

${cmd} /opt/canopsis/bin/schema2db
${cmd} /opt/canopsis/bin/canopsis-filldb --update
${cmd} /opt/canopsis/bin/canopsis-filldb --update
