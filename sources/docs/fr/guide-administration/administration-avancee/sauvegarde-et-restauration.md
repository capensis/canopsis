# Backup & Restore

## Sauvegarde de la base MongoDB courante

La commande suivante fait une sauvegarde de la base de données courante, et peut prendre quelques heures (penser à un `tmux` pour reprendre la main plus tard si nécessaire) :
```
[Vers un point de montage disposant de suffisamment d'espace pour une sauvegarde de la base]
# MONGO_DATE=$(date '+%Y%m%d') ; mkdir -p /data/backup-mongo$MONGO_DATE
# mongodump -u cpsmongo -p canopsis -d canopsis -o /data/backup-mongo$MONGO_DATE/
```

```js
> db.stats()
{
        "db" : "canopsis",
        "collections" : 35,
        "objects" : 20 827 775,
        "avgObjSize" : 1903.0303928288067,
        "dataSize" : 39635888840,
        "storageSize" : 47567834864,
        "numExtents" : 153,
        "indexes" : 125,
        "indexSize" : 26848643136,
        "fileSize" : 92230647808,
        "nsSizeMB" : 16,
        "dataFileVersion" : {
                "major" : 4,
                "minor" : 5
        },
        "extentFreeList" : {
                "num" : 0,
                "totalSize" : 0
        },
        "ok" : 1
}
```

## Restauration de la sauvegarde MongoDB

Après avoir vérifié que la commande `mongodump` précédente a fonctionné, on peut effacer puis réimporter la base, afin que MongoDB puisse purger l'espace de stockage qui était encore réservé par les données supprimées.

Suppression de la base courante :
```
# mongo admin -u admin -p admin
> use canopsis
[s'assurer de bien avoir une sauvegarde OK avant d'exécuter la commande suivante !]
> db.dropDatabase()
> exit
```

Restauration de la sauvegarde. Cette étape peut, à nouveau, prendre quelques heures :
```
# mongorestore -u admin -p admin /data/backup-mongo$MONGO_DATE/
```
