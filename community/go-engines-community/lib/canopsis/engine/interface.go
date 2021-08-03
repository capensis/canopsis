// Package engine contain implementation of canopsis engine.
package engine

//go:generate mockgen -destination=../../../mocks/lib/canopsis/engine/engine.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine Engine,Consumer,MessageProcessor,PeriodicalWorker,RPCClient,RPCMessageProcessor

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/streadway/amqp"
	"time"
)

// Engine interface is used to implement canopsis engine.
type Engine interface {
	// AddConsumer adds AMQP consumer to engine.
	AddConsumer(Consumer)
	// AddPeriodicalWorker adds periodical worker to engine.
	AddPeriodicalWorker(PeriodicalWorker)
	// Run starts goroutines for all consumers and periodical workers.
	// Engine stops if one of consumer or periodical worker return error.
	Run(context.Context) error
}

// Consumer interface is used to implement AMQP consumer of engine.
// If Consume returns error engine will be stopped.
type Consumer interface {
	Consume(context.Context) error
}

// MessageProcessor interface is used to implement AMQP message processor of consumer.
// If Process returns error engine will be stopped.
type MessageProcessor interface {
	Process(ctx context.Context, d amqp.Delivery) (newMessage []byte, err error)
}

// PeriodicalWorker interface is used to implement engine periodical worker.
// If Work returns error engine will be stopped.
type PeriodicalWorker interface {
	GetInterval() time.Duration
	Work(context.Context) error
}

// RPCClient interface is used to implement AMQP RPC client.
type RPCClient interface {
	// Consumer receives RPC responses from AMQP queue.
	Consumer
	// Call receives RPC request and publishes it to AMQP queue.
	Call(m RPCMessage) error
}

// RPCMessageProcessor interface is used to implement AMQP RPC response processor of consumer.
// If Process returns error engine will be stopped.
type RPCMessageProcessor interface {
	Process(ctx context.Context, msg RPCMessage) error
}

// RPCMessage is AMQP RPC request or response.
type RPCMessage struct {
	CorrelationID string
	Body          []byte
}

// RunInfoManager interface is used to implement engine run info storage.
type RunInfoManager interface {
	SaveInstance(ctx context.Context, info InstanceRunInfo, expiration time.Duration) error
	GetEngineQueues(ctx context.Context) ([]RunInfo, error)
	GetCacheKey(info InstanceRunInfo) string
}

// RunInfo is engine run information.
type RunInfo struct {
	Name         string `json:"name"`
	ConsumeQueue string `json:"consume_queue"`
	PublishQueue string `json:"publish_queue"`
}

// InstanceRunInfo is instance of engine run information.
type InstanceRunInfo struct {
	ID               string        `json:"_id"`
	Name             string        `json:"name"`
	ConsumeQueue     string        `json:"consume_queue"`
	PublishQueue     string        `json:"publish_queue"`
	RpcConsumeQueues []string      `json:"rpc_consume_queues"`
	RpcPublishQueues []string      `json:"rpc_publish_queues"`
	QueueLength      int           `json:"queue_length"`
	Time             types.CpsTime `json:"time"`
}

func NewInstanceRunInfo(name, consumeQueue, publishQueue string, rpcQueues ...[]string) InstanceRunInfo {
	var rpcConsumeQueues, rpcPublishQueues []string
	if len(rpcQueues) > 0 {
		if len(rpcQueues) > 2 {
			panic("too much arguments")
		}

		rpcConsumeQueues = rpcQueues[0]
		if len(rpcQueues) > 1 {
			rpcPublishQueues = rpcQueues[1]
		}
	}

	return InstanceRunInfo{
		ID:               utils.NewID(),
		Name:             name,
		ConsumeQueue:     consumeQueue,
		PublishQueue:     publishQueue,
		RpcConsumeQueues: rpcConsumeQueues,
		RpcPublishQueues: rpcPublishQueues,
	}
}
