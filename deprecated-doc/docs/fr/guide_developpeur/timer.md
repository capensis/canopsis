## Timer

 * Project : `canopsis/canopsis` -> `sources/canopsis/tools/timer`

### Description

Permet de mesurer le temps d'exécution sur une portion de code facilement.

### Utilisation

Imports :

```python
from canopsis.tools.timer import Timer
```

Utilisation basique, en mode fichier (par défaut) :

```python
from canopsis.tools.timer import Timer
from canopsis.logger import Logger, OutputNull

mylogger = Logger.get('', None, output_cls=OutputNull)
mytimer = Timer(logger=mylogger)

mytimer.start('man')
# "Starting timing for action: man"

# ... ensemble d'actions que je veux chronométrer

mytimer.stop()
# "Action : man took 42 ms to complete"
```

Remise à zéro du timer :  il peut être utile de réinitialiser le timer lorsque sa méthode stop() ne peut être appelée (en cas d'erreur par exemple).


```python
from canopsis.tools.timer import Timer
from canopsis.logger import Logger, OutputNull

mylogger = Logger.get('', None, output_cls=OutputNull)
mytimer = Timer(logger=mylogger)

try:
	mytimer.start('man')
	# "Starting timing for action: man"

	# ... L'action à chronométrer lance une exception, mytimer.stop() ne sera donc jamais appelée

	mytimer.stop()
Except Exception e:
	mytimer.reset()
	# ... N'affiche rien mais remet le timer à zéro.
```


### Paramètres

Timer a simplement besoin d'un logger dans lequel écrire.
