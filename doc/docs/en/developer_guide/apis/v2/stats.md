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
 - `above_rate`: the percentage of alarms whose ack time is above the SLA.
 - `below_rate`: the percentage of alarms whose ack time is below the SLA.

### Resolve time above or below the SLA

The statistic `resolve_time_sla` is a JSON object with the following fields:

 - `above`: the number of alarms whose resolve time is above the SLA.
 - `below`: the number of alarms whose resolve time is below the SLA.
 - `above_rate`: the percentage of alarms whose resolve time is above the SLA.
 - `below_rate`: the percentage of alarms whose resolve time is below the SLA.


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
    "ack_time_sla": 600
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

Requête:

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

Réponse:

```javascript
[
    {
        "tags": {},
        "alarms_created": 13
    }
]
```


### Multiple statistics in one request

`/api/v2/stats`

Request:

```javascript
{
    "stats": ["alarms_created", "alarms_resolved"],
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
        "alarms_created": 13,
        "alarms_created": 8,
    }
]
```
