# API TicketApiConfig

## TicketApiConfig

### Récupérer une ticketapi

```
GET /api/v2/ticketapi/<ticketapi_id>
```

Renvoie la config ticket api demandée.

```json
{
    "_id": "ticketapi_id",
    "type": "snow",
    "fields": ["Resource", "Component"],
    "regex": ".*",
    "api": {
        "base_url": "https://example.com/rest",
        "username": "",
        "password": "",
    },
    "parameters":{
        "Source": "Canopsis"
    }
}
```

### Créer une ticketapi

```json
POST /api/v2/ticketapi
{
    "_id": "ticketapi_id",
    "type": "snow",
    "fields": ["Resource", "Component"],
    "regex": ".*",
    "api": {
        "base_url": "https://example.com/rest",
        "username": "",
        "password": "",
    },
    "parameters":{
        "Source": "Canopsis"
    }
}
```

Renvoie un dictionnaire vide en cas de réussite.

### Modifier une ticketapi

```json
PUT /api/v2/ticketapi/<ticketapi_id>
{
    "_id": "ticketapi_id",
    "type": "snow",
    "fields": ["Resource", "Component"],
    "regex": ".*",
    "api": {
        "base_url": "https://example.com/rest",
        "username": "",
        "password": "",
    },
    "parameters":{
        "Source": "Canopsis"
    }
}
```

Renvoie un dictionnaire vide en cas de réussite.

### Supprimer une ticketapi

```
DELETE /api/v2/ticketapi/<ticketapi_id>
```

Renvoie un booléen.
