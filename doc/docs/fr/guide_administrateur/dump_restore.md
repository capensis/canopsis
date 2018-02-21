# Dump & Restore

Commandes à exécuter afin d’effectuer des *dumps* de base de données, collections, séries…

## Utilisation avec Docker

Pour pouvoir échanger des données avec Docker, il faut soit passer par les flux d’entrée/sortie, soit passer par les volumes.

Chemin par défaut des volumes Docker : `/var/lib/docker/volumes/<path>`.

Quand un volume a été déclaré via `docker-compose` il portera le nom du projet et son nom de volume tel que déclaré dans le fichier de configuration de `docker-compose`.

Exemple de chemin pour MongoDB et un projet nommé `uoi038` : `/var/lib/docker/volumes/uoi038_mongodbdata/_data/`

## MongoDB

### Database

Dump dans archive
```
mongodump --host <host> -u cpsmongo -p canopsis -d canopsis --gzip --archive=/data/db/dump.gz
```

Restoration d'une archive
```
mongorestore -u cpsmongo -p canopsis --host=<host> --port=<port> --authenticationDatabase=canopsis --gzip --archive=/data/db/dump.gz --drop --db canopsis
```

### Collection

Dump JSON :

```
mongoexport -u cpsmongo -p canopsis -d canopsis -c <collection_name> -o /path/to/<collection_name>.json
mongoexport -u cpsmongo -p canopsis -d canopsis -c <collection_name> | gzip > /path/to/<collection_name>.json.gz
```

Restore JSON :

```
mongoimport -u cpsmongo -p canopsis -d canopsis -c <collection_name> /path/to/<collection_name>.json
zcat /path/to/<collection_name>.json.gz | mongoimport -u cpsmongo -p canopsis -d canopsis -c <collection_name>
```

### Dump avec filtrage

Exemple en récupérant tous les documents de la collection `periodical_alarm` de la BD Canopsis avec des documents ayant un champs `v.ack` :

```
mongodump -h <host> -u cpsmongo -p canopsis --db canopsis -c periodical_alarm --archive=bad_docs.bson.gz --gzip -q '{"v.ack":{$exists: true}}'
```

## Session Web

Voir [debug front](../guide_developpeur/debug_front.md)
