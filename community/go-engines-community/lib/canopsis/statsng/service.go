package statsng

import (
	"context"
	"fmt"
	"runtime/trace"
	"time"

	libamqp "git.canopsis.net/canopsis/go-engines/lib/amqp"
	"git.canopsis.net/canopsis/go-engines/lib/api/watcherweather"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

type service struct {
	pubChannel      libamqp.Channel
	pubExchangeName string
	pubQueueName    string
	jsonEncoder     encoding.Encoder
	statsStore      watcherweather.StatsStore
	logger          zerolog.Logger
}

func NewService(pubChannel libamqp.Channel, pubExchangeName, pubQueueName string, jsonEncoder encoding.Encoder,
	statsStore watcherweather.StatsStore, logger zerolog.Logger) *service {
	s := service{
		pubChannel:      pubChannel,
		pubExchangeName: pubExchangeName,
		pubQueueName:    pubQueueName,
		jsonEncoder:     jsonEncoder,
		statsStore:      statsStore,
		logger:          logger,
	}
	return &s
}

func (s *service) setStats(ctx context.Context, timestamp types.CpsTime, state *types.AlarmStep, entityID string) error {
	stats := watcherweather.Stats{
		LastEvent: &timestamp,
	}
	if state != nil && state.Value != types.AlarmStateOK {
		stats.FailEventsCount = 1
		stats.LastFailEvent = &timestamp
	} else {
		stats.OKEventsCount = 1
	}
	return s.statsStore.SetStats(ctx, entityID, stats)
}

// ProcessAlarmChange processes an AlarmChange, and sends the corresponding
// statistic events to the statsng engine.
func (s *service) ProcessAlarmChange(ctx context.Context, change types.AlarmChange, timestamp types.CpsTime, alarm types.Alarm, entity types.Entity,
	author, eventType string) error {
	defer trace.StartRegion(ctx, "statsng.ProcessAlarmChange").End()

	var err error

	if change.Type == types.AlarmChangeTypePbhEnter {
		err = s.statsStore.ResetStats(ctx, entity.ID)
		if err != nil {
			s.logger.Warn().Err(err).Msg("reset stats")
		}
	} else if alarm.IsInActivePeriod() && eventType == types.EventTypeCheck {
		err = s.setStats(ctx, timestamp, alarm.Value.State, entity.ID)
		if err != nil {
			s.logger.Warn().Err(err).Msg("set stats")
		}
	}

	if change.Type == types.AlarmChangeTypeNone {
		return nil
	}

	switch change.Type {
	case types.AlarmChangeTypeCreate:
		err = s.statCounterInc(
			CounterAlarmsCreated,
			timestamp,
			alarm,
			entity,
			author,
		)
		if err != nil {
			s.logger.Warn().Err(err).Msg("error sending event")
			return err
		}
	case types.AlarmChangeTypeCancel:
		err = s.statCounterInc(
			CounterAlarmsCanceled,
			timestamp,
			alarm,
			entity,
			author,
		)
		if err != nil {
			s.logger.Warn().Err(err).Msg("error sending event")
			return err
		}
	case types.AlarmChangeTypeAck:
		if alarm.HasSingleAck() {
			err = s.statCounterInc(
				CounterAlarmsAcknowledged,
				timestamp,
				alarm,
				entity,
				author,
			)
			if err != nil {
				s.logger.Warn().Err(err).Msg("error sending event")
				return err
			}

			var activationDate types.CpsTime
			if alarm.Value.ActivationDate != nil {
				activationDate = *alarm.Value.ActivationDate
			} else {
				s.logger.Warn().Msg("alarm is acked but activation date is empty")
				activationDate = alarm.Value.CreationDate
			}
			err = s.statDuration(
				DurationAckTime,
				activationDate,
				timestamp.Sub(activationDate.Time),
				alarm,
				entity,
				author,
			)
			if err != nil {
				s.logger.Warn().Err(err).Msg("error sending event")
				return err
			}
		}
	}

	if alarm.Value.State != nil && change.PreviousState != alarm.Value.State.Value {
		err = s.handleStateChange(timestamp, change.PreviousState, change.PreviousStateChange, alarm, entity)
	}

	if alarm.Value.Status != nil && change.PreviousStatus != alarm.Value.Status.Value {
		err = s.handleStatusChange(timestamp, change.PreviousStatus, change.PreviousStatusChange, alarm, entity, author)
	}

	return err
}

func (s *service) handleStateChange(timestamp types.CpsTime, previousState types.CpsNumber, previousStateChange types.CpsTime, alarm types.Alarm, entity types.Entity) error {
	err := s.statStateInterval(
		StateIntervalTimeInState,
		timestamp,
		timestamp.Sub(previousStateChange.Time),
		previousState,
		alarm,
		entity,
		"",
	)
	if err != nil {
		return fmt.Errorf("error sending event: %v", err)
	}

	return nil
}

func (s *service) handleStatusChange(timestamp types.CpsTime, previousStatus types.CpsNumber, previousStatusChange types.CpsTime, alarm types.Alarm, entity types.Entity, author string) error {
	if alarm.Value.Status.Value == types.AlarmStatusFlapping {
		err := s.statCounterInc(
			CounterFlappingPeriods,
			timestamp,
			alarm,
			entity,
			author,
		)
		if err != nil {
			return fmt.Errorf("error sending event: %v", err)
		}
	}

	return nil
}

// statCounterInc sends a statcounterinc event.
func (s *service) statCounterInc(statName string, timestamp types.CpsTime, alarm types.Alarm, entity types.Entity, author string) error {
	event := types.NewEventFromAlarm(alarm)
	event.Entity = &entity
	event.Author = author

	event.EventType = types.EventTypeStatCounterInc
	event.Timestamp = timestamp
	event.StatName = statName

	return s.sendEvent(event)
}

// statDuration sends a statduration event.
func (s *service) statDuration(statName string, timestamp types.CpsTime, duration time.Duration, alarm types.Alarm, entity types.Entity, author string) error {
	event := types.NewEventFromAlarm(alarm)
	event.Entity = &entity
	event.Author = author

	duration_seconds := types.CpsNumber(duration.Seconds())
	event.EventType = types.EventTypeStatDuration
	event.Timestamp = timestamp
	event.StatName = statName
	event.Duration = &duration_seconds

	return s.sendEvent(event)
}

// statStateInterval sends a statduration event.
func (s *service) statStateInterval(statName string, timestamp types.CpsTime, duration time.Duration, state types.CpsNumber, alarm types.Alarm, entity types.Entity, author string) error {
	event := types.NewEventFromAlarm(alarm)
	event.Entity = &entity
	event.Author = author

	duration_seconds := types.CpsNumber(duration.Seconds())
	event.EventType = types.EventTypeStatStateInterval
	event.Timestamp = timestamp
	event.StatName = statName
	event.Duration = &duration_seconds
	event.State = state

	return s.sendEvent(event)
}

// sendEvent sends a statistic event.
func (s *service) sendEvent(event types.Event) error {
	jevt, err := s.jsonEncoder.Encode(event)
	if err != nil {
		return fmt.Errorf("unable to serialize stat event: %v", err)
	}

	err = s.pubChannel.Publish(
		s.pubExchangeName,
		s.pubQueueName,
		false,
		false,
		amqp.Publishing{
			Body:        jevt,
			ContentType: "application/json",
		},
	)
	if err != nil {
		return fmt.Errorf("unable to send stat event: %v", err)
	}

	return nil
}

func (s *service) ProcessResolvedAlarm(alarm types.Alarm, entity types.Entity) error {
	if alarm.Value.Resolved == nil {
		return fmt.Errorf("trying to process unresolved alarm")
	}

	err := s.statCounterInc(
		CounterAlarmsResolved,
		*alarm.Value.Resolved,
		alarm,
		entity,
		"",
	)
	if err != nil {
		return err
	}

	var activationDate types.CpsTime
	if alarm.Value.ActivationDate != nil {
		activationDate = *alarm.Value.ActivationDate
	} else {
		s.logger.Warn().Msg("alarm is resolved but activation date is empty")
		activationDate = alarm.Value.CreationDate
	}
	err = s.statDuration(
		DurationResolveTime,
		activationDate,
		alarm.Value.Resolved.Sub(activationDate.Time),
		alarm,
		entity,
		"",
	)
	if err != nil {
		return err
	}

	err = s.statStateInterval(
		StateIntervalTimeInState,
		*alarm.Value.Resolved,
		alarm.Value.Resolved.Sub(alarm.Value.State.Timestamp.Time),
		types.AlarmStateOK,
		alarm,
		entity,
		"",
	)
	if err != nil {
		return err
	}

	return nil
}
