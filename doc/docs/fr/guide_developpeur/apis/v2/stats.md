# API Statistiques

## Routes

### Calcul d'une statistique

#### URL

`POST /api/v2/stats/<nom de la statistique>`

#### Paramètres

Cette route accepte en requête un objet JSON contenant les paramètres
suivants :

 - `tstop` : un timestamp indiquant le fin de la période pour laquelle la
   statistique doit être calculées. Ce timestamp doit correspondre à une heure
   pile (e.g. 12:00, et non 12:03).
 - `duration` : la durée de la période, représentée par une chaîne
   `"<n><unité>"`, avec `<n>` un entier et `<unité>` une unité de temps (`h`,
   `d` ou `w`).
 - `mfilter` : un filtre mongodb, filtrant les entités pour lesquelles les
   statistiques doivent être calculées.
 - `parameters` : un objet contenant les paramètres spécifiques à la
   statistique calculée. Ces paramètres sont précisés dans la documentation de
   chacune des statistiques.
 - `trend` (optionnel): `true` pour calculer la tendance par rapport à la
   période précédente.
 - `sla` (optionnel): un SLA, représenté par une inégalité (e.g. `">= 0.99"`).

#### Réponse

En cas de succès, la réponse est un objet JSON contenant un champ `values`. Ce
champ est un tableau contenant les valeurs de la statistique pour chaque
entité, sous la forme suivante :

```javascript
{
    'entity': {...},  // L'entité pour laquelle la statistique a été calculée
    'value': ...,  // La valeur de la statistique
    'trend': ...,  // La tendance
    'sla': ...  // true si la valeur est conforme au SLA
}
```

#### Exemple

La requête suivante renvoie le nombre d'alarmes critiques ouvertes sur chaque
resource impactant l'entité `service`, les 18 et 19 août 2018.

```javascript
POST /api/v2/stats/alarms_created
{
    "mfilter": {
        "type": "resource",
        "impact": {
            "$in": ["service"]
        }
    },
    "tstop": 1534716000,  // 20 août à 00:00
    "duration": "2d",
    "parameters": {
        "states": [3]
    }
}
```

Le document JSON ci-dessous est un exemple de réponse à la requête précédente.

```javascript
{
    "values": [
        {
            "entity": {
                "_id": "resource1/component1",
                "type": "resource"
                "impact": [
                    "service"
                ],
                // ...
            },
            "value": 117,
            "trend": 96
        },
        {
            "entity": {
                "_id": "resource2/component1",
                "type": "resource"
                "impact": [
                    "service"
                ],
                // ...
            },
            "value": 2,
            "trend": -3
        },
        // ...
    ]
}
```


### Calcul de plusieurs statistiques

#### URL

`POST /api/v2/stats`

#### Paramètres

Cette route est similaire à la précédente, mais permet de calculer plusieurs
statistiques en une requête. Elle accepte en requête un objet JSON contenant
les paramètres suivants :

 - `tstop` : un timestamp indiquant le fin de la période pour laquelle la
   statistique doit être calculées. Ce timestamp doit correspondre à une heure
   pile (e.g.  12:00, et non 12:03).
 - `duration` : la durée de la période, représentée par une chaîne
   `"<n><unité>"`, avec `<n>` un entier et `<unité>` une unité de temps (`h`,
   `d` ou `w`).
 - `mfilter` : un filtre mongodb, filtrant les entités pour lesquelles les
   statistiques doivent être calculées.
 - `stats`: un objet contenant les statistiques à calculer. Cet objet associe
   un titre de statistique (qui sera utilisé dans la réponse) à un objet
   définissant la statistique. Cet objet contient les champs suivants :
    - `stat`: la statistique à calculer (par exemple `alarms_created`).
    - `parameters`: un objet contenant les paramètres spécifiques à la
      statistique calculée. Ces paramètres sont précisés dans la documentation
      de chacune des statistiques.
    - `trend` (optionnel): `true` pour calculer la tendance par rapport à la
      période précédente.
    - `sla` (optionnel): un SLA, représenté par une inégalité (e.g.
      `">= 0.99"`).

#### Réponse

En cas de succès, la réponse est un objet JSON contenant un champ `values`. Ce
champ est un tableau contenant les valeurs des statistiques pour chaque
entité, sous la forme suivante :

```javascript
{
    'entity': {...},  // L'entité pour laquelle la statistique a été calculée
    'titre de la statistique 1': {
        'value': ...,  // La valeur de la statistique
        'trend': ...,  // La tendance
        'sla': ...  // true si la valeur est conforme au SLA
    },
    'titre de la statistique 2': {
        'value': ...,  // La valeur de la statistique
        'trend': ...,  // La tendance
        'sla': ...  // true si la valeur est conforme au SLA
    }
}
```

#### Exemple

La requête suivante renvoie le nombre d'alarmes critiques et majeures ouvertes
sur chaque ressource impactant l'entité `service`, les 18 et 19 août 2018.

```javascript
POST /api/v2/stats
{
    "mfilter": {
        "type": "resource",
        "impact": {
            "$in": ["service"]
        }
    },
    "tstop": 1534716000,  // 20 août à 00:00
    "duration": "2d",
    "stats": {
        "Alarms critiques": {
            "stat": "alarms_created",
            "parameters": {
                "states": [3]
            },
            "trend": true,
            "sla": "<= 20"
        },
        "Alarms majeures": {
            "stat": "alarms_created",
            "parameters": {
                "states": [2]
            },
            "trend": true
        }
    }
}
```

Le document JSON ci-dessous est un exemple de réponse à la requête précédente.


```javascript
{
    "values": [
        {
            "entity": {
                "_id": "resource1/component1",
                "type": "resource"
                "impact": [
                    "service"
                ],
                // ...
            },
            "Alarmes critiques": {
                "value": 117,
                "trend": 76,
                "sla": false
            },
            "Alarmes majeures": {
                "value": 37,
                "trend": 10
            }
        },
        {
            "entity": {
                "_id": "resource2/component1",
                "type": "resource"
                "impact": [
                    "service"
                ],
                // ...
            },
            "Alarmes critiques": {
                "value": 2,
                "trend": -1
            },
            "Alarmes majeures":
                "value": 3,
                "trend": -1
            }
        },
        // ...
    ]
}
```

## Statistiques

### Compteurs d'alarmes

Les compteurs d'alarmes permettent de compter des événements sur les alarmes :

 - `alarms_created` renvoie le nombre d'alarmes créées.
 - `alarms_resolved` renvoie le nombre d'alarmes résolues.
 - `alarms_canceled` renvoie le nombre d'alarmes annulées.

Les alarmes créées alors qu'un pbehavior est actif ne sont pas prises en compte.

#### Paramètres

Ces statistiques acceptent les paramètres suivants (à indiquer dans le champ
`parameters` d'une requête).

 - `recursive` (optionnel, `true` par défaut) : `true` pour prendre en compte
   les alarmes de l'entité et de ses dépendances, `false` pour ne prendre en
   compte que les alarmes de l'entité.
 - `states` (optionnel) : Un tableau contenant les états des alarmes à prendre
   en compte (par exemple `[3]` pour ne compter que les alarmes critiques).

#### Exemple

```javascript
POST /api/v2/stats/alarms_created
{
    "mfilter": {
        "type": "resource",
        "impact": {
            "$in": ["service"]
        }
    },
    "tstop": 1534716000,
    "duration": "2d",
    "parameters": {
        "states": [3]
    }
}
```

```javascript
{
    "values": [
        {
            "entity": {
                "_id": "resource1/component1",
                "type": "resource"
                "impact": [
                    "service"
                ],
                // ...
            },
            "value": 117
        },
        // ...
    ]
}
```


### Taux conforme à un SLA

 - `ack_time_sla` renvoie le taux d'acquittement conforme au SLA
 - `resolve_time_sla` renvoie le taux de résolution conforme au SLA

Les alarmes créées alors qu'un pbehavior est actif ne sont pas prises en compte.

#### Paramètres

Ces statistiques acceptent les paramètres suivants (à indiquer dans le champ
`parameters` d'une requête).

 - `recursive` (optionnel, `true` par défaut) : `true` pour prendre en compte
   les alarmes de l'entité et de ses dépendances, `false` pour ne prendre en
   compte que les alarmes de l'entité.
 - `states` (optionnel) : Un tableau contenant les états des alarmes à prendre
   en compte (par exemple `[3]` pour ne prendre en compte que les alarmes
   critiques).
 - `sla` : le SLA, sous la forme d'une inégalité, par exemple `"<= 3600"`.

#### Exemple

```javascript
POST /api/v2/stats/resolve_time_sla
{
	"mfilter": {
		"type": "resource",
		"impact": {
			"$in": ["feeder2_80"]
		}
	},
	"tstop": 1534716000,
	"duration": "2d",
	"parameters": {
		"sla": "<= 3600"
	}
}
```

```javascript
{
    "values": [
        {
            "entity": {
                "_id": "resource1/component1",
                "type": "resource"
                "impact": [
                    "service"
                ],
                // ...
            },
            "value": 0.97
        },
        // ...
    ]
}
```


### Temps passé dans des états

 - `time_in_state` renvoie le temps passé dans des états.
 - `state_rate` renvoie la proportion du temps passé dans des états.

Les périodes pendant lesquelles un pbehavior était actif ne sont pas
prises en compte.

#### Paramètres

 - `states` : Un tableau d'états. Par exemple `[2, 3]` pour calculer la
   proportion du temps passé dans un état majeur ou critique.

#### Exemple

```javascript
POST /api/v2/stats/state_rate
{
	"mfilter": {
		"type": "resource",
		"impact": {
			"$in": ["feeder2_80"]
		}
	},
	"tstop": 1534716000,
	"duration": "2d",
	"parameters": {
        "states": [0, 1]
	}
}
```

```javascript
{
    "values": [
        {
            "entity": {
                "_id": "resource1/component1",
                "type": "resource"
                "impact": [
                    "service"
                ],
                // ...
            },
            "value": 0.94
        },
        // ...
    ]
}
```
