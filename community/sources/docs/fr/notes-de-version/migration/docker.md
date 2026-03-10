# Procédure de migration dans un environnement Docker Compose

## Sommaire

[TOC]

## Vérification MongoDB

!!! warning "Vérification"

    Avant de démarrer la procédure de mise à jour, vous devez vérifier que la valeur de `featureCompatibilityVersion` est bien positionnée à **7.0**  

```sh
CPS_EDITION=pro docker compose exec mongodb bash
mongosh -u root -p root
> db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )
> exit
```

## Arrêt de l'environnement en cours d'exécution

Vous devez prévoir une interruption du service afin de procéder à la mise à jour qui va suivre.

```sh
CPS_EDITION=pro docker compose down
```

## Mise à jour de l'applicatif Canopsis

!!! information "Information"

    Canopsis 25.10 est livré avec un nouveau jeu de configurations de référence.
    Vous devez télécharger ces configurations et y reporter vos personnalisations. 

Si vous êtes utilisateur de l'édition `community`, voici les étapes à suivre.

Télécharger le paquet de la version 25.10.0 (canopsis-community-docker-compose-25.10.0.tar.gz) disponible à cette adresse [https://git.canopsis.net/canopsis/canopsis-community/-/releases](https://git.canopsis.net/canopsis/canopsis-community/-/releases).

```sh
export CPS_EDITION=community
tar xvfz canopsis-community-docker-compose-25.10.0.tar.gz
cd canopsis-community-docker-compose-25.10.0
```

Si vous êtes utilisateur de l'édition `pro`, voici les étapes à suivre.

Télécharger le paquet de la version 25.10.0 (canopsis-pro-docker-compose-25.10.0.tar.gz) disponible à cette adresse [https://git.canopsis.net/sources/canopsis-pro-sources/-/releases](https://git.canopsis.net/sources/canopsis-pro-sources/-/releases).

```sh
export CPS_EDITION=pro
tar xvfz canopsis-pro-docker-compose-25.10.0.tar.gz
cd canopsis-pro-docker-compose-25.10.0
```

À ce stade, vous devez synchroniser les modifications réalisées sur vos anciens fichiers de configuration `docker-compose` avec les fichiers `docker-compose.yml` et/ou `docker-compose.override.yml`.

## Mise à jour de TimescaleDB

Dans cette version de Canopsis, la base de données TimescaleDB passe de la version 2.15.1 à 2.21.4.  
En plus de la mise à jour de TimescaleDB lui-même, le système de gestion de base de données PostreSQL doit être mis à jour de la version 15 à la version 17.

Deux étapes sont à suivre :

1. Mise à jour de TimescaleDB 2.15.1 vers 2.21.4
2. Mise à jour de PostgreSQL 15 vers 17

Modifiez la variable `TIMESCALEDB_TAG` du fichier `.env` de cette façon :

```diff
-TIMESCALEDB_TAG=2.21.4-pg17
+TIMESCALEDB_TAG=2.21.4-pg15
```

Démarrez le conteneur et mettez à jour l'extension TimescaleDB

```sh
CPS_EDITION=pro docker compose up -d timescaledb
CPS_EDITION=pro docker compose exec timescaledb psql postgresql://postgres:canopsis@timescaledb:5432/canopsis
canopsis=# ALTER EXTENSION timescaledb UPDATE;
canopsis=# \c canopsis_tech_metrics
canopsis=# ALTER EXTENSION timescaledb UPDATE;
```

Vérifiez que la version de l'extension soit bien mise à jour

```sh
canopsis=# \dx timescaledb
                                                List of installed extensions
    Name     | Version | Schema |                                      Description                                      
-------------+---------+--------+---------------------------------------------------------------------------------------
    timescaledb | 2.21.4  | public | Enables scalable inserts and complex queries for time-series data (Community Edition)
(1 row)
```

Sauvegarde des bases de données

```sh
CPS_EDITION=pro docker compose exec timescaledb pg_dump postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis -Ft -f /tmp/postgres_dump_archive.tar
CPS_EDITION=pro docker compose exec timescaledb pg_dump postgresql://cpspostgres_tech_metrics:canopsis@timescaledb:5432/canopsis_tech_metrics -Ft -f /tmp/postgres_dump_archive_techmetrics.tar
CPS_EDITION=pro docker compose cp timescaledb:/tmp/postgres_dump_archive.tar /tmp
CPS_EDITION=pro docker compose cp timescaledb:/tmp/postgres_dump_archive_techmetrics.tar /tmp
```

Arrêtez le conteneur et supprimez les volumes associés

```sh
CPS_EDITION=pro docker compose down -v timescaledb
```

Modifiez la variable `TIMESCALEDB_TAG` du fichier `.env` de cette façon :

```diff
-TIMESCALEDB_TAG=2.21.4-pg15
+TIMESCALEDB_TAG=2.21.4-pg17
```

Démarrer le conteneur timescaledb

```sh
CPS_EDITION=pro docker compose up -d timescaledb
```

Restaurez le dump précédemment effectué

```sh
CPS_EDITION=pro docker compose cp /tmp/postgres_dump_archive.tar timescaledb:/tmp/postgres_dump_archive.tar
CPS_EDITION=pro docker compose cp /tmp/postgres_dump_archive_techmetrics.tar timescaledb:/tmp/postgres_dump_archive_techmetrics.tar
CPS_EDITION=pro docker compose exec timescaledb pg_restore --dbname=postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis --no-owner -Ft -v /tmp/postgres_dump_archive.tar
CPS_EDITION=pro docker compose exec timescaledb pg_restore --dbname=postgresql://cpspostgres_tech_metrics:canopsis@timescaledb:5432/canopsis_tech_metrics --no-owner -Ft -v /tmp/postgres_dump_archive_techmetrics.tar
```

## Mise à jour de MongoDB

Dans cette version de Canopsis, la base de données MongoDB passe de la version 7.0 à 8.0.  

Démarrez le conteneur `mongodb` :

```sh
CPS_EDITION=pro docker compose up -d mongodb
```

Entrez ensuite à l'intérieur de ce conteneur, afin de compléter la mise à jour vers MongoDB 8.0 :

```sh
CPS_EDITION=pro docker compose exec mongodb bash
mongosh -u root -p root
> db.adminCommand( { setFeatureCompatibilityVersion: "8.0", "confirm" : true } )
```

Après avoir mis à jour mongodb, l'option de telemetry sera activée. Pour la désactiver, exécutez la commande suivante :

```sh
CPS_EDITION=pro docker compose exec mongodb bash
mongosh -u root -p root
> disableTelemetry()
exit
```

## Mise à jour de RabbitMQ

Dans cette version de Canopsis, le bus de données RabbitMQ passe de la version 4.0 à 4.1.  

Les configurations de référence fournies dans Canopsis embarquent déjà ce changement.  
Vous devez néanmoins activer l'ensemble des ["feature flags"](https://www.rabbitmq.com/docs/feature-flags) stables.  

Démarrez le conteneur `rabbitmq` :

```sh
CPS_EDITION=pro docker compose up -d rabbitmq
```

Puis activez les "feature flags" :

```sh
CPS_EDITION=pro docker compose exec rabbitmq rabbitmqctl enable_feature_flag all
```

## Lancement du provisioning `canopsis-reconfigure`

### Synchronisation du fichier de configuration `canopsis.toml` ou fichier de surcharge

Si vous avez modifié le fichier `canopsis.toml` (vous le voyez via une définition de volume dans votre fichier docker-compose.yml), vous devez vérifier qu'il soit bien à jour par rapport au fichier de référence.  

* [`canopsis.toml` pour Canopsis Community 25.10.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/25.10.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml)
* [`canopsis.toml` pour Canopsis Pro 25.10.0](https://git.canopsis.net/canopsis/canopsis-community/-/blob/25.10.0/community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml)

!!! information "Information"

    Pour éviter ce type de synchronisation fastidieuse, la bonne pratique est d'utiliser [un fichier de surcharge de cette configuration](../../../guide-administration/administration-avancee/modification-canopsis-toml/). 

Si vous avez utilisé un fichier de surcharge, alors vous n'avez rien à faire, uniquement continuer à le présenter dans un volume.

### Séparation des flux d’événements par initiateur

En version 25.04, les flags suivants avaient été dépréciés, ils sont à présent obsolètes.  
Toute référence doit être supprimée dans vos configurations. Les moteurs concernés ne démareront pas sans cela.  

* -publishQueue
* -consumeQueue
* -workers (remplacé par des flags spécifiques à chaque type de flux)

| Moteur               | Flags nouveaux     | Valeur par défaut | Flags obsolètes                             |
|----------------------|--------------------|-------------------|---------------------------------------------|
| engine-fifo          | -workers           | 10                | -publishQueue, -consumeQueue                |
| engine-che           | -externalWorkers   | 4                 | -workers, -publishQueue, -consumeQueue      |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
| engine-axe           | -externalWorkers   | 4                 | -workers, -publishQueue                     |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |
| engine-correlation   | -externalWorkers   | 4                 | -workers, -publishQueue, -consumeQueue      |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |
| engine-dynamic-infos | -externalWorkers   | 4                 | -workers, -publishQueue                     |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |
| engine-action        | -externalWorkers   | 4                 | -workers                                    |
|                      | -systemWorkers     | 4                 |                                             |
|                      | -userWorkers       | 2                 |                                             |
|                      | -rpcWorkers        | 4                 |                                             |

### Reconfiguration de Canopsis

!!! Attention

    Si vous avez personnalisé la ligne de commande de l'outil `canopsis-reconfigure`, nous vous conseillons de supprimer cette personnalisation.
    L'outil est en effet pré paramétré pour fonctionner naturellement.

```sh
CPS_EDITION=pro docker compose up -d reconfigure
```

!!! information "Information"

    Cette opération peut prendre plusieurs minutes pour s'exécuter.

Vous pouvez ensuite vérifier que le mécanisme de provisioning/reconfigure s'est correctement déroulé. Le conteneur doit présenté un "exit 0"

```sh
CPS_EDITION=pro docker compose ps -a|grep reconfigure
canopsis-pro-reconfigure-1            "/canopsis-reconfigu…"   reconfigure            exited (0)
```

## Mise à jour et démarrage final de Canopsis

Enfin, il vous reste à mettre à jour et à démarrer tous les composants applicatifs de Canopsis

```sh
CPS_EDITION=pro docker compose up -d
```

Vous pouvez ensuite vérifier que l'ensemble des conteneurs soient correctement exécutés.

```sh
CPS_EDITION=pro docker compose ps
```

Par ailleurs, le mécanisme de bilan de santé intégré à Canopsis ne doit pas présenter d'erreur.  

![Healthcheck](./img/25.10.0-healthcheck.png)