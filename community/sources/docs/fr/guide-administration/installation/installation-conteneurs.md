# Installation de Canopsis avec Docker Compose

Cette procédure décrit l'installation de Canopsis avec Docker Compose.

!!! Warning
    La présente documentation décrit un déploiement où tous les services
    (*backends* et applicatifs) seront exécutés dans des conteneurs Docker, le
    tout sur une seule machine.

    Ce type d'installation permet de monter très facilement un environnement de
    test reproductible et jetable.

    Cependant, **cette architecture n'est pas adaptée pour de la production**.
    En effet, en production, il serait recommandé d'installer les services
    *backends* – bases de données en particulier – hors de Docker, puis de
    mettre les conteneurs applicatifs sur *plusieurs* hôtes Docker afin d'avoir
    une certaine haute disponibilité (tolérance de panne).

    Le déploiement multi-nœuds fait l'objet d'un accompagnement spécifique par
    [Capensis][capensis] dans le cadre d'une souscription à l'édition Pro de
    Canopsis.

[capensis]: https://www.capensis.fr/

## Prérequis

### Utilisation de Docker Compose

[Docker Compose](https://docs.docker.com/compose/) est actuellement l'orchestrateur Docker à utiliser pour Canopsis.

!!! important
    Les conteneurs Docker et configurations Docker Compose fournies par Canopsis sont maintenues et testées seulement avec Docker Compose. La compatibilité directe avec d'autres outils d'orchestration tels que Kubernetes, Docker Swarm, Consul, OpenShift, etc. n'est à ce jour pas assurée.

### Prérequis de version du noyau Linux

Lors de l'utilisation de Docker, Canopsis nécessite **un noyau Linux 4.4 minimum sur votre système hôte**.

Vérifiez votre version du noyau à l'aide de la commande suivante :

```sh
uname -r
```

Si la version affichée est inférieure à 4.4, vous devez soit utiliser une distribution plus à jour, ou bien installer un noyau plus récent (par exemple *via* [ELRepo](https://elrepo.org/tiki/kernel-lt).

!!! important
    L'utilisation de Docker Compose avec un noyau inférieur à 4.4 n'est pas prise en charge.

## Installation de Docker et Docker Compose

Vous devez tout d'abord [installer Docker](https://docs.docker.com/engine/install/), version 20.10 minimum (se référer à page [Prérequis des versions][prereq-versions]). Veuillez utiliser les dépôts officiels de Docker, et non pas ceux proposés par votre distribution.

Veillez à installer Docker Compose à cette occasion, comme indiqué dans les instructions de la documentation d'installation officielle du Docker Engine. Le paquet se nomme `docker-compose-plugin`.

!!! information
    Dans le passé, [Docker Compose][docker-compose] était distribué sous forme de programme séparé (binaire `docker-compose` à télécharger et installer). Ce mode d'utilisation est aujourd'hui abandonné par Docker.

    Depuis plusieurs versions majeures de Docker Engine, Compose est intégré à la ligne de commande `docker` avec le *plugin* Compose.  
    Concrètement, on utilise à présent `docker compose <subcommand>...` là où les commandes `docker-compose <subcommand>...` étaient auparavant employées.

    Ces évolutions de Compose ont aussi été accompagnées de nouveaux formats et de nouvelles possibilités pour l'orchestrateur (voir [Compose Specification][compose-spec]).

    Concernant Canopsis, depuis la version 22.10 les environnements de référence fournis se basent sur ces dernières versions. L'utilisation du *plugin* Compose intégré à Docker (commandes `docker compose ...`) est maintenant la norme.

## Lancement de Canopsis avec Docker Compose

Les images Docker officielles de Canopsis sont hébergées sur leur propre registre Docker, `docker.canopsis.net`.

### Récupération de l'environnement Docker Compose

Les environnements Docker Compose de référence pour Canopsis sont disponibles via
git :

=== "Canopsis Pro"
    Pour Canopsis Pro, les fichiers sont dans le
    [dépôt git des sources de canopsis-pro][canopsis-pro-sources].

    Récupération du dépôt via Git+HTTPS :
    ```sh
    git clone https://git.canopsis.net/sources/canopsis-pro-sources.git
    ```
    Récupération du dépôt via Git+SSH :
    ```sh
    git clone git@git.canopsis.net:sources/canopsis-pro-sources.git
    ```

    Déplacez-vous ensuite dans le dossier contenant l'environnement :
    ```sh
    cd canopsis-pro-sources/pro/deployment/canopsis/docker
    ```

=== "Canopsis Community"
    Pour Canopsis Community, les fichiers sont dans le
    [dépôt git canopsis-community][canopsis-community].

    Récupération du dépôt via Git+HTTPS :
    ```sh
    git clone https://git.canopsis.net/canopsis/canopsis-community.git
    ```
    Récupération du dépôt via Git+SSH :
    ```sh
    git clone git@git.canopsis.net:canopsis/canopsis-community.git
    ```

    Déplacez-vous ensuite dans le dossier contenant l'environnement :
    ```
    cd canopsis-community/community/deployment/canopsis/docker
    ```

### Lancement de l'environnement

Récupérez les dernières images disponibles :

=== "Canopsis Pro"
    ```sh
    CPS_EDITION=pro docker compose pull
    ```

=== "Canopsis Community"
    ```sh
    CPS_EDITION=community docker compose pull
    ```

Lancez ensuite la commande suivante, afin de démarrer un environnement Canopsis
complet :

=== "Canopsis Pro"
    ```sh
    CPS_EDITION=pro docker compose up -d
    ```

=== "Canopsis Community"
    ```sh
    CPS_EDITION=community docker compose up -d
    ```

## Vérification du bon fonctionnement

=== "Canopsis Pro"
    ```sh
    CPS_EDITION=pro docker compose ps
    ```

=== "Canopsis Community"
    ```sh
    CPS_EDITION=community docker compose ps
    ```

Les services doivent être en état `Up`, `Up (healthy)` ou `Exit 0`. En fonction
des ressources de votre machine, il peut être nécessaire d'attendre quelques
minutes avant de voir l'ensemble des moteurs en état `Up`.

Vous pouvez ensuite procéder à votre [première connexion à l'interface Canopsis](premiere-connexion.md).

## Arrêt de l'environnement Docker Compose

=== "Canopsis Pro"
    ```sh
    CPS_EDITION=pro docker compose down
    ```

=== "Canopsis Community"
    ```sh
    CPS_EDITION=community docker compose down
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

Pour vérifier la validité de la configuration logrotate ajoutée, lancez la commande :

```sh
logrotate -dv /etc/logrotate.d/docker-container
```

Si vous souhaitez forcer une exécution manuelle de cette rotation sur-le-champ, vous pouvez éventuellement lancer la commande :

```sh
logrotate -fv /etc/logrotate.d/docker-container
```

[prereq-versions]: https://doc.canopsis.net/guide-administration/installation/prerequis-des-versions/#prerequis-systemes
[compose-spec]: https://docs.docker.com/compose/compose-file/
[docker-compose]: https://docs.docker.com/compose/install/#install-compose
[canopsis-pro-sources]: https://git.canopsis.net/sources/canopsis-pro-sources
[canopsis-community]: https://git.canopsis.net/canopsis/canopsis-community
