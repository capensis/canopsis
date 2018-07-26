# Statistics API

## Specification

#### URL

This API defines two routes :

 - `POST /api/v2/stats/<stat name>`: to compute the value of the statistic
   `<stat_name>`.
 - `POST /api/v2/stats`: to compute the value of multiple statistics.

#### Parameters

These routes take a JSON object with the following fields:

 - `stats` (only for `/api/v2/stats`, required) : a list containing the names
   of the statistics to compute.
 - `tstart` (optional) : a timestamp indicating the start of the interval.
 - `tstop` (optional) : a timestamp indicating the end of the interval.
 - `group_by` (optional) : a list of tags.
 - `filter` (optional) : a list of *entity groups*. The entities that are part
   of at least one of the groups will be taken into account when computing the
   statistics.
 - `parameters` (optional) : an object containing parameters for the
   statistics. For the `/api/v2/stats` route, the parameters for a statistic
   should be in `parameters.<stat_name>`.

An *entity group* is a JSON object containing `"<tag name>" : <tag filter>`
couples. An entity is part of this group if each of its tags validates the
corresponding filter.

The tag name can be used to filter according to :

 - The entity's identity, with the names `entity_id` and `entity_type`.
 - One of the entity's informations, with the names
   `entity_infos.<information_id>`. Only the information ids specified in the
   [statsng engine configuration](../../../admin_guide/statsng.md#entity-tags)
   can be used in filters.
 - The alarm, with the names `connector`, `connector_name`, `component`,
   `resource` and `alarm_state`.

The tag filter can be :

 - a string, a tag validates this filter if its value is equal to this string.
 - a list of strings, a tag validates this filter if its value is in this list.
 - an object `{"matches": "<regex>"}`, with `<regex>` a [regular
   expression](https://golang.org/pkg/regexp/syntax/), a tag validates this
   filter if its value is matched by this regular expression.

```javascript
[ // Compute statistics for entities belonging to at least one of the following groups
    { // This group contains the entities whose tags validate the following conditions
        "<tag1>": "value",                   // tag1's value is "value" ET
        "<tag2>": ["value1", "value2", ...], // tag2's value is in [...] ET
        "<tag3>": {"matches": "value\d+"}    // tag3's value is matched by the regex
    },
    // ...
]
```

#### Response

The response is a JSON array containing objects with :

 - a `tags` field, whose value is an object containing the values of the tags
   specified in `group_by` (or an empty object if `group_by` is not defined).
 - one field for each statistic that was requested.


## Statistics

### Number of alarms created

The statistic `alarms_created` is equal to the number of alarms created
during a time interval.

### Number of alarms resolved

The statistic `alarms_resolved` is equal to the number of alarms resolved
during a time interval.

### Number of alarms canceled

The statistic `alarms_canceled` is equal to the number of alarms canceled
during a time interval.

### Mean ack time

The statistic `mean_ack_time` is equal to average time taken for an alarm to
be acknowledged.

### Mean resolve time

The statistic `mean_resolve_time` is equal to average time taken for an alarm
to be resolved.

### Ack time above or below the SLA

The statistic `ack_time_sla` is a JSON object with the following fields:

 - `above`: the number of alarms whose ack time is above the SLA.
 - `below`: the number of alarms whose ack time is below the SLA.
 - `above_rate`: the ratio of alarms whose ack time is above the SLA (between 0
   and 1).
 - `below_rate`: the ratio of alarms whose ack time is below the SLA (between 0
   and 1).

The value of the SLA should be provided in seconds in the parameter `sla`.

### Resolve time above or below the SLA

The statistic `resolve_time_sla` is a JSON object with the following fields:

 - `above`: the number of alarms whose resolve time is above the SLA.
 - `below`: the number of alarms whose resolve time is below the SLA.
 - `above_rate`: the ratio of alarms whose resolve time is above the SLA
   (between 0 and 1).
 - `below_rate`: the ratio of alarms whose resolve time is below the SLA
   (between 0 and 1).

The value of the SLA should be provided in seconds in the parameter `sla`.

### Time spent in each state

The statistic `time_in_state` is a JSON object with :

 - one field for each state (between 0 and 3), containing the time spent by the
   entity in this state, in seconds
 - a `total` field, containing the total time

The intervals during which a pbehavior is active are excluded from these
values. The total time may thus be inferior to the duration of the interval
`tstop - tstart`.

### Availability

The statistic `availability` is a JSON object with the following fields:

 - `available_time`: the time during which the entity was in an available
   state, in seconds
 - `unavailable_time`: the time during which the entity was in an unavailable
   state, in seconds
 - `available_rate`: the ratio of time during which the entity was in an
   available state (between 0 and 1)
 - `unavailable_rate`: the ratio of time during which the entity was in an
   unavailable state (between 0 and 1)

The entity is considered to be available if it is in a state lower or equal to
the value of the parameter `available_state`.

The intervals during which a pbehavior is active are excluded from these
values. The total time `available_time + unavailable_time` may thus be inferior
to the duration of the interval `tstop - tstart`.

### Maintenance

The statistic `maintenance` is a JSON object with the following fields:

 - `maintenance`: the time during which the entity had an active pbehavior, in
   seconds.
 - `no_maintenance`: the time during which the entity had no active pbehavior,
   in seconds

### Mean Time Between Failures

The statistic `mtbf` is the mean time between failures, i.e. the time without
maintenance divided by the number of failures.

### Alarm List

The statistic `alarm_list` returns a list of alarms. It is a JSON array of
objects which contains the tags of the entity that created the alarm
(`entity_id`, `entity_type`, `entity_infos.<information_id>`, `connector`,
`connector_name`, `component`, `resource` and `alarm_state`), as well as the
following fields:

 - `time`: the date of creation of the alarm
 - `pbehavior`: `"True"` if there was an active pbehavior when the alarm was
   created, `"False"` otherwise.
 - `value`: the time it took for the alarm to be resolved.

### Entities impacted by the most alarms

The statistic `most_alarms_impacting` is a list containing the groups of
entities that are impacted by the largest number of alarms. The request takes
the following parameters:

 - `group_by`: the tags used to group the entities.
 - `filter` (optional): an *entity filter* used to filter the entities. This
   parameters has the same format as the main `filter` parameter.
 - `limit` (optional): the maximum number of groups to return.

It returns a list of objects ordered by descending number of alarms, with the
following fields:

 - `tags`: the tags of the entity group.
 - `value`: the number of alarms impacting this entity group.

### Entities with the worst Mean Time Between Failures

The statistic `worst_mtbf` is a list containing the groups of entities that
have the worst mtbf. The request takes the following parameters:

 - `group_by`: the tags used to group the entities.
 - `filter` (optional): an *entity filter* used to filter the entities. This
   parameters has the same format as the main `filter` parameter.
 - `limit` (optional): the maximum number of groups to return.

It returns a list of objects ordered by ascending mtbf, with the following
fields:

 - `tags`: the tags of the entity group.
 - `value`: the mtbf.

### Longest alarms

The statistic `longest_alarms` returns a list of alarms that took the longest
time to resolve. The request takes the following parameters:

 - `limit` (optional): the maximum number of groups to return.

It returns a JSON array of objects which contains the tags of the entity that
created the alarm (`entity_id`, `entity_type`, `entity_infos.<information_id>`,
`connector`, `connector_name`, `component`, `resource` and `alarm_state`), as
well as the following fields:

 - `time`: the date of creation of the alarm
 - `pbehavior`: `"True"` if there was an active pbehavior when the alarm was
   created, `"False"` otherwise.
 - `value`: the time it took for the alarm to be resolved.


## Examples

### Number of alarms created by a component

`/api/v2/stats/alarms_created`

Request:

```javascript
{
    "tstart": 1528290000,
    "tstop": 1528293000,
    "filter": [
        {
            "connector": "connector",
            "connector_name": "connector_name",
            "component": "component"
        }
    ]
}
```

Response:

```javascript
[
    {
        "tags": {},
        "alarms_created": 13
    }
]
```

### Number of alarms created by each resource of a component

`/api/v2/stats/alarms_resolved`

Request:

```javascript
{
    "tstart": 1528290000,
    "tstop": 1528293000,
    "filter": [
        {
            "connector": "connector",
            "connector_name": "connector_name",
            "component": "component"
        }
    ],
    "group_by": ["resource"]
}
```

Response:

```javascript
[
    {
        "tags": {"resource": "resource1"},
        "alarms_resolved": 4
    },
    {
        "tags": {"resource": "resource2"},
        "alarms_resolved": 3
    },
    {
        "tags": {"resource": "resource3"},
        "alarms_resolved": 1
    }
]
```

### Ack time above or below the SLA

`/api/v2/stats/ack_time_sla`

Request:

```javascript
{
    "tstart": 1528290000,
    "tstop": 1528293000,
    "filter": [
        {
            "connector": "connector",
            "connector_name": "connector_name",
            "component": "component"
        }
    ],
    "parameters": {
        "sla": 600
    }
}
```

Response:

```javascript
[
    {
        "tags": {},
        "ack_time_sla": {
            "above": 3
            "below": 9,
            "above_rate": 0.25,
            "below_rate": 0.75,
        }
    }
]
```


### Number of critical alarms created by a component

`/api/v2/stats/alarms_created`

Request:

```javascript
{
    "tstart": 1528290000,
    "tstop": 1528293000,
    "filter": [
        {
            "connector": "connector",
            "connector_name": "connector_name",
            "component": "component",
            "alarm_state": 3
        }
    ]
}
```

Response:

```javascript
[
    {
        "tags": {},
        "alarms_created": 13
    }
]
```

### Time spent by an entity in each state

`/api/v2/stats/time_in_state`

Request:

```javascript
{
    "tstart": 1528290000,
    "tstop": 1528293000,
    "filter": [
        {
            "entity_id": "watcher_0"
        }
    ]
}
```

Response:

```javascript
[
    {
        "tags": {},
        "time_in_state": {
			"total": 2454,
			"0": 1707,
			"1": 105,
			"2": 23,
			"3": 619
		}
    }
]
```

### Time during which an entity was available

`/api/v2/stats/availability`

Request:

```javascript
{
    "tstart": 1528290000,
    "tstop": 1528293000,
    "filter": [
        {
            "entity_id": "watcher_0"
        }
    ],
    "parameters": {
        "available_state": 2
    }
}
```

Response:

```javascript
[
    {
        "tags": {},
        "availability": {
			"available_time": 1835,
			"unavailable_time": 619,
			"available_rate": 0.747758761206194,
			"unavailable_rate": 0.25224123879380606
		}
    }
]
```

### For each component, the 10 resources with the worst MTBF

`/api/v2/stats/worst_mtbf`

Request:

```javascript
{
	"group_by": ["component"],
	"parameters": {
		"group_by": ["resource"],
		"limit": 10
	}
}
```

Response:

```javascript
[
    {
        "tags": {
            "component": "component_0"
        },
        "worst_mtbf": [
            {
                "tags": {
                    "resource": "resource_0"
                },
                "value": 57.333333333333336
            },
            {
                "tags": {
                    "resource": "resource_1"
                },
                "value": 106
            },
            // ...
        ]
    },
    {
        "tags": {
            "component": "component_1"
        },
        "worst_mtbf": [
            {
                "tags": {
                    "resource": "resource_0"
                },
                "value": 57.333333333333336
            },
            {
                "tags": {
                    "resource": "resource_1"
                },
                "value": 57.333333333333336
            },
            // ...
        ]
    },
    // ...
]
```

### 10 longest alarms of each component

`/api/v2/stats/longest_alarms`

Request:

```javascript
{
    "group_by": ["component"],
    "parameters": {
        "limit": 10
    }
}
```

Response:

```javascript
[
    {
        "tags": {
            "component": "component_0"
        }
        "longest_alarms": [
            {
                "alarm_state": "3",
                "connector": "...",
                "connector_name": "...",
                "entity_id": "resource_0/component_0",
                "entity_type": "resource",
                "pbehavior": "False",
                "resource": "resource_0",
                "time": 1531833020,
                "value": 3754
            },
            {
                "alarm_state": "3",
                "connector": "...",
                "connector_name": "...",
                "entity_id": "resource_1/component_0",
                "entity_type": "resource",
                "pbehavior": "False",
                "resource": "resource_1",
                "time": 1531830121,
                "value": 3562
            },
            //...
        ]
    },
    // ...
]
```

### Multiple statistics in one request

`/api/v2/stats`

Request:

```javascript
{
    "stats": ["alarms_created", "alarms_resolved", "ack_time_sla", "resolve_time_sla"],
    "tstart": 1528290000,
    "tstop": 1528293000,
    "filter": [
        {
            "connector": "connector",
            "connector_name": "connector_name",
            "component": "component"
        }
    ],
    "parameters": {
        "ack_time_sla": {
            "sla": 900
        },
        "resolve_time_sla": {
            "sla": 3600
        }
    }
}
```

Response:

```javascript
[
    {
        "tags": {},
        "alarms_created": 12,
        "alarms_resolved": 8,
        "ack_time_sla": {
            "above": 4,
            "below": 8,
            "above_rate": 0.3333333333333333,
            "below_rate": 0.6666666666666666
        },
        "ack_time_sla": {
            "above": 2
            "below": 6,
            "above_rate": 0.25,
            "below_rate": 0.75
        }
    }
]
```
