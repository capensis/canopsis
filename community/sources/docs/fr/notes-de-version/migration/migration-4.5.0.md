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

Installation des paquets :

```sh
yum makecache -y
yum --disablerepo="*" --enablerepo="timescale_timescaledb,pgdg-common,pgdg13" install timescaledb-2-loader-postgresql-13-2.5.1-0.el7 timescaledb-2-postgresql-13-2.5.1-0.el7
```

Configuration pour le système courant et désactivation de la télémétrie :

```sh
postgresql-13-setup initdb
yes | PATH=/usr/pgsql-13/bin:$PATH timescaledb-tune -pg-config /usr/pgsql-13/bin/pg_config -out-path /var/lib/pgsql/13/data/postgresql.conf -yes
echo "timescaledb.telemetry_level=off" >> /var/lib/pgsql/13/data/postgresql.conf
```

Activation et démarrage du service :

```sh
systemctl enable postgresql-13
systemctl start postgresql-13
```

!!! info "Information"
    Si vous avez besoin d'accéder à PostgreSQL/TimescaleDB depuis une autre machine, autorisez l'accès au port TCP 5432 (uniquement pour les administrateurs de la plateforme).

### Mise à jour de MongoDB

Dans cette version de Canopsis, la base de données MongoDB passe de la version 3.6 à 4.2.

Cette mise à jour doit obligatoirement être réalisée en deux étapes :

1. Mise à jour de MongoDB 3.6 à 4.0 ;
2. **Puis** mise à jour de MongoDB 4.0 à 4.2.

=== "Paquets CentOS 7"

    !!! attention
        Si vous utilisez [un Replicat Set MongoDB](https://docs.mongodb.com/manual/replication/), vous devez obligatoirement suivre une procédure différente :
    
        1. D'abord <https://docs.mongodb.com/manual/release-notes/4.0-upgrade-replica-set/>
        2. puis <https://docs.mongodb.com/manual/release-notes/4.2-upgrade-replica-set/>.

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
cd /TODO/chemin/vers/les/scripts/de/migration/intégrés
for file in $(find release4.5 -type f -name "*.js" | sort -n); do
   mongo -u cpsmongo -p canopsis canopsis < "$file"
done
```

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

### Ajustement des binaireis `engine-che` et `engine-axe` (paquets)

=== "Paquets CentOS 7"

    Cette partie s'applique seulement aux installations de paquets RPM.

    ```sh
    rm -f /opt/canopsis/bin/TODO

    # si Canopsis Community :
    ln -sf /opt/canopsis/bin/TODO-community /opt/canopsis/bin/TODO
    # OU si Canopsis Pro :
    ln -sf /opt/canopsis/bin/TODO-pro /opt/canopsis/bin/TODO
    ```

=== "Docker Compose"

    Aucune manipulation n'est nécessaire ici.

### Mise à jour de la configuration de Nginx

Plusieurs changements ont été apportés à la configuration de Nginx.

=== "Paquets CentOS 7"

    TODO : de nouvelles variables ?

    Exécutez les commandes suivantes afin de prendre en compte ces changements :

    ```sh
    cp -f /opt/canopsis/deploy-ansible/playbook/roles/canopsis/templates/nginx/cors.j2 /etc/nginx/cors.inc
    cp -f /opt/canopsis/deploy-ansible/playbook/roles/canopsis/templates/nginx/https.j2 /etc/nginx/https.inc
    sed -e 's,{{ CPS_API_URL }},http://127.0.0.1:8082,g' -e 's,{{ CPS_SERVER_NAME }},"localhost",g' /opt/canopsis/deploy-ansible/playbook/roles/canopsis/templates/nginx/default.j2 > /etc/nginx/conf.d/default.conf

    systemctl restart nginx
    ```

    !!! attention
        Si vous accédez à l'interface web de Canopsis au travers d'un nom de domaine (par exemple `canopsis.mon-si.fr`), vous devrez **obligatoirement** configurer la ligne `set $canopsis_server_name` du fichier `/etc/nginx/conf.d/default.conf` avec le nom de domaine concerné.

=== "Docker Compose"

    Si vous n'avez pas surchargé la configuration Nginx à l'aide d'un volume, vous n'avez rien à faire.

    En revanche, si vous mainteniez vos propres versions modifiées de ces fichiers de configuration, vous devez manuellement vous synchroniser avec la totalité des modifications ayant été apportées dans `/etc/nginx/`.

### Lancement du provisioning et de `canopsis-reconfigure`

Le *provisioning* doit être lancé afin de mettre à jour certaines données en base, tandis que `canopsis-reconfigure` prend en compte les changements apportés au fichier `canopsis.toml`.

=== "Paquets CentOS 7"

    Lancez les scripts de provisioning :

    ```sh
    # si vous utilisez Canopsis Community
    su - canopsis -c "canopsinit --canopsis-edition core"
    # OU si vous utilisez Canopsis Pro
    su - canopsis -c "canopsinit --canopsis-edition cat"
    ```

    Puis, lancez `canopsis-reconfigure` :

    ```bash
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    /opt/canopsis/bin/canopsis-reconfigure TODO nouvelles-options!
    ```

    TODO: Profitez-en aussi pour enlever la variable CPS_INFLUX inutile dans go-engines-vars.conf.

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
