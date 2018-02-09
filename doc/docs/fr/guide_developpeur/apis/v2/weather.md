# API weather

## Récupération de la météo (get watchers)

Cette route permet de récupérer tout ou partie de la météo.

#### Url

  `GET` /api/v2/weather/watchers/<watcher_filter>

#### GET exemple

/api/v2/weather/watchers/{}?limit=30&start=0&stop=10&pb_types=["pause", "maintenance"]

* `limit`, `start` et `stop` (optionel) permettent de paginer la réponse.
* `pb_types` (optionnel) permet de ne prendre en compte que les
pbehaviors d'un certain type (parmi une liste).

Réponse: une liste des watchers correspondants, enrichi avec les entités
```{json}
[
  {
    "alerts_not_ack": false,
    "criticity": "",
    "display_name": "let",
    "entity_id": "forever",
    "has_baseline": false,
    "hasactivepbehaviorinentities": false,
    "hasallactivepbehaviorinentities": false,
    "infos": {}
    "linklist": [],
    "mfilter": "{\"entity_id\": \"be\"}",
    "org": "",
    "pbehavior": [],
    "sla_text": "",
    "state": {
      "val": 0
    }
  },
  ...
]
```


## Récupération d'un élément de météo (weatherwatchers)

Cette route permet de récupérer la météo d'un élément (watcher).

#### Url

  `GET` /api/v2/weather/watchers/<watcher_id>

#### GET exemple

/api/v2/weather/watchers/tcb

Réponse: une liste des éléments surveillés par un watcher
```{json}
[
  {
    "entity_id": "the/sunshine",
    "infos": {},
    "linklist": [
      {
        "cat_name": "consignes",
        "links": []
      }
    ],
    "name": "underground",
    "org": "",
    "pbehavior": [],
    "sla_text": "",
    "source_type": "resource",
    "state": {
      "val": 0
    }
  },
  ...
]
```
