# Configuration avancée de la base de données MongoDB intégrée à Canopsis

La base de données [MongoDB](https://www.mongodb.com) contient la plupart des données, les vues et la configuration de l'interface Canopsis.

## Recommandations d'utilisation avancée

### Modification à effectuer sur le système

En suivant les recommendations fournies par MongoDB, des modifications dans le limiteur de ressources systèmes s'imposent.

Pour ce faire, créez le fichier `/etc/security/limits.d/mongo.conf` avec la commande

```shell
touch /etc/security/limits.d/mongo.conf
```

Puis venir y mettre la configuration suivante :

```shell
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

Il est ensuite préconisé de désactiver les `Transparent Huge Pages (THP)`

```shell
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

### Configuration de LogRotate

Afin d'éviter un surplus de log, il est nécessaire de mettre en place une règle de rétention. 

Nous recommandons en général de garder les logs 30 jours et de faire une sauvegarde régulière, cependant en fonction des ressources et de la politique interne concernant la rétention de log, ces valeurs sont adaptables.

```shell
cat > /etc/logrotate.d/mongodb.conf << EOF
/var/log/mongodb/*.log {
       daily
       rotate 30
       copytruncate
       delaycompress
       compress
       notifempty
       missingok
}
EOF
```

### Backup logique

Pour réaliser un backup de mongodb, nous recommendons d'utiliser la commande:

```shell
mongodump --uri "$CPS_MONGO_URL" --gzip --archive > /tmp/canopsis-prod-mongodb-$(date +"%Y-%m-%d")-dump.gz
```

### Optimisation    

Pour permettre un bon fonctionnement et une bonne vitesse de traitement de la plateforme, plusieurs optimisations peuvent être mises en oeuvre.

#### Taille du cache WiredTiger

Canopsis recommende l'utilisation de l'option `WiredTigerCacheSizeGB` pour fixer la consommation maximale de RAM que pourra utiliser le moteur WiredTiger avec MongoDB, cela permet d'éviter de nombreuses alertes de surconsommation de RAM. 

Pour le configurer il faut modifier le fichier `/etc/mongod.conf` et modifier la partie `storage:`

```
storage:
    ...
    wiredTiger:
       engineConfig:
           configString: cache_size=12G
```

Il est nécessaire de toujours mettre une valeur inférieure au nombre de giga disponible sur la machine en laissant toujours une petite marge.

#### Index et Collscan

Au fur et à mesure de votre utilisation de Canopsis les requêtes seront de plus en plus gourmandes en ressources matérielles et pourraient ralentir votre Canopsis. De nombreux facteurs peuvent entraîner des ralentissements: 

- Les collections qui grossissent au fur et à mesure de l'utilisation de Canopsis.
- L'enrichissement d'entitées et d'alarmes.

Exemple de log collscan:

```json
{"t":{"$date":"2024-09-23T16:15:33.418+02:00"},"s":"I",  "c":"COMMAND",  "id":51803,   "ctx":"conn431","msg":"Slow query","attr":{"type":"command","ns":"canopsis.default_entities","command":{"aggregate":"default_entities","pipeline":[{"$match":{"$or":[{"$and":[{"infos.customer.value":"CLIENT1"},{"component":{"$in":["COMPONENT1","COMPONENT2","COMPONENT3","COMPONENT4","COMPONENT5","COMPONENT6"]}}]}]}},{"$project":{"_id":1}}],"cursor":{},"lsid":{"id":{"$uuid":"ff7d8060-1f11-46bd-a186-75b2718b2277"}},"$clusterTime":{"clusterTime":{"$timestamp":{"t":1727100927,"i":3}},"signature":{"hash":{"$binary":{"base64":"81gp1Ydl+WkrIPK3eAedGx23erA=","subType":"0"}},"keyId":7414142034354634757}},"maxTimeMS":14999,"$db":"canopsis","$readPreference":{"mode":"primary"}},"planSummary":"COLLSCAN","numYields":381,"queryHash":"B0E2AF4B","planCacheKey":"589D5E1A","ok":0,"errMsg":"PlanExecutor error during aggregation :: caused by :: operation was interrupted because a client disconnected","errName":"ClientDisconnect","errCode":279,"reslen":326,"locks":{"FeatureCompatibilityVersion":{"acquireCount":{"r":382}},"Global":{"acquireCount":{"r":382}},"Mutex":{"acquireCount":{"r":1}}},"readConcern":{"level":"local","provenance":"implicitDefault"},"writeConcern":{"w":"majority","wtimeout":0,"provenance":"implicitDefault"},"remote":"IP:57832","protocol":"op_msg","durationMillis":4672}}
```

La requête a pris `4672 ms` soit `4.6 s` à s'exécuter ce qui occasionne des lenteurs dans l'utilisation de Canopsis.

Pour remédier à cette situation, il est pertinent d’ajouter des indexes sur les collections. Comme par exemple ici, sur les valeurs `infos.customer.value` et `component`

Pour créer un index, vous pouvez au travers du client `mongosh` exécuter la commande suivante:

```js
db.default_entities.createIndex({"infos.customer.value" : 1, "component" : 1})
```

Il est possible d'ajouter plusieurs indexes sur la même collection.  
Pour visualiser les indexes qui ont été créés sur une collection, éxécutez la commande:

```
db.default_entities.getIndexes()
```

## ReadPreference

Nous avons désactivé la possibilité de surcharger le [readPreference](https://www.mongodb.com/docs/manual/core/read-preference/) via l’URL de connexion Mongo.  
Le readPreference est maintenant défini par l’application elle-même et non modifiable depuis l’extérieur.

Par défaut, toutes les lectures sont effectuées depuis le noeud primaire.

Des exceptions existent dans certains composants, où la lecture depuis les secondaires est explicitement activée.

Cette décision a été prise pour éviter les erreurs fréquentes liées à un usage incorrect ou abusif du readPreference.  
Cela causait des anomalies difficiles à diagnostiquer, comme des erreurs du type : "mongo: no documents in result."

**Composants utilisant les lectures sur noeuds secondaires**

Les lectures depuis les réplicas secondaires sont actuellement activées dans :

* L'API des alarmes
* L'API d’export des alarmes
* Les moteurs engine-che et engine-dynamic-infos (dictionnaires d’info)
* Le Prometheus exporter

Cela permet de répartir la charge de manière plus équilibrée.
