# Installation de Canopsis avec Docker

Cette procédure décrit l'installation de l'édition open-source de Canopsis en mono-instance Docker.

L'ensemble des procédures décrites doivent être réalisées avec l'utilisateur `root`.

## Prérequis

### Version minimum du noyau Linux

Votre système hôte pour Docker doit disposer d'un noyau Linux suffisamment récent. Canopsis nécessite l'utilisation **d'un noyau Linux 4.4 minimum**. Le noyau installé par défaut sur CentOS 7 n'est donc pas suffisant pour héberger un environnement Docker Canopsis.

Vérifier votre version du noyau à l'aide de la commande suivante :
```sh
uname -r
```

Si vous obtenez une version inférieure à 4.4, veuillez utiliser un noyau plus récent : soit avec une distribution plus à jour, ou bien à l'aide d'[ELRepo](https://elrepo.org/tiki/kernel-lt) pour CentOS, par exemple.

### Installation de Docker CE et Docker Compose

Canopsis nécessite [l'installation de Docker CE](https://docs.docker.com/install/#supported-platforms), version 18.06 minimum. Veuillez utiliser les dépôts officiels de Docker CE, et non pas ceux proposés par votre distribution.

Une fois Docker CE installé, vous devez ensuite [installer Docker Compose](https://docs.docker.com/compose/install/#install-compose).

## Lancement de Canopsis avec Docker Compose

Les images Docker officielles de Canopsis sont hébergées sur Docker Hub : <https://hub.docker.com/u/canopsis/>.

Le [dépôt Git de Canopsis](https://git.canopsis.net/canopsis/canopsis/-/tree/develop) contient des fichiers Docker Compose d'exemple :
```sh
git clone -b develop https://git.canopsis.net/canopsis/canopsis.git && cd canopsis/docker-compose
```

La variable `CANOPSIS_IMAGE_TAG` du fichier `.env` a pour valeur par défaut `nightly`, ce qui correspond à la version la plus récente des images disponibles sur Docker Hub. Notez que cela englobe également les versions intermédiaires destinées au développement et aux tests. Aussi, le tag `nigthly` ne doit **pas** être utilisé pour un environnement de production.

Lancez ensuite la commande suivante :
```sh
docker-compose up -d
```

## Vérification du bon fonctionnement

La vérification va passer par la commande `docker-compose ps` :

```sh
docker-compose ps

             Name                           Command               State                Ports
--------------------------------------------------------------------------------------------------------
docker-compose_action_1          /engine-action                   Up
docker-compose_axe_1             /engine-axe                      Up
docker-compose_che_1             /engine-che -d -publishQue ...   Up
docker-compose_fifo_1            /engine-fifo                     Up
docker-compose_heartbeat_1       /engine-heartbeat                Up
docker-compose_influxdb_1        /entrypoint.sh influxd           Up       0.0.0.0:4444->4444/udp,
                                                                           0.0.0.0:8083->8083/tcp,
                                                                           0.0.0.0:8086->8086/tcp
docker-compose_init_1            /bin/sh -c /${_BINARY_NAME}      Exit 0
docker-compose_metric_1          /bin/sh -c /entrypoint.sh        Up       8082/tcp
docker-compose_mongodb_1         docker-entrypoint.sh --wir ...   Up       0.0.0.0:27027->27017/tcp
docker-compose_nginx_1           /bin/sh -c /entrypoint.sh        Up       0.0.0.0:80->80/tcp
docker-compose_pbehavior_1       /bin/sh -c /entrypoint.sh        Up       8082/tcp
docker-compose_provisionning_1   /bin/sh -c /entrypoint-prov.sh   Exit 0
docker-compose_rabbitmq_1        docker-entrypoint.sh rabbi ...   Up       15671/tcp,
                                                                           0.0.0.0:15672->15672/tcp,
                                                                           25672/tcp, 4369/tcp,
                                                                           5671/tcp,
                                                                           0.0.0.0:5672->5672/tcp
docker-compose_redis_1           docker-entrypoint.sh redis ...   Up       0.0.0.0:6379->6379/tcp
docker-compose_scheduler_1       /bin/sh -c /entrypoint.sh        Up       8082/tcp
docker-compose_task_importctx_1  /bin/sh -c /entrypoint.sh        Up       8082/tcp
docker-compose_watcher_1         /bin/sh -c /${_BINARY_NAME}      Up
docker-compose_webserver_1       /bin/sh -c /entrypoint.sh        Up       0.0.0.0:8082->8082/tcp
```

Les services doivent être en état `Up` ou `Exit 0`. En fonction des ressources de votre machine, il peut être nécessaire d'attendre quelques minute avant que l'ensemble des moteurs puissent passer en état `Up`.

Vous pouvez alors procéder à votre [première connexion à l'interface Canopsis](premiere-connexion.md).

## Arrêt de l'environnement Docker Compose

```sh
docker-compose down

Stopping docker-compose_nginx_1          ... done
Stopping docker-compose_webserver_1      ... done
Stopping docker-compose_metric_1         ... done
Stopping docker-compose_task_importctx_1 ... done
Stopping docker-compose_scheduler_1      ... done
Stopping docker-compose_pbehavior_1      ... done
Stopping docker-compose_mongodb_1        ... done
Stopping docker-compose_fifo_1           ... done
Stopping docker-compose_action_1         ... done
Stopping docker-compose_axe_1            ... done
Stopping docker-compose_redis_1          ... done
Stopping docker-compose_che_1            ... done
Stopping docker-compose_influxdb_1       ... done
Stopping docker-compose_heartbeat_1      ... done
Stopping docker-compose_rabbitmq_1       ... done
Stopping docker-compose_watcher_1        ... done
Removing docker-compose_nginx_1          ... done
Removing docker-compose_webserver_1      ... done
Removing docker-compose_metric_1         ... done
Removing docker-compose_provisionning_1  ... done
Removing docker-compose_task_importctx_1 ... done
Removing docker-compose_scheduler_1      ... done
Removing docker-compose_pbehavior_1      ... done
Removing docker-compose_init_1           ... done
Removing docker-compose_mongodb_1        ... done
Removing docker-compose_fifo_1           ... done
Removing docker-compose_action_1         ... done
Removing docker-compose_axe_1            ... done
Removing docker-compose_redis_1          ... done
Removing docker-compose_che_1            ... done
Removing docker-compose_influxdb_1       ... done
Removing docker-compose_heartbeat_1      ... done
Removing docker-compose_rabbitmq_1       ... done
Removing docker-compose_watcher_1        ... done
Removing network docker-compose_default
```

## Rétention des logs

La mise en place d'une politique de rétention des logs nécessite la présence du logiciel `logrotate`.

Une fois que `logrotate` est installé sur votre machine, créer le fichier `/etc/logrotate.d/docker-container` suivant :

```
/var/lib/docker/containers/*/*.log {
  rotate 7
  daily
  compress
  minsize 100M
  notifempty
  missingok
  delaycompress
  copytruncate
}
```

Pour vérifier la bonne exécution de la configuration de logrotate pour Docker, vous pouvez lancer la commande :

```sh
logrotate -fv /etc/logrotate.d/docker-container
```
