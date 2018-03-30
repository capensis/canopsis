# API événements 

## envoi d'événements

Cette route permet l'envoi d'événements canopsis.

#### Url

  `POST` /api/v2/event

#### POST exemple

/api/v2/event

json body:
```{json}
[
    {
        "component" : "component_1",
        "source_type" : "component",
        "event_type" : "check",
        "connector" : "connector",
        "connector_name" : "connector_name",
        "output" :"output",
        "state" :1
    },
    {
        "component" : "composant_2",
        "source_type" : "component",
        "event_type" : "check",
        "connector" : "connector",
        "connector_name" : "connector_name",
        "output" :"output",
        "state" :1
    }
]
```

Réponse:
```{json}
[
    {
        "component" : "component_1",
        "source_type" : "component",
        "event_type" : "check",
        "connector" : "connector",
        "connector_name" : "connector_name",
        "output" :"output",
        "state" :1
    },
    {
        "component" : "composant_2",
        "source_type" : "component",
        "event_type" : "check",
        "connector" : "connector",
        "connector_name" : "connector_name",
        "output" :"output",
        "state" :1
    }
]
```
