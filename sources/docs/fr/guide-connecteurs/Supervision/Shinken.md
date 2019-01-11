# Shinken

### Configuration

#### Pré requis

Le module broker vous permet la perte de connexion et la reconnexion avec le 
bus de messages Canopsis AMQP sans perte d'évènements.
Vous devriez envisager de définir `maxqueuelength` (le nombre maximal d'évènements à conserver 
en cas de perte de connexion). 

```
easy_install kombu
```

#### Setup

Le module broker Canopsis est présent de manière native dans la distribution Shinken. 
Vous aurez besoin d'au moins la version dev. 

Vous devez uniquement activer le module broker et au moins configurer l'adresse de l'hôte Canopsis.

Modifier le `etc/Shinken-specific.cfg` et ajouter `Canopsis` à la liste des modules activés :

```
define broker {
  modules Livestatus, Simple-log, WebUI, Canopsis
}
```

Dans le même fichier, recherchez le module `Canopsis` et définissez au moins la directive host sur l'adresse de l'hôte Canopsis :

```
define module {
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
conflit de ports MongoDB . Vous devez donc changer votre configuration Shinken et 
le fichier `/etc/mongodb.conf` :

```
    port=27018
``` 

Et redémarrer le service

```sh
systemctl restart mongod
```

Puis, éditer `shinken-specific.cfg` :

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
