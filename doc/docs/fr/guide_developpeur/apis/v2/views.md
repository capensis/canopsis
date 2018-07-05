# API Vues

## Vues

### Lister les vues

```
GET /api/v2/views
GET /api/v2/views?name=...
GET /api/v2/views?title=...
```

Renvoie toutes les vues, optionnellement filtrées par nom ou par titre.

```
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

```
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

```
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

```
POST /api/v2/views/groups/<group_id>
{
    "name": "..."
}
```

Renvoie une erreur si l'id (qui est donné par l'utilisateur) existe déjà.

### Lister les vues d'un groupe

```
GET /api/v2/views/groups/<group_id>
GET /api/v2/views/groups/<group_id>?name=...
GET /api/v2/views/groups/<group_id>?title=...
```

Renvoie une liste des vues d'un groupe, optionnellement filtrées par nom ou par
titre.

```
{
    "_id": "...",
    "name": "...",
    "views": [
        ...
    ]
}
```

### Modifier un groupe

```
PUT /api/v2/views/groups/<group_id>
{
    "name": "..."
}
```

### Supprimer un groupe

```
DELETE /api/v2/views/groups/<group_id>
```
