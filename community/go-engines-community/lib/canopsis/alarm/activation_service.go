package alarm

import (
	"context"

	amqplib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

// ActivationService checks alarm and sends activation event
// if alarm doesn't have active snooze and pbehavior.
type ActivationService interface {
	Process(ctx context.Context, alarm types.Alarm, eventReceivedTimestamp datetime.MicroTime) (bool, error)
}

type baseActivationService struct {
	encoder   encoding.Encoder
	publisher amqplib.Publisher
	queueName string
	logger    zerolog.Logger
}

func NewActivationService(
	encoder encoding.Encoder,
	publisher amqplib.Publisher,
	queueName string,
	logger zerolog.Logger,
) ActivationService {
	return &baseActivationService{
		encoder:   encoder,
		publisher: publisher,
		queueName: queueName,
		logger:    logger,
	}
}

func (s *baseActivationService) Process(ctx context.Context, alarm types.Alarm, eventRealTimestamp datetime.MicroTime) (bool, error) {
	if alarm.CanActivate() {
		err := s.sendActivationEvent(ctx, alarm, eventRealTimestamp)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (s *baseActivationService) sendActivationEvent(ctx context.Context, alarm types.Alarm, eventReceivedTimestamp datetime.MicroTime) error {
	event := types.Event{
		Connector:         alarm.Value.Connector,
		ConnectorName:     alarm.Value.ConnectorName,
		Component:         alarm.Value.Component,
		Resource:          alarm.Value.Resource,
		Timestamp:         datetime.NewCpsTime(),
		ReceivedTimestamp: eventReceivedTimestamp,
		EventType:         types.EventTypeActivate,
		Initiator:         types.InitiatorSystem,
	}
	event.SourceType = event.DetectSourceType()
	body, err := s.encoder.Encode(event)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail encode activation event")
		return err
	}

	err = s.publisher.PublishWithContext(
		ctx,
		"",
		s.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		s.logger.Error().Err(err).Msg("fail publish activation event to FIFO")
		return err
	}

	return nil
}
