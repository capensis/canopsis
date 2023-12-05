package contextgraph

import (
	"context"
	"fmt"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rmqPublisher struct {
	exchange, queue string
	encoder         encoding.Encoder
	contentType     string
	amqpPublisher   libamqp.Publisher
}

func NewEventPublisher(
	exchange, queue string,
	encoder encoding.Encoder,
	contentType string,
	amqpPublisher libamqp.Publisher,
) EventPublisher {
	return &rmqPublisher{
		exchange:      exchange,
		queue:         queue,
		encoder:       encoder,
		contentType:   contentType,
		amqpPublisher: amqpPublisher,
	}
}

func (p *rmqPublisher) SendImportResultEvent(ctx context.Context, uuid string, execTime time.Duration, state types.CpsNumber) error {
	return p.sendEvent(ctx, types.Event{
		Connector:     "taskhandler",
		ConnectorName: "task_importctx",
		Component:     "job",
		EventType:     types.EventTypeCheck,
		SourceType:    types.SourceTypeResource,
		Resource:      uuid,
		Timestamp:     datetime.NewCpsTime(),
		State:         state,
		Output:        fmt.Sprintf("Import %s failed.", uuid),
		ExtraInfos: map[string]interface{}{
			"execution_time": execTime,
		},
		Initiator: types.InitiatorSystem,
	})
}

func (p *rmqPublisher) SendPerfDataEvent(ctx context.Context, uuid string, stats importcontextgraph.Stats, state types.CpsNumber) error {
	return p.sendEvent(ctx, types.Event{
		Connector:     "Taskhandler",
		ConnectorName: "task_importctx",
		Component:     uuid,
		EventType:     types.EventTypeCheck,
		SourceType:    types.SourceTypeResource,
		Resource:      "task_importctx/report",
		State:         state,
		Output:        fmt.Sprintf("execution : %f sec, updated ent : %d, deleted ent : %d", stats.ExecTime.Seconds(), stats.Updated, stats.Deleted),
		PerfData:      fmt.Sprintf("execution_time=%ds ent_updated=%d ent_deleted=%d", int64(stats.ExecTime.Seconds()), stats.Updated, stats.Deleted),
		Initiator:     types.InitiatorSystem,
	})
}

func (p *rmqPublisher) sendEvent(ctx context.Context, event types.Event) error {
	bevent, err := p.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("error while encoding event %w", err)
	}

	return p.amqpPublisher.PublishWithContext(
		ctx,
		p.exchange,
		p.queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  p.contentType,
			Body:         bevent,
			DeliveryMode: amqp.Persistent,
		},
	)
}
