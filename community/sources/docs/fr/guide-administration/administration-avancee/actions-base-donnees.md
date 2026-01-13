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
    systemctl stop --now canopsis.service
    ```

=== "Canopsis Pro (souscription commerciale)"

    ```sh
    systemctl stop --now canopsis.service
    ```

Utilisez la commande `mongorestore`. De préférence, récupérez la sauvegarde depuis un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.mongodb.com/manual/tutorial/backup-and-restore-tools/#basic-mongorestore-operations).

```sh
mongorestore --uri mongodb://cpsmongo:canopsis@mongodb/canopsis?replicaSet=rs0 --db canopsis --gzip --archive="/chemin/vers/sauvegarde/canopsis-2026-01-05-09-06.gz"
```

Si la restauration est réussie vous pouvez redémarrer l'hyperviseur.  

=== "Canopsis Community (édition open-source)"

    ```sh
    systemctl start --now canopsis.service
    ```

=== "Canopsis Pro (souscription commerciale)"

    ```sh
    systemctl start --now canopsis.service
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
    systemctl stop --now canopsis.service
    ```

=== "Canopsis Pro (souscription commerciale)"

    ```sh
    systemctl stop --now canopsis.service
    ```

Une fois Canopsis éteint, il est nécessaire de supprimer les bases `canopsis` et/ou `canopsis_tech_metrics` avant de lancer la restauration

Pour la base `canopsis`:
```sh
echo "select 'drop table '||tablename||' cascade;' from pg_tables where schemaname = 'public'"  | psql postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis -t | psql postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis
```

Pour la base `canopsis_tech_metrics`:
```sh
echo "select 'drop table '||tablename||' cascade;' from pg_tables where schemaname = 'public'"  | psql postgresql://cpspostgres_tech_metrics:canopsis@timescaledb:5432/canopsis_tech_metrics -t | psql postgresql://cpspostgres_tech_metrics:canopsis@timescaledb:5432/canopsis_tech_metrics
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
    systemctl start --now canopsis.service
    ```

=== "Canopsis Pro (souscription commerciale)"

    ```sh
    systemctl start --now canopsis.service
    ```

## Cas d'usage 

#### Trouver les alarmes qui correspondent à un connecteur particulier

Pour récupérer l'entièreté des alarmes

```js
db.periodical_alarm.find({'v.connector':'connector_test_creation_alarmes'})
```

Pour récupérer uniquement le `display_name` des alarmes (Par exemple: Pour chercher les alarmes visuellements dans l'interface de Canopsis)

```js
db.periodical_alarm.find({'v.connector':'connector_test_creation_alarmes'})

[
  {
    _id: '019bb12c-fc43-7ef9-b3a0-afb5841d3450',
    v: { display_name: 'W0H-H09-X7U' }
  }
]
```

#### Trouver les entités qui ont beaucoup "d'Impacts" ou de "Dépendances"

Trouver les 10 entités ayant le plus `d'Impacts`
```js
db.default_entities.aggregate([{$lookup:{from:"default_entities",let:{cid:"$_id"},pipeline:[{$match:{$expr:{$and:[{$eq:["$type","resource"]},{$eq:["$connector","$$cid"]}]}}}],as:"resources_of_connector"}},{$addFields:{impactCount:{$switch:{branches:[{case:{$eq:["$type","connector"]},then:{$size:"$resources_of_connector"}},{case:{$eq:["$type","component"]},then:{$cond:[{$ifNull:["$connector",false]},1,0]}},{case:{$eq:["$type","resource"]},then:{$cond:[{$ifNull:["$component",false]},1,0]}}],default:0}}}},{$sort:{impactCount:-1,_id:1}},{$limit:10},{$project:{_id:1,name:1,type:1,impactCount:1}}]);

```

Trouver les 10 entités ayant le plus de `Dépendances`

```js
db.default_entities.aggregate([{$lookup:{from:"default_entities",let:{cid:"$_id"},pipeline:[{$match:{$expr:{$and:[{$eq:["$type","component"]},{$eq:["$connector","$$cid"]}]}}}],as:"components_of_connector"}},{$lookup:{from:"default_entities",let:{compName:"$name",compConn:"$connector"},pipeline:[{$match:{$expr:{$and:[{$eq:["$type","resource"]},{$eq:["$component","$$compName"]},{$eq:["$connector","$$compConn"]}]}}}],as:"resources_of_component"}},{$addFields:{dependencyCount:{$switch:{branches:[{case:{$eq:["$type","connector"]},then:{$size:"$components_of_connector"}},{case:{$eq:["$type","component"]},then:{$size:"$resources_of_component"}},{case:{$eq:["$type","resource"]},then:{$cond:[{$ifNull:["$connector",false]},1,0]}}],default:0}}}},{$sort:{dependencyCount:-1,_id:1}},{$limit:10},{$project:{_id:1,name:1,type:1,dependencyCount:1}}]);
```

#### Trouver les comportements périodiques qui agissent sur un nombre important d'alarmes/entités

```js
db.periodical_alarm.aggregate([{ $match:{ "v.pbehavior_info.id":{ $exists:true } } },{ $group:{ _id:"$v.pbehavior_info.id", name:{ $first:"$v.pbehavior_info.name" }, alarm_count:{ $sum:1 } } },{ $project:{ _id:0, id:"$_id", name:1, alarm_count:1 } }])

[
  {
    name: 'test',
    alarm_count: 31,
    id: '019bb14f-7eda-7dce-a3fa-d68ad3b26d11'
  }
]
```


#### Trouver les comportements périodiques qui se sont exécutés en retard
```sh

```

#### Rechercher les requêtes et les collections qui sont visés par des requêtes dont le temps d'exécution est supérieur à 10 secondes
```sh
grep "durationMillis" mongodb.log | jq -c 'select(.c == "COMMAND" and .attr.durationMillis > 10000) | {ns: .attr.ns, command: .attr.command, d: "\(.attr.durationMillis) ms"}' 
```