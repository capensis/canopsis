package importcontextgraph

import (
	"context"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
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

func (p *eventPublisher) SendEvent(ctx context.Context, event types.Event) error {
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
