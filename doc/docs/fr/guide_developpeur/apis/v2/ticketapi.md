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
    "parameters": {
        "ticket_path": "table/incident",
        "result_key": "result",
        "number_key": "number"
    },
    "payload": {
        "Source": "Club jdg",
        "Name": "{{ .Alarm.Value.Resource }}"
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
        "ticket_path": "table/incident",
        "result_key": "result",
        "number_key": "number"
    },
    "payload":{
        "Source": "Canopsis",
        "Name": "{{ .Alarm.Value.Resource }}"
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
        "username": "david",
        "password": "goodenough"
    },
    "parameters": {
        "result_key": "result",
        "ticket_path": "table/incident",
        "number_key": "number"
    },
    "payload":{
        "Source": "Canopsis NG",
        "Name": "{{ .Alarm.Value.Component }}"
    }
}
```

Renvoie un dictionnaire vide en cas de réussite.

### Supprimer une ticketapi

```
DELETE /api/v2/ticketapi/<ticketapi_id>
```

Renvoie un booléen.
