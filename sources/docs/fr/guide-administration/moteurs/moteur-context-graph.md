# Context-graph

Le moteur `context-graph` gère les dépendances entre les entités.

## Récupérer les entités du context-graph

Ce nœud final renvoie une liste d'entités de graphe de contexte, basées sur un filtre MongoDB prédéfini.

## URL

GET / api / v2 / context /? limit = 100 & start = 1 & sort = ASC | DESC


### params:

filter (optionnel): une requête MongoDB envoyée en JSON codé en URL

start   (optionnel): utilisé

limite   (facultatif): nombre maximal d'éléments à renvoyer

sort    (facultatif): ordre de tri ( ASC = croissant, DESC = décroissant). Le tri par filtre d'échantillon peut être :

{"type":"connector"}



Valeur de retour:

une liste d'entités de graphes de contexte structurées comme suit:

    [{
        "impact": [
            "Engine_context-graph/localhost",
            "localhost",
            "Engine_cleaner_events/localhost",
            "Engine_topology/localhost",
            "Engine_collectdgw/localhost",
            "Engine_linklist/localhost",
            "Engine_eventstore/localhost",
            "Engine_acknowledgement/localhost",
            "Engine_perfdata/localhost",
            "Engine_pbehavior/localhost",
            "Engine_alerts/localhost",
            "Engine_event_filter/localhost",
            "Engine_event_filter_data/localhost",
            "Engine_ticket/localhost",
            "task_importctx/localhost",
            "Engine_cleaner_alerts/localhost",
            "580063594B10017E",
            "Engine_cancel/localhost",
            "task_linklist/localhost"
        ],
        "name": "engine",
        "enable_history": [
            1500280306
        ],
        "measurements": {},
        "enabled": true,
        "depends": [],
        "infos": {
          "enabled": true,
          "enable_history": [
              1499956041
          ],
          "rk": "Engine.engine.check.resource.localhost.Engine_context-graph"
        },
        "_id": "Engine/engine",
        "type": "connector"
      },
      {
        "impact": [
            "cpu-0/localhost",
            "localhost",
            "cpu-1/localhost",
            "cpu-2/localhost",
            "cpu-3/localhost",
            "swap/localhost",
            "load/localhost",
            "disk-sda/localhost",
            "disk-sda1/localhost",
            "disk-sda2/localhost",
            "disk-sda5/localhost",
            "memory/localhost",
            "df-root/localhost",
            "df-var-lib-docker-aufs/localhost",
            "interface-lo/localhost",
            "interface-virbr0/localhost",
            "interface-br-b355daa72396/localhost",
            "interface-br-40849e9d2c75/localhost",
            "interface-docker0/localhost",
            "interface-eth0/localhost",
            "580063594B10017E",
            "canopsis_mongodb/localhost"
        ],
        "depends": [],
        "_id": "collectd/collectd2event",
        "name": "collectd2event",
        "infos": {
            "enabled": true,
            "enable_history": [
                1499956041
            ],
            "rk": "collectd.collectd2event.perf.resource.localhost.cpu-0"
        },
        "measurements": {},
        "type": "connector"
      }]
