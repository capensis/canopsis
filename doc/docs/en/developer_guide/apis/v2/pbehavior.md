# API PBehaviors


Pbehaviors are recuring calendar events that temporarily stop the monitoring of an entity during a given time (for maintenance for example).


**Note about the`tstart`et `tstop` attributes of a pbehavior**


Pbehaviors are similar to calendar events, with optional recurrence. The tstart and tstop parameters are use to define the start and end date of the first event occurrence. When the event repreats, these attributes are used to define each instance's duration, based on the first instance start and stop time.

For example, the behavior below will start à 11AM and finish 1 hour later every morning starting from 2018/06/18:

```json
{
// --snip --
	"rrule": "FREQ=WEEKLY;WKST=MO",
	"tstart": 1529312400, //Le 2018/06/18 at 11
    "tstop": 1529316000, // Le 2018/06/18 at 12
// -- snip --
}
```

**When the event is recurring, the date of the last occurrence will be stored in the `UNTIL` attribute of the `rrule`**


## Creating a pbehavior

This route Creates a new pbehavior

#### Url

  `POST /api/v2/pbehavior`

#### POST exemple

json body:

```json
{
    "name": "imagine",
    "author": "lennon",
    "filter_": {"_id": "all_the_people"},
    "rrule": "",
    "tstart": 0,
    "tstop": 1000
}
```

The body's attributes are the following:

| Name           | type    | nullable | Description                              |
|----------------|---------|----------|------------------------------------------|
| connector      | string  | No       | Identifier of the entity connector       |
| name           | string  | No       | Display name of the pbehavior            |
| author         | string  | No       | Author name                              |
| enabled        | boolean | No       | Should the pbehavior trigger or not      |
| reason         | string  | yes      | Administrative reason (optionnal)        |
| comments       | array   | yes      | Comments (option)                        |
| filter         | string  | No       | Entities filter (json)                   |
| type_          | string  | No       | Pbehavior type                           |
| connector_name | string  | No       | Display name of the entity connector     |
| rrule          | string  | yes      | Rrule (recurrence)                       |
| tstart         | integer | No       | Timestamp of the start date              |
| tstop          | integer | No       | Timestamp  end date                      |
| _id            | string  | No       | Pbheavior identifier                     |
| eids           | array   | No       | array of _ids for the impacted entities. |


Response: uid of the inserted element

```{json}
"b72e841a-d9d1-11e7-9a70-022abfd0f78f"
```

## Fetching pbehaviors for an entity

this route lists existing pbeahviors sufor an entity, identified by its eid (Entity ID)

#### URL

`GET /api/v2/pbehavior_byeid/<entityid>`

#### Parameters

* entityid <string> the ID of the target entity.


#### Response

```json
[
    {
        "connector": "canopsis",
        "name": "imagine",
        "author": "lennon",
        "enabled": true,
        "reason": "",
        "comments": null,
        "filter": "{\"_id\": \"580059AB4B100031\"}",
        "type_": "generic",
        "connector_name": "canopsis",
        "rrule": "FREQ=WEEKLY;COUNT=30;WKST=MO",
        "tstart": 1529312725,
        "tstop": 1592471125,
        "_id": "dd4cbc2c-72d6-11e8-a732-0242ac12001a",
        "isActive": true,
        "eids": [
            "580059AB4B100031"
        ]
    }
]
```

Response attributes are the following:

| Name           | type    | nullable | Description                              |
|----------------|---------|----------|------------------------------------------|
| connector      | string  | No       | Identifier of the entity connector       |
| name           | string  | No       | Display name of the pbehavior            |
| author         | string  | No       | Author name                              |
| enabled        | boolean | No       | Should the pbehavior trigger or not      |
| reason         | string  | yes      | Administrative reason (optionnal)        |
| comments       | array   | yes      | Comments (option)                        |
| filter         | string  | No       | Entities filter (json)                   |
| type_          | string  | No       | Pbehavior type                           |
| connector_name | string  | No       | Display name of the entity connector     |
| rrule          | string  | yes      | Rrule (recurrence)                       |
| tstart         | integer | No       | Timestamp of the start date              |
| tstop          | integer | No       | Timestamp  end date                      |
| _id            | string  | No       | Pbheavior identifier                     |
| eids           | array   | No       | Array of _ids for the impacted entities. |
| isActive       | boolean | No       | is the pbehavior currently active        |



## Delete a pbheavior

This route allows to remove a pbehavior

#### Url

  `DELETE /api/v2/pbehavior/<pbehavior_id>`

#### DELETE exemple



Response: a status object

```json
{
    "deletedCount": 1,
    "acknowledged": true
}
```


## Update a pbeahvior

There is currently no method to update a pbehavior in place. It is necessary to remove and recreate a pbehavior to update its content.


## Forcing pbheaviors computation

This route forces a new computation for all pbehaviors.

#### Url

  `GET` /api/v2/compute-pbehaviors

#### GET exemple

/api/v2/compute-pbehaviors

Response: has computation been trigered ?

```json
true
```
