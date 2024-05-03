# Guide de migration vers Canopsis 24.04.0

Ce guide donne les instructions vous permettant de mettre à jour Canopsis 23.10 (dernière version disponible) vers [la version 24.04.0](../24.04.0.md).

## Prérequis

L'ensemble de cette procédure doit être lu avant son exécution.

Ce document ne prend en compte que Canopsis Community et Canopsis Pro : tout développement personnalisé dont vous pourriez bénéficier ne fait pas partie du cadre de ce Guide de migration.

Les fichiers de référence qui sont mentionnés dans ce guide sont disponibles à ces adresses

| Édition           | Sources                                                                                                                              |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| Édition Community | [https://git.canopsis.net/canopsis/canopsis-community/-/releases](https://git.canopsis.net/canopsis/canopsis-community/-/releases)   |
| Édition pro       | [https://git.canopsis.net/sources/canopsis-pro-sources/-/releases](https://git.canopsis.net/sources/canopsis-pro-sources/-/releases) |

[TOC]

## Procédure de mise à jour

### Réalisation d'une sauvegarde

Des sauvegardes sont toujours recommandées, qu'elles soient régulières ou lors de modifications importantes.

La restructuration apportée dans les bases de données pour cette version de Canopsis nous amène à insister d'autant plus sur ce point. Il est donc fortement recommandé de réaliser une **sauvegarde complète** des VM hébergeant vos services Canopsis, avant cette mise à jour.

### Vérification MongoDB

!!! warning "Vérification"

    Avant de démarrer la procédure de mise à jour, vous devez vérifier que la valeur de `featureCompatibilityVersion` est bien positionnée à **5.0**  

    === "Docker Compose"
        ```sh
        CPS_EDITION=pro docker compose exec mongodb bash
        mongo -u root -p root
        > db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )
        > exit
        ```

    === "Paquets RHEL 8"

        ```sh
        mongo -u root -p root
        > db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )
        > exit
        ```

    Le retour doit être de la forme `{ "featureCompatibilityVersion" : { "version" : "5.0" }, "ok" : 1 }`
    Si ce n'est pas le cas, vous ne pouvez pas continuer la mise à jour.


### Arrêt de l'environnement en cours de lancement

Vous devez prévoir une interruption du service afin de procéder à la mise à jour qui va suivre.

=== "Docker Compose"

    ```sh
    CPS_EDITION=pro docker compose down
    ```

=== "Paquets RHEL 8"

    ```sh
    systemctl stop canopsis
    systemctl stop mongod
    systemctl stop postgresql-13
    systemctl stop rabbitmq-server
    systemctl stop redis
    ```

## Mise à jour Canopsis

!!! information "Information"

    Canopsis 24.04 est livré avec un nouveau jeu de configurations de référence.
    Vous devez télécharger ces configurations et y reporter vos personnalisations.  

=== "Docker Compose"

    Si vous êtes utilisateur de l'édition `community`, voici les étapes à suivre.

    Télécharger le paquet de la version 24.04.0 (canopsis-community-docker-compose-24.04.0.tar.gz) disponible à cette adresse [https://git.canopsis.net/canopsis/canopsis-community/-/releases](https://git.canopsis.net/canopsis/canopsis-community/-/releases).

    ```sh
    export CPS_EDITION=community
    tar xvfz canopsis-community-docker-compose-24.04.0.tar.gz
    cd canopsis-community-docker-compose-24.04.0
    ```

    Si vous êtes utilisateur de l'édition `pro`, voici les étapes à suivre.

    Télécharger le paquet de la version 24.04.0 (canopsis-pro-docker-compose-24.04.0.tar.gz) disponible à cette adresse [https://git.canopsis.net/sources/canopsis-pro-sources/-/releases](https://git.canopsis.net/sources/canopsis-pro-sources/-/releases).

    ```sh
    export CPS_EDITION=pro
    tar xvfz canopsis-pro-docker-compose-24.04.0.tar.gz
    cd canopsis-pro-docker-compose-24.04.0
    ```

    À ce stade, vous devez synchroniser les modifications réalisées sur vos anciens fichiers de configuration `docker-compose` avec les fichiers `docker-compose.yml` et/ou `docker-compose.override.yml`.

=== "Paquets RHEL 8"

    Non concerné car ces configurations sont livrées directemement dans les paquets RPM.



### Mise à jour de MongoDB

Dans cette version de Canopsis, la base de données MongoDB passe de la version 5.0 à 7.0.  
2 étapes sont nécessaires pour cette migration : 

1. De 5.0 à 6.0
2. De 6.0 à 7.0

=== "Docker Compose"

    Modifiez la variable `MONGO_TAG` du fichier `.env` de cette façon :

    ```diff
    -MONGO_TAG=7.0.8-jammy
    +MONGO_TAG=6.0.15-jammy
    ```

    Démarrez le conteneur `mongodb` :

    ```sh
    CPS_EDITION=pro docker compose up -d mongodb
    ```

    Entrez ensuite à l'intérieur de ce conteneur, afin de compléter la mise à jour vers MongoDB 6.0 :

    ```sh
    CPS_EDITION=pro docker compose exec mongodb bash
    mongosh -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "6.0" } )
    exit
    ```

    Modifiez la variable `MONGO_TAG` du fichier `.env` de cette façon :

    ```diff
    -MONGO_TAG=6.0.15-jammy
    +MONGO_TAG=7.0.8-jammy
    ```

    Démarrez le conteneur `mongodb` :

    ```sh
    CPS_EDITION=pro docker compose up -d mongodb
    ```

    Entrez ensuite à l'intérieur de ce conteneur, afin de compléter la mise à jour vers MongoDB 7.0 :

    ```sh
    CPS_EDITION=pro docker compose exec mongodb bash
    mongosh -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "7.0", "confirm" : true } )
    ```

    Après avoir mis à jour mongodb, l'option de telemetry sera activée. Pour la désactiver, exécutez la commande suivante :
    
    ```sh
    CPS_EDITION=pro docker compose exec mongodb bash
    mongosh -u root -p root
    > disableTelemetry()
    exit
    ```


=== "Paquets RHEL 8"

    !!! note
        Si vous avez mis en place des exclusions dans le fichier `/etc/yum.conf`, veillez à la désactiver le temps de cette procédure.

    Mise à jour des paquets `mongodb` :

    ```sh
    echo '[mongodb-org-6.0]
    name=MongoDB Repository
    baseurl=https://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/6.0/x86_64/
    gpgcheck=1
    enabled=1
    gpgkey=https://www.mongodb.org/static/pgp/server-6.0.asc' | tee /etc/yum.repos.d/mongodb-org-6.0.repo
    dnf makecache
    dnf install mongodb-org-6.0.15 mongodb-org-database-6.0.15 mongodb-org-server-6.0.15 mongodb-org-mongos-6.0.15 mongodb-org-tools-6.0.15
    ```

    Redémarrage de `mongodb` :

    ```sh
    systemctl start mongod
    ```

    Ensuite, complétez la mise à jour vers MongoDB 6.0 :

    ```sh
    mongosh -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "6.0" } )
    exit
    ```

    Mise à jour des paquets `mongodb` :

    !!! warning "Avertissement"

        Lors du passage de la version 6 à la version 7 l'option `storage.journal.enabled` est supprimée de mongodb.
        
        Si vous l'utilisiez, vous devez modifier le fichier `/etc/mongod.conf` et retirer les lignes suivantes:
        ```
        journal:
            enabled: true
        ```
    
    ```sh
    echo '[mongodb-org-7.0]
    name=MongoDB Repository
    baseurl=https://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/7.0/x86_64/
    gpgcheck=1
    enabled=1
    gpgkey=https://www.mongodb.org/static/pgp/server-7.0.asc' | tee /etc/yum.repos.d/mongodb-org-7.0.repo
    dnf makecache
    dnf install mongodb-org-7.0.8 mongodb-org-database-7.0.8 mongodb-org-server-7.0.8 mongodb-org-mongos-7.0.8 mongodb-org-tools-7.0.8
    ```


    Redémarrage de `mongodb` :

    ```sh
    systemctl start mongod
    ```

    Ensuite, complétez la mise à jour vers MongoDB 7.0 :

    ```sh
    mongosh -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "7.0", confirm: true } )
    exit
    ```
 
    Après avoir mis à jour mongodb, l'option de telemetry sera activée. Pour la désactiver, exécutez la commande suivante :
    
    ```sh
    mongosh -u root -p root
    > disableTelemetry()
    ```

### Mise à jour de TimescaleDB

Dans cette version de Canopsis, la base de données TimescaleDB passe de la version 2.9.3 à 2.14.2.

=== "Docker Compose"

    Relancez le conteneur `timescaledb` :

    ```sh
    CPS_EDITION=pro docker compose up -d timescaledb
    ```

    Puis mettez à jour l'extension timescaledb (La chaîne de connexion doit être adaptée à votre environnement)

    ```sh
    CPS_EDITION=pro docker compose exec timescaledb psql postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis
    canopsis=# ALTER EXTENSION timescaledb UPDATE;
    ```

    Ensuite, vérifiez que l'extension en elle-même est à présent bien à jour

    ```sh
    \dx
    ...
    timescaledb | 2.14.2   | public     | Enables scalable inserts and complex queries for time-series data (Community Edition)
    ...
    exit
    ```

=== "Paquets RHEL 8"

    Mise à jour des paquets `timescaledb` :

    ```sh
    dnf install timescaledb-2-postgresql-13-2.14.2 timescaledb-2-loader-postgresql-13-2.14.2
    ```

    Relancez le service `timescaledb` (géré par PostgreSQL) :

    ```sh
    systemctl start postgresql-13
    ```

    Puis mettez à jour l'extension timescaledb (La chaîne de connexion doit être adaptée à votre environnement)

    ```sh
    sudo -u postgres psql postgresql://cpspostgres:canopsis@localhost:5432/canopsis
    canopsis=# ALTER EXTENSION timescaledb UPDATE;
    ```

    !!! warning "Avertissement"
        Si l'opération `ALTER EXTENSION` échoue avec le message
        `ERROR:  must be owner of extension timescaledb`
        veuillez suivre les opérations suivantes :
        ```sh
        su - postgres
        psql -X
        \c canopsis
        canopsis=# ALTER EXTENSION timescaledb UPDATE;
        ```

    Ensuite, vérifiez que l'extension en elle-même est à présent bien à jour

    ```sh
    \dx
    ...
    timescaledb | 2.14.2   | public     | Enables scalable inserts and complex queries for time-series data (Community Edition)
    ...
    exit
    ```

### Mise à jour de RabbitMQ

Dans cette version de Canopsis, le bus rabbitMQ passe à la version 3.12.13.  

=== "Docker Compose"

    Il suffit de démarrer le conteneur

    ```sh
    CPS_EDITION=pro docker compose up -d rabbitmq
    ```

=== "Paquets RHEL 8"

    Passage en version 3.12 puis lancement du service `rabbitmq-server` :

    ```sh
    dnf install --repo rabbitmq_erlang --repo rabbitmq_server erlang rabbitmq-server-3.12.13
    systemctl restart rabbitmq-server
    ```

### Remise à 0 du cache Redis

Dans cette version de Canopsis, le cache de Canopsis doit repartir à 0.

=== "Docker Compose"

    ```sh
    CPS_EDITION=pro docker compose up -d redis
    CPS_EDITION=pro docker compose exec redis /usr/local/bin/redis-cli -a canopsis flushall
    OK
    ```

=== "Paquets RHEL 8"

    ```sh
    systemctl start redis
    /bin/redis-cli -a canopsis flushall
    ```

### Lancement du provisioning `canopsis-reconfigure`

#### Synchronisation du fichier de configuration `canopsis.toml` ou fichier de surcharge

Si vous avez modifié le fichier `canopsis.toml` (vous le voyez via une définition de volume dans votre fichier docker-compose.yml), vous devez vérifier qu'il soit bien à jour par rapport au fichier de référence.  

* [`canopsis.toml` pour Canopsis Community 24.04.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/24.04.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml)
* [`canopsis.toml` pour Canopsis Pro 24.04.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/24.04.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml)

!!! information "Information"

    Pour éviter ce type de synchronisation fastidieuse, la bonne pratique est d'utiliser [un fichier de surcharge de cette configuration](../../../guide-administration/administration-avancee/modification-canopsis-toml/). 


Si vous avez utilisé un fichier de surcharge, alors vous n'avez rien à faire, uniquement continuer à le présenter dans un volume.

#### Reconfiguration de Canopsis

=== "Docker Compose"

    !!! Attention

        Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette personnalisation.
        L'outil est en effet pré paramétré pour fonctionner naturellement.

    ```sh
    CPS_EDITION=pro docker compose up -d reconfigure
    ```

    !!! information "Information"

        Cette opération peut prendre plusieurs minutes pour s'exécuter.

    Vous pouvez ensuite vérifier que le mécanisme de provisioning/reconfigure s'est correctement déroulé. Le conteneur doit présenté un "exit 0"

    ```sh
    CPS_EDITION=pro docker compose ps -a|grep reconfigure
    canopsis-pro-reconfigure-1            "/canopsis-reconfigu…"   reconfigure            exited (0)
    ```

=== "Paquets RHEL 8"

    La commande `canopsis-reconfigure` doit être exécutée après mise à jour de Canopsis dans le cadre d'installation par paquets RPM.

### Mise à jour et démarrage final de Canopsis

Enfin, il vous reste à mettre à jour et à démarrer tous les composants applicatifs de Canopsis

=== "Docker Compose"

    ```sh
    CPS_EDITION=pro docker compose up -d
    ```

    Vous pouvez ensuite vérifier que l'ensemble des conteneurs soient correctement exécutés.

    ```sh
    CPS_EDITION=pro docker compose ps
    NAME                             IMAGE                                                                     COMMAND                                                                                        SERVICE           CREATED       STATUS                       PORTS
    canopsis-pro-action-1            docker.canopsis.net/docker/develop-pro/engine-action:24.04-rc3            "/engine-action"                                                                               action            3 hours ago   Up 51 minutes                
    canopsis-pro-api-1               docker.canopsis.net/docker/develop-pro/canopsis-api-pro:24.04-rc3         "/bin/sh -c /${CMD}"                                                                           api               3 hours ago   Up 51 minutes (healthy)      0.0.0.0:8082->8082/tcp, :::8082->8082/tcp
    canopsis-pro-axe-1               docker.canopsis.net/docker/develop-pro/engine-axe:24.04-rc3               "/engine-axe -publishQueue Engine_correlation"                                                 axe               3 hours ago   Up 51 minutes                
    canopsis-pro-che-1               docker.canopsis.net/docker/develop-pro/engine-che:24.04-rc3               "/engine-che"                                                                                  che               3 hours ago   Up 51 minutes                
    canopsis-pro-connector-junit-1   docker.canopsis.net/docker/develop-pro/connector-junit:24.04-rc3          "/bin/sh -c /${CMD}"                                                                           connector-junit   3 hours ago   Up 51 minutes                
    canopsis-pro-correlation-1       docker.canopsis.net/docker/develop-pro/engine-correlation:24.04-rc3       "/bin/sh -c /${CMD}"                                                                           correlation       3 hours ago   Up 51 minutes                
    canopsis-pro-dynamic-infos-1     docker.canopsis.net/docker/develop-pro/engine-dynamic-infos:24.04-rc3     "/bin/sh -c /${CMD}"                                                                           dynamic-infos     3 hours ago   Up 51 minutes                
    canopsis-pro-fifo-1              docker.canopsis.net/docker/develop-pro/engine-fifo:24.04-rc3              "/bin/sh -c /${CMD}"                                                                           fifo              3 hours ago   Up 51 minutes                
    canopsis-pro-mongodb-1           mongo:7.0.8-jammy                                                         "docker-entrypoint.sh --wiredTigerCacheSizeGB 2.5 --replSet rs0 --keyFile /data/db/.keyFile"   mongodb           4 hours ago   Up 2 hours (healthy)         0.0.0.0:27017->27017/tcp, :::27017->27017/tcp
    canopsis-pro-nginx-1             docker.canopsis.net/docker/develop-community/nginx:24.04-rc3              "/bin/sh -c /entrypoint.sh"                                                                    nginx             3 hours ago   Up 51 minutes                80/tcp, 0.0.0.0:80->8080/tcp, :::80->8080/tcp, 0.0.0.0:443->8443/tcp, :::443->8443/tcp
    canopsis-pro-pbehavior-1         docker.canopsis.net/docker/develop-community/engine-pbehavior:24.04-rc3   "/bin/sh -c /${CMD}"                                                                           pbehavior         3 hours ago   Up 51 minutes                
    canopsis-pro-rabbitmq-1          rabbitmq:3.12.13-management                                               "docker-entrypoint.sh rabbitmq-server"                                                         rabbitmq          4 hours ago   Up 52 minutes (healthy)      4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp, :::5672->5672/tcp, 15671/tcp, 15691-15692/tcp, 25672/tcp, 0.0.0.0:15672->15672/tcp, :::15672->15672/tcp
    canopsis-pro-redis-1             redis:6.2.14-bookworm                                                     "docker-entrypoint.sh /bin/sh -c 'redis-server --requirepass $REDIS_PASSWORD'"                 redis             4 hours ago   Up 52 minutes (healthy)      0.0.0.0:6379->6379/tcp, :::6379->6379/tcp
    canopsis-pro-remediation-1       docker.canopsis.net/docker/develop-pro/engine-remediation:24.04-rc3       "/bin/sh -c /${CMD}"                                                                           remediation       3 hours ago   Up 51 minutes                
    canopsis-pro-snmp-1              docker.canopsis.net/docker/develop-pro/engines-python:24.04-rc3           "/bin/sh -c /entrypoint.sh"                                                                    snmp              3 hours ago   Up 51 minutes                
    canopsis-pro-timescaledb-1       timescale/timescaledb:2.14.2-pg13                                         "docker-entrypoint.sh postgres"                                                                timescaledb       3 hours ago   Up About an hour (healthy)   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp
    canopsis-pro-webhook-1           docker.canopsis.net/docker/develop-pro/engine-webhook:24.04-rc3           "/bin/sh -c /${CMD}"                                                                           webhook           3 hours ago   Up 51 minutes         
    ```

=== "Paquets RHEL 8"

    Mise à jour de Canopsis

    ```sh
    dnf install canopsis-pro-24.04.0 canopsis-webui-24.04.0
    ```

    Reconfiguration de Canopsis

    !!! Attention

        Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette personnalisation.
        L'outil est en effet pré paramétré pour fonctionner naturellement.


    Si vous utilisez un fichier d'override du canopsis.toml, veuillez ajouter à la ligne de commande suivante l'option `-override` suivie du chemin du fichier en question.

    ```sh
    systemctl start postgresql-13
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    /opt/canopsis/bin/canopsis-reconfigure -migrate-postgres=true -migrate-mongo=true -edition pro
    ```

    !!! information "Information"

        Cette opération peut prendre plusieurs minutes pour s'exécuter.

    Vous pouvez ensuite vérifier que le mécanisme de reconfigure s'est correctement déroulé en lisant les logs sur la sortie standard de la commande.

    Redémarrage de Canopsis

    ```sh
    systemctl restart canopsis
    ```

    Vous pouvez ensuite vérifier que l'ensemble des services soient correctement exécutés.

    ```sh
    systemctl status canopsis
    ```

Par ailleurs, le mécanisme de bilan de santé intégré à Canopsis ne doit pas présenter d'erreur.  

![Healthcheck](./img/24.04.0-healthcheck.png)

