## MongoStore

 * Project : `canopsis/canopsis` -> `sources/python/common/mongo_store.py`

### Description

Permet d'instancier des pymongo.Collection simplement.

Se repose entièrement sur la lib `pymongo`.

### Utilisation

Ligne d'import :

```python
from canopsis.common.mongo_store import MongoStore
```

Utilisation minimaliste :

```python
from canopsis.common.mongo_store import MongoStore
from canopsis.confng.simpleconf import Configuration
from canopsis.confng.vendor import Ini

conf = {MongoStore.CONF_CAT: {'db': 'my_database'}}
cred_conf = Configuration.load(MongoStore.CRED_CONF_PATH, Ini)

store = MongoStore(config=conf, cred_config=cred_conf)

collection = store.get_collection('object')

# On peut maintenant utiliser 'collection' comme un object pymongo classique
collection.find_one({'_id': 'robby'})
```

### Configurations

Par compatibilité avec l'ancien code, il faut fournir deux fichiers de
configuration :

 - conf: contenant notamment le nom de la base à laquelle se connecter
 - cred_conf: contenant le nom d'utilisateur et le mot de passe d'authentification

Par défaut, MongoStore essaie de se connecter à `localhost:27017/canopsis`.

### Fonctions

 - **get_collection(name)**: retourne un objet Collection avec le nom demandé
