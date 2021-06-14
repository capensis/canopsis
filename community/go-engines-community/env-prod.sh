#!/bin/sh
if [ "${PB_CUSTOM_ENV}" = "" ]; then
    export PB_CUSTOM_ENV="go:prod"
fi

export CPS_MONGO_URL="mongodb://cpsmongo:canopsis@localhost:27027/canopsis"
export CPS_AMQP_URL="amqp://cpsrabbit:canopsis@localhost:5672/canopsis"
export CPS_REDIS_URL="redis://nouser@localhost:6379/0"
export CPS_INFLUX_URL="influxdb://cpsinflux:canopsis@localhost:28086/canopsis"
