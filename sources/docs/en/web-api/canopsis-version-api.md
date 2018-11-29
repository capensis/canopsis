# Get Canopsis version

Retrieve Canopsis version information.

**URL** : `/api/v2/version`

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None


## Success Responses

**Condition** : Data provided is valid and User is Authenticated.

**Code** : `200 OK`

**Content example** : 

```json
{
    "version": "3.4.0"
}
```

## Error Responses

**Condition** : If database error.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
   "description" : "can not retrieve the canopsis version from database, contact your administrator.",
   "name" : ""
}
```
##

**Condition** : If Canopsis version document not found.

**Code** : `404 NOT FOUND`

**Content example** :

```json
{
   "description" : "canopsis version info not found.",
   "name" : ""
}
```
