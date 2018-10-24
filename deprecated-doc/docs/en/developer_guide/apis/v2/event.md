# Event APIs

## send event

This endpoint send canopsis events.

#### Url

  `POST` /api/v2/event

#### POST example

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

Response:
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
