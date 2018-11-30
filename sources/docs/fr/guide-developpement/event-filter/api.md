# API

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


# Utilisation de l'API avec curl

## Ajouter une règle

```shell
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
    "type": "enrichment",
    "pattern": {
        "component": "192.168.0.1"
    },
    "actions": [
        {
            "type": "set_field",
            "name": "Component",
            "value": "example.com"
        }
    ],
    "priority": 101,
    "on_success": "pass",
    "on_failure": "pass"
}' 'https://<hostname>/api/v2/eventfilter/rules'
```

`<hostname>` est à remplacer par le nom de domaine ou l'adresse IP du serveur
canopsis, et `root:root` par les identifiants de l'administrateur.

Si la règle est valide, cette requête renverra l'id de la règle, sous la forme :
```json
{"_id": "cc5f3dc2-0652-42f3-be50-5e34e5bb08b7"}
```

Sinon, elle renverra une erreur, et la règle ne sera pas ajoutée à l'event-filter.

## Lister les règles

```shell
curl -X GET -u root:root 'https://<hostname>/api/v2/eventfilter/rules'
```

`<hostname>` est à remplacer par le nom de domaine ou l'adresse IP du serveur
canopsis, et `root:root` par les identifiants de l'administrateur.

## Supprimer une règle

```shell
curl -X DEL -u root:root 'https://<hostname>/api/v2/eventfilter/rules/<rule_id>'
```

`<hostname>` est à remplacer par le nom de domaine ou l'adresse IP du serveur
canopsis, `<rule_id>` par l'id de la règle à supprimer, et `root:root` par les
identifiants de l'administrateur.

## Modifier une règle

```shell
curl -X PUT -u root:root -H "Content-Type: application/json" -d '{
    "type": "enrichment",
    "pattern": {
        "component": "192.168.0.1"
    },
    "actions": [
        {
            "type": "set_field",
            "name": "Component",
            "value": "s.example.com"
        }
    ],
    "priority": 101,
    "on_success": "pass",
    "on_failure": "pass"
}' 'https://<hostname>/api/v2/eventfilter/rules/<rule_id>'
```

`<hostname>` est à remplacer par le nom de domaine ou l'adresse IP du serveur
canopsis, `<rule_id>` par l'id de la règle à modifier, et `root:root` par les
identifiants de l'administrateur.
