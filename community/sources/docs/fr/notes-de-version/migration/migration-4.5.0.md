# Guide de migration vers Canopsis 4.5.0

Ce guide donne des instructions vous permettant de mettre à jour Canopsis 4.4 vers [la version 4.5.0](../4.5.0.md).

## Prérequis

L'ensemble de cette procédure doit être lu avant son exécution.

Ce document ne prend en compte que Canopsis Community et Canopsis Pro : tout développement personnalisé dont vous pourriez bénéficier ne fait pas partie du cadre de ce Guide de migration.

### Mise à jour des prérequis Docker

En installation Docker, vous devez vous assurer d'utiliser Docker 20.10 ou une version ultérieure. Docker 19.03 (ou toute version antérieure) n'est donc plus pris en charge.

Exécutez la commande suivante afin de connaître votre version actuelle de Docker :

```sh
docker version --format '{{.Server.Version}}'
```

Si la valeur que vous obtenez commence par 20.10 ou une valeur supérieure, vous n'avez rien à faire. Si votre version est trop ancienne, [suivez la documentation officielle de Docker](https://docs.docker.com/get-docker/) afin d'installer une version à jour.

Assurez-vous aussi que Docker Compose [ne soit pas configuré en mode V2](../../guide-administration/installation/installation-conteneurs.md#installation-de-docker-et-docker-compose) avec la commande suivante :

```sh
docker-compose version --short
```

Si vous obtenez une version supérieure ou égale à 2.0.0, exécutez la commande suivante :

```sh
docker-compose disable-v2
```

## Procédure de mise à jour

### Réalisation d'une sauvegarde

Des sauvegardes sont toujours recommandées, qu'elles soient régulières ou lors de modifications importantes.

La restructuration apportée dans les bases de données pour cette version de Canopsis nous amène à insister d'autant plus sur ce point. Il est donc fortement recommandé de réaliser une **sauvegarde complète** des VM hébergeant vos services Canopsis, avant cette mise à jour.

### Arrêt de l'environnement en cours de lancement

Vous devez prévoir une interruption du service afin de procéder à la mise à jour qui va suivre.

=== "Paquets CentOS 7"

    ```sh
    canoctl stop
    ```

=== "Docker Compose"

    ```sh
    docker-compose -f 00-data.docker-compose.yml -f 01-prov.docker-compose.yml -f 02-app.docker-compose.yml down
    ```
    
    Ou bien, si vous utilisez encore l'ancien procédé :
    
    ```sh
    docker-compose down
    ```

### Suppression d'InfluxDB

InfluxDB a été totalement supprimé de Canopsis. Ce composant tiers peut donc être totalement supprimé, avec son contenu.

=== "Paquets CentOS 7"

    Exécutez les commandes suivantes :
    
    ```sh
    systemctl stop influxdb.service
    yum remove influxdb
    rm -rf /etc/influxdb /usr/lib/influxdb /var/lib/influxdb /etc/logrotate.d/influxdb
    ```

=== "Docker Compose"

    Supprimez toute variable contenant le terme `INFLUXDB` dans les fichiers `.env` et `compose.env`.
    
    Puis, enlevez toutes références au volume `influxdbdata` et au conteneur `influxdb` présentes dans votre fichier de référence Docker Compose.

Pensez aussi à révoquer toute ouverture réseau que vous auriez autorisée à destination des ports TCP 8086 et 8088.

### Ajout de TimescaleDB

=== "Paquets CentOS 7"

    Ajout des clés GPG et des dépôts TimescaleDB et PostgreSQL :
    
    ```sh
    yum install https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm
    cd /etc/pki/rpm-gpg/ && curl -L -o RPM-GPG-KEY-PKGCLOUD-TIMESCALEDB https://packagecloud.io/timescale/timescaledb/gpgkey
    
    # Add TimescaleDB repo
    cat > /etc/yum.repos.d/timescale_timescaledb.repo << EOF
    [timescale_timescaledb]
    name=timescale_timescaledb
    baseurl=https://packagecloud.io/timescale/timescaledb/el/\$releasever/\$basearch
    repo_gpgcheck=1
    # timescaledb doesn't sign all its packages
    gpgcheck=0
    enabled=1
    gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-PKGCLOUD-TIMESCALEDB
    sslverify=1
    sslcacert=/etc/pki/tls/certs/ca-bundle.crt
    metadata_expire=300
    EOF
    ```
    
    Installation des paquets TimescaleDB et des dépendances associées :
    
    ```sh
    yum makecache -y
    yum --disablerepo="*" --enablerepo="timescale_timescaledb,pgdg-common,pgdg13,base" install timescaledb-2-loader-postgresql-13-2.5.1-0.el7 timescaledb-2-postgresql-13-2.5.1-0.el7
    ```
    
    Configuration pour le système actuel, et désactivation de [la télémétrie](https://docs.timescale.com/timescaledb/latest/how-to-guides/configuration/telemetry/#telemetry-and-version-checking) :
    
    ```sh
    postgresql-13-setup initdb
    yes | PATH=/usr/pgsql-13/bin:$PATH timescaledb-tune -pg-config /usr/pgsql-13/bin/pg_config -out-path /var/lib/pgsql/13/data/postgresql.conf -yes
    echo "timescaledb.telemetry_level=off" >> /var/lib/pgsql/13/data/postgresql.conf
    ```
    
    Activation et démarrage du service :
    
    ```sh
    systemctl enable postgresql-13.service
    systemctl start postgresql-13.service
    ```
    
    Connexion à la ligne de commande PostgreSQL :
    
    ```sh
    sudo -u postgres psql
    ```
    
    Création de la base de données `canopsis` et de l'utilisateur associé, et activation de l'extension TimescaleDB sur cette base :
    
    ```sql
    postgres=# CREATE database canopsis;
    postgres=# \c canopsis
    canopsis=# CREATE EXTENSION IF NOT EXISTS timescaledb;
    canopsis=# CREATE USER cpspostgres WITH PASSWORD 'canopsis';
    canopsis=# exit
    ```
    
    !!! info "Information"
        Si vous avez besoin d'accéder à PostgreSQL/TimescaleDB depuis une autre machine, autorisez l'accès au port TCP 5432 (uniquement pour les administrateurs de la plateforme).

=== "Docker Compose"

    Dans le fichier `.env` lié à votre environnement de référence Docker Compose, ajoutez la ligne suivante :
    
    ```ini
    TIMESCALEDB_TAG=2.5.1-pg13
    ```
    
    Ajoutez ensuite les lignes suivantes au fichier `compose.env` :
    
    ```ini
    CPS_POSTGRES_URL=postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis
    POSTGRES_USER=cpspostgres
    POSTGRES_PASSWORD=canopsis
    POSTGRES_DB=canopsis
    ```
    
    Puis, ajoutez les lignes suivantes dans la section `services:`, à la suite des services RabbitMQ, MongoDB, Redis :
    
    ```yaml
    timescaledb:
      image: timescale/timescaledb:${TIMESCALEDB_TAG}
      ports:
        - "5432:5432"
      env_file:
        - compose.env
      environment:
        - TIMESCALEDB_TELEMETRY=off
      volumes:
        - timescaledata:/var/lib/postgresql/data
      restart: unless-stopped
      shm_size: 1g
    ```
    
    Enfin, déclarez un volume `timescaledata` dans la section `volumes:` du même fichier (à la suite des volumes `mongodbdata` et `rabbitmqdata`) :
    
    ```yaml
    timescaledata:
      driver: local
    ```
    
    Démarrez ensuite ce nouveau conteneur `timescaledb`.

### Mise à jour de MongoDB

Dans cette version de Canopsis, la base de données MongoDB passe de la version 3.6 à 4.2.

!!! attention
    Cette mise à jour doit **impérativement** être réalisée en deux étapes :
    
    1. Mise à jour de MongoDB 3.6 à 4.0 ;
    2. **Puis** mise à jour de MongoDB 4.0 à 4.2.

=== "Paquets CentOS 7"

    !!! important
        Si vous utilisez un [Replica Set MongoDB](https://docs.mongodb.com/manual/replication/), vous devez obligatoirement suivre une procédure différente :
    
        1. D'abord <https://docs.mongodb.com/manual/release-notes/4.0-upgrade-replica-set/> ;
        2. **puis** <https://docs.mongodb.com/manual/release-notes/4.2-upgrade-replica-set/>.
    
    Exécutez les commandes suivantes pour passer à MongoDB 4.0 :
    
    ```sh
    sed -i 's|3\.6|4.0|g' /etc/yum.repos.d/mongodb.repo
    
    yum makecache -y
    yum --disablerepo="*" --enablerepo="mongodb*" update
    
    mongo -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "4.0" } )
    > exit
    ```
    
    Puis, exécutez les commandes suivantes pour passer à MongoDB 4.2 :
    
    ```sh
    sed -i 's|4\.0|4.2|g' /etc/yum.repos.d/mongodb.repo
    
    yum makecache -y
    yum --disablerepo="*" --enablerepo="mongodb*" update
    
    mongo -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "4.2" } )
    > exit
    ```

=== "Docker Compose"

    Modifiez la variable `MONGO_TAG` du fichier `.env` de cette façon :
    
    ```diff
    -MONGO_TAG=3.6.17-xenial
    +MONGO_TAG=4.0.28-xenial
    ```
    
    Puis relancez le conteneur `mongodb` :

    ```sh
    docker-compose -f 00-data.docker-compose.yml up -d mongodb
    ```
    
    Entrez ensuite à l'intérieur de ce conteneur, afin de compléter la mise à jour vers MongoDB 4.0 :
    
    ```sh
    docker-compose -f 00-data.docker-compose.yml exec mongodb bash
    mongo -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "4.0" } )
    exit
    ```
    
    Puis, éditez à nouveau la variable `MONGO_TAG` du fichier `.env` comme suit :
    
    ```diff
    -MONGO_TAG=4.0.28-xenial
    +MONGO_TAG=4.2.18-bionic
    ```
    
    Relancez à nouveau le conteneur `mongodb` :

    ```sh
    docker-compose -f 00-data.docker-compose.yml up -d mongodb
    ```
    
    Entrez à nouveau à l'intérieur du conteneur, afin de finaliser la mise à jour vers MongoDB 4.2 :
    
    ```sh
    docker-compose -f 00-data.docker-compose.yml exec mongodb bash
    mongo -u root -p root
    > db.adminCommand( { setFeatureCompatibilityVersion: "4.2" } )
    exit
    ```

### Mise à jour de Canopsis

=== "Paquets CentOS 7"

    Appliquez la mise à jour des paquets Canopsis :
    
    ```sh
    yum --disablerepo="*" --enablerepo="canopsis*" update
    ```
    
    Note : cette mise à jour de Canopsis introduit un nouveau paquet `canopsis-webui`, identique entre Canopsis Community et Canopsis Pro.

=== "Docker Compose"

    Si et seulement si vous utilisez Canopsis Pro, vous devez maintenant utiliser les moteurs `engine-che` et `engine-axe` propres à Canopsis Pro.
    
    Pour cela, modifiez les lignes `image:` de ces 2 moteurs de la façon suivante :
    
    ```diff
       axe:
    -    image: ${DOCKER_REPOSITORY}${COMMUNITY_BASE_PATH}engine-axe:${CANOPSIS_IMAGE_TAG}
    +    image: ${DOCKER_REPOSITORY}${PRO_BASE_PATH}engine-axe:${CANOPSIS_IMAGE_TAG}
         env_file:
           - compose.env
         restart: unless-stopped
         command: /engine-axe -publishQueue Engine_correlation -withRemediation=true
    
       che:
    -    image: ${DOCKER_REPOSITORY}${COMMUNITY_BASE_PATH}engine-che:${CANOPSIS_IMAGE_TAG}
    +    image: ${DOCKER_REPOSITORY}${PRO_BASE_PATH}engine-che:${CANOPSIS_IMAGE_TAG}
         env_file:
           - compose.env
         restart: unless-stopped
         command: /engine-che -enrichContext
    ```

    Dans la partie provisioning (en-dessous de l'image `reconfigure`), ajoutez aussi les lignes suivantes (Canopsis Pro uniquement) :

    ```yaml
    migrate-metrics-meta:
      image: ${DOCKER_REPOSITORY}${PRO_BASE_PATH}migrate-metrics:${CANOPSIS_IMAGE_TAG}
      env_file:
        - compose.env
      # warning: only -onlyMeta is safe to be run by default
      command: /migrate-metrics -onlyMeta
    ```
    
    Enfin, dans tous les cas, mettez à jour la variable `CANOPSIS_IMAGE_TAG` de votre fichier `.env` pour passer à Canopsis 4.5.0 :
    
    ```ini
    CANOPSIS_IMAGE_TAG=4.5.0
    ```

### Suppression de l'option `-featureStatEvents`

L'option `-featureStatEvents` a été retirée du moteur `engine-axe`.

=== "Paquets CentOS 7"

    Lancez la commande suivante afin de savoir si cette option est utilisée :
    
    ```sh
    grep -lr "featureStatEvents" /etc/systemd/system/canopsis-engine-go@engine-axe.service.d/*
    ```
    
    Si cette commande affiche un résultat, éditez les fichiers qu'elle mentionne afin d'y retirer cette option.

=== "Docker Compose"

    Supprimez toute éventuelle utilisation de l'option `-featureStatEvents` dans vos fichiers de référence Docker Compose.

### Lancement des scripts de migration

Assurez-vous que le service MongoDB soit bien lancé et exécutez les commandes suivantes, en adaptant les identifiants MongoDB ci-dessous si nécessaire :

=== "Paquets CentOS 7"

    Sur la machine sur laquelle les paquets `canopsis*` sont installés :
    
    ```sh
    cd /opt/canopsis/share/migrations/mongodb/release4.5
    for file in $(find . -type f -name "*.js" | sort -n); do
       mongo -u cpsmongo -p canopsis canopsis < "$file"
    done
    ```

=== "Docker Compose"

    Depuis une machine qui a un client `mongo` installé et qui peut joindre le service `mongodb` d'un point de vue réseau :
    
    ```sh
    git clone --depth 1 --single-branch -b release-4.5 https://git.canopsis.net/canopsis/canopsis-community.git
    cd canopsis-community/community/go-engines-community/database/migrations
    for file in $(find release4.5 -type f -name "*.js" | sort -n); do
       mongo "mongodb://cpsmongo:canopsis@localhost:27017/canopsis" < "$file" # URI à adapter au besoin
    done
    ```
    
    Il est aussi possible de récupérer le répertoire `migrations` et de le présenter en volume dans le conteneur `mongodb` afin de réaliser le lancement du script depuis le conteneur `mongodb`.

!!! attention
    Ces scripts essaient de gérer le plus de cas d'usage possible, mais la bonne exécution de ces scripts en toute condition ne peut être garantie.

    Ils doivent obligatoirement être lancés **avant** le lancement des scripts de provisioning lors de l'étape suivante.
    
    N'hésitez pas à nous signaler tout problème d'exécution que vous pourriez rencontrer lors de cette étape.

### Synchronisation du fichier de configuration `canopsis.toml`

Vérifiez que votre fichier `canopsis.toml` soit bien à jour par rapport au fichier de référence, notamment dans le cas où vous auriez apporté des modifications locales à ce fichier :

* [`canopsis.toml` pour Canopsis Community 4.5.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/4.5.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml)
* [`canopsis.toml` pour Canopsis Pro 4.5.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/4.5.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml)

=== "Paquets CentOS 7"

    Le fichier à synchroniser est `/opt/canopsis/etc/canopsis.toml`.

=== "Docker Compose"

    Si vous n'avez pas apporté de modification locale, ce fichier est directement intégré et mise à jour dans les conteneurs, et vous n'avez donc pas de modification à apporter.
    
    Si vous modifiez ce fichier à l'aide d'un volume surchargeant `canopsis.toml`, c'est ce fichier local qui doit être synchronisé.

### Ajustements de la configuration (paquets)

Cette partie s'applique seulement aux installations de paquets RPM.

=== "Paquets CentOS 7"

    Les binaires `engine-che` et `engine-axe` sont maintenant différents entre Canopsis Community et Canopsis Pro.
    
    Dans le cadre d'une mise à jour, ce changement implique la création d'un lien symbolique, à adapter en fonction de l'édition que vous utilisez :
    
    ```sh
    rm -f /opt/canopsis/bin/engine-axe /opt/canopsis/bin/engine-che
    
    # si Canopsis Community :
    ln -sf /opt/canopsis/bin/engine-axe-community /opt/canopsis/bin/engine-axe
    ln -sf /opt/canopsis/bin/engine-che-community /opt/canopsis/bin/engine-che
    
    # OU si Canopsis Pro :
    ln -sf /opt/canopsis/bin/engine-axe-pro /opt/canopsis/bin/engine-axe
    ln -sf /opt/canopsis/bin/engine-che-pro /opt/canopsis/bin/engine-che
    ```
    
    Supprimez ensuite toute ligne `CPS_INFLUX_URL` du fichier `go-engines-vars.conf` et ajoutez-y `CPS_POSTGRES_URL` :
    
    ```sh
    sed -i '/CPS_INFLUX_URL/d' /opt/canopsis/etc/go-engines-vars.conf
    grep -q ^CPS_POSTGRES_URL= /opt/canopsis/etc/go-engines-vars.conf || echo 'CPS_POSTGRES_URL="postgresql://cpspostgres:canopsis@localhost:5432/canopsis"' >> /opt/canopsis/etc/go-engines-vars.conf
    ```

=== "Docker Compose"

    Aucune manipulation n'est nécessaire ici.

### Mise à jour de la configuration de Nginx

Plusieurs changements ont été apportés à la configuration de Nginx.

=== "Paquets CentOS 7"

    Exécutez les commandes suivantes afin de prendre en compte ces changements (en remplaçant `"localhost"` par le FQDN de votre service Canopsis si nécessaire) :
    
    ```sh
    cp -p /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf.orig
    sed \
        -e 's,{{ CPS_API_URL }},http://127.0.0.1:8082,g' \
        -e 's,{{ CPS_OLD_API_URL }},http://127.0.0.1:8081,g' \
        -e 's,{{ CPS_SERVER_NAME }},"localhost",g' \
        /opt/canopsis/deploy-ansible/playbook/roles/canopsis/templates/nginx/default.j2 > /etc/nginx/conf.d/default.conf
    
    systemctl restart nginx
    ```

=== "Docker Compose"

    Si vous n'avez pas surchargé la configuration Nginx à l'aide d'un volume, vous n'avez rien à faire.
    
    En revanche, si vous mainteniez vos propres versions modifiées de ces fichiers de configuration, vous devez manuellement vous synchroniser avec la totalité des modifications ayant été apportées au fichier `/etc/nginx/conf.d/default.conf`.

### Lancement du provisioning et de `canopsis-reconfigure`

Le provisioning doit être lancé afin de mettre à jour certaines données en base, tandis que `canopsis-reconfigure` prend en compte les changements apportés au fichier `canopsis.toml`.

=== "Paquets CentOS 7"

    Lancez les scripts de provisioning :
    
    ```sh
    # si vous utilisez Canopsis Community
    su - canopsis -c "canopsinit --canopsis-edition core"
    # OU si vous utilisez Canopsis Pro
    su - canopsis -c "canopsinit --canopsis-edition cat"
    ```
    
    Puis, lancez `canopsis-reconfigure`. Attention, cette fois-ci de nouvelles options doivent lui être données :
    
    ```bash
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    /opt/canopsis/bin/canopsis-reconfigure -migrate-postgres=true -postgres-migration-mode=up -postgres-migration-directory=/opt/canopsis/share/migrations/postgres
    ```
    
    Initialisez ensuite vos métriques TimescaleDB avec cette commande (Canopsis Pro uniquement) :
    
    ```sh
    /opt/canopsis/bin/migrate-metrics -onlyMeta
    ```

    Puis, migrez vos métriques existantes de MongoDB à TimescaleDB avec cette commande (Canopsis Pro uniquement) :

    ```sh
    # Attention : ne doit être exécuté qu'une seule fois !
    /opt/canopsis/bin/migrate-metrics
    ```

=== "Docker Compose"

    Lancez à nouveau toute la partie `data` (MongoDB, RabbitMQ, Redis, PostgreSQL…) :

    ```sh
    docker-compose -f 00-data.docker-compose.yml up -d
    ```

    !!! Attention
        Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette persionnalisation.  
        L'outil est en effet pré paramétré pour fonctionner naturellement.

    Exécutez la commande suivante :
    
    ```sh
    docker-compose -f 01-prov.docker-compose.yml up -d
    ```
    
    Ou bien, si vous utilisez encore l'ancien procédé :
    
    ```sh
    docker-compose up -d provisioning reconfigure
    # Canopsis Pro uniquement
    docker-compose up -d migrate-metrics-meta
    ```
    
    Vous pouvez ensuite migrer vos métriques existantes de MongoDB vers TimescaleDB avec cette commande (Canopsis Pro uniquement) :
    
    ```sh
    # Note : cette commande ne doit être exécutée qu'une seule fois et ne doit donc pas être ajoutée aux fichiers YAML
    # option --network à adapter si nécessaire
    docker run -it --rm --env-file compose.env --network=canopsis-pro_default docker.canopsis.net/docker/pro/migrate-metrics:4.5.0
    ```

    Le retour de la commande doit ressembler à cela

    ```sh
    WRN git.canopsis.net/canopsis/canopsis-community/community/go-engines-community@v0.0.0/lib/mongo/mongo.go:550 > MongoDB version does not support transactions, transactions are disabled
    INF git.canopsis.net/canopsis/canopsis-pro/pro/go-engines-pro/cmd/migrate-metrics/main.go:39 > entities migration finished 19.941129479s
    INF git.canopsis.net/canopsis/canopsis-pro/pro/go-engines-pro/cmd/migrate-metrics/main.go:44 > users migration finished 20.516699ms
    INF git.canopsis.net/canopsis/canopsis-pro/pro/go-engines-pro/cmd/migrate-metrics/main.go:50 > alarm metrics migration finished 1m16.718781144s
    ```

### Remise en route des moteurs et des services de Canopsis

Si et seulement si les commandes précédentes n'ont pas renvoyé d'erreur, vous pouvez relancer l'intégralité des services.

=== "Paquets CentOS 7"

    Relancez la totalité de l'environnement :
    
    ```sh
    systemctl daemon-reload
    canoctl restart
    ```

=== "Docker Compose"

    Lancez maintenant la partie `02-app`, afin de bénéficier de l'application Canopsis en elle-même :
    
    ```sh
    docker-compose -f 02-app.docker-compose.yml up -d
    ```
    
    Ou bien, si vous utilisez encore l'ancien procédé :
    
    ```sh
    docker-compose up -d
    ```

### Fin de la mise à jour

Après quelques minutes, le service devrait être à nouveau accessible sur son interface web habituelle. En cas de problème, consultez l'ensemble des logs.

Suivez la [section « Après la mise à jour »](../../guide-administration/mise-a-jour/index.md#apres-la-mise-a-jour) du Guide d'administration afin d'en savoir plus.
