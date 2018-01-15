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

In `canopsis`, use [glide](https://glide.sh/).

```
cd canopsis
glide update
```

## GoConvey

```
$GOPATH/bin/goconvey -workDir ~/go/src/git.canopsis.net/canopsis/go-revolution/
```
