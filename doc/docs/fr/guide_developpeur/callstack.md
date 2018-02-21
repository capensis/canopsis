## Callstack printer

 * Project : `canopsis/canopsis` -> `sources/python/common/callstack.py`

### Description

Callstack permet de tracer les requêtes mongo exécutées par le storage, ainsi
que la pile d'appel amenant à cette requête.

### Utilisation

La callstack est loggée depuis un engine lorsque la variable d'environnement
`CPS_CALLBACK` vaut 1.

```bash
export CPS_CALLSTACK="1"
```

Éditer le fichier supervisorctl de l'engine ciblé et modifier la ligne
`command` comme suit:

```bash
command=/bin/sh -c "CPS_CALLSTACK=\"1\" engine-launcher -e canopsis.engines.dynamic -n ...""
```

Puis redémarrer le moteur.

Les log apparaitrons dans `/opt/canopsis/var/log/callstack.log`:
```bash
tail -F ~/var/log/callstack.log
```
