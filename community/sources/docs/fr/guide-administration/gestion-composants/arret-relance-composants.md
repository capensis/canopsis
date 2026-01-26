# Arrêt et relance des composants de Canopsis

L'ensemble des commandes suivantes doit être réalisées avec l'utilisateur `root` du système.

## Gestion des composants Canopsis

La commande suivante redémarrera Canopsis :

```sh
systemctl restart canopsis
```

Les actions `start` et `stop` sont aussi disponibles, afin de respectivement démarrer et arrêter Canopsis.

## Gestion des composants liés à Canopsis

L'ensemble des composants liés à Canopsis peuvent être gérés avec la commande `systemctl` usuelle sous Linux.

### MongoDB

La base de données MongoDB peut être redémarrée avec la commande suivante :

```sh
systemctl restart mongod.service
```

### RabbitMQ

L'agent de messages RabbitMQ peut être redémarré avec la commande suivante :

```sh
systemctl restart rabbitmq-server.service
```

### Redis

Le serveur de cache Redis peut être redémarré avec la commande suivante :

```sh
systemctl restart redis.service
```

### PostgreSQL / TimescaleDB

Le serveur de cache Redis peut être redémarré avec la commande suivante :

```sh
systemctl restart postgresql-13.service
```

Veuillez noter qu'un redémarrage du service n'occasionnera pas une purge du cache existant. Ce comportement est intentionnel.

## Aller plus loin 

Pour connaître l'état de la plateforme, [rendez-vous ici](../../guide-de-depannage/etat-des-composants/index.md).
