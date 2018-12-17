# Delete Heartbeat

Delete a Heartbeat by ID.

**URL** : `/api/v2/heartbeat/<heartbeat_id>`

**Method** : `DELETE`

**Auth required** : YES

**Permissions required** : None

**Request body example** : None


## Success Responses

**Condition** : Heartbeat was removed.

**Code** : `200 OK`

**Content example** : 

```json
{
  "name": "heartbeat removed",
  "description": "<heartbeat_id>"
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
