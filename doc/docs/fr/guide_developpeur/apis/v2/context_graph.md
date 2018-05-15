# API Context Graph
# API Context Graph

## Récupérer l'id d'une entité du context

Cette route permet de récupérer l'id qui sera donné à une entité du context, selon un event.

#### Url

  `POST` /api/v2/context_graph/get_id/

#### POST exemple

/api/v2/context_graph/get_id/

json body:
```{json}
{
    "event_type": "check",
    "timestamp": 1512486748,
    "connector": "cap_kirk",
    "connector_name": "spock",
    "component": "mc_coy",
    "resource": "uhura",
    "source_type": "resource",
    "state": 2,
    "output": "NCC_1701"
}
```

Réponse: l'id de l'entité désirée
```{json}
uhura/mc_coy
```


## Supprimer une entité du context

Cette route permet de récupérer une entité dans le context.

#### Url

  `DELETE` /api/v2/context/<entity_id>

#### DELETE exemple

/api/v2/context/<entity_id>

Réponse: null si tout s'est bien passé
```{json}
null
```
