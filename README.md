## AMQP

Default bindings :

Exchange        | Queue | Routing key
--------------- | ----- | -----------
canopsis.events | axe   | #
canopsis.alerts | che   | #

## MongoDB

```
use canopsis
db.createUser({user:"cpsmongo",pwd:"canopsis",roles:["dbOwner"]})
db.periodical_alarm.createIndex({t:1, d:1})
db.periodical_alarm.createIndex({d:1})
```

## Engines

```
export CPS_AMQP_URL="amqp://cpsrabbit:canopsis@localhost/canopsis"
export CPS_MONGO_URL="mongodb://cpsmongo:canopsis@localhost/canopsis"
export CPS_REDIS_URL="redis://nouser:dbpassword@host:port/0"
export CPS_DEFAULT_CFG="$GOPATH/src/git.canopsis.net/canopsis/go-revolution/canopsis/default_configuration.toml"
```

```
cd cmd/engines-axe && go build . && ./engine-axe
cd cmd/engines-che && go build . && ./engine-che
cd cmd/engines-lifeline && go build . && ./engine-lifeline
cd cmd/feeder && go build . && ./feeder
```

## Canopsis Library

```go
import "git.canopsis.net/canopsis/go-revolution/canopsis"
```

## Go Dependencies

Use [glide](https://glide.sh/).

```
glide install
```

## Tests - GoConvey

You will need engines environment variables.

If you want to skip tests with long run times:

```
export CPS_TEST_SKIP_LONG=1
```

Then run `goconvey`:

```
goconvey -workDir ${GOPATH}/src/git.canopsis.net/canopsis/go-revolution/
```

Note: you can enable desktop notifications from the web UI to avoid checking manually.