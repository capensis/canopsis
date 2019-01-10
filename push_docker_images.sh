#!/bin/bash

# Push docker images

echo -n "Push Go images [y/N]: "
read push_go

echo -n "Push PE images too [y/N]: "
read push_pe

if [ "${CANOPSIS_TAG}" = "" ]; then
    echo "No canopsis package tag specified ; using develop..."
    export CANOPSIS_TAG=develop
fi

# Canopsis engines
for engine in {"core","cat","prov","cat-prov"}; do
    docker push canopsis/canopsis-$engine:${CANOPSIS_TAG}
done
# Go engines
if [ "${push_go}" = "Y" ]||[ "${push_go}" = "y" ]; then
    for engine in {"init", "engine-axe","engine-che","engine-heartbeat","engine-stat","engine-watcher","engine-action"}; do
        docker push canopsis/$engine:${CANOPSIS_TAG}
    done
fi

if [ "${push_pe}" = "Y" ]||[ "${push_pe}" = "y" ]; then
    docker push canopsis/init-pe:${CANOPSIS_TAG}
    docker push canopsis/canopsis-cat-pe:${CANOPSIS_TAG}
    docker push canopsis/canopsis-cat-pe:longoutput-${CANOPSIS_TAG}
    #docker push canopsis/canopsis-connector-email2canopsis-pe:${CANOPSIS_TAG}
    #docker push canopsis/canopsis-connector-snmp2canopsis-pe:${CANOPSIS_TAG}
fi
