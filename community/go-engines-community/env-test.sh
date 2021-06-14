#!/bin/sh
if [ "${PB_CUSTOM_ENV}" = "" ]; then
    export PB_CUSTOM_ENV="go:test"
fi

export CPS_MONGO_URL="mongodb://cpsmongo:canopsis@localhost:37027/canopsis"
export CPS_AMQP_URL="amqp://cpsrabbit:canopsis@localhost:5673/canopsis"
export CPS_REDIS_URL="redis://nouser@localhost:26379/0"
export CPS_INFLUX_URL="influxdb://cpsinflux:canopsis@localhost:28086/canopsis"
