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
	Process(ctx context.Context, alarm types.Alarm, eventReceivedTimestamp datetime.MicroTime, isMetaAlarmUpdated bool) (bool, error)
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
	eventRealTimestamp datetime.MicroTime,
	isMetaAlarmUpdated bool,
) (bool, error) {
	if alarm.CanActivate() {
		err := s.sendActivationEvent(ctx, alarm, eventRealTimestamp, isMetaAlarmUpdated)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (s *baseActivationService) sendActivationEvent(
	ctx context.Context,
	alarm types.Alarm,
	eventReceivedTimestamp datetime.MicroTime,
	isMetaAlarmUpdated bool,
) error {
	event := types.Event{
		Connector:          alarm.Value.Connector,
		ConnectorName:      alarm.Value.ConnectorName,
		Component:          alarm.Value.Component,
		Resource:           alarm.Value.Resource,
		Timestamp:          datetime.NewCpsTime(),
		ReceivedTimestamp:  eventReceivedTimestamp,
		IsMetaAlarmUpdated: isMetaAlarmUpdated,
		EventType:          types.EventTypeActivate,
		Initiator:          types.InitiatorSystem,
	}
	event.SourceType = event.DetectSourceType()
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
