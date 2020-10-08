# Variables d'environnement Canopsis

Ce document regroupe l'ensemble des variables d'environnement pouvant être ajustées à vos besoins, pour les différents composants de Canopsis.

!!! note
    À ce jour, ce document est encore en cours de rédaction, et n'est pas exhaustif. Il sera complété au fur et à mesure des mises à jour de la documentation.

## Modification des variables d'environnement

En installation Docker, les variables suivantes peuvent être redéfinies dans un fichier `compose.env` ou dans une section `environment:` de votre Docker Compose.

En installation par paquets, ces variables peuvent être redéfinies dans les unités systemd ou dans le shell exécutant votre commande.

## Liste des variables d'environnement

### URI de connexion aux services externes

Votre installation de Canopsis doit obligatoirement comporter les adresses et données de connexion (on parle d'URI) permettant de se connecter aux services externes Redis, MongoDB, RabbitMQ et InfluxDB.

En installation Docker, ces variables sont généralement définies dans le fichier `compose.env`. En installation par paquets, elles sont définies dans le fichier `/opt/canopsis/etc/go-engines-vars.conf`.

Ces variables concernent l'ensemble des moteurs Go, et certains binaires comme `init` ou `feeder`.

| Variable d'environnement | Valeur par défaut | Utilité |
|:-------------------------|-------------------|---------|
| `CPS_AMQP_URL` | (vide) | Une URI de connexion RabbitMQ (cf. [Spécification d'URI RabbitMQ](https://www.rabbitmq.com/uri-spec.html)) |
| `CPS_INFLUX_URL` | (vide) | Une URI de connexion InfluxDB (pas de spécification officielle) |
| `CPS_MONGO_URL` | (vide) | Une URI de connexion MongoDB (cf. [Spécification d'URI MongoDB](https://docs.mongodb.com/v3.6/reference/connection-string/)) |
| `CPS_REDIS_URL` | (vide) | Une URI de connexion Redis (cf. [Spécification d'URI Redis](https://www.iana.org/assignments/uri-schemes/prov/redis)) |

### Chemin d'accès au fichier de configuration global (`default_configuration.toml`)

Les différents moteurs et binaires Go ont besoin d'un fichier de configuration `default_configuration.toml`. La variable `CPS_DEFAULT_CFG` permet de leur indiquer le chemin où se trouve ce fichier.

En installation Docker, elle doit presque toujours valoir `/default_configuration.toml`. En installation par paquets, elle doit valoir `/opt/canopsis/etc/default_configuration.toml` (`canoctl deploy` se charge de renseigner cette valeur dans l'unité systemd `canopsis-engine-go@.service`).

Ces variables concernent l'ensemble des moteurs Go, et certains binaires comme `init` ou `feeder`.

| Variable d'environnement | Valeur par défaut | Utilité |
|:-------------------------|-------------------|---------|
| `CPS_DEFAULT_CFG` | `default_configuration.toml` (dans le répertoire courant) | Chemin d'accès vers le fichier de configuration `default_configuration.toml` |

### Tentatives de connexion aux services externes

Les variables suivantes concernent les tentatives de connexion aux services externes de Canopsis, à savoir Redis, RabbitMQ, MongoDB et InfluxDB.

Elles servent notamment à gérer le cas où les moteurs démarrent avant que ces services ne soient prêts (ce qui est essentiel pour Docker Compose).

Ces variables concernent l'ensemble des moteurs Go, et certains binaires comme `init` ou `feeder`.

| Variable d'environnement | Valeur par défaut | Utilité |
|:-------------------------|-------------------|---------|
| `CPS_MAX_RETRY` | `10` | Nombre maximum de tentatives de connexion aux services externes |
| `CPS_MAX_DELAY` | `30` | Temps maximum d'attente, en secondes, lors de chaque tentative de connexion à un service externe |
| `CPS_WAIT_FIRST_ATTEMPT` | `10` | Temps d'attente obligatoire, en secondes, avant la première tentative de connexion aux services externes |

### Durées d'exécution maximales

Les variables suivantes concernent des limites de sécurité appliquées par défaut dans Canopsis.

Dans certains cas d'utilisation, il peut être pertinent d'ajuster ces valeurs à vos besoins. Il faut néanmoins faire preuve de vigilance à ce sujet.

| Variable d'environnement | Valeur par défaut | Utilité |
|:-------------------------|-------------------|---------|
| `REGEXP2_MATCH_TIMEOUT` | `1s` | Depuis Canopsis 3.45.0. Durée (au format [`ParseDuration`](https://golang.org/pkg/time/#ParseDuration)) d'exécution maximale de l'évaluation d'une [expression régulière avancée](../../guide-utilisation/formats-et-syntaxe/format-regex.md) |

### `canopsinit`

L'outil `canopsinit` dispose lui aussi de variables d'environnement permettant de définir le temps d'attente nécessaire lors des tentatives de connexion aux services externes de Canopsis.

| Variable d'environnement | Valeur par défaut | Utilité |
|:-------------------------|-------------------|---------|
| `CPS_INIT_MAX_RETRY` | `10` | Nombre maximum de tentatives de connexion aux services externes |
| `CPS_INIT_RETRY_DELAY` | `30` | Temps maximum d'attente, en secondes, lors de chaque tentative de connexion à un service externe |
| `CPS_INIT_WAIT_FIRST_ATTEMPT` | `0` | Temps d'attente obligatoire, en secondes, avant la première tentative de connexion aux services externes |
