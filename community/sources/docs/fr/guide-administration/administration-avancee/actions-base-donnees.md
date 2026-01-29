# Nettoyage, sauvegarde, restauration et interractions avec les bases de données

## MongoDB

### Nettoyage

!!! note
    Le nettoyage des bases de données `MongoDB` et `TimescaleDB` ne se réalise pas de manière manuel.  
    Il faut pour cela utiliser la fonctionnalité [`Paramètres de stockage`](../../guide-utilisation/menu-administration/parametres-de-stockage.md) qui centralise toutes les politiques de rétention des données de Canopsis.

### Sauvegarde

Utilisez la commande `mongodump` via une tâche cron. De préférence, faites la sauvegarde sur un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.mongodb.com/manual/tutorial/backup-and-restore-tools/#basic-mongodump-operation).

!!! note
    Le mot de passe par défaut est `canopsis`, mais il peut être nécessaire d'adapter la commande selon votre contexte.

```sh
mongodump --uri mongodb://cpsmongo:canopsis@mongodb/canopsis?replicaSet=rs0 --gzip --archive="/chemin/vers/sauvegarde/canopsis-$(date +"%Y-%m-%d-%H-%M").gz"
```

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

Utilisez la commande `mongorestore`. De préférence, récupérez la sauvegarde depuis un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.mongodb.com/manual/tutorial/backup-and-restore-tools/#basic-mongorestore-operations).

```sh
mongorestore --uri mongodb://cpsmongo:canopsis@mongodb/canopsis?replicaSet=rs0 --db canopsis --gzip --archive="/chemin/vers/sauvegarde/canopsis-2026-01-05-09-06.gz"
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

## PostgreSQL (TimescaleDB)

### Sauvegarde

Utilisez la commande `pg_dump` via une tâche cron. De préférence, faites la sauvegarde sur un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.timescale.com/self-hosted/latest/backup-and-restore/logical-backup/).

!!! note
    Le mot de passe par défaut est `canopsis`, mais il peut être nécessaire d'adapter la commande selon votre contexte.

!!! Warning
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

## Cas d'usage 

#### Commandes de base pour récupérer des informations sur son Canopsis

Ces commandes seront utiles pour rapidement pouvoir récupérer des informations.

```js
db.NOM_DE_LA_COLLECTION.find() # Permet de trouver une information
db.NOM_DE_LA_COLLECTION.countDocuments() # Permet de compter le nombre de documents dans une collection
```

Pour plus d'informations sur les commandes de l'utilitaire `mongosh`, référer à la documentation officiel : [Documentation Mongosh](https://www.mongodb.com/docs/mongodb-shell/run-commands/)

Généralement, ces commandes sont utilisés avec par exemple sur les collections : `periodical_alarm`, `default_entities`, `alarm_tag`, `pbehavior`.

#### Trouver les alarmes qui correspondent à un critère particulier

Récupérer l'entièreté des alarmes actuellement ouvertes qui proviennent du même `connecteur`
```js
db.periodical_alarm.find({'v.connector':'Nom du connecteur'})
```

Récupérer l'entièreté des alarmes actuellement ouvertes qui proviennent du même `composant`
```js
db.periodical_alarm.find({'v.component':'Nom du composant'})
```

Récupérer l'entièreté des alarmes actuellement ouvertes qui possèdent une même `ressource`
```js
db.periodical_alarm.find({'v.resource':'Nom de la ressource'})
```

#### Trouver ou compter les entités liés à un tag

Trouver les alarmes actuellement ouvertes qui correspondent à un/des tags
```js
db.periodical_alarm.find({tags: { $in: ["une", "liste", "de", "tags"] }}) # Cherche à un match un ou plusieurs tags dans la liste
```

Compter les alarmes actuellement ouvertes qui correspondent à un/des tags
```js
db.periodical_alarm.countDocuments({tags: { $in: ["une", "liste", "de", "tags"] }}) # Cherche à un match un ou plusieurs tags dans la liste
```

**Trouver les alarmes qui correspondent exactement à un/des tags :**

Trouver les alarmes actuellement ouvertes qui correspondent exactement à une liste de tag:
```js
db.periodical_alarm.find({"tags": "tag" }}) # Match exactement un tag

db.periodical_alarm.find({tags: { $all: ["une", "liste", "de", "tags"] }}) # Match exactement une liste de tags
```

Compter les alarmes actuellement ouvertes qui correspondent exactement à une liste tags
```js
db.periodical_alarm.countDocuments({"tags": "tag" }}) # Match exactement un tag

db.periodical_alarm.countDocuments({tags: { $all: ["une", "liste", "de", "tags"] }}) # Match exactement une liste de tags
```

#### Trouver les connecteurs, composants, ressources ou entités qui ont le plus de dépendances

Récupérer les 10 connecteurs avec le plus de dépendances :
```js
db.default_entities.aggregate([
    {$match: {connector: {$ne: null}, type: "resource"}},
    {$group: {_id: "$connector", depends: {$sum: 1}}},
    {$sort: {depends: -1, _id: 1}},
    {$limit: 10},
]);
```

Récupérer les 10 composants avec le plus de dépendances :
```js
db.default_entities.aggregate([
    {$match: {component: {$ne: null}, type: "resource"}},
    {$group: {_id: "$component", depends: {$sum: 1}}},
    {$sort: {depends: -1, _id: 1}},
    {$limit: 10},
]);
```

Récupérer les 10 entités avec le plus de dépendances :
```js
db.default_entities.aggregate([
    {$match: {services: {$ne: null}}},
    {$project: {services: {$size: "$services"}, type: 1}},
    {$sort: {services: -1, _id: 1}},
    {$limit: 10},
]);
```

Récupérer les 10 services avec le plus de dépendances :
```js
db.default_entities.aggregate([
    {$unwind: "$services"},
    {$group: {_id: "$services", depends: {$sum: 1}}},
    {$sort: {depends: -1, _id: 1}},
    {$limit: 10},
]);
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

#### Rechercher les requêtes et les collections qui sont visés par des requêtes dont le temps d'exécution est supérieur à 1 secondes et appliquer des indexes pour les corriger.

!!! note
    Pour optimiser les temps de réponses de la base de données MongoDB, il est nécessaire d'analyser les logs du serveur mongodb pour y extraire des potentiels `COLLSCAN (Slow query)` qui pourrait ralentir votre Canopsis. Pour ce faire, il est recommandé de passer par le support Canopsis pour analyser ces logs et éventuellement ajouter des indexes permettant d'améliorer les temps de réponses.

Trouver les requêtes dont le temps d'éxécution est supérieur à 1 seconde:

```sh
grep "durationMillis" mongodb.log | jq -c 'select(.c == "COMMAND" and .attr.durationMillis > 1000) | {ns: .attr.ns, command: .attr.command, d: "\(.attr.durationMillis) ms"}' 
```

Exemple d'une requête :

```sh
{"t":{"$date":"2024-09-23T16:15:33.418+02:00"},"s":"I",  "c":"COMMAND",  "id":51803,   "ctx":"conn431","msg":"Slow query","attr":{"type":"command","ns":"canopsis.default_entities","command":{"aggregate":"default_entities","pipeline":[{"$match":{"$or":[{"$and":[{"infos.customer.value":"CLIENT1"},{"component":{"$in":["COMPONENT1","COMPONENT2","COMPONENT3","COMPONENT4","COMPONENT5","COMPONENT6"]}}]}]}},{"$project":{"_id":1}}],"cursor":{},"lsid":{"id":{"$uuid":"ff7d8060-1f11-46bd-a186-75b2718b2277"}},"$clusterTime":{"clusterTime":{"$timestamp":{"t":1727100927,"i":3}},"signature":{"hash":{"$binary":{"base64":"81gp1Ydl+WkrIPK3eAedGx23erA=","subType":"0"}},"keyId":7414142034354634757}},"maxTimeMS":14999,"$db":"canopsis","$readPreference":{"mode":"primary"}},"planSummary":"COLLSCAN","numYields":381,"queryHash":"B0E2AF4B","planCacheKey":"589D5E1A","ok":0,"errMsg":"PlanExecutor error during aggregation :: caused by :: operation was interrupted because a client disconnected","errName":"ClientDisconnect","errCode":279,"reslen":326,"locks":{"FeatureCompatibilityVersion":{"acquireCount":{"r":382}},"Global":{"acquireCount":{"r":382}},"Mutex":{"acquireCount":{"r":1}}},"readConcern":{"level":"local","provenance":"implicitDefault"},"writeConcern":{"w":"majority","wtimeout":0,"provenance":"implicitDefault"},"remote":"IP:57832","protocol":"op_msg","durationMillis":4672}}
```

La requête a pris `4672 ms` soit environs `4.6 s` à s'exécuter ce qui occasionne des lenteurs dans l'utilisation de Canopsis.  
Pour remédier à cette situation, il est pertinent d’ajouter des indexes sur les collections. Comme par exemple ici, sur les valeurs `infos.customer.value` et `component`

Créer un index :
```js
db.default_entities.createIndex({"infos.customer.value" : 1, "component" : 1})
```

Pour visualiser les indexes qui ont été créer sur une collection, éxécuter la commande:

```js
db.NOM_DE_LA_COLLECTION.getIndexes()
```




