# Go-revolution

Nouveaux moteurs pour Canopsis.

Nécessite Go 1.12 ou supérieur.

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

Exchange        | Queue            | Routing key
--------------- | ---------------- | -----------
                | Engine_action    | #
canopsis.alerts | Engine_axe       | #
canopsis.events | Engine_che       | #
canopsis.events | Engine_heartbeat | #
canopsis.events | Engine_stat      | #

### MongoDB

- [installer mongodb 3.4](https://docs.mongodb.com/v3.4/administration/install-on-linux/)
- Configurer la base :

```js
use canopsis
db.createUser({user:"cpsmongo",pwd:"canopsis",roles:["dbOwner"]})
db.periodical_alarm.createIndex({t:1, d:1})
db.periodical_alarm.createIndex({d:1})
```

### InfluxDB

- [Installer influxdb](https://portal.influxdata.com/downloads)
- configurer l'authentification : ouvrir le fichier `/etc/influxdb/influxdb.conf`

```ini
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

```sql
CREATE USER cplsinflux WITH PASSWORD 'canopsis' WITH ALL PRIVILEGES
```

- redémarrer influx

```sh
systemctl restart influxdb
```

### Redis

- installer redis

```sh
apt-get install redis-server
```

### Go setup

Installer la dernière version de go : https://golang.org/dl/

```sh
wget https://dl.google.com/go/go1.9.4.linux-amd64.tar.gz
rm -rf /usr/local/go && tar xf go1.9.4.linux-amd64.tar.gz -C /usr/local/
export PATH=$PATH:/usr/local/go/bin
```

Définir l'environnement go :

```sh
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

mkdir -p $GOPATH/{bin,src}
mkdir -p $GOPATH/src/git.canopsis.net/canopsis/go-revolution
```

Cloner le projet :

```sh
git clone https://git.canopsis.net/canopsis/go-revolution.git -b develop $GOPATH/src/git.canopsis.net/canopsis/go-revolution
```

Installer Glide: https://glide.sh/

Initialiser le projet :

```sh
cd $GOPATH/src/git.canopsis.net/canopsis/go-revolution/
make init
```

Lancer le build :

```sh
cd $GOPATH/src/git.canopsis.net/canopsis/go-revolution/
make
```

Lors du développement, il peut être utile de *builder* les binaires rapidement et des les récupérer dans un dossier partagé avec une autre machine par exemple :

```sh
make init

# réutiliser cette commande par la suite
make build BUILD_OUTPUT_DIR=/vmshare/gobin SKIP_DEPENDENCIES=true
```

### Paramétrage des moteurs

```sh
export CPS_AMQP_URL="amqp://cpsrabbit:canopsis@localhost/canopsis"
export CPS_MONGO_URL="mongodb://cpsmongo:canopsis@localhost/canopsis"
export CPS_REDIS_URL="redis://nouser:dbpassword@host:port/0"
export CPS_INFLUX_URL="http://cpsinflux:canopsis@host:8086"
export CPS_DEFAULT_CFG="$GOPATH/src/git.canopsis.net/canopsis/go-revolution/default_configuration.toml"
```

#### Paramètres spécifiques au développement

```sh
export CPS_TEST_SKIP_LONG=1
export CPS_DEBUG_PPROF_ENABLE=0
export CPS_DEBUG_PPROF_CPU=cpu.out
export CPS_DEBUG_PPROF_MEMORY=mem.out
export CPS_DEBUG_TRACE=trace.out
```

`CPS_TEST_SKIP_LONG` permet de passer les tests marqués comme long à s'exécuter.

Les variables d'environnement `CPS_DEBUG_PPROF_*` permettent d'activer le
profiling, pour voir notamment le gain de performances associé à un
développement. Les profils générés avec ces options peuvent être consultés avec
la commande `go tool pprof` (par exemple : `go tool pprof -http=":8081"
engine-axe cpu.out`).

La variable d'environnement `CPS_DEBUG_TRACE` permet de générer des traces
d'exécutions. La trace générée par cette option peut être consultée avec la
commande `go tool trace` (par exemple : `go tool trace trace.out`). L'interface
lancée par `go tool trace` n'est pour l'instant compatible qu'avec le
navigateur Chromium et ses variantes.


## Compat - Python + Go Engines

L’objectif est d’avoir le schéma suivant de communication entre les engines :

```
Exchange canopsis.events -> CHE -> Event Filter -> Axe + Autres engines...
```

### Configuration

Prendre les fichiers de conf de docker/etc et les copier dans environnement
canopsis.

Dans le cas d’une installation en `build-install` :

```sh
su - canopsis -c "supervisorctl update"
su - canopsis -c "hypcontrol start"
```

### Docker

Pour pouvoir utiliser docker compose, il faut préalablement construire l'image docker de compatibilité, ainsi que les images Docker des engines en Go.

```sh
# archive de tous les binaires construits dans Docker
make docker_images

# ou avec tag custom
make docker_images TAG=develop
```

Pour que le provisionning soit complet (reinit), il faut supprimer le volume mongo et perdre toutes les données:
```sh
docker volume rm go-revolution_mongodbdata
```

Pour ne construire qu’un seul projet :

```sh
make -C cmd/<project> -f ../../Makefile.cmd docker_image
```

Pour faire une release (archive tar):

```sh
make docker_release
```

### RabbitMQ

Pour une installation complètement cloisonnée, retirer tous les bindings des queues suivantes :

 * `Engine_alerts`
 * `Engine_cleaner_alerts`
 * `Engine_cleaner_events`

Créer et binder les queues suivantes :

 * `Engine_che` sur `canopsis.events` avec la rk `#`
 * `Engine_heartbeat` sur `canopsis.events` avec la rk `#`
 * `Engine_axe` sur `amq.direct` avec la rk `#`
 * `Engine_stat` sur `canopsis.events` avec la rk `#`

Après activation de rabbitmqadmin:

Stopper canopsis puis les engines go
```sh
/opt/canopsis/bin/canopsis-systemd stop
```

Disable les units des engines python:
```sh
systemctl disable canopsis-engine@cleaner-cleaner_alerts.service
systemctl disable canopsis-engine@dynamic-alerts.service
systemctl disable canopsis-engine@cleaner-cleaner_events.service
```

Faire les modifications sur les queues:

```sh
rabbitmqadmin -u cpsrabbit -p canopsis --vhost canopsis delete queue name=Engine_alerts
rabbitmqadmin -u cpsrabbit -p canopsis --vhost canopsis delete queue name=Engine_cleaner_alerts
rabbitmqadmin -u cpsrabbit -p canopsis --vhost canopsis delete queue name=Engine_cleaner_events

rabbitmqadmin -u cpsrabbit -p canopsis --vhost canopsis declare queue name=Engine_che durable=true auto_delete=false
rabbitmqadmin -u cpsrabbit -p canopsis --vhost canopsis declare queue name=Engine_heartbeat durable=true auto_delete=false
rabbitmqadmin -u cpsrabbit -p canopsis --vhost canopsis declare queue name=Engine_axe durable=true auto_delete=false
rabbitmqadmin -u cpsrabbit -p canopsis --vhost canopsis declare queue name=Engine_stat durable=true auto_delete=false

rabbitmqadmin -u cpsrabbit -p canopsis --vhost canopsis declare binding source=canopsis.events destination=Engine_che routing_key="#"
rabbitmqadmin -u cpsrabbit -p canopsis --vhost canopsis declare binding source=canopsis.events destination=Engine_heartbeat routing_key="#"
rabbitmqadmin -u cpsrabbit -p canopsis --vhost canopsis declare binding source=amq.direct destination=Engine_axe routing_key="#"
rabbitmqadmin -u cpsrabbit -p canopsis --vhost canopsis declare binding source=canopsis.events destination=Engine_stat routing_key="#"
```

Relancer canopsis puis les engines go
```sh
/opt/canopsis/bin/canopsis-systemd restart
```


## Paramétrage de lancement

Certains engines supportent des options au lancement, utilisez `-help` pour les voir.

À chaque *release* peut apporter des changements dans les options, veuillez y jeter un œil de temps à autres.

## Profiling intégré

Les binaires suivants permettent de lancer un *profiling* Go :

 * `engine-che`
 * `engine-axe`
 * `engine-heartbeat`
 * `engine-stat`
 * `engine-action`

Pour l’activer/désactiver globalement :

```sh
# Activation
export CPS_DEBUG_PPROF_ENABLE=1

# Désactivation
export CPS_DEBUG_PPROF_ENABLE=autrechose
```

Pour activer le profiling CPU :

```sh
export CPS_DEBUG_PPROF_CPU=/chemin/vers/trace.cpu.out
```

Pour activer le profiling Mémoire :

```sh
export CPS_DEBUG_PPROF_MEMORY=/chemin/vers/trace.mem.out
```

Ensuite lancer n’importe quel engine. Il devra être quitté proprement pour que les traces soient écrites.

Le ou les fichiers de trace sont à envoyer tels quels pour analyse.

## Procédures de purge

Dans certains cas, purger ou modifier des collections MongoDB en dehors des engines Go entraîne une incohérence dans l’état du système et ne permet pas d’avoir un comportement attendu.

Redémarrer Redis ou purger ses bases lorsque ces collections sont modifiées / purgée de manière externe :

 * `periodical_alarm`
 * `default_entities`

## Tests - GoConvey

[GoConvey](http://goconvey.co/) et `docker-compose` sont utilisés pour lancer les tests :

```sh
go get -u github.com/smartystreets/goconvey
make test
```

# Développement

## Utiliser la bibliothèque Canopsis

```go
import "git.canopsis.net/canopsis/go-revolution/canopsis"
```

## Builder les images

```sh
make docker_images TAG+[TAG]
```

## Monter un environnement de développement sous docker

A la racine du projet, exécuter :

```sh
docker-compose up
```

