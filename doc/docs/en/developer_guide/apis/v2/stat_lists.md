# Statistics Lists API

## Route

#### URL

`POST /api/v2/stats/lists/<list name>`

#### Parameters

This route takes a JSON object with the following fields:

 - `tstop`: a timestamp indicating the end of the period for which the
   list will be computed. This timestamp should be at the top of an hour
   (e.g. 12:00, not 12:03).
 - `duration`: the duration of the period, represented by a string
   `"<n><unit>"`, where `<n>` is an integer and `<unit>` a time unit (`h`, `d`,
   `w` or `m`).
 - `mfilter`: a mongodb filter, filtering the entities for which the
   list should be computed.
 - `parameters`: an object containing parameters for the computed list. See the
   documentation of each list below for the available parameters.


#### Response

If the request succeeded, the response is an array if JSON objects containing :

 - a `values` field containing the values of the list.
 - an `entity` field containing the entity for which the list was computed.


## Lists

### Intervals spent in each state

The `state_intervals` list returns a list of time intervals with the state of
the entity in each of these intervals.

#### Parameters

This statistic take the following parameters (in the `parameters` field) :

 - `states` (optional): Only the alarms whose state at the creation date is in
   this list will be taken into account.


### Example

```javascript
POST /api/v2/stats/lists/state_intervals
{
	"mfilter": {
		"type": "resource",
	},
    "tstop": 1534716000,
	"duration": "1h"
}
```

```javascript
[
    {
        "values": [
            {
                "duration": 969,
                "start": 1536242399,
                "state": 3,
                "end": 1536243119
            },
            {
                "duration": 326,
                "start": 1536243368,
                "state": 0,
                "end": 1536243694
            },
            {
                "duration": 23,
                "start": 1536243694,
                "state": 2,
                "end": 1536243708
            },
            {
                "duration": 2282,
                "start": 1536243717,
                "state": 0,
                "end": 1536245999
            }
        ],
        "entity": {
            "_id": "resource1/component1",
            "type": "resource"
            "impact": [
                "service"
            ],
        }
    },
    // ...
]
```
