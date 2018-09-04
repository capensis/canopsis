# API Vues

## Vues

### Lister les vues

```
GET /api/v2/views
GET /api/v2/views?name=...
GET /api/v2/views?title=...
```

Renvoie toutes les vues, optionnellement filtrées par nom ou par titre.

```json
{
    "groups": {
        "group1": {
            "name": "...",
            "views": [
                // ...
            ]
        },
        // ...
    }
}
```

### Créer une vue

```json
POST /api/v2/views
{
    "group_id": "<group_id>",
    "type": "...",
    "name": "...",
    "title": "..."
    // ...
}
```

Crée une vue et renvoie son id (qui est généré automatiquement).

### Récupérer une vue

```
GET /api/v2/views/<view_id>
```

### Modifier une vue

```json
PUT /api/v2/views/<view_id>
{
    "group_id": "<group_id>",
    "type": "...",
    "name": "...",
    "title": "..."
    // ...
}
```

### Supprimer une vue

```
DELETE /api/v2/views/<view_id>
```


## Groupes

### Lister les groupes

```
GET /api/v2/views/groups
GET /api/v2/views/groups?name=...
```

Renvoie la liste des groupes, optionnellement filtrées par nom.

### Créer un groupe

```json
POST /api/v2/views/groups
{
    "name": "..."
}
```

Crée un groupe et renvoie son id (qui est généré automatiquement).

### Lister les vues d'un groupe

```
GET /api/v2/views/groups/<group_id>
GET /api/v2/views/groups/<group_id>?name=...
GET /api/v2/views/groups/<group_id>?title=...
```

Renvoie une liste des vues d'un groupe, optionnellement filtrées par nom ou par
titre.

```json
{
    "_id": "...",
    "name": "...",
    "views": [
        // ...
    ]
}
```

### Modifier un groupe

```json
PUT /api/v2/views/groups/<group_id>
{
    "name": "..."
}
```

### Supprimer un groupe

```
DELETE /api/v2/views/groups/<group_id>
```
