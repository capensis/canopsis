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
  "component": "scenario",
  "connector": "Engine",
  "connector_name": "JENKINS",
  "resource": "xxxx",
  "source_type": "resource"
}
```

Réponse, the modified alarm:
```json
{
  "connector": "Engine",
  "connector_name": "JENKINS",
  "component": "scenario",
  ...
}
```
