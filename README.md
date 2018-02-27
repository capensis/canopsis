## AMQP

Default bindings :

Exchange        | Queue       | Routing key
--------------- | ----------- | -----------
canopsis.events | axe         | #
canopsis.alerts | che         | #
canopsis.events | engine_stat | #

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
export CPS_INFLUX_URL="http://cpsinflux:canopsis@host:8086"
export CPS_DEFAULT_CFG="$GOPATH/src/git.canopsis.net/canopsis/go-revolution/canopsis/default_configuration.toml"
```

```
cd cmd/engines-axe && go build . && ./engine-axe
cd cmd/engines-che && go build . && ./engine-che
cd cmd/engines-lifeline && go build . && ./engine-lifeline
cd cmd/engines-stat && go build . && ./engine-stat
cd cmd/feeder && go build . && ./feeder
```

## Canopsis Library

```go
import "git.canopsis.net/canopsis/go-revolution/canopsis"
```

## Go setup

Fetch the latest release of Go: https://golang.org/dl/

```
wget https://dl.google.com/go/go1.9.4.linux-amd64.tar.gz
rm -rf /usr/local/go && tar xf go1.9.4.linux-amd64.tar.gz -C /usr/local/
export PATH=$PATH:/usr/local/go/bin
```

Setup Go environment:

```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

mkdir -p $GOPATH/{bin,src}
mkdir -p $GOPATH/src/git.canopsis.net/canopsis
```

Clone the project if not already done:

```
git clone https://git.canopsis.net/canopsis/go-revolution.git -b develop $GOPATH/src/git.canopsis.net/canopsis/
```

Install Glide: https://glide.sh/

Build some binaries:

```
cd $GOPATH/src/git.canopsis.net/canopsis/go-revolution/
make
```

## Compat - Python + Go Engines

L’objectif est d’avoir le schéma suivant de communication entre les engines :

```
Exchange canopsis.events -> CHE -> Event Filter -> Axe + Autres engines...
```

### Configuration

 * `etc/supervisord.d/amqp2engines.conf` : retirer `engine-cleaner-alerts`, `engine-cleaner-events` et `engine-alerts`
 * `etc/amqp2engines.conf` : retirer toute occurrence des engines précédents, et ajouter `axe` dans la liste `next` de l’engine `event filter`.

```ini
[engine:event_filter]
next = axe,...
```

### RabbitMQ

Pour une installation complètement cloisonnée, retirer tous les bindings des queues suivantes :

 * `Engine_alerts`
 * `Engine_cleaner_alerts`
 * `Engine_cleaner_events`

Créer et binder les queues suivantes :

 * `Engine_che` : `canopsis.events` sur rk `#`
 * `Engine_lifeline` : `canopsis.events` sur rk `#`
 * `Engine_axe` : `amq.direct` sur rk `#`

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
