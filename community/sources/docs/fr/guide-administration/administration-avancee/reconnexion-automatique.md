# Reconnexion automatique des services et des moteurs

A partir de la version 4.1 de Canopsis, le fichier de configuration [`canopsis.toml`](./variables-environnement.md#chemin-dacces-au-fichier-de-configuration-global-canopsistoml) supporte de nouveaux paramètres permettant de configurer la reconnexion automatique en cas d'erreur.

## Configuration de la reconnexion

Les paramètres par défaut, présents dans le fichier `canopsis.toml` sont les suivants :

```ini
[Canopsis.global]
ReconnectRetries = 3
ReconnectTimeoutMilliseconds = 8
```

- `ReconnectRetries` représente le nombre de tentatives de reconnexion en cas d'erreur, 3 par défaut.
- `ReconnectTimeoutMilliseconds` est le délai **minimum** entre chaque tentative. Par défaut, il est de 8 millisecondes et on parle de délai minimum car celui-ci double à chaque tentative de reconnexion. Soit, avec la configuration par défaut, 8 ms avant le premier essai de reconnexion, 16 ms avant le second, 32 ms avant le troisième.

Ce mécanisme de reconnexion automatique est utilisé par MongoDB, Redis, RabbitMQ, les moteurs Canopsis ainsi que `canopsis-api`.

## MongoDB

Lors de l'exécution d'une commande, si une erreur de connexion est retournée, elle sera automatiquement exécutée de nouveau. Le nombre de tentatives d'exécution supplémentaires est égal à `ReconnectRetries`.

MongoDB dispose également d'options supplémentaires telles que `SocketTimeout` et `ServerSelectionTimeout` permettant de gérer les incidents de connexion. Reportez vous à la documentation MongoDB pour obtenir plus d'informations concernant l'utilisation de ces paramètres.

## Redis

Le fonctionnement est le même que pour MongoDB mais Redis n'a pas d'options internes pour gérer les incidents de connexion.

## RabbitMQ

Le comportement est identique à celui de Redis.

## Moteurs Canopsis

### [Processus périodique](../../../guide-utilisation/vocabulaire/#battement)

En cas d'incident de connexion, le processus exécute de nouveau la commande autant de fois que la valeur de `ReconnectRetries`. Si l'incident subsiste à l'issue des nouvelles tentatives, il inscrit l'erreur dans les logs et attend le prochain battement.

### Processus de travail

En cas d'erreur de connexion, le processus tente à nouveau d'exécuter la commande en fonction de la valeur de `ReconnectRetries`. Si l'incident persiste, il envoie un message de type `nack` à RabbitMQ, inscrit l'erreur dans les logs et arrête le moteur.

S'il s'agit d'une erreur d'un autre type, le processus envoie un message de type `ack` à RabbitMQ, inscrit l'erreur dans les logs et passe à la tâche suivante.

## Service `canopsis-api`

Ce service ne s'arrête jamais de fonctionner, quel que soit le type d'erreur rencontré. S'il s'agit d'une erreur de connexion, il essaie de se reconnecter indéfiniment.
