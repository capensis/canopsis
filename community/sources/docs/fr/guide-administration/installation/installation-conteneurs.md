# Installation de Canopsis avec Docker Compose

Cette procédure décrit l'installation de Canopsis avec Docker Compose.

## Prérequis

### Utilisation de Docker Compose

[Docker Compose](https://docs.docker.com/compose/) est actuellement l'orchestrateur Docker à utiliser pour Canopsis.

!!! important
    Les conteneurs Docker produits pour Canopsis ne sont pas garantis d'être compatibles avec un autre orchestrateur que Docker Compose. La compatibilité avec d'autres outils tels que Kubernetes, Docker Swarm, Consul, OpenShift, etc. n'est ainsi pas assurée.

### Prérequis de version du noyau Linux

Lors de l'utilisation de Docker, Canopsis nécessite **un noyau Linux 4.4 minimum sur votre système hôte**.

Vérifiez votre version du noyau à l'aide de la commande suivante :
```sh
uname -r
```

Si la version affichée est inférieure à 4.4, vous devez soit utiliser une distribution plus à jour, ou bien mettre à jour votre noyau à l'aide d'[ELRepo](https://elrepo.org/tiki/kernel-lt) pour CentOS, par exemple.

!!! important
    L'utilisation de Docker Compose avec un noyau inférieur à 4.4 n'est pas prise en charge.

## Installation de Docker et Docker Compose

Vous devez tout d'abord [installer Docker](https://docs.docker.com/get-docker/), version 20.10 minimum. Veuillez utiliser les dépôts officiels de Docker, et non pas ceux proposés par votre distribution.

Une fois Docker installé, vous devez ensuite [installer Docker Compose](https://docs.docker.com/compose/install/#install-compose).

!!! attention
    Dans certaines configurations, les versions les plus récentes de Docker Compose peuvent activer Compose v2 par défaut.

    Compose v2 est une réécriture importante de Docker Compose mais elle est, à ce jour, en partie incomplète et instable par rapport à Compose v1. Canopsis ne prend donc pas en charge Compose v2 pour le moment.

    Si la commande `docker-compose version --short` vous renvoie un numéro de version supérieur ou égal à `2.0.0`, vous devez désactiver Compose v2 avec la commande `docker-compose disable-v2`. Voyez [la documentation de Compose V2](https://docs.docker.com/compose/cli-command/#compose-v2-and-the-new-docker-compose-command) pour en savoir plus.

## Lancement de Canopsis avec Docker Compose

Les images Docker officielles de Canopsis sont hébergées sur leur propre registre Docker, `docker.canopsis.net`.

### Récupération de l'environnement Docker-compose

Les environnements docker compose de référence pour Canopsis sont disponible via
git :

=== "Canopsis Pro"
	Pour Canopsis Pro, elle sont dans le [dépôt git dédié a celui-ci
	](https://git.canopsis.net/sources/canopsis-pro-sources).

	Récupération du dépôt via Git+HTTPS :
	```
	git clone https://git.canopsis.net/sources/canopsis-pro-sources.git
	```
	Récupération du dépôt via Git+SSH:
	```
	git clone git@git.canopsis.net:sources/canopsis-pro-sources.git
	```

	Déplacez vous ensuite dans le dossier contenant l'environement:
	```
	cd canopsis-pro-sources/pro/deployment/canopsis/docker
	```

=== "Canopsis Community"
	Pour Canopsis Community, elle sont dans le [dépôt git dédié a celui-ci
	](https://git.canopsis.net/canopsis/canopsis-community) :
	```
	git clone https://git.canopsis.net/canopsis/canopsis-community.git
	```

	Déplacez vous ensuite dans le dossier contenant l'environement:
	```
	cd canopsis-community/community/deployment/canopsis/docker
	```

### Lancement de l'environnement

Récupérez les dernières images disponibles :
=== "Canopsis Pro"
	```sh
	CPS_EDITION=pro docker-compose pull
	```

=== "Canopsis Community"
	```sh
	CPS_EDITION=community docker-compose pull
	```


Lancez ensuite la commande suivante, afin de démarrer un environnement Canopsis
complet :
=== "Canopsis Pro"
	```sh
	CPS_EDITION=pro docker-compose up -d
	```

=== "Canopsis Community"
	```sh
	CPS_EDITION=community docker-compose up -d
	```

## Vérification du bon fonctionnement

=== "Canopsis Pro"
	```sh
	CPS_EDITION=pro docker-compose ps
	```

=== "Canopsis Community"
	```sh
	CPS_EDITION=community docker-compose ps
	```
Les services doivent être en état `Up`, `Up (healthy)` ou `Exit 0`. En fonction
des ressources de votre machine, il peut être nécessaire d'attendre quelques
minutes avant que l'ensemble des moteurs puissent passer en état `Up`.

Vous pouvez alors procéder à votre [première connexion à l'interface Canopsis](premiere-connexion.md).

## Arrêt de l'environnement Docker Compose

=== "Canopsis Pro"
	```sh
	CPS_EDITION=pro docker-compose down
	```

=== "Canopsis Community"
	```sh
	CPS_EDITION=community docker-compose down
	```


## Rétention des logs

La mise en place d'une politique de rétention des logs nécessite la présence du logiciel `logrotate`.

Une fois que `logrotate` est installé sur votre machine, créer le fichier `/etc/logrotate.d/docker-container` suivant :

```
/var/lib/docker/containers/*/*.log {
  rotate 7
  daily
  compress
  minsize 100M
  notifempty
  missingok
  delaycompress
  copytruncate
}
```

Pour vérifier la bonne exécution de la configuration de logrotate pour Docker, vous pouvez lancer la commande :

```sh
logrotate -fv /etc/logrotate.d/docker-container
```
