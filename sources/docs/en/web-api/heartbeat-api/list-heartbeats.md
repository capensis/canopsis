# List Heartbeats

Return every Heartbeats stored in database.

**URL** : `/api/v2/heartbeat/`

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

**Request body example** : None


## Success Responses

**Condition** : Read from database successfully.

**Code** : `200 OK`

**Content example** : 

```json
[
    { 
      "_id": "<heartbeat_id1>",
      "pattern": {
          "connector": "c1", 
          "connector_name": "connector1"
      }, 
      "expected_interval": "10s"
    },
    { 
      "_id": "<heartbeat_id2>",
      "pattern": {
          "connector": "c2", 
          "connector_name": "connector2"
      }, 
      "expected_interval": "15m"
    }
]
```

## Error Responses

**Condition** : If database error.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
  "description": "can not retrieve the canopsis version from database, contact your administrator."
}
```
