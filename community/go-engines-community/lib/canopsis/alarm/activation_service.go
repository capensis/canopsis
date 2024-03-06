package alarm

import (
	"context"
	"fmt"

	amqplib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ActivationService checks alarm and sends activation event
// if alarm doesn't have active snooze and pbehavior.
type ActivationService interface {
	Process(
		ctx context.Context,
		alarm types.Alarm,
		fifoAckEvent types.Event,
	) (bool, error)
}

type baseActivationService struct {
	encoder   encoding.Encoder
	publisher amqplib.Publisher
	queueName string
}

func NewActivationService(
	encoder encoding.Encoder,
	publisher amqplib.Publisher,
	queueName string,
) ActivationService {
	return &baseActivationService{
		encoder:   encoder,
		publisher: publisher,
		queueName: queueName,
	}
}

func (s *baseActivationService) Process(
	ctx context.Context,
	alarm types.Alarm,
	fifoAckEvent types.Event,
) (bool, error) {
	if alarm.CanActivate() {
		err := s.sendActivationEvent(ctx, fifoAckEvent)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (s *baseActivationService) sendActivationEvent(
	ctx context.Context,
	fifoAckEvent types.Event,
) error {
	event := types.Event{
		Connector:          fifoAckEvent.Connector,
		ConnectorName:      fifoAckEvent.ConnectorName,
		Component:          fifoAckEvent.Component,
		Resource:           fifoAckEvent.Resource,
		Timestamp:          datetime.NewCpsTime(),
		ReceivedTimestamp:  fifoAckEvent.ReceivedTimestamp,
		SourceType:         fifoAckEvent.SourceType,
		IsMetaAlarmUpdated: fifoAckEvent.IsMetaAlarmUpdated,
		EventType:          types.EventTypeActivate,
		Author:             canopsis.DefaultEventAuthor,
		Initiator:          types.InitiatorSystem,
	}

	body, err := s.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("fail encode activation event: %w", err)
	}

	err = s.publisher.PublishWithContext(
		ctx,
		"",
		s.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("fail publish activation event to FIFO: %w", err)
	}

	return nil
}
