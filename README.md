# Go-revolution

Nouveaux moteurs pour Canopsis


## Mise en place d'un environnement d'exécution des tests

### Installation d'un serveur RabbitMQ

- Installer RabbitMQ > 3.2
- installer le management plugin : `rabbitmq-plugins enable rabbitmq_management`
- [configurer rabbbitMQ pour Canopsis](https://git.canopsis.net/canopsis/canopsis/blob/develop/doc/docs/fr/guide_administrateur/package_install.md)


- créer les exchanges :


vhost    | name            | type   | Durability | auto delete | internal
---------|-----------------|--------|------------|-------------|----------
canopsis | canopsis.events | fanout | durable    | no          | no
canopsis | canopsis.alerts | fanout | durable    | no          | no


- Ajouter les bindings :  bindings :

Exchange        | Queue       | Routing key
--------------- | ----------- | -----------
canopsis.events | che         | #
canopsis.alerts | axe         | #
canopsis.events | Engine_stat | #

### MongoDB

- [installer mongodb 3.4](https://docs.mongodb.com/v3.4/administration/install-on-linux/)


- Configurer la base :

```
use canopsis
db.createUser({user:"cpsmongo",pwd:"canopsis",roles:["dbOwner"]})
db.periodical_alarm.createIndex({t:1, d:1})
db.periodical_alarm.createIndex({d:1})
```



### InfluxDB

- [Installer influxdb](https://portal.influxdata.com/downloads)

- configurer l'authentification : ouvrir le fichier `/etc/influxdb/influxdb.conf`
```
[http]
  # Determines whether HTTP endpoint is enabled.
  # enabled = true

  # The bind address used by the HTTP service.
  # bind-address = ":8086"

  # Determines whether user authentication is enabled over HTTP/HTTPS.
   auth-enabled = true

  # The default realm sent back when issuing a basic auth challenge.
  # realm = "InfluxDB"

  [...]
```

- créer l'utilisateur : `$ influx`

```
CREATE USER cplsinflux WITH PASSWORD 'canopsis' WITH ALL PRIVILEGES

```


- redémarrer influx

```
# systemctl restart influxdb
```


### redis

- installer redis

```
# apt-get install redis-server
```

### Go setup

Installer la dernière version de go : https://golang.org/dl/

```bash
wget https://dl.google.com/go/go1.9.4.linux-amd64.tar.gz
rm -rf /usr/local/go && tar xf go1.9.4.linux-amd64.tar.gz -C /usr/local/
export PATH=$PATH:/usr/local/go/bin
```

Définir l'environnement go :

```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

mkdir -p $GOPATH/{bin,src}
mkdir -p $GOPATH/src/git.canopsis.net/canopsis
```

Cloner le projet:

```bash
git clone https://git.canopsis.net/canopsis/go-revolution.git -b develop $GOPATH/src/git.canopsis.net/canopsis/
```

Installer Glide: https://glide.sh/

lancer le build:

```bash
cd $GOPATH/src/git.canopsis.net/canopsis/go-revolution/
make
```


### paramétrage des moteurs

```
export CPS_AMQP_URL="amqp://cpsrabbit:canopsis@localhost/canopsis"
export CPS_MONGO_URL="mongodb://cpsmongo:canopsis@localhost/canopsis"
export CPS_REDIS_URL="redis://nouser:dbpassword@host:port/0"
export CPS_INFLUX_URL="http://cpsinflux:canopsis@host:8086"
export CPS_DEFAULT_CFG="$GOPATH/src/git.canopsis.net/canopsis/go-revolution/canopsis/default_configuration.toml"
```





## Compat - Python + Go Engines

L’objectif est d’avoir le schéma suivant de communication entre les engines :

```
Exchange canopsis.events -> CHE -> Event Filter -> Axe + Autres engines...
```

### Configuration

 * `etc/supervisord.d/amqp2engines.conf` : retirer `engine-cleaner-alerts`, `engine-cleaner-events` et `engine-alerts`. Utile uniquement avec une installation `build-install`.
 * `etc/amqp2engines.conf` : retirer toute occurrence des engines précédents, et ajouter `axe` dans la liste `next` de l’engine `event filter`.

```ini
[engine:event_filter]
next = axe,...
```

Dans le cas d’une installation en `build-install` :

```
su - canopsis -c "supervisorctl update"
su - canopsis -c "hypcontrol start"
```

### RabbitMQ

Pour une installation complètement cloisonnée, retirer tous les bindings des queues suivantes :

 * `Engine_alerts`
 * `Engine_cleaner_alerts`
 * `Engine_cleaner_events`

Créer et binder les queues suivantes :

 * `Engine_che` : `canopsis.events` sur rk `#`
 * `Engine_heartbeat` : `canopsis.events` sur rk `#`
 * `Engine_axe` : `amq.direct` sur rk `#`

## Tests - GoConvey

You will need engines environment variables.

If you want to skip tests with long run times:

```bash
export CPS_TEST_SKIP_LONG=1
```

Then run `goconvey`:

```bash
goconvey -workDir ${GOPATH}/src/git.canopsis.net/canopsis/go-revolution/
```

Note: you can enable desktop notifications from the web UI to avoid checking manually.


# Développement


## Utiliser la bibliothèque Canopsis

```go
import "git.canopsis.net/canopsis/go-revolution/canopsis"
```