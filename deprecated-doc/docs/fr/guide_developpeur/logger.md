## Logger

 * Project : `canopsis/canopsis` -> `sources/python/logger`

### Description

Facilite l’instanciation et la récupération des loggers au sein de Canopsis, en évitant la duplication.

 * Fournis des formats en fonction du niveau de log. Un format peut être spécifié pour le logger créé.
 * La création des *handlers* est encapsulée dans des classes `Output<Kind>.make()`, permettant une plus grande clarté quant aux options à passer, et une plus grande facilité concernant les exceptions à traiter.
 * Fournis un moyen de logger dans un *handler* en mémoire, via un paramètre de `Logger.get`.

Se repose entièrement sur `logging` de Python.

Aucun besoin de généricité sur ce logger.

### Utilisation

Imports :

```python
# Import de base
from canopsis.logger import Logger

# Import optionnel des classes de création des handlers
from canopsis.logger import OutputFile, OutputStream
```

Utilisation basique, en mode fichier (par défaut) :

```python
# level: INFO
# output_cls: OutputFile
logger = Logger.get('wonderfullog', 'var/log/wonderful.log')

# message écrit dans /opt/canopsis/var/log/wonderful.log
logger.info('message')
```

Utilisation plus complète :

```python
# L'argument level supporte aussi logging.<LEVEL>
logger = Logger.get('wonderfullog', 'var/log/wonderful.log',
                    output_cls=OutputFile, level='debug')
```

OutputNull :

```python
from canopsis.logger import Logger, OutputNull

logger = Logger.get('dev_null', None, output_cls=OutputNull)

logger.info('trash')
```

Utilisation en mode *memory* :

Ce mode permet de transférer les logs à un *handler* en mémoire qui va fournir le log au *handler* final selon deux conditions :

 * La capacité du *handler* est atteinte : mettons que le *handler* accèpte de retenir 100 messages en mémoire, lorsque ce nombre est atteint, tous les logs sont transmis d’un coup au *handler* final, un fichier par exemple.
 * Le `flush_level` est atteint : vous décidez de mettre le niveau de log à `INFO`, et un message `CRITICAL` arrive. Vous voulez voir immédiatement ces messages : demandez au *handler* mémoire d’envoyer les messages au *handler* final lorsqu’un message de niveau égal au suppérieur est transmit.

```python
logger = Logger.get('wonderfullog', 'var/log/wonderful.log',
                    memory=True, memory_flushlevel='critical')

logger.info('restera en memoire')
logger.info('celui-ci aussi')
logger.warning('celui-la aussi')
logger.critical('flush les precedents ainsi que celui-ci sur le handler OutputFile')
```

### Créer une classe Output<Kind>

```python
# Ces imports sont réalisés dans le module canopsis.logger.
#
# De préférence, mettre les Output dans ce module.
import logging

from canopsis.logger import Output

class OutputKind(Output):

    def make(mandatory_param, optional_param='def_value')
        return logging.handlers.KindHandler(mandatory_param, optional_param)
```

### Paramètres optionnels

Pour passer des paramètres aux fonctions `Output<Kind>.make`, utiliser le paramètre `driver_make_args` qui transettra les options à la fonction. Cela fonctionne de la même manière qu’un `**kwargs` à la fin.

Pour savoir quels paramètres passer, voir la signature / docstring de la fonction `Output<Kind>.make()`.
