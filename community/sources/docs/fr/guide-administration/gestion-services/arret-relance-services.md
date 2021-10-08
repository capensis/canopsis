# Arrêt et relance des services liés à Canopsis

L'ensemble des commandes suivantes doit être réalisées avec l'utilisateur `root` du système.

## Gestion des services Canopsis

L'utilitaire `canoctl` permet, entre autres, de redémarrer Canopsis en lui-même. Il ne s'applique pas aux services associés, tels que RabbitMQ ou MongoDB.

La commande suivante redémarrera Canopsis, ses services et ses moteurs :

```sh
canoctl restart
```

Comme avec `systemctl`, les actions `start` et `stop` sont aussi disponibles, afin de respectivement démarrer et arrêter Canopsis.

## Gestion des services liés à Canopsis

L'ensemble des services liés à Canopsis peuvent être gérés avec la commande `systemctl` usuelle sous Linux.

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

Veuillez noter qu'un redémarrage du service n'occasionnera pas une purge du cache existant. Ce comportement est intentionnel.

## Aller plus loin 

Pour connaître l'état de votre service, [rendez-vous ici](../../guide-de-depannage/etat-des-services.md).
