# API event filter

L'API event filter permet de manipuler les règles du filtre à événements du
moteur go `che`.

## Lister les règles

```
GET /api/v2/eventfilter/rules
```

Renvoie un tableau contenant toutes les règles de l'event filter.

```json
[
    {
        "_id": "6b90880a-c4f0-4a4d-8a51-de0c7e14581e",
        "type": "drop",
        "pattern": {...},
        "priority": 100,
    },
    ...
]
```

## Récupérer une règle

```
GET /api/v2/eventfilter/rules/<rule_id>
```

Renvoie la règle dont l'id vaut `<rule_id>`, ou une erreur si celle-ci
n'existe pas.

## Créer une règle

```json
POST /api/v2/eventfilter/rules
Content-Type: "application/json"
{
    "type": "drop",
    "pattern": {...},
    "priority": 100,
}
```

Crée la règle et renvoie son id si elle est valide. Une erreur est renvoyée
si la règle est invalide.

## Supprimer une règle

```
DEL /api/v2/eventfilter/rules/<rule_id>
```

Supprime la règle dont l'id vaut `<rule_id>`, ou renvoie une erreur si
celle-ci n'existe pas.

## Modifier une règle
Content-Type:·"application/json"
```json
PUT /api/v2/eventfilter/rules/<rule_id>
Content-Type: "application/json"
{
    "type": "drop",
    "pattern": {...},
    "priority": 100,
}
```

Modifie la règle dont l'id vaut `<rule_id>`, ou renvoie une erreur si la
règle est invalide, ou si son id a été modifié.
