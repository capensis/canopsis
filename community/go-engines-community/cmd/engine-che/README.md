# Entity

You have the possibility in Canopsis to enrich the events, to create and update the context-graph.

Each contextable event launches context-graph update : 

- create new resource, component, connector entities
- links entities together
- links entities to services 

## Engine

An engine `che` loads event filters services from database, and processes incoming alarms according to
conditions. Processing matches event and entity to rules pattern and appropriately updates
event and context graph.

### Deployment
The engine assumes the presence of the `Engine_che` queue, which can be created by adding the 
following code to initialization.toml:
```
[[RabbitMQ.queues]]
name = "Engine_che"
durable = true
autoDelete = false
exclusive = false
noWait = false
# args =
```

Events get into the engine Che from the engine Fifo, through the queue `"Engine_che"`, so
`-publishQueue` startup parameter of the engine Fifo should be set as `"Engine_che"`.  
(example of engine Fifo start command: `"./engine-fifo -publishQueue Engine_che"`)

### Metrics

Use pro version of engine to support metrics.
