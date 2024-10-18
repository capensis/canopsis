# neb2canopsis : module (Event Broker) Nagios/Nagios-like pour Canopsis

## Description

Convertit des évènements de supervision Nagios/Icinga en évènements Canopsis.

Il est écrit en C.

**Note :** certains modules Nagios-like étaient parfois compatibles avec Centreon, mais ce n'est plus le cas dans ses dernières versions. Centreon nécessite dorénavant un module dédié, [Centreon Stream Connector](https://docs.centreon.com/fr/docs/integrations/data-analytics/sc-canopsis-events/).

## Installation

Installer les outils de développement pour pouvoir compiler les sources (ainsi que Git si on veut récupérer les sources de cette façon) :

```shell
# Sur Debian / Ubuntu
$ sudo apt-get install build-essential git-core
# Sur Red Hat / CentOS
$ sudo yum groupinstall "Development Tools"
$ sudo yum install git-core
```

Récupérer les sources :
```shell
# Depuis Git
$ git clone git@git.canopsis.net:canopsis-connectors/connector-neb2canopsis.git
# Ou sinon, via HTTP
$ wget https://git.canopsis.net/canopsis-connectors/connector-neb2canopsis/repository/master/archive.tar.gz -O canopsis-nagios.tgz && tar xfz canopsis-nagios.tgz
```

Compiler :
```shell
$ cd connector-neb2canopsis*
$ make
```

(**Attention :** si vous souhaitez compiler un module Nagios 4.x (**expérimental**), il doit être compilé depuis sa branche dédiée :)
```shell
# Uniquement si l'on souhaite utiliser un module Nagios 4.x (expérimental)
$ git checkout 4.x
$ make 4x
```

Installer le module dans un dossier dédié. Par exemple, pour Nagios :
```shell
$ sudo mkdir -p /usr/local/nagios/bin
$ sudo cp neb2amqp.o /usr/local/nagios/bin/
```

(Ou, pour le module 4.x :)
```shell
# Uniquement si l'on souhaite utiliser un module Nagios 4.x (expérimental)
$ sudo cp neb2amqp-4x.o /usr/local/nagios/bin/neb2amqp.o
```

## Configuration

On prend ici l'exemple de Nagios.

Ajouter les éléments suivants à la configuration Nagios (fichier `nagios.cfg`) :

```
event_broker_options=-1
broker_module=/usr/local/nagios/bin/neb2amqp.o name=Central host=<adresse IP AMQP> userid=<identifiant AMQP> password=<mot de passe AMQP>
```

Les options disponibles sont les suivantes :

```
    host =              AMQP Server (127.0.0.1)
    port =              AMQP Port (5672)
    userid =            AMQP login (guest)
    password =          AMQP password (guest)
    virtual_host =      AMQP Virtual host (canopsis)
    exchange_name =     AMQP Exchange (canopsis.events)
    name =              Poller name (Central)
    connector =         Connector name (nagios) (you can type "icinga" for icinga)
    max_size =          Maximum message size to send to the AMQP bus (8192)
    cache_file =        File in which faulty messages are stored (/usr/local/nagios/var/canopsis.cache)
                        (note: if we cannot read/create the file, the cache will
                        only run in memory)
    cache_size =        Number of messages to store in cache (1000)
    autosync =          Delay in seconds between two automatic sync of the cache into 'cache_file'.
                        If < 0 disable autosync (note: the cache will always be stored when the module
                        is unloaded). If = 0 cache every time (this is not recommended as it may consumes
                        lot of I/O) (default: 60)
    autoflush =         Delay in seconds between two automatic flush of the cache into the AMQP bus
                        if it is available (60)
    rate =              Delay in ms between two messages when depiling (5)
    flush =             Number of messages to send when depiling (-1: means it is calculated at runtime)
    purge =             If 1, purge cache at startup. /!\ This will increase Nagios' startup time. (default: 0)

    hostgroups =        If 1, send host groups on host-check events (default: 0)
    servicegroups =     If 1, send service groups on service-check events (default: 0)
    acknowledgement =   If 1, handle acknowledgement events (default: 0)
    downtime =          If 1, handle downtime events (default: 0)
    custom_variables =  If 1, add Nagios macros to event (default: 0)
    urls =              If 1, add action_url and notes_url to event (default: 0)

    amqp_wait_time =    Number of seconds before a reconnection to AMQP
```

Le service doit ensuite être relancé pour prendre en compte ces modifications :
```shell
$ sudo systemctl restart nagios.service
```

## Vérifier le bon lancement du module

Une fois le service redémarré, les logs du service doivent donner des informations sur l'état du module.

Par exemple, pour Nagios et son fichier `nagios.log` :
```
neb2amqp: Setting hostname to localhost
neb2amqp: Setting userid to cpsrabbit
neb2amqp: Setting password to canopsis
neb2amqp: Setting hostgroups to 'true'
neb2amqp: Setting custom_variables to 'true'
neb2amqp: Setting servicegroups to 'true'
neb2amqp: Setting urls to 'true'
neb2amqp: NEB2amqp 0.6-fifo by Capensis. (connector: nagios)
neb2amqp: Please visit us at http://www.canopsis.org/
neb2amqp: Initialize FIFO: /tmp/neb2amqp.cache (maximum size: 10000
neb2amqp: FIFO: Load events from file
neb2amqp: FIFO: Open fifo file '/tmp/neb2amqp.cache'
neb2amqp: FIFO: File successfully opened
neb2amqp: FIFO: Push (0)
neb2amqp: FIFO: Close file
neb2amqp: FIFO: 1 events loaded
neb2amqp: AMQP: Re-connect to amqp ...
neb2amqp: AMQP: Init connection
neb2amqp: AMQP: Creating socket
neb2amqp: AMQP: Opening socket
neb2amqp: AMQP: Login
neb2amqp: AMQP: Open channel
neb2amqp: AMQP: Successfully connected
neb2amqp: Register callbacks
neb2amqp: successfully finished initialization
Event broker module '/usr/local/nagios/bin/neb2amqp.o' initialized successfully.
```

Si le module n'arrive pas à se connecter au bus AMQP, une erreur de ce type s'affichera :
```
neb2amqp: AMQP: Re-connect to amqp ...
neb2amqp: AMQP: Init connection
neb2amqp: AMQP: Creating socket
neb2amqp: AMQP: Opening socket
neb2amqp: AMQP: Opening socket: (unknown error)
neb2amqp: FIFO: Sync 285 events to file
neb2amqp: FIFO: Remove fifo file
neb2amqp: FIFO: Open fifo file '/tmp/neb2amqp.cache'
neb2amqp: FIFO: File successfully opened
neb2amqp: FIFO: 285 events written to file
```

Cette erreur peut notamment provenir d'un serveur AMQP indisponible (adresse IP inaccessible, mauvais identifiants RabbitMQ...).

Les messages sont alors mis en attente dans un fichier temporaire, et ils seront dépilés dès que la connexion AMQP sera à nouveau fonctionnelle :
```
neb2amqp: AMQP: Re-connect to amqp ...
neb2amqp: AMQP: Init connection
neb2amqp: AMQP: Creating socket
neb2amqp: AMQP: Opening socket
neb2amqp: AMQP: Login
neb2amqp: AMQP: Open channel
neb2amqp: AMQP: Successfully connected
neb2amqp: AMQP: Shift queue, size: 734
neb2amqp: AMQP: 734/734 events shifted from Queue, new size: 0
neb2amqp: FIFO: Sync 0 events to file
neb2amqp: FIFO: Remove fifo file
neb2amqp: FIFO: Open fifo file '/tmp/neb2amqp.cache'
neb2amqp: FIFO: File successfully opened
neb2amqp: FIFO: 0 events written to file
```

## Aller plus loin

D'autres outils sont disponibles dans le répertoire `contrib/`.
