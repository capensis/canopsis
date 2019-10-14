# Installation de Canopsis avec Docker

Cette procédure décrit l'installation de l'édition open-source de Canopsis en mono-instance Docker.

L'ensemble des procédures décrites doivent être réalisées avec l'utilisateur `root`.

## Pré-requis

### Version minimum du noyau Linux

Votre système hôte pour Docker doit disposer d'un noyau Linux suffisamment récent. Canopsis nécessite l'utilisation **d'un noyau Linux 4.4 minimum**.

Le noyau installé par défaut sur CentOS 7 n'est donc **pas suffisant** pour héberger un environnement Docker Canopsis.

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

Le [dépôt Git de Canopsis](https://git.canopsis.net/canopsis/canopsis) contient des fichiers Docker Compose d'exemple :
```sh
git clone https://git.canopsis.net/canopsis/canopsis.git && cd canopsis
```

Si nécessaire, ajustez la variable `CANOPSIS_IMAGE_TAG` du fichier `.env` situé à la racine du dépôt Canopsis.

Lancez ensuite la commande suivante, afin de lever un environnement Canopsis open-core complet avec Docker :
```sh
docker-compose up -d
```

## Vérification du bon fonctionnement

La vérification va passer par la commande `docker-compose ps` :

```sh
docker-compose ps

          Name                         Command               State                          Ports                      
-----------------------------------------------------------------------------------------------------------------------
canopsis_alerts_1           /bin/sh -c /entrypoint.sh        Up       8082/tcp                                         
canopsis_cleaner_events_1   /bin/sh -c /entrypoint.sh        Up       8082/tcp                                         
canopsis_context-graph_1    /bin/sh -c /entrypoint.sh        Up       8082/tcp                                         
canopsis_event_filter_1     /bin/sh -c /entrypoint.sh        Up       8082/tcp                                         
canopsis_eventstore_1       /bin/sh -c /entrypoint.sh        Up       8082/tcp                                         
canopsis_influxdb_1         /entrypoint.sh influxd           Up       0.0.0.0:4444->4444/udp, 0.0.0.0:8083->8083/tcp,  
                                                                      0.0.0.0:8086->8086/tcp                           
canopsis_metric_1           /bin/sh -c /entrypoint.sh        Up       8082/tcp                                         
canopsis_mongodb_1          docker-entrypoint.sh --wir ...   Up       0.0.0.0:27027->27017/tcp                         
canopsis_pbehavior_1        /bin/sh -c /entrypoint.sh        Up       8082/tcp                                         
canopsis_provisionning_1    /bin/sh -c /entrypoint-prov.sh   Exit 0                                                    
canopsis_rabbitmq_1         docker-entrypoint.sh rabbi ...   Up       15671/tcp, 0.0.0.0:15672->15672/tcp, 25672/tcp,  
                                                                      4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp       
canopsis_redis_1            docker-entrypoint.sh redis ...   Up       0.0.0.0:6379->6379/tcp                           
canopsis_scheduler_1        /bin/sh -c /entrypoint.sh        Up       8082/tcp                                         
canopsis_task_importctx_1   /bin/sh -c /entrypoint.sh        Up       8082/tcp                                         
canopsis_task_mail_1        /bin/sh -c /entrypoint.sh        Up       8082/tcp                                         
canopsis_ticket_1           /bin/sh -c /entrypoint.sh        Up       8082/tcp                                         
canopsis_webserver_1        /bin/sh -c /entrypoint.sh        Up       0.0.0.0:28082->8082/tcp   
```

Les services doivent être en état `Up` ou `Exit 0`.

Vous pouvez alors procéder à votre [première connexion à l'interface Canopsis](premiere-connexion.md).

## Arrêt de l'environnement Docker Compose 

```sh
docker-compose down

Stopping canopsis_webserver_1      ... done
Stopping canopsis_metric_1         ... done
Stopping canopsis_event_filter_1   ... done
Stopping canopsis_task_importctx_1 ... done
Stopping canopsis_eventstore_1     ... done
Stopping canopsis_alerts_1         ... done
Stopping canopsis_scheduler_1      ... done
Stopping canopsis_pbehavior_1      ... done
Stopping canopsis_context-graph_1  ... done
Stopping canopsis_cleaner_events_1 ... done
Stopping canopsis_ticket_1         ... done
Stopping canopsis_task_mail_1      ... done
Stopping canopsis_rabbitmq_1       ... done
Stopping canopsis_influxdb_1       ... done
Stopping canopsis_mongodb_1        ... done
Stopping canopsis_redis_1          ... done
Removing canopsis_webserver_1      ... done
Removing canopsis_metric_1         ... done
Removing canopsis_event_filter_1   ... done
Removing canopsis_task_importctx_1 ... done
Removing canopsis_eventstore_1     ... done
Removing canopsis_alerts_1         ... done
Removing canopsis_provisionning_1  ... done
Removing canopsis_scheduler_1      ... done
Removing canopsis_pbehavior_1      ... done
Removing canopsis_context-graph_1  ... done
Removing canopsis_cleaner_events_1 ... done
Removing canopsis_ticket_1         ... done
Removing canopsis_task_mail_1      ... done
Removing canopsis_rabbitmq_1       ... done
Removing canopsis_influxdb_1       ... done
Removing canopsis_mongodb_1        ... done
Removing canopsis_redis_1          ... done
Removing network canopsis_default
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
