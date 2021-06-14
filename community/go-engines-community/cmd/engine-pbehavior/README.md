# Pbehavior

You have the possibility in Canopsis to define periods of time during which monitoring services should be delivered or when the system or its elements is on maintenance. There are 4 types of periods of changes in behavior:

 - Active state: applications that must provide a service at a given time
 - Maintenance: you want to declare entities as maintenance so that their alarms do not appear visually
 - Pause: you want to pause an application for an indefinite time
 - Inactive state: when active priods defined, the rest of the day means as inactive period.

 There are possible to add custom-defined periods based on first 3 types above with specific priority. Default and custom types have icons (except default active state) to visually identify alarms with particular period.

## Engine

A new engine `pbehavior` loads pbehavior rules from database, and process incoming alarms according to pbehavior intervals and specified priorities. Processing resolves intervals defined in rules and when period specified for the moment it marks alarms with appropriate `PbehaviorInfo` sub-document.

### Deployment
The engine assumes the presence of the `Engine_pbehavior` queue, which can be created by adding the following code to initialization.toml:
```
[[RabbitMQ.queues]]
name = "Engine_pbehavior"
durable = true
autoDelete = false
exclusive = false
noWait = false
# args =
  [RabbitMQ.queues.bind]
  key = "Engine_pbehavior"
  exchange = "amq.direct"
  noWait = false
  # args =
```

Events get into the engine Pbehavior from the engine Che, through the queue `"Engine_pbehavior"`, so
`-publishQueue` startup parameter of the engine Che should be set as `"Engine_pbehavior"`.  
(example of engine Axe start command: `"./engine-che -publishQueue Engine_pbehavior"`)
