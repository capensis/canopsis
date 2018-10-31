# Installation Conteneur

## Pré-requis

### Ports

Il faut vérifier les ports utilisés par `docker-compose.yml`, si certains ports sont utilisés sur votre machines veuillez les libérer pour le bon déroulement de l'installation.

## Installation

Pour l'installation dockerisé de Canopsis, le porcédure est la suivante :

- clôner le dépôt Canopsis : https://git.canopsis.net/canopsis/canopsis
- Dans ce dépot un fichier `docker-compose.yml` est présent. Il va servir à la création de votre Canopsis en version Dockerisé.
  
    - Rappel : [Intallation de Docker compose](https://docs.docker.com/compose/install/#install-compose)  
  
    - Troubleshotting : Si vous rencontrez une erreur de ce type lors d'un `docker-compose up -d` ou un `docker-compose --version`  
    ```
    bash: /usr/bin/docker-compose: Aucun fichier ou dossier de ce type
    ```

    - Résolution : Utiliser la commande suivante `hash docker-compose`  

- faire la commande suivante : `docker-compose up -d`

    - Troubleshotting :
    ```
    docker-compose up -d
    ERROR: Couldn't connect to Docker daemon at http+docker://localhost - is it running?
    If it's at a non-standard location, specify the URL with the DOCKER_HOST environment variable.
    ```
    - Résolution : Utiliser la commande suivante : `sudo docker-compose up -d`

## Vérification

La vérification va passer par la commande `sudo docker-compose ps` :

```
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

### Détail des colonnes

- Name : Nom des conteneurs déployés.
- Command : Commandes en rapport avec la CLI Docker.
- State : Etat du conteneur.
- Ports : Ports utilisé par le conteneur.