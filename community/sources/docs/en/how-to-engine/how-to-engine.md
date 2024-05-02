# How to create an engine?

Let's create an engine with a simple logic. For example we want to create an engine, which logs incoming `entity_id` and `alarm_id` from an event and sets the field `logged_at` to the event.

First we need to create a directory for our engine, let's create an `engine-logger` directory inside `cmd` directory. After that create `main.go` inside `engine-logger`:

```go
package main

func main() {
	
}
```

#### Create an engine instance

To create an engine we need to use a `New` function from `engine` package:

```go
func New(
	init func(ctx context.Context) error,
	deferFunc func(ctx context.Context),
	logger zerolog.Logger,
) Engine {
	return &engine{
		init:      init,
		deferFunc: deferFunc,
		logger:    logger,
	}
}
```

The function `New` return default engine struct, which implements the `Engine` interface:

```go
// Engine interface is used to implement the canopsis engine.
type Engine interface {
	// AddConsumer adds AMQP consumer to engine.
	AddConsumer(Consumer)
	// AddPeriodicalWorker adds periodical worker to engine.
	AddPeriodicalWorker(PeriodicalWorker)
	// Run starts goroutines for all consumers and periodical workers.
	// Engine stops if one of consumer or periodical worker return error.
	Run(context.Context) error
}
```

For now we're interested only in the `Run` function, which starts the engine.

So, what are the `init` and `deferFunc` functions?

Both functions are called only once. `init` function is called when the engine starts, `deferFunc` is called when the engine stops.

Let's modify `main.go` and create our new engine, which should log "init func" and "defer func"

```go
package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := log.NewLogger(false)

	engineLogger := engine.New(
		func(ctx context.Context) error {
			logger.Info().Msg("init func")

			return nil
		},
		func(ctx context.Context) {
			logger.Info().Msg("defer func")
		},
		logger,
	)

	err := engineLogger.Run(ctx)
	exitStatus := 0
	if err != nil {
		logger.Err(err).Msg("exit with error")
		exitStatus = 1
	}

	os.Exit(exitStatus)
}
```

Let's run it:

```
$ go run ./cmd/engine-logger
2021-11-23T17:37:46+07:00 INF lib\canopsis\engine\engine.go:47 > engine started consumers=0 periodical workers=0
2021-11-23T17:37:46+07:00 INF cmd\engine-logger\main.go:14 > init func
2021-11-23T17:37:46+07:00 INF cmd\engine-logger\main.go:19 > defer func
2021-11-23T17:37:46+07:00 INF lib\canopsis\engine\engine.go:99 > engine stopped
```

We can see our logs, also as you can see it stopped itself, because we don't have any workers inside it.

#### Engine workers

##### Create periodical worker
You can add two types of workers in the engine.
The first one is a periodical worker, which does some tasks periodically. It should implement the following interface:

```go
// PeriodicalWorker interface is used to implement engine periodical worker.
// If Work returns error engine will be stopped.
type PeriodicalWorker interface {
	GetInterval() time.Duration
	Work(context.Context) error
}
```

`GetInterval() time.Duration` - should return periodical interval duration.  
`Work(context.Context) error` - worker function, which does some task. If the function returns any error, then the engine will be stopped.

Let's create a periodical worker in separate file, which logs some incrementing counter:

```go
package main

import (
	"context"
	"github.com/rs/zerolog"
	"time"
)

type periodicalWorker struct {
	counter int
	interval time.Duration
	logger zerolog.Logger
}

func (w *periodicalWorker) Work(_ context.Context) error {
	w.logger.Info().Int("counter", w.counter).Msg("log counter")

	w.counter++

	return nil
}

func (w *periodicalWorker) GetInterval() time.Duration {
	return w.interval
}
```

And add it to our engine via `AddPeriodicalWorker` method:

```go
package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := log.NewLogger(false)

	engineLogger := engine.New(
		func(ctx context.Context) error {
			logger.Info().Msg("init func")

			return nil
		},
		func(ctx context.Context) {
			logger.Info().Msg("defer func")
		},
		logger,
	)

	engineLogger.AddPeriodicalWorker(&periodicalWorker{
		interval: time.Second,
		logger: logger,
	})

	err := engineLogger.Run(ctx)
	exitStatus := 0
	if err != nil {
		logger.Err(err).Msg("exit with error")
		exitStatus = 1
	}

	os.Exit(exitStatus)
}
```

If we run it, we'll see that our engine is no longer stops:

```
$ go run ./cmd/engine-logger
2021-11-23T18:18:54+07:00 INF lib\canopsis\engine\engine.go:47 > engine started consumers=0 periodical workers=1
2021-11-23T18:18:54+07:00 INF cmd\engine-logger\main.go:15 > init func
2021-11-23T18:18:55+07:00 INF cmd\engine-logger\periodical_worker.go:16 > log counter counter=0
2021-11-23T18:18:56+07:00 INF cmd\engine-logger\periodical_worker.go:16 > log counter counter=1
2021-11-23T18:18:57+07:00 INF cmd\engine-logger\periodical_worker.go:16 > log counter counter=2
2021-11-23T18:18:58+07:00 INF cmd\engine-logger\periodical_worker.go:16 > log counter counter=3
2021-11-23T18:18:59+07:00 INF cmd\engine-logger\periodical_worker.go:16 > log counter counter=4
2021-11-23T18:19:00+07:00 INF cmd\engine-logger\periodical_worker.go:16 > log counter counter=5
2021-11-23T18:19:01+07:00 INF cmd\engine-logger\periodical_worker.go:16 > log counter counter=6
```

##### Create RabbitMQ Consumer worker
The second worker type is `Consumer`, which reads messages from `RabbitMQ` queue. It should implement the following interface:

```go
// Consumer interface is used to implement AMQP consumer of engine.
// If Consume returns error engine will be stopped.
type Consumer interface {
	Consume(context.Context) error
}
```

`Consume(context.Context) error` - worker function, which does consume tasks. If the function returns any error, then the engine will be stopped.
`engine` package provides a `NewDefaultConsumer` function, which returns implementation of a rabbitmq consumer:

```go
// NewDefaultConsumer creates consumer.
func NewDefaultConsumer(
	name, queue string,
	consumePrefetchCount, consumePrefetchSize int,
	purgeQueue bool,
	nextExchange, nextQueue, fifoExchange, fifoQueue string,
	connection libamqp.Connection,
	processor MessageProcessor,
	logger zerolog.Logger,
) Consumer {
	return &defaultConsumer{
		name:                 name,
		queue:                queue,
		consumePrefetchCount: consumePrefetchCount,
		consumePrefetchSize:  consumePrefetchSize,
		purgeQueue:           purgeQueue,
		nextExchange:         nextExchange,
		nextQueue:            nextQueue,
		fifoExchange:         fifoExchange,
		fifoQueue:            fifoQueue,
		connection:           connection,
		processor:            processor,
		logger:               logger,
	}
}

// defaultConsumer implements AMQP consumer.
type defaultConsumer struct {
	// name is consumer name.
	name string
	// queue is name of AMQP queue from where consumer receives messages.
	queue                                     string
	consumePrefetchCount, consumePrefetchSize int
	purgeQueue                                bool
	// processor handles AMQP messages.
	processor MessageProcessor
	// nextQueue is name of AMQP queue to where consumer sends message after succeeded processing.
	nextQueue    string
	nextExchange string
	// fifoQueue is name of AMQP queue to where consumer sends message after failed processing
	// or if nextQueue is not defined.
	fifoQueue    string
	fifoExchange string
	// connection is AMQP connection.
	connection libamqp.Connection
	logger     zerolog.Logger
}
```

First of all we need a message processor, which will process our messages from the queue. It should implement the following interface:

```go
// MessageProcessor interface is used to implement AMQP message processor of consumer.
// If Process returns error engine will be stopped.
type MessageProcessor interface {
	Process(ctx context.Context, d amqp.Delivery) (newMessage []byte, err error)
}
```
We want to log incoming `entity_id` and `alarm_id` plus mutate an event with `logged_at` field for further engines.  
Let's add an additional field to the event structure that we want to populate with our engine:

```go

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"

type Event struct {
    //.....
    
	LoggedAt datetime.CpsTime `json:"logged_at"`
}
```

And then create a message processor in a separate file:

```go
package main

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type messageProcessor struct {
	Encoder encoding.Encoder
	Decoder encoding.Decoder
	logger  zerolog.Logger
}

func (p *messageProcessor) Process(_ context.Context, d amqp.Delivery) ([]byte, error) {
	var event types.Event

	err := p.Decoder.Decode(d.Body, &event)
	if err != nil {
		return nil, nil
	}

	if event.Entity == nil {
		return nil, errors.New("entity is nil")
	}

	if event.Alarm == nil {
		return nil, errors.New("alarm is nil")
	}

	event.LoggedAt = datetime.NewCpsTime()

	p.logger.Info().Str("entity_id", event.Entity.ID).Str("alarm_id", event.Alarm.ID).Msg("log")

	body, err := p.Encoder.Encode(event)
	if err != nil {
		return nil, nil
	}

	return body, nil
}
```

Now it's ready to be used in the `NewDefaultConsumer` function.

###### Find the engine's place in the pipeline.
To define other `NewDefaultConsumer` arguments, we need to decide where we should place our engine. Since we want to log incoming `entity_id` and `alarm_id` plus mutate an event with `logged_at` field for further engines we need to place it somewhere in the middle of the pipeline.

Community version pipeline:

fifo -> che -> pbehavior -> axe -> action

Since we need `entity_id` and `alarm_id` it should be placed after the `axe` engine, since the `axe` creates alarms. So with our new engine the pipeline will looks like this:  

fifo -> che -> pbehavior -> axe -> logs -> action

So for our `engine-logs` the previous engine is `engine-axe` and the next engine is `engine-action`.

###### Create engine's queue

Now we need to create a queue for `engine-logs`. To do this we need to add a new queue in the queues block in the `canopsis.toml` config file which is used by `canopsis-reconfigure` script.

```toml
[[RabbitMQ.queues]]
name = "Engine_logs"
durable = true
autoDelete = false
exclusive = false
noWait = false
# args =
```

After that run `canopsis-reconfigure` to create our queue

Then, since `engine-axe` should send message to the `engine-logs` we need to run `engine-axe` with flag -publishQueue Engine_logs

###### Add consumer to the engine

Now we're ready to create and add our consumer;

```go
	engineLogger.AddConsumer(
		engine.NewDefaultConsumer(
			"logs",                       // consumer name
			"Engine_logs",                // queue name
			cfg.Global.PrefetchCount,     // specific rabbitmq consumer values, just take them from global canopsis conf
			cfg.Global.PrefetchSize,      // specific rabbitmq consumer values, just take them from global canopsis conf
			false,                        // we don't want to purge our queue
			"",                           // we don't have any specific exchange for the next engine, so leave it empty
			"Engine_action",              // next queue name
			canopsis.FIFOAckExchangeName, // if message processor doesn't return message to the next queue, then if fifo exchange and queue is defined,
			canopsis.FIFOAckQueueName,    // then fifo_ack event will be sent to the engine_fifo to remove event lock. If you don't want to send fifo_ack, then leave it empty.
			amqpConnection,
			&messageProcessor{
				Encoder: json.NewEncoder(), // since we send json events
				Decoder: json.NewDecoder(), // we use json encoder and decoder
				logger:  logger,
			},
			logger,
		),
	)
```

To get canopsis config and create connections, you may use `DependencyMaker` tool from depmake package:
```go
m := depmake.DependencyMaker{}
dbClient := m.DepMongoClient(ctx)
cfg := m.DepConfig(ctx, dbClient)
amqpConnection := m.DepAmqpConnection(logger, cfg)
```

Complete main file looks like this:

```go
package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := log.NewLogger(false)

	m := depmake.DependencyMaker{}
	dbClient := m.DepMongoClient(ctx)
	cfg := m.DepConfig(ctx, dbClient)
	amqpConnection := m.DepAmqpConnection(logger, cfg)

	engineLogger := engine.New(
		func(ctx context.Context) error {
			logger.Info().Msg("init func")

			return nil
		},
		func(ctx context.Context) {
			logger.Info().Msg("defer func")
		},
		logger,
	)

	engineLogger.AddPeriodicalWorker(&periodicalWorker{
		interval: time.Second,
		logger:   logger,
	})

	engineLogger.AddConsumer(
		engine.NewDefaultConsumer(
			"logs",                       // consumer name
			"Engine_logs",                // queue name
			cfg.Global.PrefetchCount,     // specific rabbitmq consumer values, just take them from global canopsis conf
			cfg.Global.PrefetchSize,      // specific rabbitmq consumer values, just take them from global canopsis conf
			false,                        // we don't want to purge our queue
			"",                           // we don't have any specific exchange for the next engine, so leave it empty
			"Engine_action",              // next queue name
			canopsis.FIFOAckExchangeName, // if message processor doesn't return message to the next queue, then if fifo exchange and queue is defined,
			canopsis.FIFOAckQueueName,    // then fifo_ack event will be sent to the engine_fifo to remove event lock. If you don't want to send fifo_ack, then leave it empty.
			amqpConnection,
			&messageProcessor{
				Encoder: json.NewEncoder(), // since we send json events
				Decoder: json.NewDecoder(), // we use json encoder and decoder
				logger:  logger,
			},
			logger,
		),
	)

	err := engineLogger.Run(ctx)
	exitStatus := 0
	if err != nil {
		logger.Err(err).Msg("exit with error")
		exitStatus = 1
	}

	os.Exit(exitStatus)
}
```

Let's try to run it. Since we have mongo and rabbit connections, we need to specify `CPS_MONGO_URL` and `CPS_AMQP_URL` env variables:

```
CPS_MONGO_URL=mongodb://cpsmongo:canopsis@localhost:27017/canopsis CPS_AMQP_URL=amqp://cpsrabbit:canopsis@localhost:5672/canopsis go run ./cmd/engine-logger
```

And then send some events:

```json
[
    {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "source_type": "resource",
        "component": "test",
        "event_type": "check",
        "state": 1,
        "resource": "test_resource_1"
    },
    {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "source_type": "resource",
        "component": "test",
        "event_type": "check",
        "state": 1,
        "resource": "test_resource_2"
    },
    {
        "connector": "test_connector",
        "connector_name": "test_connectorname",
        "source_type": "resource",
        "component": "test",
        "event_type": "check",
        "state": 1,
        "resource": "test_resource_3"
    }
]
```

```
$ CPS_MONGO_URL=mongodb://cpsmongo:canopsis@localhost:27017/canopsis CPS_AMQP_URL=amqp://cpsrabbit:canopsis@localhost:5672/canopsis go run ./cmd/engine-logger
2021-11-24T13:35:43+07:00 INF lib\canopsis\engine\engine.go:47 > engine started consumers=1 periodical workers=1
2021-11-24T13:35:43+07:00 INF cmd\engine-logger\main.go:25 > init func
2021-11-24T13:35:46+07:00 INF cmd\engine-logger\message_processor.go:37 > log alarm_id=cdf3ee36-70e3-4ef5-bc7c-ca944b2c3d51 entity_id=test_resource_1/test
2021-11-24T13:35:46+07:00 INF cmd\engine-logger\message_processor.go:37 > log alarm_id=157df8a2-760c-4d58-9bda-1bbf57addba6 entity_id=test_resource_2/test
2021-11-24T13:35:46+07:00 INF cmd\engine-logger\message_processor.go:37 > log alarm_id=88a653ca-8486-46a1-903e-f294c87dfd82 entity_id=test_resource_3/test
```

We can see our log.

##### How to add a flag to the engine

If you need to add a flag you can use the [`flag`](https://pkg.go.dev/flag) package, which contains everything you need to create flags. Let's create a purgeQueue flag, which will be used in consumer creation.

```
var purgeQueue bool

flag.BoolVar(&purgeQueue, "purgeQueue", false, "purge consumer queue before work")
flag.Parse()
```

#### Create an RPC between two engines

Let's say we don't want to place an engine in the pipeline, but to have it isolated instead. Then we need to create an rpc connection between the engines.
We want that `engine-axe` call `engine-logs` via RabbitMQ RPC to log event and get `logged_at` value in the response. To do this we need to add an rpc-server as consumer to our `engine-logs`, and rpc-client to `engine-axe`

First of all we need to define message formats, basically we need request and response structs, let's create in `types` package file with our structs definitions:
```go
package types

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"

type LogsRpcRequest struct {
	EntityID string `json:"entity_id"`
	AlarmID  string `json:"alarm_id"`
}

type LogsRpcResponse struct {
	LoggedAt datetime.CpsTime `json:"logged_at"`
}
```

Then we need to create an rpc server, it's possible to do with `NewRPCServer` function:
```go
func NewRPCServer(
	name, queue string,
	consumePrefetchCount, consumePrefetchSize int,
	connection libamqp.Connection,
	processor MessageProcessor,
	logger zerolog.Logger,
) Consumer {
	return &rpcServer{
		name:                 name,
		queue:                queue,
		consumePrefetchCount: consumePrefetchCount,
		consumePrefetchSize:  consumePrefetchSize,
		connection:           connection,
		processor:            processor,
		logger:               logger,
	}
}

// rpcServer implements AMQP consumer of RPC requests.
type rpcServer struct {
	// name is consumer name.
	name string
	// queue is name of AMQP queue from where consumer receives messages.
	queue                                     string
	consumePrefetchCount, consumePrefetchSize int
	// processor handles AMQP messages.
	processor MessageProcessor
	// connection is AMQP connection.
	connection libamqp.Connection
	logger     zerolog.Logger
}
```

Then we need to change our message processor, now it receives `LogsRpcRequest` and returns `LogsRpcResponse`:
```go
package main

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

type messageProcessor struct {
	Encoder encoding.Encoder
	Decoder encoding.Decoder
	logger  zerolog.Logger
}

func (p *messageProcessor) Process(_ context.Context, d amqp.Delivery) ([]byte, error) {
	var request types.LogsRpcRequest
	var response types.LogsRpcResponse

	err := p.Decoder.Decode(d.Body, &request)
	if err != nil {
		return nil, nil
	}

	p.logger.Info().Str("entity_id", request.EntityID).Str("alarm_id", request.AlarmID).Msg("log")
	response.LoggedAt = datetime.NewCpsTime()

	body, err := p.Encoder.Encode(response)
	if err != nil {
		return nil, nil
	}

	return body, nil
}
```

Now we're ready to add our `RPCServer` to the engine:

```go
engineLogger.AddConsumer(engine.NewRPCServer(
	"logs",
	"Engine_logs",
	cfg.Global.PrefetchCount,
	cfg.Global.PrefetchSize,
	amqpConnection,
	&messageProcessor{
		Encoder: json.NewEncoder(), // since we send json events
		Decoder: json.NewDecoder(), // we use json encoder and decoder
		logger:  logger,
	},
	logger,
))
```

After that we need to add an `RPCClient` to the `engine-axe`.
First of all we need to create a separate queue for `engine-axe`, where `engine-logs` should send rpc responses:
```toml
[[RabbitMQ.queues]]
name = "Engine_axe_logs_rpc_client"
durable = true
autoDelete = false
exclusive = false
noWait = false
# args =
```

Do not forget to run `canopsis-reconfigure` after that.
Then we need to create an `RPCClient` and add it to the `engine-axe`, we can do it with `NewRPCClient`:
```go
// NewRPCClient creates new AMQP RPC client.
func NewRPCClient(
	name, serverQueueName, clientQueueName string,
	consumePrefetchCount, consumePrefetchSize int,
	processor RPCMessageProcessor,
	amqpChannel libamqp.Channel,
	logger zerolog.Logger,
) RPCClient {
	return &rpcClient{
		name:                 name,
		serverQueueName:      serverQueueName,
		clientQueueName:      clientQueueName,
		consumePrefetchCount: consumePrefetchCount,
		consumePrefetchSize:  consumePrefetchSize,
		processor:            processor,
		amqpChannel:          amqpChannel,
		logger:               logger,
	}
}

// rpcClient implements RPC client.
type rpcClient struct {
	// name is consumer name.
	name string
	// serverQueueName is name of AMQP queue to where client sends RPC requests.
	serverQueueName string
	// clientQueueName is name of AMQP queue from where client receives RPC response.
	clientQueueName                           string
	consumePrefetchCount, consumePrefetchSize int
	// processor handles AMQP messages.
	processor RPCMessageProcessor
	// connection is AMQP connection.
	amqpChannel libamqp.Channel
	logger      zerolog.Logger
}
```

Of course we need another message processor for that, let's create one in `engine-axe` directory, the difference from previous message processors is that it should implement `RPCMessageProcessor` interface:
```go
// RPCMessageProcessor interface is used to implement AMQP RPC response processor of consumer.
// If Process returns error engine will be stopped.
type RPCMessageProcessor interface {
	Process(ctx context.Context, msg RPCMessage) error
}
```

The difference from the simple `MessageProcessor` interface is that it receives `RPCMessage` as the second argument:
```go
type RPCMessage struct {
	CorrelationID string // the id, which helps us differentiate responses, we'll use alarm_id in our case
	Body          []byte // our response body
}
```

So our rpc message processor will look like this:

```go
package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

type logsClientMessageProcessor struct {
	Logger zerolog.Logger
	Decoder encoding.Decoder
}

func (p *logsClientMessageProcessor) Process(_ context.Context, msg engine.RPCMessage) error {
	var event types.LogsRpcResponse
	err := p.Decoder.Decode(msg.Body, &event)
	if err != nil {
		return nil
	}
	
	p.Logger.Info().Str("alarm_id", msg.CorrelationID).Int64("logged_at", event.LoggedAt.Unix()).Msg("log")

	return nil
}
```

Now add our `RPCClient` as engine's consumer in `dependencies.go` inside `engine-axe` directory, also we add our client to axe's event processor because we want to call our logs rpc when we process an event:
```go
logsRpcClient := engine.NewRPCClient(
	canopsis.AxeRPCConsumerName,
	"Engine_logs",
	"Engine_axe_logs_rpc_client",
	cfg.Global.PrefetchCount,
	cfg.Global.PrefetchSize,
	&logsClientMessageProcessor{
		Logger:  logger,
		Decoder: json.NewDecoder(),
	},
	amqpChannel,
	logger,
)
//.....
engineAxe.AddConsumer(engine.NewDefaultConsumer(
	canopsis.AxeConsumerName,
	canopsis.AxeQueueName,
	cfg.Global.PrefetchCount,
	cfg.Global.PrefetchSize,
	false,
	"",
	options.PublishToQueue,
	canopsis.FIFOAckExchangeName,
	canopsis.FIFOAckQueueName,
	amqpConnection,
	&messageProcessor{
		FeaturePrintEventOnError: options.FeaturePrintEventOnError,
		FeatureStatEvents:        options.FeatureStatEvents,
		EventProcessor: alarm.NewEventProcessor(
			alarm.NewAdapter(dbClient),
			entity.NewAdapter(dbClient),
			correlation.NewRuleAdapter(dbClient),
			alarmConfigProvider,
			m.depOperationExecutor(dbClient, alarmConfigProvider, alarmStatusService, statsService),
			alarmStatusService,
			logsRpcClient,
			redis.NewLockClient(corrRedisClient),
			logger,
		),
		StatsService:           statsService,
		RemediationRpcClient:   remediationRpcClient,
		TimezoneConfigProvider: timezoneConfigProvider,
		Encoder:                json.NewEncoder(),
		Decoder:                json.NewDecoder(),
		Logger:                 logger,
		PbehaviorAdapter:       pbehavior.NewAdapter(dbClient),
	},
	logger,
))
//..........
engineAxe.AddConsumer(logsRpcClient)
```

And add the `Call` inside  event processor `process` function:
```go
func (s *eventProcessor) Process(ctx context.Context, event *types.Event) (types.AlarmChange, error) {
	//.............

	body, err := json.Marshal(types.LogsRpcRequest{
		EntityID: event.Entity.ID,
		AlarmID:  event.Alarm.ID,
	})

	err = s.logsRpcClient.Call(engine.RPCMessage{
		CorrelationID: event.Alarm.ID,
		Body:          body,
	})
	if err != nil {
		return alarmChange, err
	}

	return alarmChange, nil
}
```

Then run engine-axe and engine-logs and send events:
```go
$ CPS_MONGO_URL=mongodb://cpsmongo:canopsis@localhost:27017/canopsis CPS_AMQP_URL=amqp://cpsrabbit:canopsis@localhost:5672/canopsis CPS_REDIS_URL=redis://localhost:6379/0 go run ./cmd/engine-axe
2021-11-24T16:46:53+07:00 INF cmd\engine-axe\main.go:51 > Statistic Events DISABLED
2021-11-24T16:46:53+07:00 INF lib\canopsis\debug\trace.go:162 > Profiling DISABLED
2021-11-24T16:46:53+07:00 INF lib\canopsis\debug\trace.go:177 > Tracing DISABLED
2021-11-24T16:46:53+07:00 INF lib\canopsis\config\config_provider.go:668 > StealthyInterval of alarm config section is used value=0s
2021-11-24T16:46:53+07:00 INF lib\canopsis\config\config_provider.go:828 > EnableLastEventDate of alarm config section is used value=false
2021-11-24T16:46:53+07:00 INF lib\canopsis\config\config_provider.go:613 > CancelAutosolveDelay of alarm config section is used value=1h0m0s
2021-11-24T16:46:53+07:00 INF lib\canopsis\config\config_provider.go:828 > DisableActionSnoozeDelayOnPbh of alarm config section is used value=false
2021-11-24T16:46:53+07:00 INF lib\canopsis\config\config_provider.go:613 > TimeToKeepResolvedAlarms of alarm config section is used value=720h0m0s
2021-11-24T16:46:53+07:00 INF lib\canopsis\config\config_provider.go:784 > DisplayNameScheme of alarm config section is used value="{{ rand_string 3 }}-{{ rand_string 3 }}-{{ rand_string 3 }}"
2021-11-24T16:46:53+07:00 INF lib\canopsis\config\config_provider.go:116 > OutputLength of alarm config section is used value=255
2021-11-24T16:46:53+07:00 INF lib\canopsis\config\config_provider.go:125 > LongOutputLength of alarm config section is used value=1024
2021-11-24T16:46:53+07:00 INF lib\canopsis\config\config_provider.go:874 > Timezone of timezone config section is used value=Europe/Paris
2021-11-24T16:46:53+07:00 INF lib\canopsis\config\config_provider.go:526 > TimeToExecute of data_storage config section is used value=Sunday,23
2021-11-24T16:46:53+07:00 INF lib\canopsis\engine\engine.go:47 > engine started consumers=5 periodical workers=5
2021-11-24T16:46:57+07:00 INF cmd\engine-axe\logs_client_message_processor.go:23 > receive rpc message alarm_id=d31115a2-2af6-43fe-b7e5-8dc21f342349 logged_at=1637747217
2021-11-24T16:46:57+07:00 INF cmd\engine-axe\logs_client_message_processor.go:23 > receive rpc message alarm_id=d8d91863-6fe8-4144-81e9-cb064685bf65 logged_at=1637747217
2021-11-24T16:46:57+07:00 INF cmd\engine-axe\logs_client_message_processor.go:23 > receive rpc message alarm_id=1666005b-90ce-4d2c-b41d-92d33de7d695 logged_at=1637747217
```
```go
$ CPS_MONGO_URL=mongodb://cpsmongo:canopsis@localhost:27017/canopsis CPS_AMQP_URL=amqp://cpsrabbit:canopsis@localhost:5672/canopsis go run ./cmd/engine-logger
2021-11-24T16:46:55+07:00 INF lib\canopsis\engine\engine.go:47 > engine started consumers=1 periodical workers=1
2021-11-24T16:46:55+07:00 INF cmd\engine-logger\main.go:33 > init func
2021-11-24T16:46:57+07:00 INF cmd\engine-logger\message_processor.go:27 > log alarm_id=d31115a2-2af6-43fe-b7e5-8dc21f342349 entity_id=test_resource_10/test
2021-11-24T16:46:57+07:00 INF cmd\engine-logger\message_processor.go:27 > log alarm_id=d8d91863-6fe8-4144-81e9-cb064685bf65 entity_id=test_resource_20/test
2021-11-24T16:46:57+07:00 INF cmd\engine-logger\message_processor.go:27 > log alarm_id=1666005b-90ce-4d2c-b41d-92d33de7d695 entity_id=test_resource_30/test
```

As we can see rpc requests are sent and rpc responses are received successfully.
