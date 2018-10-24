## Models

 * Projet : `canopsis/canopsis` -> `sources/python/models`

### Description

Les classes de ce module permettent de gérer les modèles de base de Canopsis,
sous forme d'objet (et non plus de dict).

### Utilisation

```python
from canopsis.models.entity import Entity

entities = storage._backend.find({})
entities = [Entity(**Entity.convert_keys(e)) for e in entities]
```

Ou encore mieux, avec les MongoStore / MongoCollection :

```python
from canopsis.models.entity import Entity
from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection

conf = {
    MongoStore.CONF_CAT: {
        'db': 'canopsis',
        'user': '*',
        'pwd': '*'
    }
}

store = MongoStore(config=conf)
collection = MongoCollection(store.get_collection('default_entities'))

entities = collection.find({})
entities = [Entity(**Entity.convert_keys(e)) for e in entities]
```
