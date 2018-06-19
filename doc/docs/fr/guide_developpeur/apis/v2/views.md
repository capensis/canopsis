# API Vues

## Lister les vues

```
GET /api/v2/views[?name=...&title=...]
```

Renvoie une liste de toutes les vues, optionnellement filtrées par nom ou par
titre.

## Créer une vue

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

Renvoie l'id de la vue, qui est généré automatiquement.

## Supprimer, modifier ou récupérer une vue

```
DELETE,PUT,GET /api/v2/views/<view_id>
```


# Groupes

## Lister les groupes

```
GET /api/v2/views/groups
```

## Créer un groupe

```
POST /api/v2/views/groups/<group_id>
{
    "name": "..."
}
```

Renvoie une erreur si l'id (qui est donné par l'utilisateur) existe déjà.

## Lister les vues d'un groupe

```
GET /api/v2/views/groups/<group_id>[?name=...&title=...]
```

Renvoie une liste des vues d'un groupe, optionnellement filtrées par nom ou par
titre.

```
{
    "name": "...",
    "views": [
        ...
    ]
}
```

## Supprimer ou modifier un groupe

```
DELETE,PUT /api/v2/views/groups/<group_id>
```
