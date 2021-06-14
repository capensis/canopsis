package alarm

import (
	amqplib "git.canopsis.net/canopsis/go-engines/lib/amqp"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"time"
)

// AlarmActivationService checks alarm and sends activation event
// if alarm doesn't have active snooze and pbehavior
type ActivationService interface {
	Process(*types.Alarm) (bool, error)
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

func (s *baseActivationService) Process(alarm *types.Alarm) (bool, error) {
	if !alarm.IsActivated() && !alarm.IsSnoozed() && alarm.Value.PbehaviorInfo.IsActive() {
		err := s.sendActivationEvent(alarm)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (s *baseActivationService) sendActivationEvent(alarm *types.Alarm) error {
	event := types.Event{
		Connector:     alarm.Value.Connector,
		ConnectorName: alarm.Value.ConnectorName,
		Component:     alarm.Value.Component,
		Resource:      alarm.Value.Resource,
		Timestamp:     types.CpsTime{Time: time.Now()},
		EventType:     types.EventTypeActivate,
	}

	if alarm.Value.Resource == "" {
		event.SourceType = types.SourceTypeComponent
	} else {
		event.SourceType = types.SourceTypeResource
	}

	body, err := s.encoder.Encode(event)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail encode activation event")
		return err
	}

	err = s.publisher.Publish(
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
