# Installation de paquets Canopsis sur Red Hat Enterprise Linux 8

!!! Information
    Si vous souhaitez réaliser une mise à jour, la procédure est décrite dans le [guide de mise à jour](../mise-a-jour/index.md).

Cette procédure décrit l'installation de Canopsis en mono-instance à l'aide de paquets RHEL 8. Les binaires sont compilés pour l'architecture x86-64.

!!! Warning
    L'installation de tous les services sur un seul nœud, décrite ici, permet
    de monter aisément un environnement de test.  
    Cependant, **ce type d'installation ne saurait être conseillé pour de la
    production**.

    Une plate-forme Canopsis de production repose généralement sur plusieurs
    nœuds, avec les mécanismes de haute disponibilité appropriés. Ce type de
    déploiement plus avancé est supporté uniquement dans le cadre
    d'une souscription à l'édition Pro de Canopsis.

    [Capensis][capensis] peut vous accompagner dans la mise en place de
    l'architecture Canopsis distribuée et des solutions de haute disponibilité.

[capensis]: https://www.capensis.fr/

Toutes les commandes dans cette procédure doivent être invoquées avec
l'utilisateur `root`.

## Prérequis

Assurez-vous d'avoir suivi les [prérequis réseau et de sécurité](../administration-avancee/configuration-parefeu-et-selinux.md), notamment concernant la désactivation de SELinux.

Pour vérifier l'état de SELinux :

```sh
sestatus
```

L'installation nécessite l'ajout de dépôts RPM tiers, ainsi qu'un accès HTTP et HTTPS pour le téléchargement de diverses dépendances. Plus de détails dans la [matrice des flux réseau](../matrice-des-flux-reseau/index.md).

!!! Information
    Notez que les versions de MongoDB, RabbitMQ, Redis et TimescaleDB dont l'installation est décrite ici sont les seules validées pour fonctionner avec Canopsis.

    Plus de détails sur les [prérequis des versions](prerequis-des-versions.md).

## Dépendances

### Mise à jour système

Assurez-vous que le système est à jour (l'installation sur RHEL 8 suppose que le système est relié à des dépôts à jour de la distribution, et en particulier pas figé dans une ancienne version mineure 8.x.) :

```sh
dnf update
```

### Configuration système

Vous pouvez vérifier les limites de ressources systèmes avec la commande suivante :

```sh
ulimit -a
```

Pour appliquer la [configuration recommandée par le projet MongoDB](https://www.mongodb.com/docs/v7.0/reference/ulimit/), créez le fichier `/etc/security/limits.d/mongo.conf` :

```sh
cat << EOF > /etc/security/limits.d/mongo.conf
#<domain>      <type>  <item>         <value>
mongo           soft    fsize           unlimited
mongo           soft    cpu             unlimited
mongo           soft    as              unlimited
mongo           soft    memlock         unlimited
mongo           hard    nofile          64000
mongo           hald    nproc           64000
EOF
```

Désactivez la gestion des `Transparent Huge Pages (THP)` selon la [préconisation MongoDB](https://www.mongodb.com/docs/manual/tutorial/transparent-huge-pages/)

```sh
cat << EOF > /etc/systemd/system/disable-transparent-huge-pages.service
[Unit]
Description=Disable Transparent Huge Pages (THP)
DefaultDependencies=no
After=sysinit.target local-fs.target
Before=mongod.service
[Service]
Type=oneshot
ExecStart=/bin/sh -c 'echo never | tee /sys/kernel/mm/transparent_hugepage/enabled > /dev/null'
[Install]
WantedBy=basic.target
EOF
```

```sh
systemctl daemon-reload
systemctl enable --now disable-transparent-huge-pages
```

### Ajout des dépôts tiers

Ajout du dépôt pour PostgreSQL :

```sh
dnf install https://download.postgresql.org/pub/repos/yum/reporpms/EL-8-x86_64/pgdg-redhat-repo-latest.noarch.rpm
```

Ajout du dépôt pour MongoDB :

```sh
cat << EOF > /etc/yum.repos.d/mongodb-org-7.0.repo
[mongodb-org-7.0]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/\$releasever/mongodb-org/7.0/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-7.0.asc
EOF
```

Ajout du dépôt pour RabbitMQ :

```sh
cat << EOF > /etc/yum.repos.d/rabbitmq.repo
##
## Zero dependency Erlang RPM
##

[modern-erlang]
name=modern-erlang-el8
# Use a set of mirrors maintained by the RabbitMQ core team.
# The mirrors have significantly higher bandwidth quotas.
baseurl=https://yum1.novemberain.com/erlang/el/8/$basearch
        https://yum2.novemberain.com/erlang/el/8/$basearch
repo_gpgcheck=1
enabled=1
gpgkey=https://github.com/rabbitmq/signing-keys/releases/download/3.0/cloudsmith.rabbitmq-erlang.E495BB49CC4BBE5B.key
gpgcheck=1
sslverify=1
sslcacert=/etc/pki/tls/certs/ca-bundle.crt
metadata_expire=300
pkg_gpgcheck=1
autorefresh=1
type=rpm-md

##
## RabbitMQ Server
##

[rabbitmq-el8-noarch]
name=rabbitmq-el8-noarch
baseurl=https://yum2.novemberain.com/rabbitmq/el/8/noarch
        https://yum1.novemberain.com/rabbitmq/el/8/noarch
repo_gpgcheck=1
enabled=1
# Cloudsmith's repository key and RabbitMQ package signing key
gpgkey=https://github.com/rabbitmq/signing-keys/releases/download/3.0/cloudsmith.rabbitmq-server.9F4587F226208342.key
       https://github.com/rabbitmq/signing-keys/releases/download/3.0/rabbitmq-release-signing-key.asc
gpgcheck=1
sslverify=1
sslcacert=/etc/pki/tls/certs/ca-bundle.crt
metadata_expire=300
pkg_gpgcheck=1
autorefresh=1
type=rpm-md
EOF
```

Ajout du dépôt pour TimescaleDB :

```sh
cat << EOF > /etc/yum.repos.d/timescale_timescaledb.repo
[timescale_timescaledb]
name=timescale_timescaledb
baseurl=https://packagecloud.io/timescale/timescaledb/el/8/\$basearch
repo_gpgcheck=1
# TimescaleDB doesn’t sign all its packages
gpgcheck=0
enabled=1
gpgkey=https://packagecloud.io/timescale/timescaledb/gpgkey
sslverify=1
sslcacert=/etc/pki/tls/certs/ca-bundle.crt
metadata_expire=300
EOF
```

### Configuration des dépôts

Exécuter la commande suivante et vérifier dans la sortie que les dépôts ajoutés
sont bien contactés (accepter les cléfs des différents dépôts lorsque demandé) :

```sh
dnf makecache
```

Désactiver le module PostgreSQL ([requis pour l'installation de TimescaleDB sur RHEL 8](https://docs.timescale.com/install/latest/self-hosted/installation-redhat/)) :

```sh
dnf module disable postgresql
```

Activer le module Nginx 1.20.* :

```sh
dnf module disable php
dnf module enable nginx:1.20
```

Activer le module Redis 6.0.* :

```sh
dnf module enable redis:6
```

### Installation

```sh
dnf install logrotate socat mongodb-org nginx redis timescaledb-2-postgresql-13-2.14.2 timescaledb-2-loader-postgresql-13-2.14.2
dnf install erlang rabbitmq-server
```

Pour éviter un upgrade automatique des dépendances de Canopsis, vous pouvez épingler leurs paquets en ajoutant la directive suivante dans le fichier `/etc/yum.conf` :

```
exclude=mongodb-org,mongodb-org-server,mongodb-org-shell,mongodb-org-mongos,mongodb-org-tools,nginx,nginx-filesystem,erlang,rabbitmq-server,redis,timescaledb-2-postgresql-13,timescaledb-2-loader-postgresql-13
```

### Ouverture des ports

Pratiquer les ouvertures de ports nécessaires à l'accès au service.

Les commandes données couvrent le cas standard où le pare-feu système `firewalld` est utilisé, et servent surtout à rappeler les ports ou services à ouvrir. (cf. [matrice des flux réseau](../matrice-des-flux-reseau/index.md))

```sh
firewall-cmd --add-port=5672/tcp --add-port=15672/tcp --permanent
firewall-cmd --add-port=8080/tcp --permanent
firewall-cmd --add-port=27017/tcp --permanent
firewall-cmd --add-service=postgresql --permanent
firewall-cmd --add-service=redis --permanent
firewall-cmd --reload
```

### Configuration de MongoDB

Une clef d'authentification doit être créée pour sécuriser le futur *replicaset*
MongoDB. Pour cela, exécuter les commandes suivantes qui vont générer le
fichier `/etc/mongodb-keyfile` et lui attribuer les bonnes permissions et
appartenances :

```sh
openssl rand -base64 756 > /etc/mongodb-keyfile
chmod 400 /etc/mongodb-keyfile
chown mongod:root /etc/mongodb-keyfile
```

Ensuite, éditer le fichier `/etc/mongod.conf`. Il faut mettre les contenus
suivants dans les sections `security` et `replication` pour :

- protéger l'instance (authentification obligatoire)
- activer le mode *replicaset*

```yaml
security:
  authorization: enabled
  keyFile: /etc/mongodb-keyfile

replication:
  oplogSizeMB: 1024
  replSetName: rs0
```

!!! information
    Le nom de *replicaset* `rs0` est un exemple. Un autre nom peut être choisi,
    du moment que ce nom est répercuté plus tard dans l'URL MongoDB au sein du
    fichier `/opt/canopsis/etc/go-engines-vars.conf`.

On peut à présent activer et démarrer le service :

```sh
systemctl enable --now mongod.service
```

L'instance MongoDB étant démarrée, il reste à la configurer.

On se connecte dans un shell `mongo` et on désactive la télémétrie :

```sh
mongosh
> disableTelemetry()
```

On initialise le *replicaset* :

```sh
> rs.initiate()
```

L'état du *replicaset* peut être vérifié avec la commande `rs.status()` :

```sh
> rs.status()
```

Au bout de quelques secondes, le prompt du shell `mongo` doit faire apparaître
que le nœud est PRIMARY :

```sh
rs0:PRIMARY>
```

Lorsque c'est le cas, le *replicaset* est prêt. On poursuit avec la création
des utilisateurs MongoDB `root` puis `canopsis`, toujours dans le shell
`mongo` :

```sh
> use admin
> db.createUser({user: "root", pwd: "UNMOTDEPASSEFORT", roles: [ { role: "root", db: "admin" }]})
> exit
```

On se reconnecte avec le shell `mongo`, cette fois-ci en s'authentifiant en tant
que `root` MongoDB :

```sh
mongosh -u root -p UNMOTDEPASSEFORT
> use canopsis
> db.createUser({user: "cpsmongo", pwd: "canopsis", roles: [ { role: "dbOwner", db: "canopsis" }, { role: "clusterMonitor", db: "admin"}]})
> exit
```

Les manipulations d'installation dans MongoDB sont terminées.

### Configuration de TimescaleDB

Initialiser l'instance PostgreSQL puis initialiser TimescaleDB (cf. [documentation de l'outil de règlage](https://docs.timescale.com/timescaledb/latest/how-to-guides/configuration/timescaledb-tune/) de TimescaleDB) :

```sh
postgresql-13-setup initdb
timescaledb-tune -yes --pg-config=/usr/pgsql-13/bin/pg_config
echo "timescaledb.telemetry_level=off" >> /var/lib/pgsql/13/data/postgresql.conf
```

Activer et démarrer le service :

```sh
systemctl enable --now postgresql-13.service
```

Se connecter à l'instance PostgreSQL avec l'identité du superuser `postgres` :

```sh
sudo -u postgres psql
```

Créer la base de données `canopsis` et l'utilisateur associé dans l'instance PostgreSQL :

```sql
postgres=# CREATE database canopsis;
postgres=# \c canopsis
canopsis=# CREATE EXTENSION IF NOT EXISTS timescaledb;
canopsis=# SET password_encryption = 'scram-sha-256';
canopsis=# CREATE USER cpspostgres WITH PASSWORD 'canopsis';
canopsis=# exit
```

### Configuration de RabbitMQ

Activer et démarrer le service :

```sh
systemctl enable --now rabbitmq-server.service
```

Activer le plugin qui apporte l'interface de management RabbitMQ, puis redémarrer le service :

```sh
rabbitmq-plugins enable rabbitmq_management
systemctl restart rabbitmq-server.service
```

Créer les objets RabbitMQ (vhost, utilisateur avec les bons droits) utiles à l'application Canopsis :

```sh
rabbitmqctl add_vhost canopsis
rabbitmqctl add_user cpsrabbit canopsis
rabbitmqctl set_user_tags cpsrabbit administrator
rabbitmqctl set_permissions --vhost canopsis cpsrabbit '.*' '.*' '.*'
```

### Démarrage de Redis

Ajouter un mot de passe ( ici `canopsis`)

```sh
sed -i 's/^# requirepass.*/requirepass canopsis/' /etc/redis.conf
```

Activer et démarrer le service :

```sh
systemctl enable --now redis
```

## Installation de Canopsis Community ou Pro

Canopsis est disponible dans une édition « Community », open-source et gratuitement accessible à tous, et une édition « Pro », souscription commerciale ajoutant des fonctionnalités supplémentaires. Voyez [le site officiel de Canopsis](https://www.capensis.fr/canopsis/) pour en savoir plus.

Notez que l'édition Pro de Canopsis était auparavant connue sous le nom de « CAT » et que certains éléments peuvent encore la désigner sous ce nom.

Cliquez sur l'un des onglets « Community » ou « Pro » suivants, en fonction de l'édition choisie.

=== "Canopsis Community (édition open-source)"

    Ajout du dépôt de paquets Canopsis pour RHEL 8 :

    ```sh
    cat << EOF > /etc/yum.repos.d/canopsis.repo
    [canopsis]
    name = canopsis
    baseurl=https://nexus.canopsis.net/repository/canopsis/el8/community/
    gpgcheck=0
    enabled=1
    EOF
    ```

    Installation de l'édition open-source de Canopsis :

    ```sh
    dnf makecache
    dnf install canopsis
    ```

=== "Canopsis Pro (souscription commerciale)"

    !!! attention
        L'édition Pro nécessite une souscription commerciale ainsi qu'une authentification d'accès aux repos à renseigner dans `baseurl` du fichier `/etc/yum.repos.d/canopsis-pro.repo`.

    Ajout des dépôts de paquets Canopsis pour RHEL 8 :

    ```sh
    cat << EOF > /etc/yum.repos.d/canopsis.repo
    [canopsis]
    name = canopsis
    baseurl=https://nexus.canopsis.net/repository/canopsis/el8/community/
    gpgcheck=0
    enabled=1
    EOF

    cat << EOF > /etc/yum.repos.d/canopsis-pro.repo
    [canopsis-pro]
    name = canopsis-pro
    baseurl=https://user:password@nexus.canopsis.net/repository/canopsis-pro/el8/pro/
    gpgcheck=0
    enabled=1
    EOF
    ```

    Installation de Canopsis Pro :

    ```sh
    dnf makecache
    dnf install canopsis-pro
    ```

## Initialisation de Canopsis

Le fichier de configuration est `/opt/canopsis/etc/go-engines-vars.conf`, qui
est normalement dans l'état suivant :

```sh
CPS_MONGO_URL="mongodb://cpsmongo:canopsis@localhost:27017/canopsis?replicaSet=rs0"
CPS_AMQP_URL="amqp://cpsrabbit:canopsis@localhost:5672/canopsis"
CPS_POSTGRES_URL="postgresql://cpspostgres:canopsis@localhost:5432/canopsis"
CPS_REDIS_URL="redis://localhost:6379/0"
CPS_API_URL="http://localhost:8082"
#CPS_OLD_API_URL="http://localhost:8081"
CPS_POSTGRES_TECH_URL="postgresql://cpspostgres:canopsis@localhost:5432/canopsis_tech_metrics"
```

!!! Note
    On remarque dans `CPS_MONGO_URL` le nom de *replicaset*, ici `rs0`,
    qui doit être cohérent avec ce qui a été mis lors de la
    [configuration de MongoDB](#configuration-de-mongodb).

Il est maintenant temps d'initialiser l'instance Canopsis ;
cliquez sur l'un des onglets « Community » ou « Pro » suivants, en fonction de
l'édition choisie.

=== "Canopsis Community (édition open-source)"

    Provisionner Canopsis :

    ```sh
    set -o allexport; source /opt/canopsis/etc/go-engines-vars.conf; /opt/canopsis/bin/canopsis-reconfigure -migrate-postgres=true -edition community
    ```

    Activer et démarrer les services :

    ```sh
    systemctl enable --now canopsis-engine-go@engine-action.service \
                           canopsis-engine-go@engine-axe.service \
                           canopsis-engine-go@engine-che.service \
                           canopsis-engine-go@engine-fifo.service \
                           canopsis-engine-go@engine-pbehavior.service \
                           canopsis-service@canopsis-api.service \
                           canopsis.service
    ```

=== "Canopsis Pro (souscription commerciale)"

    Provisionner Canopsis :

    ```sh
    set -o allexport; source /opt/canopsis/etc/go-engines-vars.conf; /opt/canopsis/bin/canopsis-reconfigure -migrate-postgres=true -edition pro
    ```

    Activer et démarrer les services :

    ```sh
    systemctl enable --now canopsis-engine-go@engine-action.service \
                           canopsis-engine-go@engine-axe.service \
                           canopsis-engine-go@engine-che.service \
                           canopsis-engine-go@engine-correlation.service \
                           canopsis-engine-go@engine-dynamic-infos.service \
                           canopsis-engine-go@engine-fifo.service \
                           canopsis-engine-go@engine-pbehavior.service \
                           canopsis-engine-go@engine-remediation.service \
                           canopsis-engine-go@engine-webhook.service \
                           canopsis-service@canopsis-api.service \
                           canopsis-engine-python-snmp.service \
                           canopsis.service
    ```

Tester un envoi d'évènement :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
  "event_type": "check",
  "connector": "connector_test_creation_alarmes",
  "connector_name": "test",
  "component": "component_test_creation_alarmes",
  "resource": "resource_test_creation_alarmes",
  "source_type": "resource",
  "author": "QA_canopsis",
  "state": 2,
  "debug": true,
  "output": "Test création alarmes Canopsis"
}' 'http://localhost:8082/api/v4/event'
```

## Lancement de la Web UI de Canopsis

Installer le paquet :

```sh
dnf install canopsis-webui
```

Activer et démarrer Nginx :

```sh
systemctl enable --now nginx.service
```

Une fois cette commande terminée, vous pouvez alors réaliser votre [première connexion à l'interface Canopsis](premiere-connexion.md).

Si vous souhaitez réaliser une mise à jour, la procédure est décrite dans le [Guide de mise à jour](../mise-a-jour/index.md).

