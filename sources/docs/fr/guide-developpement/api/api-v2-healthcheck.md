# Healthcheck

L'API healthcheck permet d'obtenir l'état de fonctionnement de Canopsis ; elle permet, par exemple, de savoir si les services nécessaires sont disponibles.

## Configuration

La configuration de l'API healthcheck est dans le fichier `/opt/canopsis/etc/healthcheck/manager.conf`.

Sa structure est la suivante :

```ini
[HEALTHCHECK]

check_amqp_limit_size = 100000
check_amqp_queues = Engine_alerts,Engine_cleaner_events,Engine_context-graph,Engine_event_filter,Engine_pbehavior,task_importctx
check_collections = default_entities,periodical_alarm
check_engines = cleaner-cleaner_events,dynamic-alerts,dynamic-context-graph,dynamic-pbehavior,dynamic-watcher,event_filter-event_filter,task_importctx-task_importctx
check_ts_db = canopsis
check_webserver = canopsis-webserver
systemctl_engine_prefix = canopsis-engine@
```

Les paramètres sont :

- `check_amqp_limit_size` : Le nombre maximum de messages dans une file RabbitMQ ; au delà, la file est considérée comme surchargée.
- `check_amqp_queues` : La liste, séparée par des virgules et sans espaces, des files RabbitMQ qui seront surveillées.
- `check_collections` : La liste, séparée par des virgules et sans espaces, des collections MongoDB dont l'existence sera surveillée.
- `check_engines` : La liste, séparée par des virgules et sans espaces, des moteurs, Python comme Go, qui seront surveillés.
- `check_ts_db` : Le nom de la base de données de statistiques.
- `check_webserver` : Le nom utilisé dans le système pour le serveur web de Canopsis.
- `systemctl_engine_prefix` : Le préfixe utilisé dans `systemctl` pour les différents moteurs de Canopsis.

### Récupérer l'état global

Renvoie un résumé de l'état.

**URL** : `/api/v2/healthcheck/`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer l'état général :

```sh
curl -X GET -u root:root 'http://localhost:8082/api/v2/healthcheck/'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "amqp": "",
    "cache": "",
    "database": "",
    "engines": "",
    "time_series": "",
    "timestamp": 1542795713,
    "overall": true
}
```

Le bon fonctionnement général est annoncé par la clef **overall** qui doit être à *true*.

### Récupérer l'état des services

Récupère l'état de services en particulier.

**URL** : `/api/v2/healthcheck/?criticals=<id_des_services>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer l'état des services `amqp`, `cache`, `database` et `engines` :

```sh
curl -X GET -u root:root 'http://localhost:8082/api/v2/healthcheck/?criticals=amqp,cache,database,engines'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "amqp": "",
    "cache": "",
    "database": "",
    "engines": "",
    "timestamp": 1541774674,
    "overall": true
}
```

### Description du résumé

La réponse contient les champs suivants :

* amqp, cache, database, engines, timeseries : un message d'erreur associé à chaque service ; rien s'il n'y a pas d'erreur détectée ;
* timestamp : le moment de création du résultat ;
* overall : un booléen pour savoir si l'état global est bon ou mauvais.

Concernant `overall`, par défaut, tous les services sont pris en compte pour calculer l'état global. Il est toutefois possible de sélectionner les services à considérer comme indispensables en utilisant le paramètre `criticals` dans l'URL GET.
`criticals` est une liste de services séparés par des virgules.

### Ce qui est vérifié

* amqp : la connexion et le canal sont ouverts, et il est possible de publier un message. On vérifie aussi qu'une liste de files existe bien, qu'il y a au moins un Consumer dessus, et que la file est active et non saturée (> 100 000 messages en attente) ;
* cache : la connexion est fonctionnelle et il est possible de faire un ECHO ;
* database : la connexion fonctionne et il est possible de lire dans une liste de collections ;
* engines : hors Docker, vérifie avec `systemctl` que les engines de base (Python) sont en état *running* ;
* time series : vérifie que la base existe et que l'on peut y lire des *measurements*.
