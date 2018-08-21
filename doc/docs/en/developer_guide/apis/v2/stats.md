# Statistics API

## Routes

### Compute a single statistic

#### URL

`POST /api/v2/stats/<stat name>`

#### Parameters

This route takes a JSON object with the following fields:

 - `tstop`: a timestamp indicating the end of the last period for which the
   statistics will be computed. This timestamp should be at the top of an hour
   (e.g. 12:00, not 12:03).
 - `duration`: the duration of the period, represented by a string
   `"<n><unit>"`, where `<n>` is an integer and `<unit>` a time unit (`h`, `d`
   ou `w`).
 - `mfilter`: a mongodb filter, filtering the entities for which the
   statistics should be computed.
 - `parameters`: an object containing parameters for the computed statistic.
   See the documentation of each statistic below for the available parameters.

#### Response

If the request succeeded, the response is a JSON object containing a `values`
field. This field is a table containing the values of the statistic for each
entity, as follows:

```javascript
{
    'entity': {...},  // The entity for which the statistic was computed
    'value': ...  // The value of the statistic
}
```

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
            "value": 117
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
            "value": 2
        },
        // ...
    ]
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

#### Response

If the request succeeded, the response is a JSON object containing a `values`
field. This field is a table containing the values of the statistics for each
entity, as follows:

```javascript
{
    'entity': {...},  // The entity for which the statistic was computed
    'title of the 1st statistic': ...,  // The value of the statistic
    'title of the 2nd statistic': ...  // The value of the statistic
}
```

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
            }
        },
        "Major alarms": {
            "stat": "alarms_created",
            "parameters": {
                "states": [2]
            }
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
            "Critical alarms": 117,
            "Major alarms": 37
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
            "Critical alarms": 2,
            "Major alarms": 3
        },
        // ...
    ]
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

The periods during which an pbehavior was actie are not taken into account.

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
