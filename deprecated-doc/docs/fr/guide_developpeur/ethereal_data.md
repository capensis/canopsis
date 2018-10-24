## EtherealData

 * Project : `canopsis/canopsis` -> `sources/python/common/ethereal_data.py`

### Description

Permet de lire des données de configuration depuis la base de données mongo au 
travers d'un cache paramétrable.

Aucunes dépendances externes.

### Utilisation

Ligne d'import :

```python
from canopsis.common.ethereal_data import EtherealData
```

Utilisation minimaliste :

```python
from pymongo import MongoClient
from canopsis.common.ethereal_data import EtherealData

my_collection = MongoClient().my_database.my_collection

ed = EtherealData(collection=my_collection,
                  filter_={'_id': 'doc_id'},
                  timeout=30)
# Les données sont dans la collection 'my_collection', le document où sont 
# présentes les données s'appel 'doc_id' et l'on veut que les données aient 
# moins de 30 secondes d'ancienneté

ed.get('Mario', 'Luigi')  # = "Luigi"
ed.set('Mario', 'bros')
ed.get('Mario', 'Luigi')  # = "bros"
```

### Fonctions

 - **get(value, default=None)**: lit une donnée dans la base, ou renvoie une valeur par défaut si elle n'existe pas
 - **set(key, value)**: met à jour la donnée en base
