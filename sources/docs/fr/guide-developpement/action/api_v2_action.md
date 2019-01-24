# API Action

## Action

### Récupérer une action

```
GET /api/v2/actions/<action_id>
```

Renvoie l'action demandée.

```json
{
    "_id": "action_id",
    "type": "pbehavior",
    "fields": ["Resource", "Component"],
    "regex": ".*",
    "parameters": {
        "name": "pbehavior_name",
        "author": "System",
        "type": "Pause",
        "rrule": "",
        "reason": "",
        "tstart": 0,
        "tstop": 253402297199
    }
}
```

### Créer une action

```json
POST /api/v2/actions
{
    "_id": "action_id",
    "type": "pbehavior",
    "fields": ["Resource", "Component"],
    "regex": ".*",
    "parameters": {
        "name": "pbehavior_name",
        "author": "System",
        "type": "Pause",
        "rrule": "",
        "reason": "",
        "tstart": 0,
        "tstop": 253402297199
    }
}
```

Renvoie un dictionnaire vide en cas de réussite.

### Modifier une action

```json
PUT /api/v2/actions/<action_id>
{
    "_id": "action_id",
    "type": "pbehavior",
    "fields": ["Resource", "Component"],
    "regex": ".*",
    "parameters": {
        "name": "pbehavior_name",
        "author": "Myself",
        "type": "Pause",
        "rrule": "",
        "reason": "",
        "tstart": 0,
        "tstop": 253402297199
    }
}
```

Renvoie un dictionnaire vide en cas de réussite.

### Supprimer une action

```
DELETE /api/v2/actions/<action_id>
```

Renvoie un booléen.
