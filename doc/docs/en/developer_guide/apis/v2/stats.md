# Statistics API

## Routes

### Compute a single statistic

#### URL

`POST /api/v2/stats/<stat name>`

#### Parameters

This route takes a JSON object with the following fields:

 - `tstop`: a timestamp indicating the end of the period for which the
   statistic will be computed. This timestamp should be at the top of an hour
   (e.g. 12:00, not 12:03).
 - `duration`: the duration of the period, represented by a string
   `"<n><unit>"`, where `<n>` is an integer and `<unit>` a time unit (`h`, `d`,
   `w` or `m`).
 - `mfilter`: a mongodb filter, filtering the entities for which the
   statistic should be computed.
 - `parameters`: an object containing parameters for the computed statistic.
   See the documentation of each statistic below for the available parameters.
 - `trend` (optional): `true` to compute the trend with the previous period.
 - `sla` (optional): a SLA, represented by an inequality (e.g. `">= 0.99"`).
 - `aggregate` (optional): an array containing the names of aggregations
   functions used to aggregate the values of the statistic (`"sum"` is the only
   available function for now).
 - `sort_order` (optional): `"desc"` to sort the results by descending value,
   `"asc"` to sort them by ascending value. The results are not sorted by
   default.
 - `limit` (optional): the maximum number of values to return. All values are
   returned by default.


#### Response

If the request succeeded, the response is a JSON object containing :

 - a `values` field. This field is a table containing the values of the
   statistic for each entity, as follows:

   ```javascript
   {
       'entity': {...},  // The entity for which the statistic was computed
       'value': ...,  // The value of the statistic
       'trend': ...,  // The trend
       'sla': ...  // true if the value is conform to the SLA
   }
   ```
 - an `aggregations` field. This field is an object containing the values of
   the aggregations, using the functions given in the `aggregate` parameter.

#### Example

The following request returns the number of critical alarms opened on each
resource impacting the entity `service`, on the 18th and 19th of August 2018.

```javascript
POST /api/v2/stats/alarms_created
{
    "mfilter": {
        "type": "resource",
        "impact": {
            "$in": ["service"]
        }
    },
    "tstop": 1534716000,  // August 20th at 00:00
    "duration": "2d",
    "parameters": {
        "states": [3]
    },
    "trend": true,
    "aggregate": ["sum"]
}
```

The JSON document below is an example of a response to the previous request.

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

### Compute multiple statistics in one request

#### URL

`POST /api/v2/stats`

#### Parameters

This route is similar to the previous one, but allows to compute multiple statistics in a single request. It takes a JSON object with the following fields:

 - `tstop`: a timestamp indicating the end of the last period for which the
   statistics will be computed. This timestamp should be at the top of an hour
   (e.g. 12:00, not 12:03).
 - `duration`: the duration of the period, represented by a string
   `"<n><unit>"`, where `<n>` is an integer and `<unit>` a time unit (`h`, `d`
   ou `w`).
 - `mfilter`: a mongodb filter, filtering the entities for which the
   statistics should be computed.
 - `stats`: an object containing the statistics to compute. This objects maps a
   title (which will be used in the response) to an object defining the
   statistic, which has the following fields:
    - `stat`: the statistic (for example `alarms_crated`).
    - `parameters`: an object containing parameters for the computed statistic.
      See the documentation of each statistic below for the available
      parameters.
    - `trend` (optional): `true` to compute the trend with the previous period.
    - `sla` (optional): a SLA, represented by an inequality (e.g. `">= 0.99"`).
    - `aggregate` (optional): an array containing the names of aggregations
      functions used to aggregate the values of the statistic (`"sum"` is the
      only available function for now).
 - `sort_column` (optional): the title of the statistic whose values will be
   used to sort the results.
 - `sort_order` (optional): `"desc"` to sort the results by descending value,
   `"asc"` to sort them by ascending value. The results are not sorted by
   default.
 - `limit` (optional): the maximum number of values to return. All values are
   returned by default.

#### Response

If the request succeeded, the response is a JSON object containing:

 - a `values` field. This field is a table containing the values of the
   statistics for each entity, as follows:

   ```javascript
   {
       'entity': {...},  // The entity for which the statistic was computed
       'title of the 1st statistic': {
           'value': ...,  // The value of the statistic
           'trend': ...,  // The trend
           'sla': ...  // true if the value is conform to the SLA
       },
       'title of the 2nd statistic': {
           'value': ...,  // The value of the statistic
           'trend': ...,  // The trend
           'sla': ...  // true if the value is conform to the SLA
       }
   }
   ```
 - an `aggregation` field. This field is an object containing the values of the
   aggregations for each statistic.

#### Example

The following request returns the number of critical and major alarms opened on
each resource impacting the entity `service`, on the 18th and 19th of August
2018.

```javascript
POST /api/v2/stats
{
    "mfilter": {
        "type": "resource",
        "impact": {
            "$in": ["service"]
        }
    },
    "tstop": 1534716000,
    "duration": "2d",
    "stats": {
        "Critical alarms": {
            "stat": "alarms_created",
            "parameters": {
                "states": [3]
            },
            "trend": true,
            "sla": "<= 20",
            "aggregate": ["sum"]
        },
        "Major alarms": {
            "stat": "alarms_created",
            "parameters": {
                "states": [2]
            },
            "trend": true
        }
    },
    "sort_column": "Critical alarms",
    "sort_order": "desc"
}
```

The JSON document below is an example of a response to the previous request.


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
            "Critical alarms": {
                "value": 117,
                "trend": 76,
                "sla": false
            },
            "Major alarms": {
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
            "Critical alarms": {
                "value": 2,
                "trend": -1,
                "sla": true
            },
            "Major alarms": {
                "value": 3,
                "trend": -1
            }
        },
        // ...
    ],
    "aggregations": {
        "Critical alarms": {
            "sum": 253
        }
    }
}
```

### Compute statistics on multiple periods

#### URL

`POST /api/v2/stats/evolution`

#### Parameters

This route is similar to the previous one, but allows to compute statistics on
multiple periods. It takes a JSON object with the following fields:

 - `tstop`: a timestamp indicating the end of the last period for which the
   statistics will be computed. This timestamp should be at the top of an hour
   (e.g. 12:00, not 12:03).
 - `duration`: the duration of the period, represented by a string
   `"<n><unit>"`, where `<n>` is an integer and `<unit>` a time unit (`h`, `d`
   ou `w`).
 - `mfilter`: a mongodb filter, filtering the entities for which the
   statistics should be computed.
 - `stats`: an object containing the statistics to compute. This objects maps a
   title (which will be used in the response) to an object defining the
   statistic, which has the following fields:
    - `stat`: the statistic (for example `alarms_crated`).
    - `parameters`: an object containing parameters for the computed statistic.
      See the documentation of each statistic below for the available
      parameters.
    - `trend` (optional): `true` to compute the trend with the previous period.
    - `sla` (optional): a SLA, represented by an inequality (e.g. `">= 0.99"`).
    - `aggregate` (optional): an array containing the names of aggregations
      functions used to aggregate the values of the statistic (`"sum"` is the
      only available function for now).
 - `periods`: the number of periods.

#### Response

If the request succeeded, the response is a JSON object containing:

 - a `values` field. This field is a table containing the values of the
   statistics for each entity, as follows:

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
 - an `aggregations` field. This field is an object containing the values of
   the aggregations, using the functions given in the `aggregate` parameter.

#### Example

The following request returns the number of critical and major alarms opened on
each resource impacting the entity `service`, on the 18th and 19th of August
2018.

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
        "Critical alarms": {
            "stat": "alarms_created",
            "parameters": {
                "states": [3]
            },
            "trend": true,
            "sla": "<= 20",
            "aggregate": ["sum"]
        },
        "Major alarms": {
            "stat": "alarms_created",
            "parameters": {
                "states": [2]
            },
            "trend": true
        }
    }
}
```

The JSON document below is an example of a response to the previous request.


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
            "Critical alarms": [
                {
                    "start": 1534543200,  // August 18 at 0:00
                    "end": 1534629600,
                    "value": 19,
                    "trend": 12,
                    "sla": false
                },
                {
                    "start": 1534629600,  // August 19 at 00:00
                    "end": 1534716000,
                    "value": 98,
                    "trend": 79,
                    "sla": false
                }
            ],
            "Major alarms": [
                {
                    "start": 1534543200,  // August 18 at 00:00
                    "end": 1534629600,
                    "value": 11,
                    "trend": -1
                },
                {
                    "start": 1534629600,  // August 19 at 00:00
                    "end": 1534716000,
                    "value": 26,
                    "trend": 15
                }
            ]
        },
        // ...
    ],
    "aggregations": {
        "Critical alarms": [
            {
                "start": 1534543200,  // August 18 at 00:00
                "end": 1534629600,
                "value": 37
            },
            {
                "start": 1534629600,  // August 19 at 00:00
                "end": 1534716000,
                "value": 136
            }
        ]
    }
}
```


## Statistics

### Alarm counters

The alarms counters count events on the alarms :

 - `alarms_created` returns the number of alarms created.
 - `alarms_resolved` returns the number of alarms resolved.
 - `alarms_canceled` returns the number of alarms canceled.

The alarms created while a pbehavior was active are not taken into account.

#### Parameters

These statistics take the following parameters (in the `parameters` field) :

 - `recursive` (optional, `true` by default): `true` to get the value of the
   counter on the entity and its dependencies, `false` to get the value only
   for the entity itself.
 - `states` (optional): Only the alarms whose state at the creation date is in
   this list will be taken into account.
 - `authors` (optional): Only the events whose author is in this list will be
   taken into account.

#### Example

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

### Proportion conform to a SLA

 - `ack_time_sla` returns the proportion of ack times that are conform to a
   SLA.
 - `resolve_time_sla` returns the proportion of resolve times that are conform
   to a SLA.

The alarms created while a pbehavior was active are not taken into account.

#### Parameters

These statistics take the following parameters (in the `parameters` field):

 - `recursive` (optional, `true` by default): `true` to get the value of the
   counter on the entity and its dependencies, `false` to get the value only
   for the entity itself.
 - `states` (optional): Only the alarms whose state at the creation date is in
   this list will be taken into account.
 - `authors` (optional): Only the events whose author is in this list will be
   taken into account.
 - `sla`: the SLA as an inequality, for example `"<= 3600"`.

#### Example

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

### Time spent in state

 - `time_in_state` returns the time spent in some states.
 - `state_rate` returns the proportion of time spent in some states.

The periods during which a pbehavior was active are not taken into account.

#### Parameters

 - `states`: An array of stats. For example `[2, 3]` to compute the proportion
   of the time spent in a major or critical state.

#### Example

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


### Mean Time Between Failures

The `mtbf` statistic returns the Mean Time Between Failures, i.e. the available
time divided by the number of alarms.

The periods during which a pbehavior was active are not taken into account.

#### Paramètres

This statistic does not take any parameters.

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
`


### Current state

The `current_state` statistic returns the current state of an entity (at the
time the request was made). This statistic does not take into account the
`tstop` and `duration` parameters.

#### Paramètres

This statistic does not take any parameters.

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


### Ongoing Alarms

The following statistics compute the numver of ongoing alarms.

 - `ongoing_alarms` returns the number of ongoing alarms during a period.
 - `current_ongoing_alarms` returns the number of ongoing alarms at the time of
   the request. This statistic does not take into account the `tstop` and
   `duration` parameters.

The alarms created while a pbehavior was active are not taken into account.

#### Parameters

These statistics take the following parameters (in the `parameters` field) :

 - `states` (optional): Only the alarms whose state at the creation date is in
   this list will be taken into account.

#### Example

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
