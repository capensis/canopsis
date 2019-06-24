# Actions sur la base de données

## MongoDB

### Purge

!!! attention
    Cette manipulation à un impact métier important et ne doit être réalisée que par une personne compétente.

Cette section va lister différentes commandes pour purger des collections de la base de données. La connexion à la base et la purge peuvent se faire via la ligne de commande (`mongo canopsis -u cpsmongo -p MOT_DE_PASSE --host XXXX`) ou bien via [Robo3T](https://robomongo.org). Dans les sous-sections suivantes, les commandes ont été réalisées en ligne de commande.

#### Purge simple d'une collection

Pour vider simplement les documents d'une collection, vous pouvez utiliser la commande `db.<nom de la collection>.remove({})`. La fonction `remove` en paramètre une requête, qui est ici la requête `{}` qui va matcher tous les documents de la collection.

Au moment de la purge, un message va indiquer le nombre d'éléments supprimés. Vous pouvez ensuite vérifier que `db.<nom de la collection>.find({})` ne retourne aucun résultat.

```bash
> db.periodical_alarm.remove({})
WriteResult({ "nRemoved" : 235 })
> db.entities.remove({})
WriteResult({ "nRemoved" : 17 })
> db.default_pbehavior.remove({})
WriteResult({ "nRemoved" : 6 })
> db.periodical_alarm.find({})
>
```

#### Purge d'une collection avec filtre

La fonction `remove` sur une collection prend en paramètre une requête. On peut donc filtrer sur les documents des collections.

Avant de supprimer ces documents, vous pouvez toujours vérifier la liste des documents concernés avec `find(<requête>)` et voir leur nombre `count(<requête>)`.

Pour les requêtes sur les dates, vous pouvez vous aider de sites comme [epochconverter.com](https://www.epochconverter.com/) pour convertir les dates en timestamp UNIX.

##### Alarmes

Voici une liste non exhaustive des requêtes portant sur différentes propriétés de la collection des alarmes, `periodical_alarm`.

| Type d'alarmes                                                            | Requête                                                                                 |
|:--------------------------------------------------------------------------|:----------------------------------------------------------------------------------------|
| Alarmes résolues                                                          | `db.periodical_alarm.find(`**`{"v.resolved":{$ne:null}}`**`)`                           |
| Alarmes non résolues                                                      | `db.periodical_alarm.find(`**`{"v.resolved":null}`**`)`                                 |
| Alarmes associées à l'entité `XXX/ZZZ`                                    | `db.periodical_alarm.find(`**`{"v.component" : "ZZZ", "v.resource" : "XXX"}`**`)`       |
| Alarmes non mises à jour depuis le 1er janvier 2019 00:00:00 GMT          | `db.periodical_alarm.find(`**`{"v.last_update_date":{$lte:1546300800}}`**`)`            |

##### Entités

Voici une liste non exhaustive des requêtes portant sur différentes propriétés de la collection des entités, `default_entities`.

### Sauvegarde

Utilisez la commande `mongodump` via une tâche cron. De préférence, faites la sauvegarde sur un système de fichier externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.mongodb.com/manual/tutorial/backup-and-restore-tools/#basic-mongodump-operation).

!!! note
    Le mot de passe par défaut est "canopsis" mais il peut être nécessaire d'adapter la commande selon votre contexte.

```bash
mongodump --username cpsmongo --password votre_password --db canopsis --out /path/to/backup
```

### Restauration

!!! attention
    Cette manipulation à un impact métier important et ne doit être réalisée que par une personne compétente. La restauration de la base de donnée ne doit être effectuée que si celle-ci est endommagée, pour corriger l'incident.

Avant de procéder à la restauration, arrêtez l'hyperviseur.
```shell
/opt/canopsis/bin/canopsis-systemd stop
```

Utilisez la commande `mongorestore`. De préférence, récupérez la sauvegarde depuis un système de fichier externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.mongodb.com/manual/tutorial/backup-and-restore-tools/#basic-mongorestore-operations).

!!! note
    Le mot de passe par défaut est "canopsis" mais il peut être nécessaire d'adapter la commande selon votre contexte.

```shell
mongorestore --username cpsmongo --password votre_password --db canopsis /path/to/backup
```

!!! note
    Lors du dump de la base, la commande créé un sous dossier dans `/path/to/backup` pour y stocker les fichiers. Ce sous-dossier doit être ajouté au `path` dans la commande `mongorestore`.

Si la restauration est réussie vous pouvez redémarrer l'hyperviseur.
```shell
/opt/canopsis/bin/canopsis-systemd start
```

## InfluxDB

### Sauvegarde

Utilisez la commande `influxd backup` via une tâche cron. De préférence, faites la sauvegarde sur un système de fichier externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.influxdata.com/influxdb/v1.7/administration/backup_and_restore/#backup).

```bash
influxd backup -portable -database canopsis /path/to/backup
```

### Restauration

!!! attention
    Cette manipulation à un impact métier important et ne doit être réalisée que par une personne compétente. La restauration de la base de donnée ne doit être effectuée que si celle-ci est endommagée, pour de corriger l'incident.

Avant de procéder à la restauration, arrêtez l'hyperviseur.
```shell
/opt/canopsis/bin/canopsis-systemd stop
```

Utilisez la commande `influxd restore`. De préférence, récupérez la sauvegarde depuis un système de fichier externe à la machine (NAS, SAN).Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.influxdata.com/influxdb/v1.7/administration/backup_and_restore/#restore).

```shell
influxd restore -portable /path/to/backup
```

!!! note
    Il est possible que la commande retourne un message d'erreur :

    ```
    error updating meta: DB metadata not changed. database may already exist
    restore: DB metadata not changed. database may already exist
    ```

    Il s'agit uniquement des metadatas qui sont déjà présentes dans Influx et n'ont pas changé. Le contenu de la table canopsis a bien été restauré.

Si la restauration est réussie vous pouvez redémarrer l'hyperviseur.

```shell
/opt/canopsis/bin/canopsis-systemd start
```
