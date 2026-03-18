# Guide de migration vers Canopsis 25.10.0

Ce guide donne les instructions vous permettant de mettre à jour Canopsis 25.04 (dernière version disponible) vers [la version 25.10.0](../25.10.0.md).

## Sommaire

[TOC]

## Prérequis

L'ensemble de cette procédure doit être lue avant son exécution.

Ce document ne prend en compte que Canopsis Community et Canopsis Pro : tout développement personnalisé dont vous pourriez bénéficier ne fait pas partie du cadre de ce Guide de migration.

Les fichiers de référence qui sont mentionnés dans ce guide sont disponibles à ces adresses

| Édition           | Sources                                                                                                                              |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| Édition Community | [https://git.canopsis.net/canopsis/canopsis-community/-/releases](https://git.canopsis.net/canopsis/canopsis-community/-/releases)   |
| Édition pro       | [https://git.canopsis.net/sources/canopsis-pro-sources/-/releases](https://git.canopsis.net/sources/canopsis-pro-sources/-/releases) |

## Réalisation d'une sauvegarde

Des sauvegardes sont toujours recommandées, qu'elles soient régulières ou lors de modifications importantes.

## Mise à jour de MongoDB

### Vérifications MongoDB

!!! warning "Vérification"

    Avant de démarrer la procédure de mise à jour, vous devez vérifier que la valeur de `featureCompatibilityVersion` est bien positionnée à **7.0**  

```sh
mongosh -u root -p root
> db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )
> exit
```

### Mise à jour

=== "RPM"

    Pour commencer, il faut couper le service

    ```sh
    systemctl stop mongod.service
    ```

    Installer le dépôt pour la version `8.0` de MongoDB

    ```sh
    cat << EOF > /etc/yum.repos.d/mongodb-org-8.0.repo
    [mongodb-org-8.0]
    name=MongoDB Repository
    baseurl=https://repo.mongodb.org/yum/redhat/\$releasever/mongodb-org/8.0/x86_64/
    gpgcheck=1
    enabled=1
    gpgkey=https://www.mongodb.org/static/pgp/server-8.0.asc
    EOF
    ```

    Lancer l'upgrade

    ```sh
    dnf upgrade mongodb-org mongodb-org-database mongodb-org-server mongodb-org-mongos mongodb-org-tools -y
    ```

    Une fois mis à jour, le service doit être relancé
    ```sh
    systemctl start mongod.service
    ```

    Il faut ensuite se connecter avec l'utilisateur `root` et terminer la mise à jour
    ```sh
    mongosh -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "8.0", confirm: true } )
    exit
    ```

    Après avoir mis à jour mongodb, l'option de telemetry sera activée. Pour la désactiver, exécutez la commande suivante :

    ```sh
    mongosh -u root -p root
    > disableTelemetry()
    ```

=== "Docker Compose"

    Dans cette version de Canopsis, la base de données MongoDB passe de la version 7.0 à 8.0.  

    Démarrez le conteneur `mongodb` :

    ```sh
    CPS_EDITION=pro docker compose up -d mongodb
    ```

    Entrez ensuite à l'intérieur de ce conteneur, afin de compléter la mise à jour vers MongoDB 8.0 :

    ```sh
    CPS_EDITION=pro docker compose exec mongodb bash
    mongosh -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "8.0", "confirm" : true } )
    ```

    Après avoir mis à jour mongodb, l'option de telemetry sera activée. Pour la désactiver, exécutez la commande suivante :

    ```sh
    CPS_EDITION=pro docker compose exec mongodb bash
    mongosh -u root -p root
    > disableTelemetry()
    exit
    ```

=== "Helm"

    !!! warning Attention
            Ce bloc est réservé uniquement aux environnements impliquant MongoDB exécuté dans un environnement Kubernetes.
            
            Si ce n'est pas votre cas, référez-vous au bloc [Paquets RPM](#__tabbed_1_2)

    Dump de la base de données Canopsis :
    ```sh
    kubectl exec -n canopsis canopsis-mongodb-0 -- mongodump --uri="mongodb://cpsmongo:canopsis@localhost:27017/canopsis" --gzip --out /tmp/dump_canopsis.gz
    ```

    Récupération en local du dump :
    ```sh
    kubectl cp canopsis/canopsis-mongodb-0:/tmp/dump_canopsis.gz .
    ```

    Arrêt des pods MongoDB :
    ```sh
    kubectl get statefulset --no-headers=true |grep mongodb| awk {'print $1'}| xargs kubectl delete statefulset
    ```

    Suppresion des PVCs MongoDB :
    ```sh
    kubectl get pvc --no-headers=true | awk '{print $1}' | grep mongodb | xargs kubectl delete pvc
    ```

    !!! warning Attention
        Veillez à bien adapter la commande ci-dessous avec vos paramètres présent dans votre fichier de surcharges, par exemple customer-values.yml.

    Mise à jour de MongoDB :
    ```sh
    helm repo update

    helm upgrade canopsis canopsis/mongodb \
    --set enabled=true \
    --set settings.rootUsername=root \
    --set settings.rootPassword=root \
    --set replicaSet.enabled=true \
    --set replicaSet.name=rs0 \
    --set replicaSet.secondaries=0 \
    --set replicaSet.key="c29tZXJhbmRvbXN0cmluZzEyMzQ1Ng==" \
    --set userDatabase.name=canopsis \
    --set userDatabase.user=cpsmongo \
    --set userDatabase.password=canopsis \
    --set storage.keepPvc=true \
    --set storage.requestedSize=8Gi \
    --version 1.1.0
    ```

    Lorsque le POD est UP, copie du dump de la DB sur l'instance 0 de MongoDB : 
    ```sh
    kubectl cp ./canopsis canopsis-mongodb-0:/tmp/
    ```

    Restauration du dump :
    ```sh
    kubectl exec -n canopsis canopsis-mongodb-0 -- mongorestore -u cpsmongo --password canopsis --gzip --db canopsis /tmp/canopsis
    ``` 

    Suppression du Statefulset :
    ```sh
    kubectl get statefulset --no-headers=true | awk '{print $1}' | grep mongo | xargs kubectl delete statefulset
    ```

## Mise à jour TimescaleDB

Dans cette version de Canopsis, la base de données TimescaleDB passe de la version 2.15.1 à 2.21.4.  
En plus de la mise à jour de TimescaleDB lui-même, le système de gestion de base de données PostreSQL doit être mis à jour de la version 15 à la version 17.

Deux étapes sont à suivre :

1. Mise à jour de TimescaleDB 2.15.1 vers 2.21.4
2. Mise à jour de PostgreSQL 15 vers 17

=== "RPM"

    Dans un premier temps, on sauvegarde les bases de données `canopsis` et `canopsis_tech_metrics`

    ```sh
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    sudo -iu postgres pg_dump $(eval echo "$CPS_POSTGRES_URL") --no-owner -Fc -v -f /tmp/canopsis-$(date +"%Y-%m-%d")-canopsis-dump.sql.gz
    sudo -iu postgres pg_dump $(eval echo "$CPS_POSTGRES_TECH_URL") --no-owner -Fc -v -f /tmp/canopsis-$(date +"%Y-%m-%d")-canopsis_tech_metrics-dump.sql.gz
    ```

    On peut ensuite arrêter le service et le désactiver

    ```sh
    systemctl stop postgresql-15.service
    systemctl disable postgresql-15.service
    ```

    On peut maintenant supprimer l'éventuel `versionlock` relatif à cette brique, et mettre à jour `timescaledb` avant de poser un `versionlock` sur la nouvelle version

    ```sh
    dnf versionlock delete timescaledb-2-loader-postgresql-15 timescaledb-2-postgresql-15
    dnf install timescaledb-2-postgresql-17-2.21.4 timescaledb-2-loader-postgresql-17-2.21.4 -y
    dnf versionlock add --raw timescaledb-2-loader-postgresql-17 timescaledb-2-postgresql-17
    ```

    Il faut ensuite initialiser la nouvelle instance PostgreSQL 17

    ```sh
    postgresql-17-setup initdb
    timescaledb-tune -yes --pg-config=/usr/pgsql-17/bin/pg_config
    echo "timescaledb.telemetry_level=off" >> /var/lib/pgsql/17/data/postgresql.conf
    ```

    Définir la zone de temps de la base de donnée à `UTC`: (Nécessaire pour le bon fonctionnement de Canopsis)
    ```sh
    sed -i "s/^#\?timezone.*/timezone = 'UTC'/" /var/lib/pgsql/17/data/postgresql.conf
    sed -i "s/^#\?log_timezone.*/log_timezone = 'UTC'/" /var/lib/pgsql/17/data/postgresql.conf
    ```

    La réactiver au boot et on vérifie son bon démarrage
    ```sh
    systemctl enable --now postgresql-17.service
    systemctl status postgresql-17.service
    ```

    Il faut ensuite recréer les bases de données ainsi que les utilisateurs associés

    ```sql 
    sudo -iu postgres psql
    postgres=# CREATE database canopsis;
    postgres=# \c canopsis
    canopsis=# CREATE EXTENSION IF NOT EXISTS timescaledb;
    canopsis=# SET password_encryption = 'scram-sha-256';
    canopsis=# CREATE USER cpspostgres WITH PASSWORD 'canopsis';
    canopsis=# GRANT ALL ON DATABASE canopsis TO cpspostgres;
    canopsis=# ALTER DATABASE canopsis OWNER TO cpspostgres;
    canopsis=# exit
    ```

    ```sql
    sudo -iu postgres psql
    postgres=# CREATE database canopsis_tech_metrics;
    postgres=# \c canopsis_tech_metrics
    canopsis_tech_metrics=# CREATE EXTENSION IF NOT EXISTS timescaledb;
    canopsis_tech_metrics=# SET password_encryption = 'scram-sha-256';
    canopsis_tech_metrics=# CREATE USER cpspostgres_tech_metrics WITH PASSWORD 'canopsis';
    canopsis_tech_metrics=# GRANT ALL ON DATABASE canopsis_tech_metrics TO cpspostgres_tech_metrics;
    canopsis_tech_metrics=# ALTER DATABASE canopsis_tech_metrics OWNER TO cpspostgres_tech_metrics;
    canopsis_tech_metrics=# exit        
    ```

    Puis vérifier que l'extension Timescaledb est bien à jour
    ```sql
    postgres=# \c canopsis
    canopsis=# \dx timescaledb
                                                    List of installed extensions
        Name     | Version | Schema |                                      Description                                      
    -------------+---------+--------+---------------------------------------------------------------------------------------
    timescaledb | 2.21.4  | public | Enables scalable inserts and complex queries for time-series data (Community Edition)
    (1 row)
    ```

    On modifie les droits de l'utilisateur `cpspostgres` et `cpspostgres_tech_metrics` pour réaliser les imports

    ```sql
    sudo -iu postgres psql
    postgres=# ALTER ROLE cpspostgres WITH LOGIN SUPERUSER CREATEDB CREATEROLE REPLICATION BYPASSRLS;
    postgres=# ALTER ROLE cpspostgres_tech_metrics WITH LOGIN SUPERUSER CREATEDB CREATEROLE REPLICATION BYPASSRLS;
    ```

    Et on réimporte les données

    ```sh
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    sudo -iu postgres pg_restore --no-owner -Fc -v -d $(eval echo "$CPS_POSTGRES_URL") /tmp/canopsis-YYYY-MM-DD-canopsis-dump.sql.gz
    sudo -iu postgres pg_restore --no-owner -Fc -v -d $(eval echo "$CPS_POSTGRES_TECH_URL") /tmp/canopsis-YYYY-MM-DD-canopsis_tech_metrics-dump.sql.gz
    ```

    Réinitialiser les droits des utilisateurs
    ```sql
    sudo -iu postgres psql
    postgres=# ALTER ROLE cpspostgres WITH LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS;
    postgres=# ALTER ROLE cpspostgres_tech_metrics WITH LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS;
    ```

=== "Docker Compose"

    Modifiez la variable `TIMESCALEDB_TAG` du fichier `.env` de cette façon :

    ```diff
    -TIMESCALEDB_TAG=2.21.4-pg17
    +TIMESCALEDB_TAG=2.21.4-pg15
    ```

    Démarrez le conteneur et mettez à jour l'extension TimescaleDB

    ```sh
    CPS_EDITION=pro docker compose up -d timescaledb
    CPS_EDITION=pro docker compose exec timescaledb psql postgresql://postgres:canopsis@timescaledb:5432/canopsis
    canopsis=# ALTER EXTENSION timescaledb UPDATE;
    canopsis=# \c canopsis_tech_metrics
    canopsis=# ALTER EXTENSION timescaledb UPDATE;
    ```

    Vérifiez que la version de l'extension soit bien mise à jour

    ```sh
    canopsis=# \dx timescaledb
                                                    List of installed extensions
        Name     | Version | Schema |                                      Description                                      
    -------------+---------+--------+---------------------------------------------------------------------------------------
        timescaledb | 2.21.4  | public | Enables scalable inserts and complex queries for time-series data (Community Edition)
    (1 row)
    ```

    Sauvegarde des bases de données

    ```sh
    CPS_EDITION=pro docker compose exec timescaledb pg_dump postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis -Ft -f /tmp/postgres_dump_archive.tar
    CPS_EDITION=pro docker compose exec timescaledb pg_dump postgresql://cpspostgres_tech_metrics:canopsis@timescaledb:5432/canopsis_tech_metrics -Ft -f /tmp/postgres_dump_archive_techmetrics.tar
    CPS_EDITION=pro docker compose cp timescaledb:/tmp/postgres_dump_archive.tar /tmp
    CPS_EDITION=pro docker compose cp timescaledb:/tmp/postgres_dump_archive_techmetrics.tar /tmp
    ```

    Arrêtez le conteneur et supprimez les volumes associés

    ```sh
    CPS_EDITION=pro docker compose down -v timescaledb
    ```

    Modifiez la variable `TIMESCALEDB_TAG` du fichier `.env` de cette façon :

    ```diff
    -TIMESCALEDB_TAG=2.21.4-pg15
    +TIMESCALEDB_TAG=2.21.4-pg17
    ```

    Démarrer le conteneur timescaledb

    ```sh
    CPS_EDITION=pro docker compose up -d timescaledb
    ```

    Restaurez le dump précédemment effectué

    ```sh
    CPS_EDITION=pro docker compose cp /tmp/postgres_dump_archive.tar timescaledb:/tmp/postgres_dump_archive.tar
    CPS_EDITION=pro docker compose cp /tmp/postgres_dump_archive_techmetrics.tar timescaledb:/tmp/postgres_dump_archive_techmetrics.tar
    CPS_EDITION=pro docker compose exec timescaledb pg_restore --dbname=postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis --no-owner -Ft -v /tmp/postgres_dump_archive.tar
    CPS_EDITION=pro docker compose exec timescaledb pg_restore --dbname=postgresql://cpspostgres_tech_metrics:canopsis@timescaledb:5432/canopsis_tech_metrics --no-owner -Ft -v /tmp/postgres_dump_archive_techmetrics.tar
    ```

=== "Helm"

    Sauvegarder la base de données Canopsis :

    ```sh
    kubectl exec canopsis-timescaledb-0 -- pg_dump postgresql://cpspostgres:canopsis@canopsis-timescaledb:5432/canopsis -Ft -f /tmp/postgres_canopsis_dump.tar

    kubectl cp canopsis-timescaledb-0:/tmp/postgres_canopsis_dump.tar postgres_canopsis_dump.tar
    ```

    Supprimer le Statefulset ainsi que le PVC :

    ```sh
    kubectl get statefulset --no-headers=true | awk '{print $1}' | grep timescaledb | xargs kubectl delete statefulset

    kubectl get pvc --no-headers=true | awk '{print $1}' | grep timescaledb | xargs kubectl delete pvc
    ```

    Déploiement de PostgreSQL 17 - TimescaleDB 2.21.4

    ```sh
    kubectl apply -f - <<EOF
    apiVersion: apps/v1
    kind: StatefulSet
    metadata:
        name: canopsis-timescaledb
    spec:
        selector:
        matchLabels:
            app.kubernetes.io/name: canopsis-pro
        serviceName: canopsis-timescaledb-headless
        updateStrategy:
        type: RollingUpdate
        template:
        metadata:
            labels:
            app.kubernetes.io/name: canopsis-pro
        spec:
            containers:
            - name: timescaledb
                image: docker.io/timescale/timescaledb:2.21.4-pg17
                ports:
                - containerPort: 5432
                env:
                - name: TIMESCALEDB_TELEMETRY
                    value: "off"
                - name: POSTGRES_DB
                    value: "canopsis"
                - name: POSTGRES_USER
                    value: "cpspostgres"
                - name: POSTGRES_PASSWORD
                    valueFrom:
                    secretKeyRef:
                        name: canopsis-timescaledb
                        key: timescaledb-password
                readinessProbe:
                exec:
                    command:
                    - /bin/bash
                    - -c
                    - pg_isready -d \$POSTGRES_DB -U \$POSTGRES_USER
                initialDelaySeconds: 5
                periodSeconds: 10
                timeoutSeconds: 5
                volumeMounts:
                - name: datadir
                    mountPath: /var/lib/postgresql/data
            imagePullSecrets:
            - name: canopsisregistry
        volumeClaimTemplates:
        - metadata:
            name: datadir
            annotations:
                helm.sh/resource-policy: "keep"
            spec:
            accessModes:
                - ReadWriteOnce
            resources:
                requests:
                storage: 8Gi
    EOF
    ```

    Restauration du dump :

    ```sh
    kubectl cp postgres_canopsis_dump.tar canopsis-timescaledb-0:/tmp

    kubectl exec canopsis-timescaledb-0 -- pg_restore --dbname=postgresql://cpspostgres:canopsis@canopsis-timescaledb-0:5432/canopsis --no-owner -Ft -v /tmp/postgres_canopsis_dump.tar
    ```

    Suppression du Statefulset :

    ```sh
    kubectl get statefulset --no-headers=true | awk '{print $1}' | grep timescaledb | xargs kubectl delete statefulset
    ```

## Mise à jour RabbitMQ

=== "RPM"

    Pour mettre à jour RabbitMQ, on retire en premier lieu le `versionlock` de la version 4.0 de `rabbitmq-server`, puis on le réapplique pour la version `4.1` :

    ```sh
    dnf versionlock delete 'rabbitmq-server-4.0*'
    dnf versionlock add --raw 'rabbitmq-server-4.1*'
    ```

    Il faut couper le service

    ```sh
    systemctl stop rabbitmq-server.service
    ```

    Une fois fait, on peut lancer la mise à jour : 

    ```sh
    dnf upgrade rabbitmq-server -y
    ```

    Le service doit être démarrer

    ```sh
    systemctl start rabbitmq-server.service
    ```

    On peut ensuite vérifier la version de `rabbitmq` :

    ```sh
    rabbitmqctl --version
    4.1.5
    ```

=== "Docker Compose"

    Dans cette version de Canopsis, le bus de données RabbitMQ passe de la version 4.0 à 4.1.  

    Les configurations de référence fournies dans Canopsis embarquent déjà ce changement.  
    Vous devez néanmoins activer l'ensemble des ["feature flags"](https://www.rabbitmq.com/docs/feature-flags) stables.  

    Démarrez le conteneur `rabbitmq` :

    ```sh
    CPS_EDITION=pro docker compose up -d rabbitmq
    ```

    Puis activez les "feature flags" :

    ```sh
    CPS_EDITION=pro docker compose exec rabbitmq rabbitmqctl enable_feature_flag all
    ```

=== "Helm"

    Dans cette version de Canopsis, le bus de données RabbitMQ passe de la version 4.0 à 4.1.  

    Supprimer le volume associé à RabbitMQ

    ```sh
    kubectl get pvc --no-headers=true | awk '{print $1}' | grep rabbitmq | xargs kubectl delete pvc
    ```

## Mise à jour Valkey (Helm seulement)

=== "Helm"

    Valkey ne s’appuie plus sur les charts Bitnami, mais sur nos propres charts maintenus en interne.
    Il est donc nécessaire de supprimer le StatefulSet ainsi que le PVC associé.

    Supprimer le Statefulset Valkey

    ```sh
    kubectl get statefulset |grep valkey| awk {'print $1'}| xargs kubectl delete statefulset
    ```

    Supprimer le PVC Valkey

    ```sh
    kubectl get pvc |grep valkey| awk {'print $1'}| xargs kubectl delete pvc
    ```

## Mise à jour de Canopsis

### Synchronisation du fichier de configuration `canopsis.toml` ou fichier de surcharge

Si vous avez modifié le fichier `canopsis.toml` (vous le voyez via une définition de volume dans votre fichier docker-compose.yml), vous devez vérifier qu'il soit bien à jour par rapport au fichier de référence.  

* [`canopsis.toml` pour Canopsis Community 25.10.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/25.10.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml)
* [`canopsis.toml` pour Canopsis Pro 25.10.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/25.10.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml)

!!! information "Information"

    Pour éviter ce type de synchronisation fastidieuse, la bonne pratique est d'utiliser [un fichier de surcharge de cette configuration](../../../guide-administration/administration-avancee/modification-canopsis-toml/). 

Si vous avez utilisé un fichier de surcharge, alors vous n'avez rien à faire, uniquement continuer à le présenter dans un volume.

### Séparation des flux d’événements par initiateur

En version 25.04, les flags suivants avaient été dépréciés, ils sont à présent obsolètes.  
Toute référence doit être supprimée dans vos configurations. Les moteurs concernés ne démareront pas sans cela.  

* -publishQueue
* -consumeQueue
* -workers (remplacé par des flags spécifiques à chaque type de flux)

| Moteur               | Flags nouveaux     | Valeur par défaut | Flags obsolètes                             |
|----------------------|--------------------|-------------------|---------------------------------------------|
| engine-fifo          | -workers           | 10                | -publishQueue, -consumeQueue                |
| engine-che           | -externalWorkers   | 4                 | -workers, -publishQueue, -consumeQueue      |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
| engine-axe           | -externalWorkers   | 4                 | -workers, -publishQueue                     |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |
| engine-correlation   | -externalWorkers   | 4                 | -workers, -publishQueue, -consumeQueue      |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |
| engine-dynamic-infos | -externalWorkers   | 4                 | -workers, -publishQueue                     |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |
| engine-action        | -externalWorkers   | 4                 | -workers                                    |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |

### Lancement du provisioning `canopsis-reconfigure` et terminer la mise à jour.

!!! Attention

    Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette personnalisation.
    L'outil est en effet pré paramétré pour fonctionner naturellement.

=== "RPM"

    Si vous utilisez un fichier d'override du canopsis.toml, veuillez ajouter à la ligne de commande suivante l'option `-override` suivie du chemin du fichier en question.

    Vérifiez que les services `postgresql-17.service` et `mongod.service` sont bien démarrés puis lancer le reconfigure
    ```sh
    systemctl status postgresql-17 mongod
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    /opt/canopsis/bin/canopsis-reconfigure -migrate-postgres=true -migrate-tech-postgres=true -migrate-mongo=true -edition pro
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

    ![Healthcheck](./img/25.10.0-healthcheck.png)

=== "Docker Compose"

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

    Enfin, il vous reste à mettre à jour et à démarrer tous les composants applicatifs de Canopsis

    ```sh
    CPS_EDITION=pro docker compose up -d
    ```

    Vous pouvez ensuite vérifier que l'ensemble des conteneurs soient correctement exécutés.

    ```sh
    CPS_EDITION=pro docker compose ps
    ```

    Par ailleurs, le mécanisme de bilan de santé intégré à Canopsis ne doit pas présenter d'erreur.  

    ![Healthcheck](./img/25.10.0-healthcheck.png)

=== "Helm"

    Définir le nom de votre instance

    ```sh
    export RELEASE_NAME="canopsis-prod"
    ```

    Mise à jour des repos helm

    ```sh
    helm repo update
    ```

    Mise à jour de Canopsis

    ```sh
    helm upgrade ${RELEASE_NAME} canopsis/canopsis-pro -f customer-values.yaml
    ```

    Par ailleurs, le mécanisme de bilan de santé intégré à Canopsis ne doit pas présenter d'erreur.  

    ![Healthcheck](./img/25.10.0-healthcheck.png)