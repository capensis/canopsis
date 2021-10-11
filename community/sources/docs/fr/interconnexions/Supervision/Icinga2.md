# Connector Livestatus2Canopsis

<!-- XXX: documentation de très faible qualité -->

## Description

Convertit des évènements de supervision Icinga 2 en évènements Canopsis.

## Structure

```
├── application.py
├── doc
├── install
│   ├── apache-jmeter-2.13.tgz
│   ├── config.example.json
│   ├── lib-python-skeleton.install
│   ├── lib-python-specific.install
│   ├── skeleton.ansible
│   ├── specific.ansible
│   ├── uninstall.ansible
│   ├── virtualenv.py
│   └── virtualenv.pyc
└── README.md
```

## Installation && Configuration VENV

```
git clone https://git.canopsis.net/canopsis-connectors/connector-livestatus2canopsis.git
cd connector-livestatus2canopsis
git checkout rc
ansible-playbook install/skeleton.ansible
ansible-playbook install/specific.ansible
```

## Supprimer VENV

```
cd connector-livestatus2canopsis
ansible-playbook install/uninstall.ansible
```

## Application.py

Ce fichier contient le connecteur de classe. Dans cette classe, il existe quelques fonctions:

* **Init** initialise ConnectorTool, CanopsisPublisher, NagiosPublisher

* **Lister** Fonction qui exécute un écouteur de socket avec la configuration d'installation.

* **Lance** Fonction qui initie le fonctionnement du connecteur (Socket, One Exec, Daemon).

* **Processing** convertit les données d'entrée en données Canopsis mais publie également un évènement.

## Configuration
La configuration de ce connecteur ce fait dans `./etc/config.json`

* `wait_loop` :  le temps d'attente pour fermer la connexion lors de la demande de livestatus.

## Attribu :
* `state_changed` : est défini lorsque `state` est différent de `last_state` et peut être utiliser pour limiter l'utilisation de la bande passante lorsqu'il utilisé pour envoyer un évènement via Canopsis2Canopsis (Pro).

# Installation

## Requirements

Debian like (Debian, Ubuntu ...):

```
apt-get install build-essential git-core python supervisord
```

Redhat like (Centos ..):

```
yum groupinstall "Development Tools"
yum install git-core python supervisord
```

## Download & Build

via git:

```
    mkdir /opt/canopsis-connectors/
    git clone https://git.canopsis.net/canopsis-connectors/connector-livestatus2canopsis.git
    git clone -b develop git@git.canopsis.net:canopsis-connectors/connector-libs.git connector-libs
    cd connector-livestatus2canopsis
```

## InstallH


### Ansible H> v2

```
    ansible-playbook install/skeleton.ansible
    ansible-playbook install/specific.ansible
```

### Manuellement

```
    python install/virtualenv.py .
    for pkg in `cat install/lib-python-skeleton.install`; do bin/pip install $pkg; done
    for pkg in `cat install/lib-python-specific.install`; do bin/pip install $pkg; done
    mkdir etc
    cp install/config.example.json etc/config.json
    cp -r ../connector-libs lib/
    bin/python lib/connector-libs/setup.py install
    mv lib/connector-libs/connector_libs lib/python{2.6|2.7}/site-packages/
```

## Configuration

### livestatus2canopsis

Edit etc/config.json

| cat        | Param            | Type    | Description                           |
|------------|------------------|---------|---------------------------------------|
| connector  | loglevel         | Boolean | Enable debug mode                     |
| connector  | publish2nagios   | Boolean | Publish scenario result to nagios     |
| connector  | write2json       | Boolean | DUMP JSON DOC TO FILE                 |
| connector  | wait_loop        | Integer | Wait time during pulling              |
| amqp       | host             | IP      | AMQP IP                               |
| amqp       | port             | Integer | AMQP Port (default: 5672)             |
| amqp       | user             | String  | AMQP User (default: "guest")          |
| amqp       | pass             | String  | AMQP Pass (default: "guest")          |
| amqp       | vhost            | String  | AMQP VHost (default: "canopsis")      |
| livestatus | threadProcessing | Integer | How many thread processing the events |
| livestatus | socket           | String  | Livestatus socket                     |

### supervisord

Créer le fichier : /etc/supervisord.d/livestatus2canopsis.ini

```
    [program:connector-livestatus2canopsis]
    command=/opt/canopsis-connectors/connector-livestatus2canopsis/bin/python /opt/canopsis-connectors/connector-livestatus2canopsis/application.py
    directory=/opt/canopsis-connectors/connector-livestatus2canopsis
    process_name=%(program_name)s
    stdout_logfile=/var/log/%(program_name)s_out.log
    stderr_logfile=/var/log/%(program_name)s_err.log
    autostart=true
    autorestart=true
```
Créer les fichier de logs

```
    touch /var/log/connector-livestatus2canopsis_out.log
    touch /var/log/connector-livestatus2canopsis_err.log
    chmod 777 /var/log/connector-livestatus2canopsis_out.log
    chmod 777 /var/log/connector-livestatus2canopsis_err.log
```

# Use

## Start / Stop

Comme ce connecteur est un script python autonome, il n’a pas de mod deadmon, c’est pourquoi nous utilisons Supervord pour gérer le mode Deamon.

Et pour démarrer le processus de supervision global.

```
    service supervisord start
```

Et pour arrêter le processus de supervision global.

```
    service supervisord stop
```

Pour démarrer seulement un processus.

```
	service supervisord status
	supervisord start "Nom inside the status return"
```

Pour stopper seulement un processus.

```
	service supervisord status
	supervisord stop "Nom inside the status return"
```
