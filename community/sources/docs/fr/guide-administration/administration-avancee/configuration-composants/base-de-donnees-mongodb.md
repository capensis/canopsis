# Configuration avancée de la base de données MongoDB intégrée à Canopsis

La base de données [MongoDB](https://www.mongodb.com) contient la plupart des données, les vues et la configuration de l'interface Canopsis.

## Recommandations d'utilisation avancée

Par défaut, la commande `canoctl deploy` installe MongoDB sur la même machine que le reste des composants de Canopsis. En dehors d'une utilisation pour un périmètre réduit, il est recommandé que MongoDB soit installé sur une machine virtuelle ou physique dédiée à cette tâche.

La [mise en place d'un *replica set*](https://docs.mongodb.com/manual/replication/) peut vous permettre d'apporter une redondance et une disponibilité supplémentaires à votre base de données. Elle permet aussi de bénéficier [des transactions](https://docs.mongodb.com/manual/core/transactions/) (depuis Canopsis 4.5.0) et donc d'une meilleure sûreté des données de production.

Canopsis ne propose ces configurations que par le biais de certaines souscriptions Canopsis Pro.

## Optimisations système

Si vous constatez que votre utilisation de MongoDB est importante dans votre environnement Canopsis, et si MongoDB a déjà été déplacé sur une instance dédiée, il peut être utile d'appliquer les recommandations officielles de MongoDB sur [les limites système](https://docs.mongodb.com/v4.2/reference/ulimit/#linux-distributions-using-systemd) et [les Transparent Huge Pages](https://docs.mongodb.com/v4.2/tutorial/transparent-huge-pages/).

Par défaut, MongoDB recommande aussi l'utilisation du système de fichiers XFS pour le dossier `/var/lib/mongodb`. Il s'agit du système de fichiers par défaut de CentOS.

## ReadPreference

https://www.mongodb.com/docs/manual/core/read-preference/

Nous avons désactivé la possibilité de surcharger le [readPreference](https://www.mongodb.com/docs/manual/core/read-preference/) via l’URL de connexion Mongo.  
Le readPreference est maintenant défini par l’application elle-même et non modifiable depuis l’extérieur.

Par défaut, toutes les lectures sont effectuées depuis le noeud primaire.

Des exceptions existent dans certains composants, où la lecture depuis les secondaires est explicitement activée.

Cette décision a été prise pour éviter les erreurs fréquentes liées à un usage incorrect ou abusif du readPreference.  
Cela causait des anomalies difficiles à diagnostiquer, comme des erreurs du type : "mongo: no documents in result."

**Composants utilisant les lectures sur noeuds secondaires**

Les lectures depuis les réplicas secondaires sont actuellement activées dans :

* L'API des alarmes
* L'API d’export des alarmes
* Les moteurs engine-che et engine-dynamic-infos (dictionnaires d’info)
* Le Prometheus exporter

Cela permet de répartir la charge de manière plus équilibrée.
