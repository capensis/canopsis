# Configuration avancée de la base de données MongoDB intégrée à Canopsis

La base de données [MongoDB](https://www.mongodb.com) contient l'ensemble du contenu, des vues et de la configuration de Canopsis.

## Recommandations d'utilisation avancée

Par défaut, la commande `canoctl deploy` installe MongoDB sur la même machine que le reste des composants de Canopsis. En dehors d'une utilisation pour un périmètre réduit, il est recommandé que MongoDB soit installé sur une machine virtuelle ou physique dédiée à cette tâche.

La [mise en place d'un *replica set*](https://docs.mongodb.com/manual/replication/) peut vous permettre d'apporter une redondance et une disponibilité supplémentaires à votre base de données.

Canopsis ne propose ces configurations que par le biais de certaines souscriptions Canopsis Pro.

## Optimisations système

Si vous constatez que votre utilisation de MongoDB est importante dans votre environnement Canopsis, et si MongoDB a déjà été déplacé sur une instance dédiée, il peut être utile d'appliquer les recommandations officielles de MongoDB sur [les limites système](https://docs.mongodb.com/v4.2/reference/ulimit/#linux-distributions-using-systemd) et [les Transparent Huge Pages](https://docs.mongodb.com/v4.2/tutorial/transparent-huge-pages/).

Par défaut, MongoDB recommande aussi l'utilisation du système de fichiers XFS pour le dossier `/var/lib/mongodb`. Il s'agit du système de fichiers par défaut de CentOS.
