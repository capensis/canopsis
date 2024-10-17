# Reconnexion automatique des services et des moteurs

À partir de Canopsis 4.1.0, le fichier de configuration [`canopsis.toml`](modification-canopsis-toml.md) prend en charge de nouveaux paramètres permettant de configurer la reconnexion automatique en cas d'erreur.

## Configuration de la reconnexion

Les paramètres présents par défaut dans le fichier `canopsis.toml` pour la reconnexion automatique sont les suivants :

```ini
[Canopsis.global]
ReconnectRetries = 3
ReconnectTimeoutMilliseconds = 8
```

- `ReconnectRetries` représente le nombre de tentatives de reconnexion en cas d'erreur, 3 par défaut.
- `ReconnectTimeoutMilliseconds` est le délai **minimum** entre chaque tentative. Par défaut, il est de 8 millisecondes et on parle de délai minimum, car celui-ci double à chaque tentative de reconnexion. Soit, avec la configuration par défaut, 8 ms avant le premier essai de reconnexion, 16 ms avant le second, 32 ms avant le troisième.

Ce mécanisme de reconnexion automatique est utilisé par MongoDB, Redis, RabbitMQ, les moteurs Canopsis ainsi que `canopsis-api`.

Ce mécanisme de reconnexion automatique est utilisé par MongoDB, Redis, RabbitMQ, les moteurs Canopsis ainsi que `canopsis-api`.

!!! note
    Toute modification d'une de ces valeurs implique de suivre de le [Guide de modification du fichier `canopsis.toml`](modification-canopsis-toml.md).

## Perte de connexion à MongoDB

Lors de l'exécution d'une commande, si une erreur de connexion est reçue, la commande sera automatiquement exécutée de nouveau. Le nombre de tentatives d'exécution supplémentaires est égal à `ReconnectRetries`.

MongoDB dispose également d'options supplémentaires telles que `SocketTimeout` et `ServerSelectionTimeout` permettant de gérer les incidents de connexion. Reportez-vous à la documentation MongoDB pour obtenir plus d'informations concernant l'utilisation de ces paramètres.

## Perte de connexion à Redis

Le fonctionnement est le même que pour MongoDB mais Redis n'a pas d'options internes pour gérer les incidents de connexion.

## Perte de connexion à RabbitMQ

Le comportement est identique à celui de Redis.

## Comportement de la reconnexion dans les moteurs et services Canopsis

### Processus périodique

En cas d'incident de connexion, le processus exécute de nouveau la commande autant de fois que la valeur de `ReconnectRetries`. Si l'incident subsiste à l'issue des nouvelles tentatives, il inscrit l'erreur dans les logs et attend le prochain [battement](../../guide-utilisation/vocabulaire/index.md#battement).

### Processus de travail

En cas d'erreur de connexion, le processus tente de nouveau d'exécuter la commande en fonction de la valeur de `ReconnectRetries`. Si l'incident persiste, il envoie un message de type `nack` à RabbitMQ, inscrit l'erreur dans les logs et arrête le moteur.

S'il s'agit d'une erreur d'un autre type, le processus envoie un message de type `ack` à RabbitMQ, inscrit l'erreur dans les logs et passe à la tâche suivante.

## Service `canopsis-api`

Ce service ne s'arrête jamais de fonctionner, quel que soit le type d'erreur rencontré. S'il s'agit d'une erreur de connexion, il essaie de se reconnecter indéfiniment.
