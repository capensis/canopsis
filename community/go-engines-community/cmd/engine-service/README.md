# Service

You have the possibility in Canopsis to define service entity which may represent 
an application, component, or team which is set in order to explore alarms and states of 
corresponding entities and other services. Service is linked to alarm which state is resolved
by states of dependent entities according to defined method. 

## Engine

An engine `service` loads services from database, and processes incoming alarms 
to resolve and update state of service alarms.

### Deployment
The engine assumes the presence of the `Engine_service` queue, which can be created by adding the 
following code to initialization.toml:
```
[[RabbitMQ.queues]]
name = "Engine_service"
durable = true
autoDelete = false
exclusive = false
noWait = false
# args =
  [RabbitMQ.queues.bind]
  key = "Engine_service"
  exchange = "amq.direct"
  noWait = false
  # args =
```

Events get into the engine Service from the engine Axe or Correlation in Cat, through the queue `"Engine_service"`, so
`-publishQueue` startup parameter of the engine Axe/Correlation should be set as `"Engine_service"`.  
(example of engine Axe start command: `"./engine-axe -publishQueue Engine_service"`)

Other engines can communicate with Service via RabbitMQ-RPC. RabbitMQ-RPC requires additional 
rabbitmq queues, the configuration should be updated for:

```
[[RabbitMQ.queues]]
name = "Engine_service_rpc_server"
durable = true
autoDelete = false
exclusive = false
noWait = false
# args =

[[RabbitMQ.queues]]
name = "Engine_axe_service_rpc_client"
durable = true
autoDelete = false
exclusive = false
noWait = false
# args =
```
