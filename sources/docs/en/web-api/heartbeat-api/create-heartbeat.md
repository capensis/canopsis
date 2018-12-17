# Create Heartbeat

Create a new heartbeat. Read the body of the request to extract the heartbeat as a json.

**URL** : `/api/v2/heartbeat`

**Method** : `POST`

**Auth required** : YES

**Permissions required** : None

**Request body example** :
```json
{
  "pattern": {
      "connector": "c1", 
      "connector_name": "connector1"
  }, 
  "expected_interval": "10s"
}
```


## Success Responses

**Condition** : Heartbeat successfully created.

**Code** : `200 OK`

**Content example** : 

```json
{
  "name": "heartbeat created",
  "description": "<Heartbeat ID>"
}
```

## Error Responses

**Condition** : If invalid Heartbeat payload.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
  "description": "invalid heartbeat payload."
}
```
##

**Condition** : If Heartbeat pattern already exist.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
  "description": "heartbeat pattern already exists"
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
