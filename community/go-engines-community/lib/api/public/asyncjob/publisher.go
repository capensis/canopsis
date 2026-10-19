package asyncjob

import (
	"context"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher signals a Worker that a job is ready to be processed.
type Publisher interface {
	Publish(ctx context.Context, jobType, jobID string) error
}

type publisher struct {
	publisher libamqp.Publisher
	encoder   encoding.Encoder
}

func NewPublisher(p libamqp.Publisher, encoder encoding.Encoder) Publisher {
	return &publisher{
		publisher: p,
		encoder:   encoder,
	}
}

func (p *publisher) Publish(ctx context.Context, jobType, jobID string) error {
	body, err := p.encoder.Encode(Message{Type: jobType, JobID: jobID})
	if err != nil {
		return err
	}

	return p.publisher.PublishWithContext(
		ctx,
		canopsis.DefaultExchangeName,
		canopsis.PublicApiJobQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}
