# Présentation

Vous trouverez ici toute la documentation nécessaire au développement sur canopsis.

!!! tip "Note"
    Cette page contient le plan de la documentation de développement, qui est en cours d'écriture.

# Process de développement
## Organisation des dépôts
## Process de release
## Nomenclature des messages de commit
<!--  - specification des segments de canopsis (alerts, action, …) -->

# Installation d'un environnement de développement
## Python
### VM
### LXC
## Python et Go
### Docker
### VM
### LXC

# Backend
## Python
### Installation de nouvelle source python

```bash
source ~/doc/pyvenv/canopsis/bin/activate
pip install -U ${HOME}/canopsis/sources/canopsis/
```

Cette technique est valide pour n’importe quel environnement.

Quand vous utiliserez LXC pour déployer vos développements, vous pourrez utiliser la même technique, à quelques détails près :

```bash
ssh root@<ip_conteneur>
chown canopsis:canopsis /opt/canopsis -R 
su - canopsis
pip install -U /mnt/canopsis/sources/canopsis
```

Le fait de faire `su - canopsis` va vous "connecter" à l’utilisateur `canopsis` et les fichiers `.bashrc` et `.bash_profile` sont déjà configurés pour *sourcer* le *venv* automatiquement.

##### Plus vite

Dans le cas où vous **savez** que aucune dépendance n’a été mise à jour, vous pouvez exécuter cette commande :

```bash
pip install -U --no-deps canopsis/sources/canopsis
```

Cela va vous faire gagner du temps.

### Structure du projet
<!--
  - organisation des packages
  - architecture à mettre en place : modele, adapter, api
-->
### Création d'engines
### Création d'API

# Golang
## Compilation

Cloner le projet :

```bash
git clone https://git.canopsis.net/canopsis/go-revolution.git -b develop $GOPATH/src/git.canopsis.net/canopsis/go-revolution
```

Initialiser le projet :

```bash
cd $GOPATH/src/git.canopsis.net/canopsis/go-revolution/
make init
```

Lancer le build :

```bash
cd $GOPATH/src/git.canopsis.net/canopsis/go-revolution/
make
```

Lors du développement, il peut être utile de *builder* les binaires rapidement et des les récupérer dans un dossier partagé avec une autre machine par exemple :

```
make init

# réutiliser cette commande par la suite
make build BUILD_OUTPUT_DIR=/vmshare/gobin SKIP_DEPENDENCIES=true
```

Ouvrez les ports de Mongo, Rabbit, Redis et Influx de la machine/vm où se trouve votre canopsis, bindez les adresses dans les différents fichiers de configuration, notemment pour Mongo et Redis, à `0.0.0.0`, puis exportez ces valeurs, en prenant bien soin de changer les couples `host:port` par ceux de votre machine/vm canopsis.


```bash
export CPS_AMQP_URL="amqp://cpsrabbit:canopsis@host:post/canopsis"
export CPS_MONGO_URL="mongodb://cpsmongo:canopsis@host:port/canopsis"
export CPS_REDIS_URL="redis://nouser:dbpassword@host:port/0"
export CPS_INFLUX_URL="http://cpsinflux:canopsis@host:port"
export CPS_DEFAULT_CFG="$GOPATH/src/git.canopsis.net/canopsis/go-revolution/default_configuration.toml"
```

Repris en partie du readme go-revolution :  

```bash
cd $GOPATH/src/git.canopsis.net/canopsis/go-revolution/
make init
make build
cd build
./init -conf ../cmd/init/initialisation.toml
```

`init` devrait soulever tous les problèmes de ports fermés, adresses mal bindées et vous permettre de débloquer la situation.

Ensuite, lancez votre engine, par exemple : 

```bash
./engine-axe
```

## Architecture du projet
## Création de moteurs
## Implémentation de source de données externes (pour l'event-filter)

# Base de données
<!--
## default_entities
### Présentation générale
### Présentation de la structure d'un document.
## periodical_alarms
### Présentation générale
### Présentation de la structure d'un document.
-->

# Front-end
## Mise en place de l'environnement de développement
## Technologies utilisées
## Structure du projet
## Règles de style
## Les mixins, helpers et filters
## Le store Vuex
## Guides de création nouvelle fonctionnalité
### Modal
### Vue
### Widget (+ Paramètres du widget)

# API

[Présentation de toutes les routes disponibles](API.md)

  * Pbehavior
  * Event-filter
  * [Healthcheck](./healthcheck/api_v2_healthcheck.md)

# Gestion de la documentation
