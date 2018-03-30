# Canopsis Root Path

Ou la racine de travail de Canopsis.

Afin de faciliter les choses, `canopsis.common` contient la variable `root_path` qui est initialisée au premier chargement de ce module.

## Détection

Elle va contenir l’une de ces valeurs (la première trouvée gagne) :

 * `CPS_PREFIX` : variable d’environnement si disponible
 * `sys.prefix` : le préfix tel que découvert par Python.
 * `/opt/canopsis` : chemin en dur du dernier espoir.

## Impossibilité de détecter la racine

Pour définir si l’un de ces chemin est correct, la fonction `_root_path()` va tenter de trouver le répertoire `os.path.join(root_path, 'etc')`.

**Si aucun chemin ne réussi**, l’exception **canopsis.common.CanopsisUnsupportedEnvironment** est levée et aucun processus Canopsis ne doit être lancé.

## Utilisation

Exemple d’ouverture du fichier de configuration `mongo_store.conf` situé dans `etc/common` :

```python
import os

from canopsis.common import root_path

fpath = os.path.join(root_path, 'etc', 'common', 'mongo_store.conf')
with open(fpath, 'r') as fh:
    ...
```
