# Initialisation.

## Introduction

La commande *canopsis-reconfigure* permet d'initialiser les queues et les exchanges RabbitMQ à
l'aide d'une configuration.

Il est à terme de prévu d'initialiser les différents services dont canopsis à
besoin pour fonctionner.

## Utilisation
### Docker-compose

Pour utiliser *canopsis-reconfigure* dans un environnement docker, il suffit juste de
lancer canopsis à l'aide de docker-compose et du fichier docker-compose.yml
à la racine du projet.

```sh
docker-compose up -d
```

### Canopsis classique

Si vous n'utilisez pas docker-compose pour démarrer canopsis, vous pouvez
lancer la commande manuellement.

```sh
CPS_AMQP_URL=amqp://cpsrabbit:canopsis@rabbitmq:5672/canopsis ./canopsis-reconfigure
```

Ou si vous souhaitez utiliser une configuration particulière

```sh
CPS_AMQP_URL=amqp://cpsrabbit:canopsis@rabbitmq:5672/canopsis ./canopsis-reconfigure -conf=./autre_initialisation.toml
```

## Lancement manuel

Après avoir construit le binaire comme tous les autres :

```
export CPS_MAX_RETRY=10
export CPS_MAX_DELAY=1
export CPS_WAIT_FIRST_ATTEMPT=10
./canopsis-reconfigure
```
