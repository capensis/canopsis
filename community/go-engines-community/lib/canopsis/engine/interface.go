// engine contain implementation of canopsis engine.
package engine

//go:generate mockgen -destination=../../../mocks/lib/canopsis/engine/engine.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine Engine,Consumer,MessageProcessor,PeriodicalWorker,RPCClient,RPCMessageProcessor

import (
	"context"
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

// PeriodicalWork interface is used to implement engine periodical worker.
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
	Process(RPCMessage) error
}

// RPCMessage is AMQP RPC request or response.
type RPCMessage struct {
	CorrelationID string
	Body          []byte
}

// RunInfoManager interface is used to implement engine run info storage.
type RunInfoManager interface {
	Save(ctx context.Context, info RunInfo, expiration time.Duration) error
	Get(ctx context.Context, engineName string) (*RunInfo, error)
	GetAll(ctx context.Context) ([]RunInfo, error)
	GetGraph(ctx context.Context) (*RunInfoGraph, error)
	ClearAll(ctx context.Context) error
}

// RunInfo is engine run information to detect engines order.
type RunInfo struct {
	Name            string `json:"name"`
	ConsumeQueue    string `json:"input_queue"`
	PublishQueue    string `json:"output_queue"`
	PublishExchange string `json:"output_exchange,omitempty"`
}

type RunInfoGraph struct {
	Nodes []RunInfo `json:"nodes"`
	Edges []Edge    `json:"edges"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}
