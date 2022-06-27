package importcontextgraph

import (
	"fmt"
	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type eventPublisher struct {
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
	return &eventPublisher{
		exchange:      exchange,
		queue:         queue,
		encoder:       encoder,
		contentType:   contentType,
		amqpPublisher: amqpPublisher,
	}
}

func (p *eventPublisher) SendUpdateEntityServiceEvent(serviceId string) error {
	return p.sendEvent(types.Event{
		EventType:     types.EventTypeRecomputeEntityService,
		Connector:     types.ConnectorEngineService,
		ConnectorName: types.ConnectorEngineService,
		Component:     serviceId,
		Timestamp:     types.CpsTime{Time: time.Now()},
		Author:        canopsis.DefaultEventAuthor,
		SourceType:    types.SourceTypeService,
	})
}

func (p *eventPublisher) sendEvent(event types.Event) error {
	bevent, err := p.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("error while encoding event %+v", err)
	}

	return p.amqpPublisher.Publish(
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
