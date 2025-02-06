# Nettoyage, sauvegarde et restauration des bases de données

## MongoDB

### Nettoyage

Cette section va lister différentes commandes pour purger des collections de la base de données. La connexion à la base et le nettoyage peuvent se faire via la ligne de commande (`mongosh canopsis -u cpsmongo -p MOT_DE_PASSE --host XXXX`) ou bien via [Robo3T](https://robomongo.org). Dans les sous-sections suivantes, les commandes ont été réalisées en ligne de commande.

!!! attention
    Cette manipulation a une incidence métier importante et ne doit être réalisée que par une personne compétente. **Avant toute opération, il est vivement conseillé de faire une [sauvegarde de la base MongoDB](#sauvegarde)**.

Avant de supprimer des documents, vous pouvez toujours vérifier la liste des documents concernés avec `db.<nom de la collection>.find(<requête>)` et voir leur nombre `db.<nom de la collection>.count(<requête>)`. Ces fonctions prennent en paramètre une requête, qui va filtrer sur les documents de la collection.

Une fois que vous avez vérifié que les documents correspondent à ce que vous voulez supprimer, vous pouvez utiliser la commande `db.<nom de la collection>.remove(<requête>)`. Au moment de la suppression, un message va indiquer le nombre d'éléments supprimés.

```js
> db.periodical_alarm.remove({"t" : 1537894605})
WriteResult({ "nRemoved" : 3 })
> db.entities.remove({"name": "foldable"})
WriteResult({ "nRemoved" : 17 })
```

!!! attention
    La requête vide `{}` va s'appliquer à tous les documents de la collection. Par conséquent, **`db.<nom de la collection>.remove({})` va vider complètement la collection**. Pensez donc à ne jamais avoir `{}` comme paramètre, sauf si vous voulez vider complètement la collection.

Le tableau suivant montre plusieurs examples de requêtes sur les collections d'objets Canopsis, avec la collection MongoDB en _italique_ et le filtre en **gras**. Pour rappel, les entités sont stockées dans la collection `default_entities` tandis que les alarmes sont stockées dans `periodical_alarm`, les comportements périodiques dans `default_pbehavior` et les vues dans `views`.

| Type d'objets                                                             | Requête MongoDB                                                                               |
|:--------------------------------------------------------------------------|:----------------------------------------------------------------------------------------------|
| Action de type `snooze`                                                   | `db.`_`default_action`_`.find(`**`{"type":"snooze"}`**`)`                                     |
| Alarmes résolues                                                          | `db.`_`periodical_alarm`_`.find(`**`{"v.resolved":{$ne:null}}`**`)`                           |
| Alarmes non résolues                                                      | `db.`_`periodical_alarm`_`.find(`**`{"v.resolved":null}`**`)`                                 |
| Alarmes associées à l'entité `XXX/ZZZ`                                    | `db.`_`periodical_alarm`_`.find(`**`{"v.component" : "ZZZ", "v.resource" : "XXX"}`**`)`       |
| Alarmes non mises à jour depuis le 1er janvier 2019 00:00:00 GMT          | `db.`_`periodical_alarm`_`.find(`**`{"v.last_update_date":{$lte:1546300800}}`**`)`            |
| Expression régulière sur l'attribut `client` dans l'entité                | `db.`_`default_entities`_`.find(`**`{"infos.client.value":{$regex:'.*SSBU.*',$options:'i'}}`**`)`|
| Comportements périodiques créés par `emile-zola`                                         | `db.`_`default_pbehavior`_`.find(`**`{"author":"emile-zola"}`**`)`                            |
| Comoportements périodiques avec un `tstart` placé dans le futur                           | `db.`_`default_pbehavior`_`.find(`**`{"tstart" : {$gt : Math.floor(Date.now() / 1000)}})`**`)`|
| Vues désactivées                                                          | `db.`_`views`_`.find(`**`{"enabled":false}`**`)`                                              |

Pour les requêtes sur les dates, vous pouvez vous aider de sites comme [epochconverter.com](https://www.epochconverter.com/) pour convertir les dates en timestamp UNIX. Le timestamp correspondant au temps courant est `Math.floor(Date.now() / 1000)`. Vous pouvez également vous servir des objets `Date()` pour les requêtes temporelles (voir la [documentation officielle de MongoDB sur `Date()`](https://docs.mongodb.com/manual/reference/method/Date/index.html)) qu'il faudra convertir en timestamp pour le filtrage. L'exemple suivant montre l'affichage des alarmes non mises à jour depuis un mois.

```js
> var d = new Date();
> d.setMonth(d.getMonth() - 1);
> oneMonthAgo = Math.floor(d / 1000);
> db.periodical_alarm.find({"v.last_update_date":{$lte:oneMonthAgo}})
```

Il est également possible de filtrer grâce aux expressions régulières en utilisant l'opérateur `$regex` (voir la [documentation officielle de MongoDB sur `$regex`](https://docs.mongodb.com/manual/reference/operator/query/regex/index.html)). Les deux lignes ci-dessous sont équivalentes, elles vont afficher les alarmes dont la ressource correspond à la regex .`[0-9a-fA-F]+`.

```js
> db.periodical_alarm.find({"v.resource":{$regex:'[0-9a-f]+',$options:'i'}}) // L'option i rend la regex insensible à la casse
> db.periodical_alarm.find({"v.resource":{$regex:/[0-9a-f]+/i}})             // On retrouve ici également l'option i
```

### Sauvegarde

Utilisez la commande `mongodump` via une tâche cron. De préférence, faites la sauvegarde sur un système de fichiers externe à la machine (NAS, SAN). Vous pouvez consulter la documentation de la commande en suivant ce [lien](https://docs.mongodb.com/manual/tutorial/backup-and-restore-tools/#basic-mongodump-operation).

!!! note
    Le mot de passe par défaut est `canopsis`, mais il peut être nécessaire d'adapter la commande selon votre contexte.

```sh
mongodump --username cpsmongo --password canopsis --db canopsis --out /chemin/vers/sauvegarde
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
mongorestore --username cpsmongo --password canopsis --db canopsis /chemin/vers/sauvegarde
```

!!! note
    Lors de la sauvegarde de la base, la commande crée un sous-dossier dans `/chemin/vers/sauvegarde` pour y stocker les fichiers. Ce sous-dossier doit être ajouté au chemin dans la commande `mongorestore`.

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
                           canopsis-service@canopsis-api.service \
                           canopsis-engine-python-snmp.service \
                           canopsis.service
    ```
