# Guide de migration vers Canopsis 24.10.0

Ce guide donne les instructions vous permettant de mettre à jour Canopsis 24.04 (dernière version disponible) vers [la version 24.10.0](../24.10.0.md).

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

    Avant de démarrer la procédure de mise à jour, vous devez vérifier que la valeur de `featureCompatibilityVersion` est bien positionnée à **7.0**  

    === "Docker Compose"
        ```sh
        CPS_EDITION=pro docker compose exec mongodb bash
        mongosh -u root -p root
        > db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )
        > exit
        ```

    === "Paquets RHEL 8"

        ```sh
        mongosh -u root -p root
        > db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )
        > exit
        ```

    === "Helm"

        A venir

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

=== "Helm"

    A venir

## Mise à jour Canopsis

!!! information "Information"

    Canopsis 24.10 est livré avec un nouveau jeu de configurations de référence.
    Vous devez télécharger ces configurations et y reporter vos personnalisations.  

=== "Docker Compose"

    Si vous êtes utilisateur de l'édition `community`, voici les étapes à suivre.

    Télécharger le paquet de la version 24.10.0 (canopsis-community-docker-compose-24.10.0.tar.gz) disponible à cette adresse [https://git.canopsis.net/canopsis/canopsis-community/-/releases](https://git.canopsis.net/canopsis/canopsis-community/-/releases).

    ```sh
    export CPS_EDITION=community
    tar xvfz canopsis-community-docker-compose-24.10.0.tar.gz
    cd canopsis-community-docker-compose-24.10.0
    ```

    Si vous êtes utilisateur de l'édition `pro`, voici les étapes à suivre.

    Télécharger le paquet de la version 24.10.0 (canopsis-pro-docker-compose-24.10.0.tar.gz) disponible à cette adresse [https://git.canopsis.net/sources/canopsis-pro-sources/-/releases](https://git.canopsis.net/sources/canopsis-pro-sources/-/releases).

    ```sh
    export CPS_EDITION=pro
    tar xvfz canopsis-pro-docker-compose-24.10.0.tar.gz
    cd canopsis-pro-docker-compose-24.10.0
    ```

    À ce stade, vous devez synchroniser les modifications réalisées sur vos anciens fichiers de configuration `docker-compose` avec les fichiers `docker-compose.yml` et/ou `docker-compose.override.yml`.

=== "Paquets RHEL 8"

    Non concerné car ces configurations sont livrées directemement dans les paquets RPM.

=== "Helm"

    A venir

### Mise à jour de TimescaleDB

Dans cette version de Canopsis, la base de données TimescaleDB passe de la version 2.14.2 à 2.15.1.  
En plus de la mise à jour de TimescaleDB lui-même, le système de gestion de base de données PostreSQL doit être mis à jour de la version 13 à la version 15.

Deux étapes sont à suivre :

1. Mise à jour de PostgreSQL de 13 vers 15
2. Mise à jour de TimescaleDB 14.2 à 2.15.1


=== "Docker Compose"


    !!! warning "Cas de la base de données des métriques techniques"
         
        Certaines commandes ci-dessous concernent la fonctionalité [métriques techniques](../../../guide-de-depannage/metriques-techniques/). Si vous n'avez pas mis en oeuvre cette fonctionnalité pour le moment, vous pouvez ignorer ces commandes.


    Modifiez la variable `TIMESCALEDB_TAG` du fichier `.env` de cette façon :

    ```diff
    -TIMESCALEDB_TAG=2.15.1-pg15
    +TIMESCALEDB_TAG=2.14.2-pg13
    ```

    Sauvegarde de la base de données

    ```sh
    CPS_EDITION=pro docker compose up -d timescaledb
    CPS_EDITION=pro docker compose exec timescaledb pg_dump postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis -Ft -f /tmp/postgres_dump_archive.tar
    CPS_EDITION=pro docker compose exec timescaledb pg_dump postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis_tech_metrics -Ft -f /tmp/postgres_dump_archive_techmetrics.tar
    CPS_EDITION=pro docker compose cp timescaledb:/tmp/postgres_dump_archive.tar /tmp
    CPS_EDITION=pro docker compose cp timescaledb:/tmp/postgres_dump_archive_techmetrics.tar /tmp
    ```

    Arrêtez le conteneur et supprimez les volumes associés

    ```sh
    CPS_EDITION=pro docker compose down -v timescaledb
    ```

    Modifiez la variable `TIMESCALEDB_TAG` du fichier `.env` de cette façon :

    ```diff
    -TIMESCALEDB_TAG=2.14.2-pg13
    +TIMESCALEDB_TAG=2.14.2-pg15
    ```

    Démarrer le conteneur timescaledb

    ```sh
    CPS_EDITION=pro docker compose up -d timescaledb
    ```

    Restaurez le dump précédemment effectué

    ```sh
    CPS_EDITION=pro docker compose cp /tmp/postgres_dump_archive.tar timescaledb:/tmp/postgres_dump_archive.tar
    CPS_EDITION=pro docker compose cp /tmp/postgres_dump_archive_techmetrcis.tar timescaledb:/tmp/postgres_dump_archive_techmetrics.tar
    CPS_EDITION=pro docker compose exec timescaledb pg_restore --dbname=postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis --no-owner -Ft -v /tmp/postgres_dump_archive.tar
    CPS_EDITION=pro docker compose exec timescaledb pg_restore --dbname=postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis_tech_metrics --no-owner -Ft -v /tmp/postgres_dump_archive_techmetrics.tar
    ```

    Arrêtez le conteneur

    ```sh
    CPS_EDITION=pro docker compose stop timescaledb
    ```

    Modifiez la variable `TIMESCALEDB_TAG` du fichier `.env` de cette façon :

    ```diff
    -TIMESCALEDB_TAG=2.14.2-pg15
    +TIMESCALEDB_TAG=2.15.1-pg15
    ```

    Démarrer le conteneur timescaledb et mettez à jour l'extension

    ```sh
    CPS_EDITION=pro docker compose up -d timescaledb
    CPS_EDITION=pro docker compose exec timescaledb psql postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis
    canopsis=# ALTER EXTENSION timescaledb UPDATE;
    ```

    ```sh
    CPS_EDITION=pro docker compose up -d timescaledb
    CPS_EDITION=pro docker compose exec timescaledb psql postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis_tech_metrics
    canopsis=# ALTER EXTENSION timescaledb UPDATE;
    ```


    Ensuite, vérifiez que l'extension en elle-même est à présent bien à jour

    ```sh
    \dx
    ...
    timescaledb | 2.15.1   | public     | Enables scalable inserts and complex queries for time-series data (Community Edition)
    ...
    exit
    ```


=== "Paquets RHEL 8"

    Arrêter Canopsis

    ```sh
    systemctl stop canopsis-engine-go@engine-action.service \
                       canopsis-engine-go@engine-axe.service \
                       canopsis-engine-go@engine-che.service \
                       canopsis-engine-go@engine-fifo.service \
                       canopsis-engine-go@engine-pbehavior.service \
                       canopsis-service@canopsis-api.service \
                       canopsis.service
    ```

    Sauvegarder les bases de données

    ```sh
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    sudo -iu postgres pg_dump $(eval echo "$CPS_POSTGRES_URL") --no-owner -Fc -v -f /tmp/canopsis-$(date +"%Y-%m-%d")-canopsis-dump.sql.gz
    sudo -iu postgres pg_dump $(eval echo "$CPS_POSTGRES_TECH_URL") --no-owner -Fc -v -f /tmp/canopsis-$(date +"%Y-%m-%d")-canopsis_tech_metrics-dump.sql.gz
    ```

    Stopper l'instance postgresql

    ```sh
    systemctl stop postgresql-13.service
    systemctl disable postgresql-13.service
    ```

    Autoriser dnf à réaliser la montée de version

    ```sh
    dnf versionlock delete timescaledb-2-loader-postgresql-13 timescaledb-2-postgresql-13
    ```

    Installer la version 15 de postgreSQL et la version 2.14 de TimescaleDB

    ```sh
    dnf install timescaledb-2-postgresql-15-2.14.2 timescaledb-2-loader-postgresql-15-2.14.2
    ```

    Initialiser la nouvelle base

    ```sh
    postgresql-15-setup initdb
    timescaledb-tune -yes --pg-config=/usr/pgsql-15/bin/pg_config
    echo "timescaledb.telemetry_level=off" >> /var/lib/pgsql/15/data/postgresql.conf
    systemctl enable --now  postgresql-15.service
    ```

    Créer la base de donnée `cpspostgres` avec les mêmes informations d'identification.

    ```sh
    sudo -iu postgres psql
    postgres=# CREATE database canopsis;
    postgres=# \c canopsis
    canopsis=# CREATE EXTENSION IF NOT EXISTS timescaledb;
    canopsis=# SET password_encryption = 'scram-sha-256';
    canopsis=# CREATE USER cpspostgres WITH PASSWORD 'canopsis';
    canopsis=# exit
    ```

    Créer la base de donnée `canopsis_tech_metrics` avec les mêmes informations d'identification.

    ```sh
    sudo -iu postgres psql
    postgres=# CREATE database canopsis_tech_metrics;
    postgres=# \c canopsis_tech_metrics
    canopsis_tech_metrics=# CREATE EXTENSION IF NOT EXISTS timescaledb;
    canopsis_tech_metrics=# SET password_encryption = 'scram-sha-256';
    canopsis_tech_metrics=# exit
    ```

    Modifier les droits de l'utilisateur cpspostgres pour réaliser les imports

    ```
    sudo -iu postgres psql
    postgres=# ALTER ROLE cpspostgres WITH LOGIN SUPERUSER CREATEDB CREATEROLE REPLICATION BYPASSRLS;
    ```

    Importer les dumps

    ```sh
    sudo -iu postgres pg_restore --no-owner -Fc -v -d $(eval echo "$CPS_POSTGRES_URL") /tmp/canopsis-YYYY-MM-DD-canopsis-dump.sql.gz
    sudo -iu postgres pg_restore --no-owner -Fc -v -d $(eval echo "$CPS_POSTGRES_TECH_URL") /tmp/canopsis-YYYY-MM-DD-canopsis_tech_metrics-dump.sql.gz
    ```

    Réinitialiser les droits des utilisateurs

    ```
    sudo -iu postgres psql
    postgres=# ALTER ROLE cpspostgres WITH LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS;
    ```

    Installer la version `2.15.1` de TimescaleDB

    ```sh
    dnf install timescaledb-2-postgresql-15-2.15.1 timescaledb-2-loader-postgresql-15-2.15.1
    ```

    Se connecter sur l'instance pgsql pour mettre à jour l'extention TimescaleDB:

    ```sh
    sudo -iu postgres psql -X
    postgres=# \c canopsis
    postgres=# ALTER EXTENSION timescaledb UPDATE;
    postgres=# \c canopsis_tech_metrics
    postgres=# ALTER EXTENSION timescaledb UPDATE
    ```

    vérouiller la version pour éviter des mises à jour non souhaitées

    ```sh
    dnf versionlock add timescaledb-2-loader-postgresql-15 timescaledb-2-postgresql-15
    ```

    Supprimer la version 13 de TimescaleDB

    ```sh
    dnf remove timescaledb-2-loader-postgresql-13-2.14.2 timescaledb-2-postgresql-13-2.14.2
    ```

=== "Helm"

    A venir

### Lancement du provisioning `canopsis-reconfigure`

#### Synchronisation du fichier de configuration `canopsis.toml` ou fichier de surcharge

Si vous avez modifié le fichier `canopsis.toml` (vous le voyez via une définition de volume dans votre fichier docker-compose.yml), vous devez vérifier qu'il soit bien à jour par rapport au fichier de référence.  

* [`canopsis.toml` pour Canopsis Community 24.10.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/24.10.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml)
* [`canopsis.toml` pour Canopsis Pro 24.10.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/24.10.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml)

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

=== "Helm"

    A venir

### Mise à jour et démarrage final de Canopsis

Enfin, il vous reste à mettre à jour et à démarrer tous les composants applicatifs de Canopsis

=== "Docker Compose"

    ```sh
    CPS_EDITION=pro docker compose up -d
    ```

    Vous pouvez ensuite vérifier que l'ensemble des conteneurs soient correctement exécutés.

    ```sh
    CPS_EDITION=pro docker compose ps
    NAME                             IMAGE                                                                       COMMAND                  SERVICE           CREATED          STATUS                    PORTS
    canopsis-pro-action-1            docker.canopsis.net/docker/develop-pro/engine-action:24.10.0            "/engine-action"         action            41 seconds ago   Up 28 seconds             
    canopsis-pro-api-1               docker.canopsis.net/docker/develop-pro/canopsis-api-pro:24.10.0         "/bin/sh -c /${CMD}"     api               41 seconds ago   Up 28 seconds (healthy)   0.0.0.0:8082->8082/tcp, :::8082->8082/tcp
    canopsis-pro-axe-1               docker.canopsis.net/docker/develop-pro/engine-axe:24.10.0               "/engine-axe -publis…"   axe               41 seconds ago   Up 28 seconds             
    canopsis-pro-che-1               docker.canopsis.net/docker/develop-pro/engine-che:24.10.0               "/engine-che"            che               41 seconds ago   Up 28 seconds             
    canopsis-pro-connector-junit-1   docker.canopsis.net/docker/develop-pro/connector-junit:24.10.0          "/bin/sh -c /${CMD}"     connector-junit   41 seconds ago   Up 28 seconds             
    canopsis-pro-correlation-1       docker.canopsis.net/docker/develop-pro/engine-correlation:24.10.0       "/bin/sh -c /${CMD}"     correlation       41 seconds ago   Up 28 seconds             
    canopsis-pro-dynamic-infos-1     docker.canopsis.net/docker/develop-pro/engine-dynamic-infos:24.10.0     "/bin/sh -c /${CMD}"     dynamic-infos     41 seconds ago   Up 28 seconds             
    canopsis-pro-fifo-1              docker.canopsis.net/docker/develop-pro/engine-fifo:24.10.0              "/bin/sh -c /${CMD}"     fifo              41 seconds ago   Up 28 seconds             
    canopsis-pro-mongodb-1           mongo:7.0.14-jammy                                                          "docker-entrypoint.s…"   mongodb           2 minutes ago    Up 2 minutes (healthy)    0.0.0.0:27017->27017/tcp, :::27017->27017/tcp
    canopsis-pro-nginx-1             docker.canopsis.net/docker/develop-community/nginx:24.10.0              "/bin/sh -c /entrypo…"   nginx             41 seconds ago   Up 28 seconds             80/tcp, 0.0.0.0:80->8080/tcp, [::]:80->8080/tcp, 0.0.0.0:443->8443/tcp, [::]:443->8443/tcp
    canopsis-pro-pbehavior-1         docker.canopsis.net/docker/develop-community/engine-pbehavior:24.10.0   "/bin/sh -c /${CMD}"     pbehavior         41 seconds ago   Up 28 seconds             
    canopsis-pro-rabbitmq-1          rabbitmq:3.12.13-management                                                 "docker-entrypoint.s…"   rabbitmq          2 minutes ago    Up 2 minutes (healthy)    4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp, :::5672->5672/tcp, 15671/tcp, 15691-15692/tcp, 25672/tcp, 0.0.0.0:15672->15672/tcp, :::15672->15672/tcp
    canopsis-pro-recorder-1          docker.canopsis.net/docker/develop-pro/events-recorder:24.10.0          "/bin/sh -c /${CMD}"     recorder          41 seconds ago   Up 28 seconds             
    canopsis-pro-redis-1             redis:6.2.14-bookworm                                                       "docker-entrypoint.s…"   redis             2 minutes ago    Up 2 minutes (healthy)    0.0.0.0:6379->6379/tcp, :::6379->6379/tcp
    canopsis-pro-remediation-1       docker.canopsis.net/docker/develop-pro/engine-remediation:24.10.0       "/bin/sh -c /${CMD}"     remediation       41 seconds ago   Up 28 seconds             
    canopsis-pro-snmp-1              docker.canopsis.net/docker/develop-pro/engines-python:24.10.0           "/bin/sh -c /entrypo…"   snmp              41 seconds ago   Up 28 seconds             
    canopsis-pro-timescaledb-1       timescale/timescaledb:2.15.1-pg15                                           "docker-entrypoint.s…"   timescaledb       2 minutes ago    Up 2 minutes (healthy)    0.0.0.0:5432->5432/tcp, :::5432->5432/tcp
    canopsis-pro-webhook-1           docker.canopsis.net/docker/develop-pro/engine-webhook:24.10.0           "/bin/sh -c /${CMD}"     webhook           41 seconds ago   Up 28 seconds
    ```

=== "Paquets RHEL 8"

    Mise à jour de Canopsis

    ```sh
    dnf install canopsis-pro-24.10.0 canopsis-webui-24.10.0
    ```

    Reconfiguration de Canopsis

    !!! Attention

        Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette personnalisation.
        L'outil est en effet pré paramétré pour fonctionner naturellement.


    Si vous utilisez un fichier d'override du canopsis.toml, veuillez ajouter à la ligne de commande suivante l'option `-override` suivie du chemin du fichier en question.

    ```sh
    systemctl start postgresql-15
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

=== "Helm"

    A venir

Par ailleurs, le mécanisme de bilan de santé intégré à Canopsis ne doit pas présenter d'erreur.  

![Healthcheck](./img/24.10-healthcheck.png)

