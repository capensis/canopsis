# Context-graph APIs



## Fetch context-graph entities


This endpoint returns a list of context-graph entities, based on a predefined MongoDB filter.


#### url
  `GET` /api/v2/context/<filter>?limit=100&start=1&sort=ASC|DESC

#### params:

- **filter** (optional) : a mongoDB request sent as an urlencoded json
- **start**  (optional) : used
- **limit**  (optional) : max number of items to return
- **sort**   (optional) : sort order (__ASC__ = ascending, __DESC__ = descending). Sorts by
â sample filter can be :


      `{"type":"connector"} `


#### Return value:

a list of context-graph entites structured as follows :

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
