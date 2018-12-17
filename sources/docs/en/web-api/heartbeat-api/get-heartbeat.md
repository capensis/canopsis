# Get Heartbeat

Get a Heartbeat by ID.

**URL** : `/api/v2/heartbeat/<heartbeat_id>`

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

**Request body example** : None


## Success Responses

**Condition** : Heartbeat was found.

**Code** : `200 OK`

**Content example** : 

```json
{ 
  "_id": "<heartbeat_id>",
  "pattern": {
      "connector": "c1", 
      "connector_name": "connector1"
  }, 
  "expected_interval": "10s"
}
```

## Error Responses

**Condition** : If Heartbeat not found.

**Code** : `404 NOT FOUND`

**Content example** :

```json
{
  "description" : "heartbeat not found",
  "name" : "<heartbeat_id>"
}
```
##

**Condition** : If database error.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
  "description": "can not retrieve the canopsis version from database, contact your administrator."
}
```
