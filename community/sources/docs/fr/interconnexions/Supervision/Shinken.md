# Shinken

Convertit des évènements de supervision Shinken en évènements Canopsis.

## Prérequis

- Installez Shinken Entreprise 02.07.06 ou supérieur en suivant la documentation de l'éditeur.
- Installez également la bibliothèque kombu soit avec `pip` soit en passant par le gestionnaire de paquets de votre distribution `yum install python2-kombu`/`apt-get install python-kombu`

## Installation

Le module broker Canopsis sera bientôt présent de manière native dans la distribution Shinken Entreprise. Pour le moment notre support peut vous fournir les fichiers nécessaires.

Copiez le dossier `broker-exporter-canopsis` vers `/var/lib/shinken/modules/`

Puis copiez le fichier `broker-exporter-canopsis.cfg` vers `/etc/shinken/modules/` et modifiez les droits du fichier avec la commande suivante :
```sh
chown shinken:shinken /etc/shinken/modules/broker-exporter-canopsis.cfg
```

## Configuration

Éditez le fichier `/etc/shinken/brokers/broker-master.cfg` et ajoutez `broker-exporter-canopsis` à la liste des modules activés.
```apacheconf
#======== Modules to enable for this daemon =========
    # Available:
    # - Simple-log            : save all logs into a common file
    # - WebUI                 : visualisation interface
    # - Graphite-Perfdata     : save all metrics into a graphite database
    # - sla                   : save sla into a database
    # - Livestatus            : TCP API to query element state, used by nagios external tools like NagVis or Thruk
    # - event-manager-writer  : save events for events manager (do not forget to activate the module in your webui to see data)
    modules                   Simple-log, WebUI, Graphite-Perfdata, sla, event-manager-writer, Livestatus, broker-exporter-canopsis
```

Personnalisez ensuite le contenu du fichier `/etc/shinken/modules/broker-exporter-canopsis.cfg`, par exemple :
```apacheconf
    #======== Canopsis address =========
    # Canopsis host to connect to
    host                 192.168.0.123
    # rabbitmq port
    port                 5672
    # must be changed
    user                 cpsrabbit
    password             canopsis
    virtual_host         canopsis
    exchange_name        canopsis.events
    # need a unique identifier because there should be more than one shinken in canopsis
    identifier           shinken-1
    # maximum event stored in queue when connection with canopsis is lost
    maxqueuelength       50000
    # frequency (in seconds) on which the queue is saved for retention
    queue_dump_frequency 300
```

Le module broker vous permet la perte de connexion et la reconnexion avec le bus de messages Canopsis AMQP sans perte d'évènements.
Vous pouvez modifier la valeur de `maxqueuelength` (le nombre maximal d'évènements à conserver en cas de perte de connexion) selon vos besoins.

Enfin redémarrez Shinken.

```sh
service shinken restart
```

Aucune configuration n'est nécessaire dans l'interface de Shinken, le connecteur enverra automatiquement les évènements vers Canopsis.
