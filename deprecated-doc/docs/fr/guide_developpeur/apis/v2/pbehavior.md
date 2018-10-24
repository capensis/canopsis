# API PBehaviors

Les pbehaviors sont des évènements de calendrier récurrents qui permettent de mettre en pause la surveillance d'une alarme pendant une période donnée (pour des maintenances ou des astreintes par exemple).


**Note sur les attributs `tstart`et `tstop` d'un pbehavior**

un pbehavior peut être vu comme un élément de calendrier, avec une récurrence optionnelle. Les paramètres tstart et tstop sont utilisés pour définir la date et l'heure de début et de fin de l'évènement. En cas de récurrence, ils sont utilisés pour que chaque nouvelle instance de l'évènement démarre et se termine à la même heure que la première instance.

Par exemple, le pbehavior suivant durera 1 heure de 11H à midi tous les lundi à partir du 18/06/18 :

```json
{
// --snip --
	"rrule": "FREQ=WEEKLY;WKST=MO",
	"tstart": 1529312400, //Le 18/06/2018 à 11H
    "tstop": 1529316000, // Le 18/06/2018 à midi
// -- snip --
}
```

**En cas de récurrence, la date de dernière occurrence est déterminée par l'attribut `UNTIL` de la `rrule`**


## Creation d'un pbheavior

Cette route permet de créer un pbehavior.

#### Url

  `POST /api/v2/pbehavior`

#### POST exemple

json body:

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

Les attributs possibles de l'objet sont les suivants :


| Nom            | type    | nullable | Description                                                                           |
|----------------|---------|----------|---------------------------------------------------------------------------------------|
| connector      | string  | Non      | identifiant du connecteur de l'entité                                                 |
| name           | string  | Non      | Nom d'affichage du pbehavior                                                          |
| author         | string  | Non      | Nom de l'auteur du pbehavior                                                          |
| enabled        | boolean | Non      | Le pbehavior est il autorisé à se déclencher ou non                                   |
| reason         | string  | Oui      | Motif (optionnel)                                                                     |
| comments       | array   | Oui      | Commentaires (optionnels)                                                             |
| filter         | string  | Non      | Filtre d'entités (json)                                                               |
| type_          | string  | Non      | Type de pbehavior                                                                     |
| connector_name | string  | Non      | Nom d'affichage du connecteur de l'entité                                             |
| rrule          | string  | Oui      | Rrule (récurrence)                                                                    |
| tstart         | integer | Non      | Timestamp de la date de début                                                         |
| tstop          | integer | Non      | Timestamp de la date de fin (en cas de récurrence, seule l'heure est prise en compte) |
| _id            | string  | Non      | Identifiant du pbehavior                                                              |
| eids           | array   | Non      | Tableau d'identifiants des entités concernées.                                        |

Réponse: l'uid de l'élément inséré

```json
"b72e841a-d9d1-11e7-9a70-022abfd0f78f"
```

## Récupérer les pbehaviors d'une entité

Cette route permet de lister les pbeahviors présents sur une entité, définie par son eid (Entity ID)

#### URL
`GET /api/v2/pbehavior_byeid/<entityid>`

#### Paramètres
* entityid <string> the ID of the target entity.


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

Les attributs de la réponse sont les suivants:

| Nom            | type    | nullable | Description                                                                           |
|----------------|---------|----------|---------------------------------------------------------------------------------------|
| connector      | string  | Non      | Identifiant du connecteur de l'entité                                                 |
| name           | string  | Non      | Nom d'affichage du pbehavior                                                          |
| author         | string  | Non      | Nom de l'auteur du pbehavior                                                          |
| enabled        | boolean | Non      | Le pbehavior est il autorisé à se déclencher ou non                                   |
| reason         | string  | Oui      | Motif (optionnel)                                                                     |
| comments       | array   | Oui      | Commentaires (optionnels)                                                             |
| filter         | string  | Non      | Filtre d'entités (json)                                                               |
| type_          | string  | Non      | Type de pbehavior                                                                     |
| connector_name | string  | Non      | Nom d'affichage du connecteur de l'entité                                             |
| rrule          | string  | Oui      | Rrule (récurrence)                                                                    |
| tstart         | integer | Non      | Timestamp de la date de début                                                         |
| tstop          | integer | Non      | Timestamp de la date de fin (en cas de récurrence, seule l'heure est prise en compte) |
| _id            | string  | Non      | Identifiant du pbehavior                                                              |
| isActive       | boolean | Non      | Le pbehavior est-il dans sa période d'activité                                        |
| eids           | array   | Non      | Tableau d'identifiants des entités concernées.                                        |



## Supprimer un pbheavior

Cette route permet de supprimer un pbehavior.

#### Url

  `DELETE /api/v2/pbehavior/<pbehavior_id>`

#### DELETE exemple

/api/v2/pbehavior/<pbehavior_id>

Réponse: un dictionnnaire de status

```json
{
    "deletedCount": 1,
    "acknowledged": true
}
```


## Mettre à jour un pbeahvior

Il n'existe pas à ce jour de méthode de mise à jour d'un pbheavior. Il est donc nécessaire de supprimer puis re-créer le pbehavior pour en mettre à jour le contenu.

## Forcer le recacul des pbheaviors

Cette route permet de forcer le recalcul de tous les pbehaviors.

#### Url

  `GET` /api/v2/compute-pbehaviors

#### GET exemple

/api/v2/compute-pbehaviors

Réponse: le recalcule a-t'il bien été lancé ?

```json
true
```
