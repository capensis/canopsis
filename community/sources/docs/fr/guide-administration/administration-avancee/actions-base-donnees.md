# Nettoyage, sauvegarde et restauration des bases de données

## MongoDB

### Nettoyage

!!! note
    Depuis la version `VERSION` de Canopsis, le nettoyage des bases de données `MongoDB` et `TimescaleDB` ne se réalise plus de manière manuel.  
    Il faut pour cela utiliser la fonctionnalité [`Paramètres de stockage`](../../guide-utilisation/menu-administration/parametres-de-stockage.md) qui centralise toutes les politiques de rétention des données de Canopsis.

### Sauvegarde

Utilisez la commande `mongodump` via une tâche cron. De préférence, faites la sauvegarde sur un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.mongodb.com/manual/tutorial/backup-and-restore-tools/#basic-mongodump-operation).

!!! note
    Le mot de passe par défaut est `canopsis`, mais il peut être nécessaire d'adapter la commande selon votre contexte.

```sh
mongodump --uri mongodb://cpsmongo:canopsis@mongodb/canopsis?replicaSet=rs0 --gzip --archive="/chemin/vers/sauvegarde/canopsis-$(date +"%Y-%m-%d-%H-%M").gz"
```

### Restauration

!!! attention
    Cette manipulation a une incidence métier importante et ne doit être réalisée que par une personne compétente. La restauration de la base de données ne doit être effectuée que si celle-ci est endommagée, pour corriger l'incident.

Avant de procéder à la restauration, arrêtez l'hyperviseur.  

=== "Canopsis Community (édition open-source)"

    ```sh
    systemctl stop --now canopsis-engine-go@engine-action.service \
                           canopsis-engine-go@engine-axe.service \
                           canopsis-engine-go@engine-che.service \
                           canopsis-engine-go@engine-fifo.service \
                           canopsis-engine-go@engine-pbehavior.service \
                           canopsis-service@canopsis-api.service \
                           canopsis.service
    ```

=== "Canopsis Pro (souscription commerciale)"

    ```sh
    systemctl stop --now canopsis-engine-go@engine-action.service \
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

Utilisez la commande `mongorestore`. De préférence, récupérez la sauvegarde depuis un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.mongodb.com/manual/tutorial/backup-and-restore-tools/#basic-mongorestore-operations).

```sh
mongorestore --uri mongodb://cpsmongo:canopsis@mongodb/canopsis?replicaSet=rs0 --db canopsis --gzip --archive="/chemin/vers/sauvegarde/canopsis-2026-01-05-09-06.gz"
```

Si la restauration est réussie vous pouvez redémarrer l'hyperviseur.  

=== "Canopsis Community (édition open-source)"

    ```sh
    systemctl start --now canopsis-engine-go@engine-action.service \
                            canopsis-engine-go@engine-axe.service \
                            canopsis-engine-go@engine-che.service \
                            canopsis-engine-go@engine-fifo.service \
                            canopsis-engine-go@engine-pbehavior.service \
                            canopsis-service@canopsis-api.service \
                            canopsis.service
    ```

=== "Canopsis Pro (souscription commerciale)"

    ```sh
    systemctl start --now canopsis-engine-go@engine-action.service \
                            canopsis-engine-go@engine-axe.service \
                            canopsis-engine-go@engine-che.service \
                            canopsis-engine-go@engine-correlation.service \
                            canopsis-engine-go@engine-dynamic-infos.service \
                            canopsis-engine-go@engine-fifo.service \
                            canopsis-engine-go@engine-pbehavior.service \
                            canopsis-engine-go@engine-remediation.service \
                            canopsis-engine-go@engine-webhook.service \
                            canopsis-engine-go@engine-events-recorder.service \
                            canopsis-service@canopsis-api.service \
                            canopsis-engine-python-snmp.service \
                            canopsis.service
    ```

## PostgreSQL (TimescaleDB)

### Sauvegarde

Utilisez la commande `pg_dump` via une tâche cron. De préférence, faites la sauvegarde sur un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.timescale.com/self-hosted/latest/backup-and-restore/logical-backup/).

!!! note
    Le mot de passe par défaut est `canopsis`, mais il peut être nécessaire d'adapter la commande selon votre contexte.

!!! attention
    Il est nécessaire de réaliser cette action une fois pour chaque base de données, `canopsis` et `canopsis_tech_metrics`.

Pour la base `canopsis`:
```sh
pg_dump postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis --no-owner -Fc -v -f /tmp/canopsis-$(date +"%Y-%m-%d-%H-%M")-canopsis-dump.sql.gz
```

Pour la base `canopsis_tech_metrics`:
```sh
pg_dump postgresql://cpspostgres_tech_metrics:canopsis@timescaledb:5432/canopsis_tech_metrics --no-owner -Fc -v -f /tmp/canopsis-$(date +"%Y-%m-%d-%H-%M")-canopsis_tech_metrics-dump.sql.gz
```

### Restauration

!!! attention
    Cette manipulation a une incidence métier importante et ne doit être réalisée que par une personne compétente. La restauration des bases de données ne doit être effectuée que si celles-ci sont endommagées, pour corriger l'incident.

Avant de procéder à la restauration, arrêtez l'hyperviseur.  

=== "Canopsis Community (édition open-source)"

    ```sh
    systemctl stop --now canopsis-engine-go@engine-action.service \
                            canopsis-engine-go@engine-axe.service \
                            canopsis-engine-go@engine-che.service \
                            canopsis-engine-go@engine-fifo.service \
                            canopsis-engine-go@engine-pbehavior.service \
                            canopsis-service@canopsis-api.service \
                            canopsis.service
    ```

=== "Canopsis Pro (souscription commerciale)"

    ```sh
    systemctl stop --now canopsis-engine-go@engine-action.service \
                            canopsis-engine-go@engine-axe.service \
                            canopsis-engine-go@engine-che.service \
                            canopsis-engine-go@engine-correlation.service \
                            canopsis-engine-go@engine-dynamic-infos.service \
                            canopsis-engine-go@engine-fifo.service \
                            canopsis-engine-go@engine-pbehavior.service \
                            canopsis-engine-go@engine-remediation.service \
                            canopsis-engine-go@engine-webhook.service \
                            canopsis-engine-go@engine-events-recorder.service \
                            canopsis-service@canopsis-api.service \
                            canopsis-engine-python-snmp.service \
                            canopsis.service
    ```

Utilisez la commande `pg_restore`. De préférence, récupérez la sauvegarde depuis un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.timescale.com/self-hosted/latest/backup-and-restore/logical-backup/).

Tout d'abord, il vous faut vous connecter à la base postgresql
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

Vous devrez ensuite entrer votre base en mode restauration
```sql
SELECT timescaledb_pre_restore();
```

Une fois la base en mode restauration, vous pouvez importer vos dumps

Pour la base `canopsis`
```sh
pg_restore -Fc -d "postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis" canopsis-YYYY-mm-dd-HH-MM-canopsis-dump.sql.gz
```

Pour la base `canopsis_tech_metrics`:
```sh
pg_restore -Fc -d "postgresql://cpspostgres_tech_metrics:canopsis@timescaledb:5432/canopsis_tech_metrics" canopsis-YYYY-mm-dd-HH-MM-canopsis_tech_metrics-dump.sql.gz
```

Une fois les dumps importés, vous pouvez sortir du mode restauration
```sql
SELECT timescaledb_post_restore();
```

Si la restauration est réussie vous pouvez redémarrer l'hyperviseur.  

=== "Canopsis Community (édition open-source)"

    ```sh
    systemctl start --now canopsis-engine-go@engine-action.service \
                            canopsis-engine-go@engine-axe.service \
                            canopsis-engine-go@engine-che.service \
                            canopsis-engine-go@engine-fifo.service \
                            canopsis-engine-go@engine-pbehavior.service \
                            canopsis-service@canopsis-api.service \
                            canopsis.service
    ```

=== "Canopsis Pro (souscription commerciale)"

    ```sh
    systemctl start --now canopsis-engine-go@engine-action.service \
                            canopsis-engine-go@engine-axe.service \
                            canopsis-engine-go@engine-che.service \
                            canopsis-engine-go@engine-correlation.service \
                            canopsis-engine-go@engine-dynamic-infos.service \
                            canopsis-engine-go@engine-fifo.service \
                            canopsis-engine-go@engine-pbehavior.service \
                            canopsis-engine-go@engine-remediation.service \
                            canopsis-engine-go@engine-webhook.service \
                            canopsis-engine-go@engine-events-recorder.service \
                            canopsis-service@canopsis-api.service \
                            canopsis-engine-python-snmp.service \
                            canopsis.service
    ```