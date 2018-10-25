# Guide Connecteurs

SOURCE : https://canopsis.readthedocs.io/en/readthedocs/_sources/canopsis-connectors/connector-shinken2canopsis/setup.txt

## Section : Supervision / Shinken

### Configuration

#### Pré requis

The broker module allow dynamic connection loss and reconnect with
canopsis amqp message bus without loss of events. You should consider to
set ``maxqueuelength`` (the maximum number of events that should be kept
in case of connection loss).

Le module broker vous permet la perte de connexion et la reconnexion avec le 
bus de messages canopsis amqp sans perte d'événements.
Vous devriez envisager de définir `` maxqueuelength`` (le nombre maximal d'événements à conserver 
en cas de perte de connexion). 

```
easy_install kombu
```

#### Setup

Le module broker canopsis est présent de manière native dans la distribution shinken. 
Vous aurez besion d'au moins la version dev. 

Vous devez uniquement activer le module broker et au moins configurer l'adresse de l'hôte canopsis.

Modifier le `` etc / Shinken-specific.cfg`` et ajouter `` Canopsis`` à la liste des modules activés :

```
define broker {
  modules Livestatus, Simple-log, WebUI, Canopsis
}
```

Dans le même fichier, recherchez le module `` Canopsis`` et définissez au moins la directive host sur l'adresse de l'hôte canopsis :

```
define module{
       module_name          Canopsis
       module_type          canopsis
       host                 xxx.xxx.xxx.xxx
       port                 5672
       user                 guest
       password             guest
       virtual_host         canopsis
       exchange_name        canopsis.events
       identifier           shinken-1
       maxqueuelength       50000
       queue_dump_frequency 300
}
```

Lorsque vous souhaitez connecter Shinken sur Canopsis, il existe un 
conflit de ports Mongodb . Vous devez donc changer votre configuration Shinken et 
le fichier ``/etc/mongodb.conf``

```
    port=27018
``` 

Et redémarrer le service

```
    /etc/init.d/mongodb restart
```

Puid éditer ``shinken-specific.cfg``

```
    define module {
      module_name Mongodb
      module_type mongodb
      uri mongodb://localhost:27018/?safe=true
      database shinken
    }
    define module {
      module_name mongologs
      module_type logstore_mongodb
      mongodb_uri mongodb://localhost:27018/?safe=true
    }
    define module {
      module_name MongodbRetention
      module_type mongodb_retention
      uri mongodb://localhost:27018/?safe=true
      database shinken
    }
```