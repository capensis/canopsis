# Variables d'environnement Canopsis

Ce document regroupe l'ensemble des variables d'environnement pouvant être ajustées à vos besoins, pour les différents composants de Canopsis.

!!! note
    À ce jour, ce document est encore en cours de rédaction, et n'est pas exhaustif. Il sera complété au fur et à mesure des mises à jour de la documentation.

## Modification des variables d'environnement

En installation Docker, les variables suivantes peuvent être définies dans un fichier `compose.env` ou dans une section `environment:` de votre Docker Compose.

En installation par paquets, ces variables peuvent être définies dans `/opt/canopsis/etc/go-engines-vars.conf` si elles doivent s'appliquer à l'ensemble des moteurs, ou dans la section `[Environment]` d'une unité systemd dédiée si un moteur en particulier est ciblé.

## Variables d'environnement des composants Go

Ces variables concernent l'ensemble des moteurs Go, et certains outils comme `canopsis-reconfigure`.

### URI de connexion aux services externes

Votre installation de Canopsis doit obligatoirement comporter les adresses et données de connexion (on parle d'URI) permettant de se connecter aux services externes Redis, MongoDB et RabbitMQ.

| Variable d'environnement | Valeur par défaut | Utilité |
|:-------------------------|-------------------|---------|
| `CPS_AMQP_URL` | (vide) | Une URI de connexion RabbitMQ (cf. [Spécification d'URI RabbitMQ](https://www.rabbitmq.com/uri-spec.html)) |
| `CPS_API_URL` | (vide) | Une URI de connexion à l'API Canopsis |
| `CPS_MONGO_URL` | (vide) | Une URI de connexion MongoDB (cf. [Spécification d'URI MongoDB](https://docs.mongodb.com/v4.2/reference/connection-string/)) |
| `CPS_OLD_API_URL` | (vide) | URI de connexion à l'ancienne API Gunicorn de Canopsis |
| `CPS_POSTGRES_URL` | (vide) | URI de connexion PostgreSQL/TimescaleDB (cf. [Spécification d'URI PostgreSQL](https://www.postgresql.org/docs/13/libpq-connect.html#LIBPQ-CONNSTRING)) |
| `CPS_REDIS_URL` | (vide) | Une URI de connexion Redis (cf. [Spécification d'URI Redis](https://www.iana.org/assignments/uri-schemes/prov/redis)) |

### Tentatives de connexion aux services externes

Les variables suivantes concernent les tentatives de connexion aux services externes de Canopsis, à savoir Redis, RabbitMQ et MongoDB.

Elles servent notamment à gérer le cas où les moteurs démarrent avant que ces services ne soient prêts (ce qui est essentiel pour Docker Compose).

| Variable d'environnement | Valeur par défaut | Utilité |
|:-------------------------|-------------------|---------|
| `CPS_MAX_RETRY` | `10` | Nombre maximum de tentatives de connexion aux services externes |
| `CPS_MAX_DELAY` | `30` | Temps maximum d'attente, en secondes, lors de chaque tentative de connexion à un service externe |
| `CPS_WAIT_FIRST_ATTEMPT` | `10` | Temps d'attente obligatoire, en secondes, avant la première tentative de connexion aux services externes |

### Utilisation d'un proxy HTTP ou HTTPS

Certains moteurs permettent l'utilisation d'un proxy HTTP ou HTTPS pour les accès à des ressources web externes. Si les deux types de proxy sont activés, le proxy HTTPS sera généralement privilégié.

| Variable d'environnement | Valeur par défaut | Utilité |
|:-------------------------|-------------------|---------|
| `HTTP_PROXY` | (vide) | URL vers un proxy HTTP (ex : `http://utilisateur:motdepasse@192.168.0.253:3128/`) |
| `HTTPS_PROXY` | (vide) | URL vers un proxy HTTPS (ex : `https://utilisateur:motdepasse@192.168.0.253:3128/`) |
| `NO_PROXY` | (vide) | Liste des adresses pour lesquelles le proxy ne doit **pas** s'appliquer, séparées par des virgules (ex : `127.0.0.1,localhost`) |

### Durées d'exécution maximales

Les variables suivantes concernent des limites de sécurité appliquées par défaut dans Canopsis.

Dans certains cas d'utilisation, il peut être pertinent d'ajuster ces valeurs à vos besoins. Il faut néanmoins faire preuve de vigilance à ce sujet.

| Variable d'environnement | Valeur par défaut | Utilité |
|:-------------------------|-------------------|---------|
| `REGEXP2_MATCH_TIMEOUT` | `1s` | Depuis Canopsis 3.45.0. Durée (au format [`ParseDuration`](https://golang.org/pkg/time/#ParseDuration)) d'exécution maximale de l'évaluation d'une [expression régulière avancée](../../guide-utilisation/formats-et-syntaxe/format-regex.md) |

## Variables d'environnement propres à certains outils

### `canopsinit`

L'outil `canopsinit` dispose lui aussi de variables d'environnement permettant de définir le temps d'attente nécessaire lors des tentatives de connexion aux services externes de Canopsis.

| Variable d'environnement | Valeur par défaut | Utilité |
|:-------------------------|-------------------|---------|
| `CPS_INIT_MAX_RETRY` | `10` | Nombre maximum de tentatives de connexion aux services externes |
| `CPS_INIT_RETRY_DELAY` | `30` | Temps maximum d'attente, en secondes, lors de chaque tentative de connexion à un service externe |
| `CPS_INIT_WAIT_FIRST_ATTEMPT` | `0` | Temps d'attente obligatoire, en secondes, avant la première tentative de connexion aux services externes |

## Variables d'environnement propres à Docker

Les variables suivantes ne sont disponibles que dans un environnement Canopsis reposant sur Docker. Elles sont sans effet sur les autres méthodes d'installation.

| Variable d'environnement | Valeur par défaut | Utilité |
|:-------------------------|-------------------|---------|
| `CPS_LOGGING_LEVEL` | `info` | Permet de surcharger le niveau de verbosité des moteurs et tasks Python. `debug` permet d'obtenir davantage d'informations |
