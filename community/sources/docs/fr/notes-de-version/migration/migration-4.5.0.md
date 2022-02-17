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

### Réalisation d'une sauvegarde de vos machines virtuelles

TODO

### Arrêt de l'environnement en cours de lancement

Vous devez prévoir une interruption du service afin de procéder à la mise à jour qui va suivre.

=== "Paquets CentOS 7"

    ```sh
    canoctl stop
    ```

=== "Docker Compose"

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
    yum --disablerepo="*" --enablerepo="timescale_timescaledb,pgdg-common,pgdg13" install timescaledb-2-loader-postgresql-13-2.5.1-0.el7 timescaledb-2-postgresql-13-2.5.1-0.el7
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

    TODO

### Mise à jour de MongoDB

Dans cette version de Canopsis, la base de données MongoDB passe de la version 3.6 à 4.2.

Cette mise à jour doit obligatoirement être réalisée en deux étapes :

1. Mise à jour de MongoDB 3.6 à 4.0 ;
2. **Puis** mise à jour de MongoDB 4.0 à 4.2.

=== "Paquets CentOS 7"

    !!! attention
        Si vous utilisez [un Replicat Set MongoDB](https://docs.mongodb.com/manual/replication/), vous devez obligatoirement suivre une procédure différente :
    
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

    TODO

### Mise à jour de Canopsis

=== "Paquets CentOS 7"

    Appliquez la mise à jour des paquets Canopsis :

    ```sh
    yum --disablerepo="*" --enablerepo="canopsis*" update
    ```

    Note : cette mise à jour de Canopsis introduit un nouveau paquet `canopsis-webui`, identique entre Canopsis Community et Canopsis Pro.

=== "Docker Compose"

    <https://git.canopsis.net/canopsis/canopsis-pro/-/tree/release-4.5/pro/deployment/canopsis/docker>

    TODO

### Lancement des scripts de migration

Assurez-vous que le service MongoDB soit bien lancé et exécutez les commandes suivantes, en adaptant les identifiants MongoDB ci-dessous si nécessaire :

```sh
cd /opt/canopsis/share/migrations/mongodb/release4.5
for file in $(find . -type f -name "*.js" | sort -n); do
   mongo -u cpsmongo -p canopsis canopsis < "$file"
done
```

!!! attention
    Ces scripts essaient de gérer le plus de cas d'usage possible, mais la bonne exécution de ces scripts en toute condition ne peut être garantie.

    Ils doivent obligatoirement être lancés **avant** le lancement des scripts de provisioning lors de l'étape suivante.

    N'hésitez pas à nous signaler tout problème d'exécution que vous pourriez rencontrer lors de cette étape.

### Synchronisation du fichier de configuration `canopsis.toml`

Vérifiez que votre fichier `canopsis.toml` soit bien à jour par rapport au fichier de référence, notamment dans le cas où vous auriez apporté des modifications locales à ce fichier :

<!-- ces liens seront fonctionnels lorsque la version 4.5.0 sera publiée -->
* [`canopsis.toml` pour Canopsis Community 4.5.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/4.5.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml)
* [`canopsis.toml` pour Canopsis Pro 4.5.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/4.5.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml)

=== "Paquets CentOS 7"

    Le fichier à synchroniser est `/opt/canopsis/etc/canopsis.toml`.

=== "Docker Compose"

    Si vous n'avez pas apporté de modification locale, ce fichier est directement intégré et mise à jour dans les conteneurs, et vous n'avez donc pas de modification à apporter.

    Si vous modifiez ce fichier à l'aide d'un volume surchargeant `canopsis.toml`, c'est ce fichier local qui doit être synchronisé.

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

    En revanche, si vous mainteniez vos propres versions modifiées de ces fichiers de configuration, vous devez manuellement vous synchroniser avec la totalité des modifications ayant été apportées dans `/etc/nginx/`.

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

    Vous pouvez ensuite migrer vos métriques existantes de MongoDB vers TimescaleDB avec les commandes suivantes :

    ```sh
    /opt/canopsis/bin/migrate-metrics -onlyMeta
    /opt/canopsis/bin/migrate-metrics
    ```

=== "Docker Compose"

    Exécutez la commande suivante :

    ```sh
    docker-compose -f 01-prov.docker-compose.yml up -d
    ```

    Ou bien, si vous utilisez encore l'ancien procédé :

    ```sh
    docker-compose up -d provisioning reconfigure
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
