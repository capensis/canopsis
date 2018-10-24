## Tracer

 * Projet : `canopsis/canopsis` -> `sources/python/tracer`

### Description

Permet de tracer un des actions réalisées sur un ensemble d'entités du ContextGraph.

Une trace s’enregistre dans la collection `default_tracer` en respectant ce format :

```json
{
    "_id": "uniq_id_string",
    "impact_entities": ["list", "of", "entity", "ids"],
    "triggered_by": "string",
    "extra": {
        "custom": "dict"
    }
}
```

L’action réalisée se renseigne dans le champs `triggered_by`, qui prend le même nom dans l’API du Manager (ne pas créer ses traces à la main).

### Utilisation

```python
from canopsis.tracer.manager import TracerManager

TM = TracerManager()

help(TM)
```

### Créer une trace

Utiliser l’api du manager pour créer les traces.

L’information `triggered_by` doit être une *string*, mais veillez à bien la créer. Par exemple, si le déclencheur est `baseline`, il conviendra de peupler le champs de cette manière : `baseline.<baseline_name>`.

### Ajouter des entités

Une trace concernera des entités du ContextGraph.

Utiliser l’api du manager.
