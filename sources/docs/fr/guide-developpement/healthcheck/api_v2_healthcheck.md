# API healthcheck

L'API healthcheck permet d'obtenir l'état de fonctionnement de Canopsis ; par exemple elle permet de savoir si les services nécessaires sont disponibles.


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

- `check_amqp_limit_size` : Le nombre maximum de messages dans une queue RabbitMQ, au delà la queue est considérée comme surchargée.
- `check_amqp_queues` : La liste, séparée par des virgules et sans espaces, des queues de RabbitMQ qui seront surveillées.
- `check_collections` : La liste, séparée par des virgules et sans espaces, des collections MongoDB dont l'existence sera surveillée.
- `check_engines` : La liste, séparée par des virgules et sans espaces, des moteurs, Python comme Go, qui seront surveillés.
- `check_ts_db` : Le nom de la base de données de stats.
- `check_webserver` : Le nom utilisé dans le système pour le webserver de Canopsis.
- `systemctl_engine_prefix` : Le préfixe utilisé dans `systemctl` pour les différents moteurs de Canopsis.

### Récupérer l'état global

```
GET /api/v2/healthcheck/
```

Renvoie un résumé de l'état.

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

```
GET /api/v2/healthcheck/?criticals=amqp,cache,database,engines,time_series
```

Renvoie un résumé de l'état.

```json
{
    "amqp": "",
    "cache": "",
    "database": "",
    "engines": "",
    "time_series": "",
    "timestamp": 1541774674,
    "overall": true
}
```

### Description du résumé

La réponse contient les champs suivants :
* amqp, cache, database, engines, timeseries : un message d'erreur associé à chacque service ; rien s'il n'y a pas d'erreur détectée ;
* timestamp : le moment de création du résultat ;
* overall : un booléen pour savoir si l'état global est bon ou mauvais.

Concernant `overall`, par défaut, tous les services sont pris en compte pour calculer l'état global. Il est toutefois possible de sélectionner les services à considérer comme indispensable en utilisant le paramètre `criticals` dans l'url GET.
`criticals` est une liste de services séparés par des virgules.

### Ce qui est vérifié

* amqp : la connexion est ouverte, le channel aussi, et il est possible de publier un message. On vérifie aussi qu'une liste de queues existe bien, qu'il y a au moins un Consumer dessus, et que la queue est active et non saturée (> 100 000 messages en attente) ;
* cache: la connexion est fonctionnelle et il est possible de faire un ECHO ;
* database: la connexion fonctionne et il est possible de lire dans une liste de collections ;
* engines : hors docker, vérifie par systemctl que les engines de bases (python) sont "running" ;
* time series : vérifie que la database existe et que l'on peut lire des measurements.
