# Statistics API

## Routes

### Compute a single statistic

#### URL

`POST /api/v2/stats/<stat name>`

#### Parameters

This route takes a JSON object with the following fields:

 - `tstart` (optional): a timestamp indicating the start of the last period
   for which the statistics will be computed.
 - `tstop` (optional): a timestamp indicating the end of the last period for
   which the statistics will be computed.
 - `periods` (optional): the number of periods for which the statistics will
   be computed.
 - `group_by` (optional): a list of tags used to group the results. The
   available tags or the same as the ones used in the *alarm groups*, and are
   defined below.
 - `filter` (optional): a list of *alarm groups*. The alarms that are part
   of at least one of the groups will be taken into account when computing the
   statistics.
 - `parameters` (optional): an object containing parameters for the computed
   statistic. See the documentation of each statistic below for the available
   parameters.

An *alarm group* is a JSON object containing `"<tag name>": <tag filter>`
couples. An alarm is part of this group if each of its tags validates the
corresponding filter.

The tag name can be used to filter according to:

 - The identity of the entity that created the alarm, with the names
   `entity_id` and `entity_type`.
 - One of the entity's informations, with the names
   `entity_infos.<information_id>`. Only the information ids specified in the
   [statsng engine configuration](../../../admin_guide/statsng.md#entity-tags)
   can be used in filters.
 - The alarm, with the names `connector`, `connector_name`, `component`,
   `resource` and `alarm_state`.

The tag filter can be:

 - a string, a tag validates this filter if its value is equal to this string.
 - a list of strings, a tag validates this filter if its value is in this list.
 - an object `{"matches": "<regex>"}`, with `<regex>` a [regular
   expression](https://golang.org/pkg/regexp/syntax/), a tag validates this
   filter if its value is matched by this regular expression.

```javascript
[ // Compute statistics for alarms belonging to at least one of the following groups
    { // This group contains the alarms whose tags validate the following conditions
        "<tag1>": "value",                   // tag1's value is "value" ET
        "<tag2>": ["value1", "value2", ...], // tag2's value is in [...] ET
        "<tag3>": {"matches": "value\d+"}    // tag3's value is matched by the regex
    },
    // ...
]
```

#### Response

If the request succeeded, the response is a JSON array containing the groups
obtained by grouping with the tags defined in `group_by`.

Each group is an object containing the following fields:

 - `tags`: an object containing the balues of the tags defined in `group_by`.
 - `periods`: an array containing the statistics for each period. The period
   or ordered chronologically.

Each period is an object containing the following fields:

 - `tstart`: a timestamp indicating the start of the period.
 - `tstop`: a timestamp indicating the end of the period.
 - `<stat name>`: the value of the statistic. The type of this value depends
   on the computed statistic.

#### Example

The following request returns the proportion of alarms per day whose resolve
time is lower than a SLA for each resource of the component `c`, between the
23rd and the 29th of July 2018.

```javascript
POST /api/v2/stats/resolve_time_sla
{
    "filter": [{
        "component": "c",
        "entity_type": "resource"
    }],
    "group_by": ["resource"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "periods": 7,
    "parameters": {
        "sla": 3600
    }
}
```

The alarms whose component is `c` and created by an entity of type `resource`
are taken into account.

The statistic is computed for seven consecutive periods of same duration. The
last of these periods starts at `tstart` (July 29th at 00:00) and ends at
`tstop` (July 30th at 00:00).

The value of the SLA is defined in the `parameters` field, since it is a
parameter that is spectific to the `resolve_time_sla` statistic.

The JSON document below is an example of a response to the previous request.
The value of the statistic is a dictionnary containing multiple values. See the
documentation of the `resolve_time_sla` statistic for more details.

```javascript
[ // Array of groups
    {
        "tags": { // Tags of the group
            "resource": "resource1"
        },
        "periods": [ // Array of periods
            {
                "tstart": 1532296800, // July 23rd at 00:00
                "tstop": 1532383200, // July 24th at 00:00
                "resolve_time_sla": { // Value of the statistic
                    "above": 10,
                    "below": 90,
                    "above_rate": 0.1,
                    "below_rate": 0.9
                }
            },
            {
                "tstart": 1532383200, // July 24th at 00:00
                "tstop": 1532469600, // July 25th at 00:00
                "resolve_time_sla": {
                    "above": 0,
                    "below": 47,
                    "above_rate": 0,
                    "below_rate": 1,
                }
            },
            // ...
            {
                "tstart": 1532815200, // July 29th at 00:00
                "tstop": 1532901600, // July 30th at 00:00
                "resolve_time_sla": {
                    "above": 4,
                    "below": 28,
                    "above_rate": 0.125,
                    "below_rate": 0.875,
                }
            }
        ]
    },
    {
        "tags": {
            "resource": "resource2"
        },
        "periods": [
            // ...
        ]
    },
    // ...
]
```


### Compute multiple statistics in one request

#### URL

`POST /api/v2/stats`

#### Parameters

This route takes a JSON object with the same parameters as the previous route,
with two exceptions:

 - A new `stats` field (required) containing a list of the statistics to
   compute.
 - The `parameters`field containing an object associating to each statistic its
   parameters.

#### Response

The response has the same format as the previous route. Each period contains
the value of multiple statistics.

#### Example

The following request returns the ratio of alarms per day whose resolve time is
lower than a SLA and the alarms created per day for each resource of the
component `c`, between the 23rd and the 29th of July 2018.

```javascript
POST /api/v2/stats
{
    "stats": ["resolve_time_sla", "alarms_created"],
    "filter": [{
        "component": "c",
        "entity_type": "resource"
    }],
    "group_by": ["resource"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "periods": 7,
    "parameters": {
        "resolve_time_sla": {
            "sla": 3600
        }
    }
}
```

The body of this request is the same as the previous example, with two
exceptions:

 - The list of the statistics to compute has been added to the `stats` field.
 - The `sla` parameter which was in the `parameters` field has been moved to
   `parameters.resolve_time_sla`. The `alarms_created` statistic does not take
   parameters. If it did, they would have to be defined in
   `parameters.alarms_created`.

The JSON document below is an example of a response to the previous request.

```javascript
[ // Array of groups
    {
        "tags": { // Tags of the group
            "resource": "resource1"
        },
        "periods": [ // Array of periods
            {
                "tstart": 1532296800, // July 23rd at 00:00
                "tstop": 1532383200, // July 24th at 00:00
                // Values of the two statistics
                "alarms_created": 100,
                "resolve_time_sla": {
                    "above": 10,
                    "below": 90,
                    "above_rate": 0.1,
                    "below_rate": 0.9
                }
            },
            // ...
        ]
    },
    // ...
]
```


## Statistics

### Number of alarms created

The `alarms_created` statistic returns the number of alarms created. The alarms
created while a pbehavior was active are not taken into account. It does not
take any parameters.

#### Example

Request:

```javascript
POST /api/v2/stats/alarms_created
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Response:

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "alarms_created": 100
            }
        ]
    },
    // ...
]
```

### Number of alarms impacting an entity

The `alarms_impacting` statistic returns the number of alarms impacting an
entity. The alarms created while a pbehavior was active are not taken into
account. It does not take any parameters.

#### Example

Request:

```javascript
POST /api/v2/stats/alarms_impacting
{
    "group_by": ["entity_id"],
    "filter": [{
        "entity_type": "component"
    }],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Response:

```javascript
[
    {
        "tags": {
            "entity_id": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "alarms_impacting": 100
            }
        ]
    },
    // ...
]
```

### Number of alarms resolved

The `alarms_resolved` statistic returns the number of alarms resolved. The
alarms *created* while a pbehavior was active are not taken into account. It
does not take any parameters.

#### Example

Request:

```javascript
POST /api/v2/stats/alarms_resolved
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Response:

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "alarms_resolved": 86
            }
        ]
    },
    // ...
]
```

### Number of alarms canceled

The `alarms_canceled` statistic returns the number of alarms canceled. The
alarms *created* while a pbehavior was active are not taken into account. It
does not take any parameters.

#### Example

Request:

```javascript
POST /api/v2/stats/alarms_canceled
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Response:

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "alarms_canceled": 7
            }
        ]
    },
    // ...
]
```

### Mean acknowledgement time

The `mean_ack_time` rstatistic eturns the average acknowledgement time. The
alarms *created* while a pbehavior was active are not taken into account. It
does not take any parameters.

The acknowledgement time is the difference between the date of the *first*
acknowledgement and the date of creation of the alarm.

#### Example

Request:

```javascript
POST /api/v2/stats/mean_ack_time
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Response:

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "mean_ack_time": 426
            }
        ]
    },
    // ...
]
```

### Mean resolve time

The `mean_resolve_time` statistic returns the average resolve time. The alarms
*created* while a pbehavior was active are not taken into account. It does not
take any parameters.

The resolve time is the difference between the date of the resolution and the
date of creation of the alarm.

#### Example

Request:

```javascript
POST /api/v2/stats/mean_resolve_time
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Response:

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "mean_resolve_time": 2687
            }
        ]
    },
    // ...
]
```

### Acknowledgement time above or below the SLA

The `ack_time_sla` statistic returns the numbers and rates of acknowledgement
times above or below a SLA. The alarms *created* while a pbehavior was active
are not taken into account. The statistic takes a `sla` parameter whose value
is the SLA in seconds, and returns a JSON object containing the following
fields:

 - `above`: the number of alarms whose acknowledgement time is above the SLA.
 - `below`: the number of alarms whose acknowledgement time is below the SLA.
 - `above_rate`: the ratio of alarms whose acknowledgement time is above the
   SLA (between 0 and 1).
 - `below_rate`: the ratio of alarms whose acknowledgement time is below the
   SLA (between 0 and 1).

The acknowledgement time is the difference between the date of the *first*
acknowledgement and the date of creation of the alarm.

#### Example

Request:

```javascript
POST /api/v2/stats/resolve_time_sla
{
    "tstart": 1532815200,
    "tstop": 1532901600,
    "parameters": {
        "sla": 3600
    }
}
```

Response:

```javascript
[
    {
        "tags": {},
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "ack_time_sla": {
                    "above": 10
                    "below": 90,
                    "above_rate": 0.1,
                    "below_rate": 0.9,
                }
            }
        ]
    }
]
```

### Resolve time above or below the SLA

The `resolve_time_sla` statistic returns the numbers and rates of resolve times
above or below a SLA. The alarms *created* while a pbehavior was active are not
taken into account. The statistic takes a `sla` parameter whose value is the
SLA in seconds, and returns a JSON object containing the following fields:

 - `above`: the number of alarms whose resolve time is above the SLA.
 - `below`: the number of alarms whose resolve time is below the SLA.
 - `above_rate`: the ratio of alarms whose resolve time is above the SLA
   (between 0 and 1).
 - `below_rate`: the ratio of alarms whose resolve time is below the SLA
   (between 0 and 1).

The resolve time is the difference between the date of the resolution and the
date of creation of the alarm.

#### Example

Request:

```javascript
POST /api/v2/stats/resolve_time_sla
{
    "tstart": 1532815200,
    "tstop": 1532901600,
    "parameters": {
        "sla": 3600
    }
}
```

Response:

```javascript
[
    {
        "tags": {},
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "ack_time_sla": {
                    "above": 10
                    "below": 90,
                    "above_rate": 0.1,
                    "below_rate": 0.9,
                }
            }
        ]
    }
]
```

### Time spent in each state

The `time_in_state` statistic returns a JSON object with:

 - one field for each state (between 0 and 3), containing the time spent by the
   entity in this state, in seconds
 - a `total` field, containing the total time

The intervals during which a pbehavior was active are excluded from these
values. The total time may thus be inferior to the duration of the interval
`tstop - tstart`.

This statistic can only be computed for groups containing only one entity. It
is necessary to ensure that each group only contains one, for example by adding
`entity_id` to the `group_by`parameter.

#### Example

Request:

```javascript
POST /api/v2/stats/time_in_state
{
    "filter": [{
        "component": "c"
    }],
    "group_by": ["entity_id"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "periods": 2
}
```

Response:

```javascript
[
    {
        "tags": {
            "entity_id": "ressource1/c"
        },
        "periods": [
            {
                "tstart": 1532728800,
                "tstop": 1532815200,
                "time_in_state": {
                    0: 48159,
                    1: 34051,
                    2: 2203,
                    3: 1387,
                    "total": 85800
                }
            },
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "time_in_state": {
                    0: 52563,
                    1: 28465,
                    2: 4245,
                    3: 527,
                    "total": 85800
                }
            }
        ]
    },
    // ...
]
```

### Availability

The `availability` statistic returns the times and rates of availability and
unavailability. It takes an `available_state` parameter whose value is the
state until which an entity is considered to be available. It returns a JSON
object with the following fields:

 - `available`: the time during which the entity was in an available state
   (lower or equal to `available_state`), in seconds.
 - `unavailable`: the time during which the entity was in an unavailable state
   (strictly higher than `available_state`), in seconds
 - `available_rate`: the ratio of time during which the entity was in an
   available state (lower or equal to `available_state`). This value is between
   0 and 1.
 - `unavailable_rate`: the ratio of time during which the entity was in an
   unavailable state (strictly higher than `available_state`). This value is
   between 0 and 1.

The intervals during which a pbehavior was active are excluded from these
values. The total time `available_time + unavailable_time` may thus be inferior
to the duration of the interval `tstop - tstart`.

This statistic can only be computed for groups containing only one entity. It
is necessary to ensure that each group only contains one, for example by adding
`entity_id` to the `group_by`parameter.

#### Example

Request:

```javascript
POST /api/v2/stats/availability
{
    "filter": [{
        "component": "c"
    }],
    "group_by": ["entity_id"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "parameters": {
        "available_state": 1
    }
}
```

Response:

```javascript
[
    {
        "tags": {
            "entity_id": "ressource1/c"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "time_in_state": {
                    "available": 81028,
                    "unavailable": 4772,
                    "available_rate": 0.9443822843822843,
                    "unavailable_rate": 0.05561771561771562
                }
            }
        ]
    },
    // ...
]
```

### Maintenance

The `maintenance` statistic returns the time during which a pbehavior was or
was not active on an entity. It does not take any parameters and returns a JSON
object with the following fields:

 - `maintenance`: the time during which the entity had an active pbehavior, in
   seconds.
 - `no_maintenance`: the time during which the entity had no active pbehavior,
   in seconds.

This statistic can only be computed for groups containing only one entity. It
is necessary to ensure that each group only contains one, for example by adding
`entity_id` to the `group_by`parameter.

#### Example

Request:

```javascript
POST /api/v2/stats/maintenance
{
    "filter": [{
        "component": "c"
    }],
    "group_by": ["entity_id"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Response:

```javascript
[
    {
        "tags": {
            "entity_id": "ressource1/c"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "maintenance": {
                    "maintenance": 600,
                    "no_maintenance": 85800
                }
            }
        ]
    },
    // ...
]
```

### Mean Time Between Failures

The `mtbf` statistic returns the mean time between failures, i.e. the time
without maintenance divided by the number of failures.

This statistic can only be computed for groups containing only one entity. It
is necessary to ensure that each group only contains one, for example by adding
`entity_id` to the `group_by`parameter.

#### Example

Request:

```javascript
POST /api/v2/stats/mtbf
{
    "filter": [{
        "component": "c"
    }],
    "group_by": ["entity_id"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Response:

```javascript
[
    {
        "tags": {
            "entity_id": "ressource1/c"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "mtbf": 31.931522143654632
            }
        ]
    },
    // ...
]
```

### Alarm List

The `alarm_list` statistic returns a list of alarms. It does not take any
parameters, and returns a JSON array of objects which contains the tags of the
entity that created the alarm (`entity_id`, `entity_type`,
`entity_infos.<information_id>`, `connector`, `connector_name`, `component`,
`resource` and `alarm_state`), as well as the following fields:

 - `time`: the date of creation of the alarm
 - `pbehavior`: `"True"` if there was an active pbehavior when the alarm was
   created, `"False"` otherwise.
 - `value`: the time it took for the alarm to be resolved.

Only the resolved alarms are taken into account.

#### Example

Request:

```javascript
POST /api/v2/stats/alarm_list
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Response:

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "alarm_list": [
                    {
                        "time": 1532815202,
                        "entity_id": "resource1/component1",
                        "entity_type": "resource",
                        "connector": "connector",
                        "connector_name": "connector_name",
                        "resource": "resource1",
                        "alarm_state": "3",
                        "pbehavior": "False",
                        "value": 157
                    },
                    {
                        "time": 1532815325,
                        "entity_id": "resource2/component1",
                        "entity_type": "resource",
                        "connector": "connector",
                        "connector_name": "connector_name",
                        "resource": "resource2",
                        "alarm_state": "1",
                        "pbehavior": "False",
                        "value": 849
                    },
                    // ...
                ]
            }
        ]
    },
    // ...
]
```

### State list

The `state_list` statistic returns a list containing time intervals during
which an entity was in a certain state. It does not take any parameters, and
returns an array of JSON objects containing the following fields:

 - `start` : the date of the start of the interval.
 - `stop` : the date of the end of the interval.
 - `duration` : the duration of the interval.
 - `state` : the state of the entity during this interval.

This statistic can only be computed for groups containing only one entity. It
is necessary to ensure that each group only contains one, for example by adding
`entity_id` to the `group_by`parameter.

#### Example

Request:

```javascript
POST /api/v2/stats/state_list
{
    "group_by": ["entity_id"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Response:

```javascript
[
    {
        "tags": {
            "entity_id": "c/r"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "state_list": [
                    {
                        "start": 1532815200,
                        "end": 1532875424,
                        "duration": 60224,
                        "state: 0
                    },
                    {
                        "start": 1532875424,
                        "end": 1532878559,
                        "duration": 3135,
                        "state: 3
                    },
                    {
                        "start": 1532878559,
                        "end": 1532901600,
                        "duration": 23041,
                        "state: 0
                    }
                ]
            }
        ]
    },
    // ...
]
```

### Entities impacted by the most alarms

The `most_alarms_impacting` statistic returns a list containing the groups of
entities that are impacted by the largest number of alarms. The request takes
the following parameters:

 - `group_by` (required): the tags used to group the entities.
 - `filter` (optional): an entity filter. This parameters has the same format
   as the main `filter` parameter.
 - `limit` (optional): the maximum number of groups to return.

The parameters `group_by` and `filter` have to be defined in the `parameters`
field (or in `parameters.most_alarms_impacting` for the `/api/v2/stats` route),
and are distinct from the main `group_by` and `filter` fields. For example, to
get a the resources impacted by the most alarms grouped by components, the
`parameters.group_by` field should be set to `resource` (to compute the number
of alarms per resources), the `parameters.filter`field should contain
`"entity_type: "resource"` (to compute the number of alarms only for
resources), and the `group_by` field should be set to `component` (to group the
results by component). See the example below for the full request.

The request returns a list of objects ordered by descending number of alarms,
with the following fields:

 - `tags`: the tags of the group.
 - `value`: the number of alarms impacting this group.

#### Example

Request:

```javascript
POST /api/v2/stats/most_alarms_impacting
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "parameters": {
        "group_by": ["resource"],
        "filter": [{
            "entity_type": "resource"
        }],
        "limit": 2
    }
}
```

Response:

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "most_alarms_impacting": [
                    {
                        "tags": {
                            "resource": "resource3"
                        },
                        "value": 451
                    },
                    {
                        "tags": {
                            "resource": "resource1"
                        },
                        "value": 210
                    }
                ]
            }
        ]
    },
    // ...
]
```

### Entities creating the most alarms

The `most_alarms_created` statistic returns a list containing the groups of
entities that created the largest number of alarms. The request takes the
following parameters:

 - `group_by` (required): the tags used to group the entities.
 - `filter` (optional): an entity filter. This parameters has the same format
   as the main `filter` parameter.
 - `limit` (optional): the maximum number of groups to return.

The parameters `group_by` and `filter` have to be defined in the `parameters`
field (or in `parameters.most_alarms_created` for the `/api/v2/stats` route),
and are distinct from the main `group_by` and `filter` fields. For example, to
get a the resources impacted by the most alarms grouped by components, the
`parameters.group_by` field should be set to `resource` (to compute the number
of alarms per resources), the `parameters.filter`field should contain
`"entity_type: "resource"` (to compute the number of alarms only for
resources), and the `group_by` field should be set to `component` (to group the
results by component). See the example below for the full request.

The request returns a list of objects ordered by descending number of alarms,
with the following fields:

 - `tags`: the tags of the group.
 - `value`: the number of alarms created by this group.

#### Example

Request:

```javascript
POST /api/v2/stats/most_alarms_created
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "parameters": {
        "group_by": ["resource"],
        "filter": [{
            "entity_type": "resource"
        }],
        "limit": 2
    }
}
```

Response:

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "most_alarms_created": [
                    {
                        "tags": {
                            "resource": "resource3"
                        },
                        "value": 451
                    },
                    {
                        "tags": {
                            "resource": "resource1"
                        },
                        "value": 210
                    }
                ]
            }
        ]
    },
    // ...
]
```

### Entities with the worst Mean Time Between Failures

The `worst_mtbf` statistic returns a list containing the groups of entities
that have the worst mtbf. The request takes the following parameters:

 - `group_by` (required): the tags used to group the entities.
 - `filter` (optional): an entity filter. This parameters has the same format
   as the main `filter` parameter.
 - `limit` (optional): the maximum number of groups to return.

The parameters `group_by` and `filter` have to be defined in the `parameters`
field (or in `parameters.most_alarms_impacting` for the `/api/v2/stats` route),
and are distinct from the main `group_by` and `filter` fields. For example, to
get a the resources impacted by the most alarms grouped by components, the
`parameters.group_by` field should be set to `resource` (to compute the number
of alarms per resources), the `parameters.filter`field should contain
`"entity_type: "resource"` (to compute the number of alarms only for
resources), and the `group_by` field should be set to `component` (to group the
results by component). See the example below for the full request.

The request returns a list of objects ordered by descending number of alarms,
with the following fields:

 - `tags`: the tags of the group.
 - `value`: the mtbf.

#### Example

Request:

```javascript
POST /api/v2/stats/worst_mtbf
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "parameters": {
        "group_by": ["resource"],
        "filter": [{
            "entity_type": "resource"
        }],
        "limit": 2
    }
}
```

Response:

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "worst_mtbf": [
                    {
                        "tags": {
                            "resource": "resource3"
                        },
                        "value": 45
                    },
                    {
                        "tags": {
                            "resource": "resource1"
                        },
                        "value": 157
                    }
                ]
            }
        ]
    },
    // ...
]
```

### Longest alarms

The `longest_alarms` statistic returns a list of alarms that took the longest
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

#### Example

Request:

```javascript
POST /api/v2/stats/longest_alarms
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Response:

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "alarm_list": [
                    {
                        "time": 1532895472,
                        "entity_id": "resource2/component1",
                        "entity_type": "resource",
                        "connector": "connector",
                        "connector_name": "connector_name",
                        "resource": "resource2",
                        "alarm_state": "1",
                        "pbehavior": "False",
                        "value": 4892
                    },
                    {
                        "time": 1532854763,
                        "entity_id": "resource1/component1",
                        "entity_type": "resource",
                        "connector": "connector",
                        "connector_name": "connector_name",
                        "resource": "resource1",
                        "alarm_state": "3",
                        "pbehavior": "False",
                        "value": 3542
                    },
                   // ...
                ]
            }
        ]
    },
    // ...
]
```
