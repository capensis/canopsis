## Engine HeartBeat

Create an `Engine_heartbeat` queue bound to `canopsis.events` with routing key `#` in RabbitMQ.

```
cd cmd/engine-heartbeat

export CPS_MONGO_URL="mongodb://cpsmongo:canopsis@canopsis:27017/canopsis"
export CPS_AMQP_URL="amqp://cpsrabbit:canopsis@canopsis:5672/canopsis"
export CPS_REDIS_URL="redis://localhost:6379/0"

go build . && ./engine-heartbeat
```

### Configuration

A documents in the `heartbeat` collection:

```json
{
    "_id": "<HeartbeatID>",
    "pattern" : {"connector" : "myconnector"},
    "expected_interval" : "10s"
}
```
