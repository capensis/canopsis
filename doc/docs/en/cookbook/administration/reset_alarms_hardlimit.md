# Resetting hard-limit alarsm

** Warning: This process destroys data in the database. We urge you to perform a full backup of the database before proceeding**


This procedure resets all alarms that have reached their hard limit in order to make them editable again.


2 steps are necessary 

- Database backup
- Alarms reset

Resetting alarms destroys all history of all the alarms that have reached their hard limit


## Database backup 


- Log into the MongoDB master server
- from a bash sell, start the backup:

```bash
$ mongodump -u cpsmongo -p canopsis -d canopsis --gzip --archive=dump.gz
```

Store the output file`dump.gz`. It will be used if a restore is needed

## Resetting the alarms

- Log into the MongoDB master server
- open a mongo shell: 

```bash
$ mongo -u cpsmongo -p canopsis canopsis
```

Check that alarms need to be reset :

```javascript

db.getCollection('periodical_alarm').find({"v.hard_limit": {$ne: null}})

```
If this request resturns no result, there is no need to continue: all alarms are editable. 

Otherwise, run the following query

```javascript
db.getCollection('periodical_alarm').updateMany({"v.hard_limit": {$ne: null}}, {$set: {"v.hard_limit": null,"v.steps": []}})
```
The expexted result displays the number of affected rows

```json
{
    "acknowledged" : true,
    "matchedCount" : 6.0,
    "modifiedCount" : 6.0
}
```

The attribute `matchedCount` indicates the number of alarms that match the query executed

 The attibute`modifiedCount` indicates the number of alarms that have been updated by the query

Run again the following query to ensure that all alarms have been updated

```javascript
db.getCollection('periodical_alarm').find({"v.hard_limit": {$ne: null}})

```
The request should return no result. If some results remain, run the process again, then contact the Support


Alarms are now editable, but their whole history is now empty.


## In case of issue: Restoring the database

To restore the database, execute the following commands: 

- Log in on the MongoDB Master 

Navigate to the folder where the file `dump.gz` is stored, then run the following command:

```bash
$  mongorestore -u cpsmongo -p canopsis -d canopsis --drop --gzip --archive=dump.gz

```
