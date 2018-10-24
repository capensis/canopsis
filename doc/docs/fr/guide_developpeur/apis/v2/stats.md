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
   `d`, `w` ou `m`).
 - `mfilter` : un filtre mongodb, filtrant les entités pour lesquelles les
   statistiques doivent être calculées.
 - `parameters` : un objet contenant les paramètres spécifiques à la
   statistique calculée. Ces paramètres sont précisés dans la documentation de
   chacune des statistiques.
 - `trend` (optionnel) : `true` pour calculer la tendance par rapport à la
   période précédente.
 - `sla` (optionnel) : un SLA, représenté par une inégalité (e.g. `">= 0.99"`).
 - `aggregate` (optionnel) : un tableau contenant les noms de fonctions
   d'aggrégation à appliquer aux valeurs de la statistique (`"sum"` est la
   seule fonction disponible pour le moment).
 - `sort_order` (optionnel) : `"desc"` pour trier les résultats par valeur
   décroissante, `"asc"` pour les trier par valeur croissante. Les résultats ne
   sont pas triés par défaut.
 - `limit` (optionnel) : le nombre maximal de valeurs à renvoyer. Toutes les
   valeurs sont renvoyées par défaut.


#### Réponse

En cas de succès, la réponse est un objet JSON contenant :

 - un champ `values`. Ce champ est un tableau contenant les valeurs de la
   statistique pour chaque entité, sous la forme suivante :

   ```javascript
   {
       'entity': {...},  // L'entité pour laquelle la statistique a été calculée
       'value': ...,  // La valeur de la statistique
       'trend': ...,  // La tendance
       'sla': ...  // true si la valeur est conforme au SLA
   }
   ```
 - un champs `aggregations`. Ce champ est un objet contenant la valeurs des
   aggregations, obtenues en utilisant les fonctions données en paramètre
   `aggregate`.

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
    },
    "trend": true,
    "aggregate": ["sum"]
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
    ],
    "aggregations": {
        "sum": 253  // 117 + 2 + ...
    }
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
    - `aggregate` (optionnel) : un tableau contenant les noms de fonctions
      d'aggrégation à appliquer aux valeurs de la statistique (`"sum"` est la
      seule fonction disponible pour le moment).
 - `sort_column` (optionnel) : le titre de la statistique dont la valeur sera
   utilisée pour trier les résultats.
 - `sort_order` (optionnel) : `"desc"` pour trier les résultats par valeur
   décroissante, `"asc"` pour les trier par valeur croissante. Les résultats ne
   sont pas triés par défaut.
 - `limit` (optionnel) : le nombre maximal de valeurs à renvoyer. Toutes les
   valeurs sont renvoyées par défaut.


#### Réponse

En cas de succès, la réponse est un objet JSON contenant :

 - un champ `values`. Ce champ est un tableau contenant les valeurs des
   statistiques pour chaque entité, sous la forme suivante :

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
 - un champ `aggregations`. Ce champs est un objet contenant les valeurs des
   aggrégations pour chaque statistiques.

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
            "sla": "<= 20",
            "aggregation": ["sum"]
        },
        "Alarms majeures": {
            "stat": "alarms_created",
            "parameters": {
                "states": [2]
            },
            "trend": true
        }
    },
    "sort_column": "Alarmes critiques",
    "sort_order": "desc"
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
                "trend": -1,
                "sla": true
            },
            "Alarmes majeures":
                "value": 3,
                "trend": -1
            }
        },
        // ...
    ],
    "aggregations": {
        "Alarmes critiques": {
            "sum": 253
        }
    }
}
```

### Calcul de statistiques sur plusieurs périodes

#### URL

`POST /api/v2/stats/evolution`

#### Paramètres

Cette route est similaire à la précédente, mais permet de calculer les
statistiques sur plusieurs périodes. Elle accepte en requête un objet JSON
contenant les paramètres suivants :

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
    - `aggregate` (optionnel) : un tableau contenant les noms de fonctions
      d'aggrégation à appliquer aux valeurs de la statistique (`"sum"` est la
      seule fonction disponible pour le moment).
 - `periods`: le nombre de périodes pour lesquelles les statistiques doivent
   être calculéees.

Les paramètres `sort_column`, `sort_order` et `limit` ne sont pas disponibles
pour cette route.

#### Réponse

En cas de succès, la réponse est un objet JSON contenant ;

 - un champ `values`. Ce champ est un tableau contenant les valeurs des
   statistiques pour chaque entité, sous la forme suivante :

   ```javascript
   {
       'entity': {...},  // L'entité pour laquelle la statistique a été calculée
       'titre de la statistique 1': [
           {
               'start': ...,  // Timestamp du début de la période
               'end': ...,  // Timestamp du fin de la période
               'value': ...,  // La valeur de la statistique
               'trend': ...,  // La tendance
               'sla': ...  // true si la valeur est conforme au SLA
           },
           {
               'start': ...,  // Timestamp du début de la période
               'end': ...,  // Timestamp du fin de la période
               'value': ...,  // La valeur de la statistique
               'trend': ...,  // La tendance
               'sla': ...  // true si la valeur est conforme au SLA
           },
           // ...
       ],
       'titre de la statistique 2': [
           // ...
       ]
   }
   ```
 - un champ `aggregations`. Ce champs est un objet contenant les valeurs des
   aggrégations pour chaque statistiques, sur chaque période.

#### Exemple

La requête suivante renvoie le nombre d'alarmes critiques et majeures ouvertes
sur chaque ressource impactant l'entité `service`, les 18 et 19 août 2018.

```javascript
POST /api/v2/stats/evolution
{
    "mfilter": {
        "type": "resource",
        "impact": {
            "$in": ["service"]
        }
    },
    "tstop": 1534716000,  // 20 août à 00:00
    "duration": "1d",
    "periods": 2,
    "stats": {
        "Alarms critiques": {
            "stat": "alarms_created",
            "parameters": {
                "states": [3]
            },
            "trend": true,
            "sla": "<= 20",
            "aggregation": ["sum"]
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
            "Alarmes critiques": [
                {
                    "start": 1534543200,  // 18 août à 00:00
                    "end": 1534629600,
                    "value": 19,
                    "trend": 12,
                    "sla": false
                },
                {
                    "start": 1534629600,  // 19 août à 00:00
                    "end": 1534716000,
                    "value": 98,
                    "trend": 79,
                    "sla": false
                }
            ],
            "Alarmes majeures": [
                {
                    "start": 1534543200,  // 18 août à 00:00
                    "end": 1534629600,
                    "value": 11,
                    "trend": -1
                },
                {
                    "start": 1534629600,  // 19 août à 00:00
                    "end": 1534716000,
                    "value": 26,
                    "trend": 15
                }
            ]
        },
        // ...
    ],
    "aggregations": {
        "Alarmes critiques": [
            {
                "start": 1534543200,  // 18 août à 00:00
                "end": 1534629600,
                "value": 37
            },
            {
                "start": 1534629600,  // 19 août à 00:00
                "end": 1534716000,
                "value": 136
            }
        ]
    }
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
 - `authors` (optionnel) : Un tableau contenant les auteurs des événements à
   prendre en compte.

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
 - `authors` (optionnel) : Un tableau contenant les auteurs des événements à
   prendre en compte.
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


### Indice de fiabilité

La statistique `mtbf` (Mean Time Between Failures) renvoie l'indice de
fiabilité, ou temps moyen entre panne, c'est-à-dire le temps de disponibilité
divisé par le nombre d'indisponibilités.

Les périodes pendant lesquelles un pbehavior était actif ne sont pas
prises en compte.

#### Paramètres

Cette statistique ne prend pas de paramètres.

#### Exemple

```javascript
POST /api/v2/stats/mtbf
{
	"mfilter": {
		"type": "resource",
		"impact": {
			"$in": ["feeder2_80"]
		}
	},
	"tstop": 1534716000,
	"duration": "2d"
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
            "value": 406.56
        },
        // ...
    ]
}
```


### État courant

La statistique `current_state` renvoie l'état courant des entités (à l'instant
où la requête est effectuée). Cette statistique ne prends pas en compte les
paramètres `tstop` et `duration`.

#### Paramètres

Cette statistique ne prend pas de paramètres.

#### Exemple

```javascript
POST /api/v2/stats/current_state
{
	"mfilter": {
		"type": "resource",
		"impact": {
			"$in": ["feeder2_80"]
		}
	},
	"tstop": 1534716000,
	"duration": "2d"
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
            "value": 3
        },
        // ...
    ]
}
```


### Alarmes en cours

Les statistiques suivantes permettent de calculer le nombre d'alarmes en cours.

 - `ongoing_alarms` renvoie le nombre d'alarmes en cours pendant une période.
 - `current_ongoing_alarms` renvoie le nombre d'alarmes en cours lorsque la
   requête est effectuée. Cette statistique ne prend pas en compte les
   paramètres `tstop` et `duration`.

Les alarmes créées alors qu'un pbehavior est actif ne sont pas prises en compte.

#### Paramètres

Ces statistiques acceptent les paramètres suivants (à indiquer dans le champ
`parameters` d'une requête).

 - `states` (optionnel) : Un tableau contenant les états des alarmes à prendre
   en compte (par exemple `[3]` pour ne compter que les alarmes critiques).

#### Exemple

```javascript
POST /api/v2/stats/ongoing_alarms
{
	"mfilter": {
		"type": "resource",
		"impact": {
			"$in": ["feeder2_80"]
		}
	},
	"tstop": 1534716000,
	"duration": "2d"
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
            "value": 3
        },
        // ...
    ]
}
```
