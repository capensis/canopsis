# connector-centreon-engine : module (Event Broker) Centreon pour Canopsis

## Description

Ce dépôt contient un module Broker Centreon pour l'envoi d'évènements à Canopsis. Il est écrit en C++, et des modules pré-compilés sont disponibles pour certains environnements.

Ce module vient s'ajouter dans une instance Centreon existante, et doit ensuite être configuré depuis l'interface web de Centreon afin de rediriger le trafic vers Canopsis.

**Pré-requis réseau :** la transmission de flux réseau doit être permise entre Centreon et Canopsis (port 5672 par défaut).

## Installation d'un module pré-compilé

### Récupération du module

Clôner le dépôt Git contenant le module (**note :** activer votre proxy HTTP si nécessaire) :
```shell
# Installation de Git sur Debian / Ubuntu
$ sudo apt-get install git-core
# Installation de Git sur Red Hat / CentOS
$ sudo yum install git-core

# Clône du dépôt
$ git clone https://VOTRE-UTILISATEUR-GITLAB-ICI@git.canopsis.net/canopsis-connectors/connector-centreon-engine.git
$ cd connector-centreon-engine
```

### Installation du module

Des modules pré-compilés sont disponibles dans le répertoire `precompiled/` de ce dépôt Git.

Il faut, pour cela, connaître votre environnement (CentOS 6, CentOS 7…) et votre version du Centreon Broker sur le système cible :
```shell
$ cat /etc/centos-release
CentOS Linux release 7.4.1708 (Core)
$ cbd -v
[1513786864] info:    Centreon Broker 3.0.11
```

Dans cet exemple, on est sur une CentOS 7 avec un Centreon Broker (CBD) 3.0.11. Le module qui nous intéresse est donc `precompiled/Centos7/85-amqp-cbd-3.0.11.so`.

Il faut ensuite l'envoyer dans le répertoire d'installation des modules Centreon (attention : le nom `85-amqp.so` est attendu en destination) :
```shell
# À adapter en fonction du système cible !
$ sudo cp precompiled/Centos7/85-amqp-cbd-3.0.11.so /usr/share/centreon/lib/centreon-broker/85-amqp.so
```

### Installation de l'extension web

Il faut ensuite ajouter l'extension `centreon-extension/connector-centreon-canopsis` présente dans ce dépôt Git dans l'installation Centreon, afin de pouvoir finaliser l'installation du module.

```shell
$ sudo cp -r centreon-extension/connector-centreon-canopsis/ /usr/share/centreon/www/modules/
$ sudo chown -R apache:apache /usr/share/centreon/www/modules/connector-centreon-canopsis/
```

On peut alors installer le module depuis les menus suivants de l'interface web Centreon (Administration > Extensions > Modules > connector-centreon-canopsis et cliquer sur le bouton Action, sur la droite du tableau) :

![Installation du module depuis l'interface web : étape 1](img/webextension_install.png)

Puis, valider l'installation de ce module en cliquant sur « Install Module » :

![Installation du module depuis l'interface web : étape 2](img/webextension_install2.png)

## Configuration du module

### Configuration

**Attention :** la transmission de flux réseau entre Centreon et Canopsis doit être permise sur vos équipements réseau (port 5672 par défaut, à ajuster en fonction de votre configuration Canopsis/Centreon).

Aller dans Configuration > Pollers > Broker Configuration > central-broker-master.

![Configuration du module AMQP Canopsis : étape 1](img/module_parameters.png)

Puis, dans la nouvelle page qui apparaît, aller dans l'onglet Output, choisir le module « AMQP - Canopsis AMQP bus » dans le menu déroulant, et cliquer sur le bouton « Add ».

![Configuration du module AMQP Canopsis : étape 2](img/module_parameters2.png)

Des options de configuration « Canopsis AMQP bus » apparaissent alors en bas de page. Il faut alors renseigner les informations de connexion à l'instance AMQP Canopsis voulue (adresse, port, identifiants, nom de l'Exchange et du Virtual Host...). Valider ces changements avec le bouton « Save ».

![Configuration du module AMQP Canopsis : étape 3](img/module_parameters3.png)

### Redémarrage

**Attention :** les redémarrages suivants occasionnent une interruption de service le temps du redémarrage du Broker et des Engines Centron.

On redémarre ensuite le service pour s'assurer du bon chargement de la nouvelle configuration. Pour cela, aller dans Configuration > Pollers > cocher les éléments concernés > cliquer sur « Export configuration ».

![Redémarrage du service : étape 1](img/module_restart1.png)

Sur la nouvelle page qui s'affiche, il faut ensuite cocher les cases « Move Export Files » et « Restart Monitoring Engine », puis choisir la méthode « Restart » dans le menu déroulant, et enfin cliquer sur le bouton « Export ».

![Redémarrage du service : étape 2](img/module_restart1.png)

**ATTENTION :** Il faut bien faire un `restart` et non pas un simple `reload` ! Sans quoi vous risquez des problèmes de cohérence sur les évènements échangés avec Canopsis.

## Aller plus loin

[Compilation manuelle du module](Compilation-Manuelle-Module-Centreon.md) : peut être nécessaire s'il n'existe pas de binaire pré-compilé pour votre environnement.
