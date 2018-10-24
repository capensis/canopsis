# MongoDB

La configuration de l’accès à MongoDB se trouve dans `etc/common/mongo_store.conf` **uniquement**. Les fichiers `etc/mongo/storage.conf` et `etc/cstorage.conf`, même si remplis, ne sont plus utilisés.

```ini
[DATABASE]
host = HOST
port = PORT
db = DATABASE
user = USER
pwd = PASSWORD
```

## ReplicaSet

```ini
[DATABASE]
host = HOST1:PORT1,HOST2:PORT2,HOST3
port = PORT3
replicaset = RS_NAME
read_preference = PYMONGO_CONST
```

La valeur `RS_NAME` doit être le nom configuré du *replicaset* auquel on va se connecter.

La valeur `PYMONGO_CONST` va prendre une des valeurs listée pour `MongoClient` : https://api.mongodb.com/python/2.8/api/pymongo/index.html#pymongo.read_preferences.ReadPreference

Voici l’extrait de la doc :

 * `PRIMARY`: Queries are sent to the primary of a shard.
 * `PRIMARY_PREFERRED`: Queries are sent to the primary if available, otherwise a secondary.
 * `SECONDARY`: Queries are distributed among shard secondaries. An error is raised if no secondaries are available.
 * `SECONDARY_PREFERRED`: Queries are distributed among shard secondaries, or the primary if no secondary is available.
 * `NEAREST`: Queries are distributed among all members of a shard.

### ReplicaSet As ENVVAR

ReplicaSet is handled by env2cfg script.

The replicaSet url format is as following :
```
mongodb://[username:password@]host[:port]/canopsis?replicaSet=[replicaset_name]
```
