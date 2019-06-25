# Actions sur la base de données

## MongoDB

### Purge

Cette section va lister différentes commandes pour purger des collections de la base de données. La connexion à la base et la purge peuvent se faire via la ligne de commande (`mongo canopsis -u cpsmongo -p MOT_DE_PASSE --host XXXX`) ou bien via [Robo3T](https://robomongo.org). Dans les sous-sections suivantes, les commandes ont été réalisées en ligne de commande.

!!! attention
    Cette manipulation à un impact métier important et ne doit être réalisée que par une personne compétente. **Avant toute purge, il est vivement conseillé de faire une [sauvegarde de la base Mongo](#Sauvegarde)** en utilisant `mongorestore`.


Avant de supprimer des documents, vous pouvez toujours vérifier la liste des documents concernés avec `db.<nom de la collection>.find(<requête>)` et voir leur nombre `db.<nom de la collection>.count(<requête>)`. Ces fonctions prennent en paramètre une requête, qui va filtrer sur les documents de la collection.

Une fois que vous avez vérifié que les documents correspondent à ce que vous voulez supprimez, vous pouvez `db.<nom de la collection>.remove(<requête>)`. Au moment de la purge, un message va indiquer le nombre d'éléments supprimés.

```bash
> db.periodical_alarm.remove({"t" : 1537894605})
WriteResult({ "nRemoved" : 3 })
> db.entities.remove({"name": "foldable"})
WriteResult({ "nRemoved" : 17 })
```

!!! attention
    La requête vide `{}` va matcher tous les documents de la collection. Par conséquent, **`db.<nom de la collection>.remove({})` va vider complètement la collection**. Pensez donc à ne jamais avoir `{}` comme paramètre, sauf si vous voulez vider complètement la collection.


Les sous-sections suivantes vont montrer plusieurs examples de requêtes sur les collections d'objets Canopsis. Pour les requêtes sur les dates, vous pouvez vous aider de sites comme [epochconverter.com](https://www.epochconverter.com/) pour convertir les dates en timestamp UNIX.

#### Alarmes

Voici une liste non exhaustive des requêtes portant sur différentes propriétés de la collection des alarmes, `periodical_alarm`.

| Type d'alarme                                                             | Requête                                                                                 |
|:--------------------------------------------------------------------------|:----------------------------------------------------------------------------------------|
| Alarmes résolues                                                          | `db.periodical_alarm.find(`**`{"v.resolved":{$ne:null}}`**`)`                           |
| Alarmes non résolues                                                      | `db.periodical_alarm.find(`**`{"v.resolved":null}`**`)`                                 |
| Alarmes associées à l'entité `XXX/ZZZ`                                    | `db.periodical_alarm.find(`**`{"v.component" : "ZZZ", "v.resource" : "XXX"}`**`)`       |
| Alarmes non mises à jour depuis le 1er janvier 2019 00:00:00 GMT          | `db.periodical_alarm.find(`**`{"v.last_update_date":{$lte:1546300800}}`**`)`            |

#### Entités

Voici une liste non exhaustive des requêtes portant sur différentes propriétés de la collection des entités, `default_entities`.

| Type d'entité                                                             | Requête                                                                                 |
|:--------------------------------------------------------------------------|:----------------------------------------------------------------------------------------|
| Expression régulière sur l'attribut `client` dans l'entité                | `db.default_entities.find(`**`{"infos.client.value":{$regex:'.*SSBU.*',$options:'i'}}`**`)`|

#### Pbehavior

Voici une liste non exhaustive des requêtes portant sur différentes propriétés de la collection des pbehaviors, `default_pbehavior`.

| Type de pbehavior                                                         | Requête                                                                                 |
|:--------------------------------------------------------------------------|:----------------------------------------------------------------------------------------|
| Pbehaviors créés par `emile-zola`                                         | `db.default_pbehavior.find(`**`{"author":"emile-zola"}`**`)`                            |

#### Vues

Voici une liste non exhaustive des requêtes portant sur différentes propriétés de la collection des vues, `default_pbehavior`.

| Type de vue                                                               | Requête                                                                                 |
|:--------------------------------------------------------------------------|:----------------------------------------------------------------------------------------|
| Vues désactivées                                                          | `db.views.find(`**`{"enabled":false}`**`)`                                              |

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
