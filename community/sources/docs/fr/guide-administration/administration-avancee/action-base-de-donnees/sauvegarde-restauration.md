# Sauvegarde et restauration des bases de données

### Sauvegarde

Utilisez la commande `mongodump` via une tâche cron. De préférence, faites la sauvegarde sur un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.mongodb.com/manual/tutorial/backup-and-restore-tools/#basic-mongodump-operation).

!!! note
    Le mot de passe par défaut est `canopsis`, mais il peut être nécessaire d'adapter la commande selon votre contexte.

```sh
mongodump --uri ${CPS_MONGO_URL} --gzip --archive="/chemin/vers/sauvegarde/canopsis-$(date +"%Y-%m-%d-%H-%M").gz"
```

La commande `mongodump` peut-être réaliser depuis un noeud ou un poste utilisateur à condition que les [`database-tools`](https://www.mongodb.com/docs/database-tools/installation/) soient installés.

### Restauration

!!! Warning
    Cette manipulation a une incidence métier importante et ne doit être réalisée que par une personne compétente. La restauration de la base de données ne doit être effectuée que si celle-ci est endommagée, pour corriger l'incident.

Avant de procéder à la restauration, arrêtez l'hyperviseur.  

=== "RPM (EL8/EL9)"

    ```sh
    systemctl stop canopsis
    ```
       
=== "Docker Compose"

    ```sh
    docker compose down
    ```

=== "Helm"

    ```sh
    kubectl delete deployments --all
    ```

Avant de restaurer votre base de données, assurez-vous de la vider :
```sh
mongosh ${CPS_MONGO_URL} --eval 'db.getCollectionNames().forEach(c => db[c].drop())'
```

Utilisez la commande `mongorestore`. De préférence, récupérez la sauvegarde depuis un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.mongodb.com/manual/tutorial/backup-and-restore-tools/#basic-mongorestore-operations).

```sh
mongorestore --uri ${CPS_MONGO_URL} --db canopsis --gzip --archive="/chemin/vers/sauvegarde/canopsis-2026-01-05-09-06.gz"
```

La commande `mongorestore` peut-être réaliser depuis un noeud ou un poste utilisateur à condition que les [`database-tools`](https://www.mongodb.com/docs/database-tools/installation/) soient installés.

Si la restauration est réussie vous pouvez redémarrer l'hyperviseur.  

=== "RPM (EL8/EL9)"

    ```sh
    systemctl start canopsis
    ```
       
=== "Docker Compose"

    ```sh
    docker compose up -d
    ```

=== "Helm"

    ```sh
    export RELEASE_NAME="canopsis-prod"
    helm upgrade ${RELEASE_NAME} canopsis/canopsis-pro -f customer-values.yaml
    ```

## PostgreSQL (TimescaleDB)

### Sauvegarde

Utilisez la commande `pg_dump` via une tâche cron. De préférence, faites la sauvegarde sur un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.timescale.com/self-hosted/latest/backup-and-restore/logical-backup/).

!!! note
    Le mot de passe par défaut est `canopsis`, mais il peut être nécessaire d'adapter la commande selon votre contexte.

!!! Warning
    Il est nécessaire de réaliser cette action une fois pour chaque base de données, `canopsis` et `canopsis_tech_metrics`.

Pour la base `canopsis`:
```sh
pg_dump ${CPS_POSTGRES_URL} --no-owner -Fc -v -f /tmp/canopsis-$(date +"%Y-%m-%d-%H-%M")-canopsis-dump.sql.gz
```

Pour la base `canopsis_tech_metrics`:
```sh
pg_dump ${CPS_POSTGRES_TECH_URL} --no-owner -Fc -v -f /tmp/canopsis-$(date +"%Y-%m-%d-%H-%M")-canopsis_tech_metrics-dump.sql.gz
```

La commande `pg_dump` peut-être réaliser depuis un noeud ou un poste utilisateur à condition que les [`outils supplémentaires de PostgreSQL`](https://www.postgresql.org/download/) soient installés.

### Restauration

!!! Warning
    Cette manipulation a une incidence métier importante et ne doit être réalisée que par une personne compétente. La restauration des bases de données ne doit être effectuée que si celles-ci sont endommagées, pour corriger l'incident.

Avant de procéder à la restauration, arrêtez l'hyperviseur.  

=== "RPM (EL8/EL9)"

    ```sh
    systemctl stop canopsis
    ```
       
=== "Docker Compose"

    ```sh
    docker compose down
    ```

=== "Helm"

    ```sh
    kubectl delete deployments --all
    ```

Une fois Canopsis éteint, il est nécessaire de supprimer les tables des bases `canopsis` et/ou `canopsis_tech_metrics` avant de lancer la restauration

Pour la base `canopsis`:
```sh
echo "select 'drop table '||tablename||' cascade;' from pg_tables where schemaname = 'public'"  | psql ${CPS_POSTGRES_URL} -t | psql ${CPS_POSTGRES_URL}
```

Pour la base `canopsis_tech_metrics`:
```sh
echo "select 'drop table '||tablename||' cascade;' from pg_tables where schemaname = 'public'"  | psql ${CPS_POSTGRES_TECH_URL} -t | psql ${CPS_POSTGRES_TECH_URL}
```

La commande `psql` peut-être réaliser depuis un noeud ou un poste utilisateur à condition que les [`outils supplémentaires de PostgreSQL`](https://www.postgresql.org/download/) soient installés.

Utilisez la commande `pg_restore`. De préférence, récupérez la sauvegarde depuis un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.timescale.com/self-hosted/latest/backup-and-restore/logical-backup/).

Tout d'abord, il faut se connecter à la base postgresql
```sh
sudo -u postgres psql
```

Puis créer les bases de Canopsis
```sql 
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
postgres=# CREATE database canopsis_tech_metrics;
postgres=# \c canopsis_tech_metrics
canopsis_tech_metrics=# CREATE EXTENSION IF NOT EXISTS timescaledb;
canopsis_tech_metrics=# SET password_encryption = 'scram-sha-256';
canopsis_tech_metrics=# CREATE USER cpspostgres_tech_metrics WITH PASSWORD 'canopsis';
canopsis_tech_metrics=# GRANT ALL ON DATABASE canopsis_tech_metrics TO cpspostgres_tech_metrics;
canopsis_tech_metrics=# ALTER DATABASE canopsis_tech_metrics OWNER TO cpspostgres_tech_metrics;
canopsis_tech_metrics=# exit        
```

Vous devrez ensuite passer votre base en mode restauration
```sql
SELECT timescaledb_pre_restore();
```

Une fois la base en mode restauration, vous pouvez importer vos dumps

Pour la base `canopsis`
```sh
pg_restore -Fc -d ${CPS_POSTGRES_URL} canopsis-YYYY-mm-dd-HH-MM-canopsis-dump.sql.gz
```

Pour la base `canopsis_tech_metrics`:
```sh
pg_restore -Fc -d ${CPS_POSTGRES_TECH_URL} canopsis-YYYY-mm-dd-HH-MM-canopsis_tech_metrics-dump.sql.gz
```

La commande `pg_restore` peut-être réaliser depuis un noeud ou un poste utilisateur à condition que les [`outils supplémentaires de PostgreSQL`](https://www.postgresql.org/download/) soient installés.

Une fois les dumps importés, vous pouvez sortir du mode restauration
```sql
SELECT timescaledb_post_restore();
```

Si la restauration est réussie vous pouvez redémarrer l'hyperviseur.  

=== "RPM (EL8/EL9)"

    ```sh
    systemctl start canopsis
    ```
       
=== "Docker Compose"

    ```sh
    docker compose up -d
    ```

=== "Helm"

    ```sh
    export RELEASE_NAME="canopsis-prod"
    helm upgrade ${RELEASE_NAME} canopsis/canopsis-pro -f customer-values.yaml
    ```
