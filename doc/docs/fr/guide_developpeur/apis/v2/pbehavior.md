# API événements

## Creation d'un pbheavior

Cette route permet de créer un pbehavior.

#### Url

  `POST` /api/v2/pbehavior

#### POST exemple

/api/v2/pbehavior

json body:
```{json}
{
    "name": "imagine",
    "author": "lennon",
    "filter_": {"_id": "all_the_people"},
    "rrule": "",
    "tstart": 0,
    "tstop": 1000
}
```

Réponse: l'uid de l'élément inséré
```{json}
"b72e841a-d9d1-11e7-9a70-022abfd0f78f"
```

## Supprimer un pbheavior

Cette route permet de supprimer un pbehavior.

#### Url

  `DELETE` /api/v2/pbehavior/<pbehavior_id>

#### DELETE exemple

/api/v2/pbehavior/<pbehavior_id>

Réponse: un dictionnnaire de status
```{json}
{
    "deletedCount": 1,
    "acknowledged": true
}
```

## Forcer le recacule des pbheaviors

Cette route permet de forcer le recalcule de tous les pbehaviors.

La route n'est appelable qu'une fois toutes les 10 secondes.

#### Url

  `GET` /api/v2/compute-pbehaviors

#### GET exemple

/api/v2/compute-pbehaviors

Réponse: le recalcule a-t'il bien été lancé ?
```{json}
true
```
