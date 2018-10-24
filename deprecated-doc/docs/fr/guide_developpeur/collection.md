## MongoCollection

 * Project : `canopsis/canopsis` -> `sources/python/common/collection.py`

### Description

Facilite l’utilisation de la lib pymongo

 * Permet de faire les opérations de base sur une collection (find, insert, update, remove)
 * Permet de vérifier qu'une transaction s'est bien déroulée
 * Et tout est bien loggé

Se repose entièrement sur la lib `pymongo`.

### Utilisation

Ligne d'import :

```python
from canopsis.common.collection import MongoCollection
```

Utilisation basique :

```python
from canopsis.common.collection import MongoCollection
from canopsis.middleware.core import Middleware

storage = Middleware.get_middleware_by_uri('storage-default-macollection://')
mc = MongoCollection(collection=storage._backend)

# On peut insérer/updater des objects directement dans la collection
id_ = mc.insert({'Eric': 'Idle'})
result = mc.update({'_id': id_}, {'Graham': 'Chapman'})

# Et vérifier que tout s'est bien passé
MongoCollection.is_successfull(result)  # == True
```

### Notes

 * Attention ! La fonction update de mongo remplace complètement le document (càd qu'il n'y a pas de fusion entre l'ancienne donnée et la nouvelle). Pour faire un update partiel, il faut utiliser « $set »:
```python
mc.update({'_id': id_}, {'$set': {'John': 'Cleese'}})
mc.find_one({'_id': id_})
# => {'_id': 'xx-xx', 'Graham': 'Chapman', 'John': 'Cleese'}
```
