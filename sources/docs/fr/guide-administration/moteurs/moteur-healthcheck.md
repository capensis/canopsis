# Healthcheck

## Configuration

La configuration du moteur `healthcheck` est dans le fichier `/opt/canopsis/etc/healthcheck/manager.conf`.

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
- `check_collections` : La liste, séparée par des virgules et sans espaces, des collections MongoDB qui seront surveillées.
- `check_engines` : La liste, séparée par des virgules et sans espaces, des moteurs, Python comme Go, qui seront surveillés.
- `check_ts_db` : Le nom de la table utilisée par la base de données de stats.
- `check_webserver` : Le nom utilisé dans le système pour le webserver de Canopsis.
- `systemctl_engine_prefix` : Le préfixe utilisé dans `systemctl` pour les différents moteurs de Canopsis.