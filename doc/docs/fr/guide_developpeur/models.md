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
