# API alerts

## Done

Cette route permet l'envoi d'une action done.

#### Url

  `POST` /api/v2/alerts/done

#### POST exemple

/api/v2/alerts/done

json body:
```json
{
  "author": "root",
  "comment": "",
  "component": "gtk",
  "connector": "the",
  "connector_name": "knack",
  "resource": "my_sharona",
  "source_type": "resource"
}
```

Réponse, the modified alarm:
```json
{
  "connector": "the",
  "connector_name": "knack",
  "component": "gtk",
  ...
}
```
