# Backup & Restore

## MongoDB

### Sauvegarde de la base

On utilise la commande `mongodump` via une tâche cron. La sauvegarde sera, de préférence, effectuée sur un filesystem externe à la machine (NAS, SAN).  
Pour obtenir plus d'information à propos de cette commande vous pouvez consulter cette [page](https://docs.mongodb.com/manual/tutorial/backup-and-restore-tools/#basic-mongodump-operations).

```bash
$ mongodump --db canopsis --out <path-to-backup>
```

### Restoration de la base

On utilise la commande `mongorestore` via une tâche cron. La sauvegarde sera, de préférence, effectuée sur un filesystem externe à la machine (NAS, SAN).  
Pour obtenir plus d'information à propos de cette commande vous pouvez consulter cette [page](https://docs.mongodb.com/manual/tutorial/backup-and-restore-tools/#basic-mongorestore-operations).

```shell
$ mongorestore <path-to-backup>
```

## InfluxDB

Pour obtenir plus d'information à propos des commandes utilisées dans cette section vous pouvez consulter cette [page](https://docs.influxdata.com/influxdb/v1.7/administration/backup_and_restore/).

### Sauvegarde de la base

On utilise la commande `influxd backup` via une tâche cron. La sauvegarde sera, de préférence, effectuée sur un filesystem externe à la machine (NAS, SAN).  

```bash
$ influxd backup -portable -database canopsis <path-to-backup>
```

### Restoration de la base

On utilise la commande `influxd restore` via une tâche cron. La sauvegarde sera, de préférence, effectuée sur un filesystem externe à la machine (NAS, SAN).  

```shell
$ influxd restore -portable <path-to-backup>
```
