# Procédure de migration dans un environnement RPM

## Sommaire

[TOC]

## Mise à jour de dépôts

Certaines briques logicielles nécessitent un changement de dépôts

### MongoDB:

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

## Mise à jour des paquets de l'applicatif Canopsis

Pour réaliser la mise à jour de Canopsis il faut dans un premier temps arrêter l'application

```sh
systemctl stop canopsis.service
```

Une fois le service correctement arrêté, il faut lever le `versionlock` éventuellement présent pour mettre à jour vers la version `25.10`

```sh
dnf versionlock delete 'canopsis-pro-25.04.*'
dnf versionlock delete 'canopsis-webui-25.04.*'

dnf versionlock add --raw 'canopsis-pro-25.10.*'
dnf versionlock add --raw 'canopsis-webui-25.10.*'
```

La mise à jour des packages peut ensuite commencer

```sh
dnf upgrade canopsis-pro canopsis-webui -y
```

Le paquet `canopsis-pro` a une dépendance vers le client `mongodb`, de ce fait ce paquet sera installé durant cette installation.


## Mise à jour TimescaleDB

Dans cette version de Canopsis, la base de données TimescaleDB passe de la version 2.15.1 à 2.21.4.  
En plus de la mise à jour de TimescaleDB lui-même, le système de gestion de base de données PostreSQL doit être mis à jour de la version 15 à la version 17.

Deux étapes sont à suivre :

1. Mise à jour de TimescaleDB 2.15.1 vers 2.21.4
2. Mise à jour de PostgreSQL 15 vers 17

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

## Mise à jour MongoDB

### Vérifications MongoDB

!!! warning "Vérification"

    Avant de démarrer la procédure de mise à jour, vous devez vérifier que la valeur de `featureCompatibilityVersion` est bien positionnée à **7.0**  

```sh
mongosh -u root -p root
> db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )
> exit
```

### Mise à jour

Pour commencer, il faut couper le service

```sh
systemctl stop mongod.service
```

Une fois le [dépôt mise à jour vers la version `8.0`](#mise-a-jour-de-depots), on peut lancer l'upgrade

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

## Mise à jour RabbitMQ

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

## Lancement du provisioning `canopsis-reconfigure`

Reconfiguration de Canopsis

!!! Attention

    Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette personnalisation.
    L'outil est en effet pré paramétré pour fonctionner naturellement.


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