
# Guide de migration vers Canopsis 23.04.0

Ce guide donne des instructions vous permettant de mettre à jour Canopsis 22.10 vers [la version 23.04.0](../23.04.0.md).

## Prérequis

L'ensemble de cette procédure doit être lu avant son exécution.

Ce document ne prend en compte que Canopsis Community et Canopsis Pro : tout développement personnalisé dont vous pourriez bénéficier ne fait pas partie du cadre de ce Guide de migration.

Les fichiers de référence qui sont mentionnés dans ce guide sont disponibles à ces adresses

| Edition           | Sources                                                          |
| ----------------- | ---------------------------------------------------------------- |
| Edition Community | https://git.canopsis.net/canopsis/canopsis-community/-/releases  |
| Edition pro       | https://git.canopsis.net/sources/canopsis-pro-sources/-/releases |

## Procédure de mise à jour

### Réalisation d'une sauvegarde

Des sauvegardes sont toujours recommandées, qu'elles soient régulières ou lors de modifications importantes.

La restructuration apportée dans les bases de données pour cette version de Canopsis nous amène à insister d'autant plus sur ce point. Il est donc fortement recommandé de réaliser une **sauvegarde complète** des VM hébergeant vos services Canopsis, avant cette mise à jour.


### Vérification MongoDB

!!! warning "Vérification"

    Avant de démarrer la procédure de mise à jour, vous devez vérifier que la valeur de `featureCompatibilityVersion` est bien positionnée à **4.4**  

    === "Docker Compose"
        ```sh
        CPS_EDITION=pro docker-compose exec mongodb bash
        mongo -u root -p root
        > db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )
        > exit
        ```

    === "Paquets RH8"

        ```sh
        mongo -u root -p root
        > db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )
        > exit
        ```
    
    Le retour doit être de la forme `{ "featureCompatibilityVersion" : { "version" : "4.4" }, "ok" : 1 }`
    Si ce n'est pas le cas, vous ne pouvez pas continuer la mise à jour.

### Arrêt de l'environnement en cours de lancement

Vous devez prévoir une interruption du service afin de procéder à la mise à jour qui va suivre.

=== "Docker Compose"

    ```sh
    CPS_EDITION=pro docker-compose down
    ```

=== "Paquets RH8"

    ```sh
    systemctl stop canopsis
    systemctl stop mongod
    systemctl stop postgresql-13
    systemctl stop rabbitmq-server
    systemctl stop redis
    ```

### Mise à jour Canopsis

!!! information "Information"

    Canopsis 23.04 est livré avec un nouveau jeu de configurations de référence.
    Vous devez télécharger ces configurations et y reporter vos personnalisations.  


=== "Docker Compose"

    Si vous êtes utilisateur de l'édition `community`, voici les étapes à suivre.

    ```sh
    export CPS_EDITION=community
    mkdir 23.04
    cd 23.04
    git clone https://git.canopsis.net/canopsis/canopsis-community.git -b release-23.04
    cd canopsis-community/community/deployment/canopsis/docker/
    ```

    Si vous êtes utilisateur de l'édition `pro`, voici les étapes à suivre.

    ```sh
    export CPS_EDITION=pro
    mkdir 23.04
    cd 23.04
    git clone https://git.canopsis.net/sources/canopsis-pro-sources.git -b release-23.04
    cd canopsis-pro-sources/pro/deployment/canopsis/docker/
    ```

    A ce stade, vous devez synchroniser les modifications réalisées sur vos anciens fichiers de configuration `docker-compose` avec les fichiers `docker-compose.yml` et/ou `docker-compose-override.yml`.

=== "Paquets RH8"

    A venir

### Mise à jour de MongoDB

Dans cette version de Canopsis, la base de données MongoDB passe de la version 4.4 à 5.0.

=== "Docker Compose"

    Démarrez le conteneur `mongodb` :

    ```sh
    CPS_EDITION=pro docker-compose up -d mongodb
    ```

    Entrez ensuite à l'intérieur de ce conteneur, afin de compléter la mise à jour vers MongoDB 5.0 :

    ```sh
    CPS_EDITION=pro docker-compose exec mongodb bash
    mongo -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "5.0" } )
    exit
    ```

=== "Paquets RH8"

    !!! note
        Si vous avez mis en place des exclusions dans le fichier `/etc/yum.conf`, veillez à la désactiver le temps de cette procédure.

    Mise à jour des paquets `mongodb` :

    ```sh
    echo '[mongodb-org-5.0]
    name=MongoDB Repository
    baseurl=https://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/5.0/x86_64/
    gpgcheck=1
    enabled=1
    gpgkey=https://www.mongodb.org/static/pgp/server-5.0.asc' | tee /etc/yum.repos.d/mongodb-org-5.0.repo
    dnf makecache
    dnf install mongodb-org-5.0.16 mongodb-org-database-5.0.16 mongodb-org-server-5.0.16 mongodb-org-mongos-5.0.16 mongodb-org-tools-5.0.16 mongodb-org-shell-5.0.16
    ```

    Redémarrage de `mongodb` :

    ```sh
    systemctl start mongod
    ```

    Ensuite, complétez la mise à jour vers MongoDB 5.0 :

    ```sh
    mongo -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "5.0" } )
    exit
    ```

### Mise à jour de TimescaleDB

Dans cette version de Canopsis, la base de données TimescaleDB passe de la version 2.7.2 à 2.9.3.

=== "Docker Compose"

    Relancez le conteneur `timescaledb` :

    ```sh
    CPS_EDITION=pro docker-compose up -d timescaledb
    ```

    Puis mettez à jour l'extension timescaledb (La chaîne de connexion doit être adaptée à votre environnement)

    ```sh
    CPS_EDITION=pro docker-compose exec timescaledb psql postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis
    canopsis=# ALTER EXTENSION timescaledb UPDATE;
    ```

    Ensuite, vérifiez que l'extension en elle-même est à présent bien à jour

    ```sh
    \dx
    ...
    timescaledb | 2.9.3   | public     | Enables scalable inserts and complex queries for time-series data
    ...
    exit
    ```
   

=== "Paquets RH8"

    Mise à jour des paquets `timescaledb` :

    ```sh
    dnf install timescaledb-2-postgresql-13-2.9.3 timescaledb-2-loader-postgresql-13-2.9.3
    ```

    Relancez le service `timescaledb` :

    ```sh
    systemctl start timescaledb
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
    timescaledb | 2.9.3   | public     | Enables scalable inserts and complex queries for time-series data
    ...
    exit
    ```

### Mise à jour de RabbitMQ

Dans cette version de Canopsis, le bus rabbitMQ passe à la version 3.11.11.  

=== "Docker Compose"
    
    Un passage par la version 3.10.20 est nécessaire lorsque votre installation part de la version 3.7.28.

    Passage en version 3.10.20 puis lancement du conteneur `rabbitmq` :

    ```sh
    sed -i "s/RABBITMQ_TAG=.*-management$/RABBITMQ_TAG=3.10.20-management/g" .env
    CPS_EDITION=pro docker-compose up -d rabbitmq
    ```

    Il faut à présent activer la fonctionnalité `FEATURE_FLAGS` de rabbitmq :

    ```sh
    CPS_EDITION=pro docker-compose exec rabbitmq /bin/bash
    rabbitmqctl enable_feature_flag all
    exit
    ```

    Passage en version 3.11.11 puis lancement du conteneur `rabbitmq` :

    ```sh
    sed -i "s/RABBITMQ_TAG=.*-management$/RABBITMQ_TAG=3.11.11-management/g" .env
    CPS_EDITION=pro docker-compose up -d rabbitmq
    ```

=== "Paquets RH8"

    Il faut commencer par activer la fonctionnalité `FEATURE_FLAGS` de rabbitmq :

    ```sh
    systemctl start rabbitmq-server
    rabbitmqctl enable_feature_flag all
    ```
   
    Passage en version 3.11.11 puis lancement du service `rabbitmq-server` :

    ```sh
    dnf install --repo rabbitmq_erlang --repo rabbitmq_server erlang rabbitmq-server-3.11.11
    systemctl restart rabbitmq-server
    ```

### Remise à 0 du cache Redis

Dans cette version de Canopsis, le cache de Canopsis doit repartir à 0.

=== "Docker Compose"

    ```sh
    CPS_EDITION=pro docker-compose up -d redis
    CPS_EDITION=pro docker-compose exec redis /usr/local/bin/redis-cli flushall
    ```

=== "Paquets RH8"

    ```sh
    systemctl start redis
    /bin/redis-cli flushall
    ```

### Lancement du provisioning `canopsis-reconfigure`

#### Synchronisation du fichier de configuration `canopsis.toml` ou fichier de surcharge

Si vous avez modifié le fichier `canopsis.toml` (vous le voyez via une définition de volume dans votre fichier docker-compose.yml), vous devez vérifier qu'il soit bien à jour par rapport au fichier de référence.  

* [`canopsis.toml` pour Canopsis Community 23.04.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/23.04.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml)
* [`canopsis.toml` pour Canopsis Pro 23.04.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/23.04.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml)

!!! information "Information"

    Pour éviter ce type de synchronisation fastidieuse, la bonne pratique est d'utiliser [un fichier de surcharge de cette configuration](../../../guide-administration/administration-avancee/modification-canopsis-toml/). 


Si vous avez utilisé un fichier de surcharge, alors vous n'avez rien à faire, uniquement continuer à le présenter dans un volume.

#### Reconfiguration de Canopsis

=== "Docker Compose"

    !!! Attention
    
    Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette personnalisation.
    L'outil est en effet pré paramétré pour fonctionner naturellement.

    ```sh
    CPS_EDITION=pro docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d reconfigure
    ```

    !!! information "Information"

        Cette opération peut prendre plusieurs minutes pour s'exécuter.

    Vous pouvez ensuite vérifier que le mécanisme de provisioning/reconfigure s'est correctement déroulé. Le conteneur doit présenté un "exit 0"

    ```sh
    CPS_EDITION=pro docker-compose -f docker-compose.yml -f docker-compose.override.yml ps
    canopsis-pro-reconfigure-1            "/canopsis-reconfigu…"   reconfigure            exited (0)
    ```

=== "Paquets RH8"

    La commande `canopsis-reconfigure` doit être exécutée après mise à jour de Canopsis dans le cadre d'installation par paquets RPM.

#### Migration des statistiques de remédiation

Historiquement, les statistiques de remédiation étaient stockées dans la base mongoDB.  
A présent, c'est TimescaleDB qui porte ces statistiques.

Une migration des données est donc nécessaire.

=== "Docker Compose"

    Exécution de la commande `migrate-instruction-metrics`

    ```sh
    A CORRIGER avec un patch
    docker run -e EXTERNAL_API_USERNAME='[testuser]' -e EXTERNAL_API_PASSWORD='[testpassword]' -e CPS_MONGO_URL='mongodb://cpsmongo:canopsis@mongodb/canopsis' -e CPS_AMQP_URL='amqp://cpsrabbit:canopsis@rabbitmq/canopsis' -e CPS_POSTGRES_URL='postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis' --network=canopsis-pro_default -it --rm docker.canopsis.net/docker/pro/migrate-instruction-metrics:23.04-preview3
    ```

=== "Paquets RH8"

    Exécution de la commande `migrate-instruction-metrics`

    ```sh
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    /opt/canopsis/bin/migrate-instruction-metrics
    ```

#### Mise à jour et démarrage final de Canopsis

Enfin, il vous reste à mettre à jour et à démarrer tous les composants applicatifs de Canopsis

=== "Docker Compose"

    ```sh
    CPS_EDITION=pro docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d
    ```

    Vous pouvez ensuite vérifier que l'ensemble des conteneurs soient correctement exécutés.

    ```sh
    CPS_EDITION=pro docker-compose -f docker-compose.yml -f docker-compose.override.yml ps
    NAME                                  COMMAND                  SERVICE                STATUS              PORTS
    canopsis-pro-action-1                 "/engine-action -wit…"   action                 running             
    canopsis-pro-api-1                    "/canopsis-api -docs"    api                    running (healthy)   0.0.0.0:8082->8082/tcp, :::8082->8082/tcp
    canopsis-pro-axe-1                    "/engine-axe -publis…"   axe                    running             
    canopsis-pro-che-1                    "/engine-che"            che                    running             
    canopsis-pro-connector-junit-1        "/bin/sh -c /${CMD}"     connector-junit        running             
    canopsis-pro-correlation-1            "/bin/sh -c /${CMD}"     correlation            running             
    canopsis-pro-dynamic-infos-1          "/bin/sh -c /${CMD}"     dynamic-infos          running             
    canopsis-pro-fifo-1                   "/bin/sh -c /${CMD}"     fifo                   running             
    canopsis-pro-migrate-metrics-meta-1   "/bin/true /migrate-…"   migrate-metrics-meta   exited (0)          
    canopsis-pro-mongodb-1                "docker-entrypoint.s…"   mongodb                running (healthy)   0.0.0.0:27027->27017/tcp, :::27027->27017/tcp
    canopsis-pro-nginx-1                  "/bin/sh -c /entrypo…"   nginx                  running             80/tcp, 0.0.0.0:80->8080/tcp, :::80->8080/tcp, 0.0.0.0:443->8443/tcp, :::443->8443/tcp
    canopsis-pro-pbehavior-1              "/bin/sh -c /${CMD}"     pbehavior              running             
    canopsis-pro-rabbitmq-1               "docker-entrypoint.s…"   rabbitmq               running (healthy)   4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp, :::5672->5672/tcp, 15671/tcp, 15691-15692/tcp, 25672/tcp, 0.0.0.0:15672->15672/tcp, :::15672->15672/tcp
    canopsis-pro-reconfigure-1            "/canopsis-reconfigu…"   reconfigure            exited (0)          
    canopsis-pro-redis-1                  "docker-entrypoint.s…"   redis                  running (healthy)   0.0.0.0:6379->6379/tcp, :::6379->6379/tcp
    canopsis-pro-remediation-1            "/bin/sh -c /${CMD}"     remediation            running             
    canopsis-pro-service-1                "/engine-service -pu…"   service                running             
    canopsis-pro-timescaledb-1            "docker-entrypoint.s…"   timescaledb            running (healthy)   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp
    ```

=== "Paquets RH8"

    Mise à jour de Canopsis 

    ```sh
    dnf install canopsis-pro-23.04.0 canopsis-webui-23.04.0
    ```

    Reconfiguration de Canopsis 

    !!! Attention
    
        Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette personnalisation.
        L'outil est en effet pré paramétré pour fonctionner naturellement.
        

    Si vous utilisez un fichier d'override du canopsis.toml, veuillez ajouter à la ligne de commande suivante l'option `-override` suivie du chemin du fichier en question.

    ```sh
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

![Healthcheck](./img/23.04.0-healthcheck.png)

### Moteur SNMP et Linkbuilder

Toutes les parties Python de Canopsis ont été supprimées dans la version 22.10.  
Les fonctionnalités SNMP et Linkbuilder n'ont pas encore été migrées en GO et réintégrées dans Canopsis.  
Pour continuer à les utiliser, vous devez utiliser les packages des versions 4.x.  
Cela est temporaire mais nécessaire.

Que ce soit sur une installation par paquets RPM ou par image Docker, nous vous livrons des configurations à exécuter avec Docker.

Vous devez simplement définir la variable d'environnement `CPS_OLD_API` dans `/opt/canopsis/etc/go-engines-vars.conf` pour une installation paquets ou `compose.env` en installation docker.

!!! warning "Avertissement"

    Les configurations `docker-compose` livrées ne doivent absolument pas être modifiées, notamment les versions d'images utilisées.  
    Par ailleurs, si vous utilisiez déjà la partie SNMP dans votre installation alors la collection `schema` existe très certainement dans mongoDB.  
    Si ce n'était pas le cas, vous allez devoir la créer avec ces instructions.  
    ```
    docker-compose exec snmp /bin/bash
    schema2db
    ```

```sh
cd ../../../mock/external-services/snmp/docker
docker-compose up -d 
```
