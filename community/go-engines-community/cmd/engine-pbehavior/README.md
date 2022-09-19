# Pbehavior

You have the possibility in Canopsis to define periods of time during which monitoring services should be delivered or 
when the system or its elements is on maintenance. There are 4 types of periods of changes in behavior:

 - Active state: applications that must provide a service at a given time
 - Maintenance: you want to declare entities as maintenance so that their alarms do not appear visually
 - Pause: you want to pause an application for an indefinite time
 - Inactive state: when active priods defined, the rest of the day means as inactive period.

 There are possible to add custom-defined periods based on first 3 types above with specific priority. Default and custom 
 types have icons (except default active state) to visually identify alarms with particular period.

## Engine

The `engine-pbehavior` loads pbehavior rules from database and process alarms and entities according to pbehavior intervals 
and specified priorities. Processing resolves intervals defined in rules and when period specified for the moment it marks
alarms and entities with appropriate `PbehaviorInfo` sub-document.

### Deployment
The engine assumes the presence of the `Engine_pbehavior_rpc_server` queue, which can be created by adding the following
code to initialization.toml:
```
[[RabbitMQ.queues]]
name = "Engine_pbehavior_rpc_server"
durable = true
autoDelete = false
exclusive = false
noWait = false
# args =
```
