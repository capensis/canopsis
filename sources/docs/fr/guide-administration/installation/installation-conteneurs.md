# Installation Conteneur

## Pré-requis

### Ports

Il faut vérifier les ports utilisés par `docker-compose.yml`.  
Veuillez effectuer les adaptations nécessaires si certains ports sont déjà en utilisation sur le noeud d'installation.

## Installation

- Cloner le dépôt Canopsis : https://git.canopsis.net/canopsis/canopsis
- Actuellement, les conteneurs sont gérés dans le Docker Hub officiel : `https://hub.docker.com/u/canopsis/`
- Dans ce dépôt, un fichier `docker-compose.yml` est présent. Il va servir à la création de votre Canopsis en version Dockerisée.
  
    - Rappel : [Installation de Docker compose](https://docs.docker.com/compose/install/#install-compose)  
  
    - Troubleshooting : Si vous rencontrez une erreur de ce type lors d'un `docker-compose up -d` ou un `docker-compose --version`  
    ```bash
    bash: /usr/bin/docker-compose: Aucun fichier ou dossier de ce type
    ```

    - Résolution : Utiliser la commande suivante `hash docker-compose`  

- Exécuter la commande suivante : `docker-compose up -d`

    - Troubleshooting :
    ```sh
    docker-compose up -d
    ERROR: Couldn't connect to Docker daemon at http+docker://localhost - is it running?
    If it's at a non-standard location, specify the URL with the DOCKER_HOST environment variable.
    ```
    - Résolution : Utiliser la commande suivante : `sudo docker-compose up -d`

## Vérification

La vérification va passer par la commande `sudo docker-compose ps` :

```sh
sudo docker-compose ps

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

## Arrêt et suppression du docker-compose

```sh
sudo docker-compose down

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
