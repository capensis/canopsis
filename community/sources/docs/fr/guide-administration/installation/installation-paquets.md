# Installation de paquets Canopsis sur CentOS 7

Cette procédure décrit l'installation de Canopsis en mono-instance à l'aide de paquets CentOS 7. Les binaires sont compilés pour l'architecture x86-64.

L'ensemble des commandes suivantes doivent être réalisées avec l'utilisateur `root`.

## Prérequis

Assurez-vous d'avoir suivi les [prérequis réseau et de sécurité](../administration-avancee/configuration-parefeu-et-selinux.md), notamment concernant la désactivation de SELinux.

L'installation nécessite l'ajout de dépôts RPM tiers, ainsi qu'un accès HTTP et HTTPS pour le téléchargement de diverses dépendances. Plus de détails dans la [matrice des flux réseau](../matrice-des-flux-reseau/index.md).

!!! information
    Notez que Canopsis installe ses propres versions de MongoDB, Nginx, Redis et RabbitMQ, et que seules ces versions sont validées pour fonctionner avec Canopsis. Veillez à ne pas remplacer ces versions par vos propres versions, que ce soit de façon intentionnelle, ou par l'ajout de dépôt tiers qui pourraient écraser celles fournies par Canopsis.

    Plus de détails sur les [prérequis des versions](prerequis-des-versions.md).

## Ajout de dépôts tiers et de dépendances

Activation d'EPEL et application des dernières mises à jour du système :
```sh
yum install yum-utils epel-release
yum update
```

Activation de SCL et installation d'une version corrigée de Python 2.7 (pour quelques composants historiques de Canopsis nécessitant cette version) :
```sh
yum install centos-release-scl
yum install python27
```

## Installation de Canopsis Community ou Pro

Canopsis est disponible dans une édition « Community », open-source et gratuitement accessible à tous, et une édition « Pro », souscription commerciale ajoutant des fonctionnalités supplémentaires. Voyez [le site officiel de Canopsis](https://www.capensis.fr/canopsis/) pour en savoir plus.

Notez que l'édition Pro de Canopsis était auparavant connue sous le nom de « CAT » et que certains éléments peuvent encore la désigner sous ce nom.

Cliquez sur l'un des onglets « Community » ou « Pro » suivants, en fonction de l'édition choisie.

=== "Canopsis Community (édition open-source)"

    Ajout du dépôt de paquets Canopsis pour CentOS 7 :
    ```sh
    echo "[canopsis]
    name = canopsis
    baseurl=https://repositories.canopsis.net/pulp/repos/centos7-canopsis4/
    gpgcheck=0
    enabled=1" > /etc/yum.repos.d/canopsis.repo
    ```

    Installation de l'édition open-source de Canopsis :
    ```sh
    yum install canopsis-core
    ```

=== "Canopsis Pro (souscription commerciale)"

    !!! attention
        L'édition Pro nécessite une souscription commerciale, ainsi que d'une demande d'autorisation d'accès au dépôt `canopsis-cat`.

    Ajout des dépôts de paquets Canopsis pour CentOS 7 :
    ```sh
    echo "[canopsis]
    name = canopsis
    baseurl=https://repositories.canopsis.net/pulp/repos/centos7-canopsis4/
    gpgcheck=0
    enabled=1" > /etc/yum.repos.d/canopsis.repo

    echo "[canopsis-cat]
    name = canopsis-cat
    baseurl=https://repositories.canopsis.net/pulp/repos/centos7-canopsis4-cat/
    gpgcheck=0
    enabled=1" > /etc/yum.repos.d/canopsis-cat.repo
    ```

    Installation de Canopsis Pro :
    ```sh
    yum install canopsis-cat
    ```

## Initialisation de Canopsis

Vous devez ensuite initialiser l'environnement Canopsis à l'aide de la commande suivante (qui procédera à l'installation des dépendances, la finalisation des fichiers de configuration, l'activation des moteurs…).

```sh
canoctl deploy
```

Cette commande peut prendre quelques minutes. Elle ne doit être exécutée qu'à la **première installation** de Canopsis.

Une fois cette commande terminée, vous pouvez alors réaliser votre [première connexion à l'interface Canopsis](premiere-connexion.md). 

Si vous souhaitez réaliser une mise à jour, la procédure est décrite dans le [Guide de mise à jour](../mise-a-jour/index.md).
