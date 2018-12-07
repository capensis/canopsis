# Utilisation

# API PBehaviors

Pbehaviors sont des évènements de calendrier récurrents qui arrêtent temporairement la surveillance d'une entité pendant un temps donné (pour la maintenance par exemple).

**Note sur les attributs tstart et tstop d'un pbehavior**

Les comportements sont similaires aux évènements de calendrier, avec une récurrence facultative.
Les paramètres tstart et tstop servent à définir les dates de début et de fin de la première occurrence d'évènement.
Lorsque l'évènement est répété, ces attributs sont utilisés pour définir la durée de chaque instance, en fonction des heures de début et de fin de la première instance.

Par exemple, le comportement ci-dessous commence à 11 heures et se termine une heure plus tard tous les matins à partir du 2018/06/18:

```js
{
   // --snip --
   "rrule": "FREQ=WEEKLY;WKST=MO",
   "tstart": 1529312400, //Le 2018/06/18 at 11
   "tstop": 1529316000, // Le 2018/06/18 at 12
   // -- snip --
}
```

**Lorsque l'évènement est récurrent, la date de la dernière occurrence est stockée dans l'attribut `UNTIL` de l' évènement.rrule**

## Créer un pbehavior

#### Url

  `POST /api/v2/pbehavior`

#### POST exemple

json :

```json
{
    "name": "imagine",
    "author": "lennon",
    "filter_": {"_id": "all_the_people"},
    "rrule": "",
    "tstart": 0,
    "tstop": 1000
}
```

Les attributs du corps sont les suivants:

| Name            | type    | nullable | Description                                  |
|-----------------|---------|----------|----------------------------------------------|
| connector       | string  | No       | Identifiant du connecteur d'entité           |
| name            | string  | No       | Display name du pbehavior                    |
| author          | string  | No       | Nom de l'auteur                              |
| enabled         | boolean | No       | Si le pbehavior est déclenché ou non         |
| reason          | string  | yes      | motif d'administration (optionnel)           |
| comments        | array   | yes      | Commentaires (option)                        |
| filter          | string  | No       | filtre d'entité (JSON)                       |
| type\_          | string  | No       | type de Pbehavior                            |
| connector\_name | string  | No       | Display name du connector                    |
| rrule           | string  | yes      | Rrule (récurrence)                           |
| tstart          | integer | No       | Timestamp de la date de départ               |
| tstop           | integer | No       | Timestamp de la date de fin                  |
| \_id            | string  | No       | indentifiant du Pbehavior                    |
| eids            | array   | No       | tableau du \_ids pour les entités impactées. |


Réponse : UID de l'élément inséré

```{json}
"b72e841a-d9d1-11e7-9a70-022abfd0f78f"
```

## Récupérer les pbehaviors d'une entité

Cette route répertorie les pbehaviors existant sur une entité, identifiée par son eid (Entity ID)

#### URL

`GET /api/v2/pbehavior_byeid/<entityid>`

#### Parameters

* entityid l'ID de l'entité cible.

#### Réponse

```json
[
    {
        "connector": "canopsis",
        "name": "imagine",
        "author": "lennon",
        "enabled": true,
        "reason": "",
        "comments": null,
        "filter": "{\"_id\": \"580059AB4B100031\"}",
        "type_": "generic",
        "connector_name": "canopsis",
        "rrule": "FREQ=WEEKLY;COUNT=30;WKST=MO",
        "tstart": 1529312725,
        "tstop": 1592471125,
        "_id": "dd4cbc2c-72d6-11e8-a732-0242ac12001a",
        "isActive": true,
        "eids": [
            "580059AB4B100031"
        ]
    }
]
```

Les attributs de réponse sont les suivants :

| Name            | type    | nullable | Description                              |
|-----------------|---------|----------|------------------------------------------|
| connector       | string  | No       | Identifier of the entity connector       |
| name            | string  | No       | Display name of the pbehavior            |
| author          | string  | No       | Author name                              |
| enabled         | boolean | No       | Should the pbehavior trigger or not      |
| reason          | string  | yes      | Administrative reason (optionnal)        |
| comments        | array   | yes      | Comments (option)                        |
| filter          | string  | No       | Entities filter (json)                   |
| type\_          | string  | No       | Pbehavior type                           |
| connector\_name | string  | No       | Display name of the entity connector     |
| rrule           | string  | yes      | Rrule (recurrence)                       |
| tstart          | integer | No       | Timestamp of the start date              |
| tstop           | integer | No       | Timestamp  end date                      |
| \_id            | string  | No       | Pbehavior identifier                     |
| eids            | array   | No       | Array of _ids for the impacted entities. |
| isActive        | boolean | No       | is the pbehavior currently active        |

## Supprimer un pbehavior

Cette route permet de supprimer un pbehavior

#### Url

  `DELETE /api/v2/pbehavior/<pbehavior_id>`

#### Exemple de suppression

Response: a status object

```json
{
    "deletedCount": 1,
    "acknowledged": true
}
```


## Update un pbeahvior

Il n'y a actuellement aucune méthode pour mettre à jour un comportement en place. Il est nécessaire de supprimer et de recréer un comportement pour mettre à jour son contenu.

## Forcer le calcul de pbehaviors

Cette route impose un nouveau calcul pour tous les comportements.


#### URL

GET / api / v2 / compute-pbehaviors


#### GET exemple

/ api / v2 / compute-pbehaviors

Réponse: les calculs ont-ils été traités?


```json
true
```
