package action

import (
	"context"
	"errors"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

const MaxRetries = 5

type service struct {
	alarmAdapter            libalarm.Adapter
	scenarioInputChannel    chan<- ExecuteScenariosTask
	delayedScenarioManager  DelayedScenarioManager
	executionStorage        ScenarioExecutionStorage
	encoder                 encoding.Encoder
	decoder                 encoding.Decoder
	fifoChan                libamqp.Channel
	fifoExchange, fifoQueue string
	activationService       libalarm.ActivationService
	logger                  zerolog.Logger
}

// NewService gives the correct action adapter.
func NewService(
	alarmAdapter libalarm.Adapter,
	scenarioInputChan chan<- ExecuteScenariosTask,
	delayedScenarioManager DelayedScenarioManager,
	storage ScenarioExecutionStorage,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
	fifoChan libamqp.Channel,
	fifoExchange string,
	fifoQueue string,
	activationService libalarm.ActivationService,
	logger zerolog.Logger,
) Service {
	service := service{
		alarmAdapter:           alarmAdapter,
		scenarioInputChannel:   scenarioInputChan,
		delayedScenarioManager: delayedScenarioManager,
		executionStorage:       storage,
		fifoChan:               fifoChan,
		encoder:                encoder,
		decoder:                decoder,
		fifoExchange:           fifoExchange,
		fifoQueue:              fifoQueue,
		activationService:      activationService,
		logger:                 logger,
	}

	return &service
}

func (s *service) ListenScenarioFinish(parentCtx context.Context, channel <-chan ScenarioResult) {
	ctx, cancel := context.WithCancel(parentCtx)
	s.logger.Debug().Msg("start listen scenario results")

	go func() {
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				s.logger.Debug().Msg("scenario results listener is cancelled")
				return
			case result, ok := <-channel:
				if !ok {
					s.logger.Debug().Msg("scenario results channel closed")
					return
				}

				s.logger.Debug().Msgf("scenario for alarm = %s finished", result.Alarm.ID)
				// Fetch updated alarm from storage since task manager returns
				// updated alarm after one scenario and not after all scenarios.
				alarm, err := s.alarmAdapter.GetAlarmByAlarmId(ctx, result.Alarm.ID)
				if err != nil {
					s.logger.Error().Err(err).Msg("failed to fetch alarm")
					break
				}

				fifoAckEvent := result.FifoAckEvent
				if fifoAckEvent.EventType == "" {
					fifoAckEvent = types.Event{
						Connector:     alarm.Value.Connector,
						ConnectorName: alarm.Value.ConnectorName,
						Component:     alarm.Value.Component,
						Resource:      alarm.Value.Resource,
						SourceType:    result.EntityType,
						Alarm:         &alarm,
					}
				}

				activationSent := false
				if !alarm.IsResolved() && (result.Err == nil ||
					(result.Err != nil && len(result.ActionExecutions) > 0 &&
						result.ActionExecutions[len(result.ActionExecutions)-1].Action.Type == types.ActionTypeWebhook)) {
					// Send activation event
					ok, err = s.activationService.Process(ctx, alarm, fifoAckEvent)
					if err != nil {
						s.logger.Error().Err(err).Msg("failed to send activation")
						break
					}

					if ok {
						activationSent = true
					}
				}

				if !activationSent {
					s.sendEventToFifoAck(ctx, fifoAckEvent)
				}
			}
		}
	}()
}

func (s *service) Process(ctx context.Context, event *types.Event) error {
	if event.Alarm == nil || event.Entity == nil {
		s.sendEventToFifoAck(ctx, *event)
		return nil
	}

	fifoAckEvent := types.Event{
		EventType:          event.EventType,
		Connector:          event.Connector,
		ConnectorName:      event.ConnectorName,
		Component:          event.Component,
		Resource:           event.Resource,
		SourceType:         event.SourceType,
		Timestamp:          event.Timestamp,
		ReceivedTimestamp:  event.ReceivedTimestamp,
		Author:             event.Author,
		UserID:             event.UserID,
		Initiator:          event.Initiator,
		IsMetaAlarmUpdated: event.IsMetaAlarmUpdated,
	}

	alarm := *event.Alarm
	entity := *event.Entity

	switch event.AlarmChange.Type {
	case types.AlarmChangeTypePbhEnter, types.AlarmChangeTypePbhLeave,
		types.AlarmChangeTypePbhLeaveAndEnter:
		var err error
		if event.PbehaviorInfo.IsActive() {
			err = s.delayedScenarioManager.ResumeDelayedScenarios(ctx, alarm)
		} else {
			err = s.delayedScenarioManager.PauseDelayedScenarios(ctx, alarm)
		}
		if err != nil {
			return err
		}
	}

	if event.EventType == types.EventTypeRunDelayedScenario {
		additionalData := AdditionalData{}
		err := s.decoder.Decode([]byte(event.DelayedScenarioData), &additionalData)
		if err != nil {
			s.logger.Err(err).Msg("invalid additional data for delayed scenario")
		}
		s.scenarioInputChannel <- ExecuteScenariosTask{
			Alarm:             alarm,
			Entity:            entity,
			DelayedScenarioID: event.DelayedScenarioID,
			AdditionalData:    additionalData,
			FifoAckEvent:      fifoAckEvent,
		}

		return nil
	}

	triggers := event.AlarmChange.GetTriggers()
	if len(triggers) == 0 {
		var activated bool
		var err error
		if event.AlarmChange.Type != types.AlarmChangeTypeNone {
			activated, err = s.activationService.Process(ctx, alarm, *event)
			if err != nil {
				return err
			}
		}

		if !activated {
			s.sendEventToFifoAck(ctx, *event)
		}

		return nil
	}

	s.scenarioInputChannel <- ExecuteScenariosTask{
		Triggers: triggers,
		Alarm:    alarm,
		Entity:   entity,
		AdditionalData: AdditionalData{
			AlarmChangeType: string(event.AlarmChange.Type),
			Author:          event.Author,
			User:            event.UserID,
			Initiator:       event.Initiator,
			Output:          event.Output,
		},
		IsMetaAlarmUpdated:   event.IsMetaAlarmUpdated,
		IsInstructionMatched: event.IsInstructionMatched,
		FifoAckEvent:         fifoAckEvent,
	}

	return nil
}

func (s *service) ProcessAbandonedExecutions(ctx context.Context) error {
	abandonedExecutions, err := s.executionStorage.GetAbandoned(ctx)
	if err != nil {
		return err
	}

	for _, execution := range abandonedExecutions {
		alarm, err := s.alarmAdapter.GetOpenedAlarmByAlarmId(ctx, execution.AlarmID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				s.logger.Warn().Str("execution", execution.GetCacheKey()).Msg("Alarm for scenario execution doesn't exist or resolved. Execution will be removed")
				err = s.executionStorage.Del(ctx, execution.GetCacheKey())
				if err != nil {
					return err
				}
			}

			continue
		}

		completed := execution.ActionExecutions[len(execution.ActionExecutions)-1].Executed

		if completed {
			s.logger.Debug().Str("execution", execution.GetCacheKey()).Msg("Execution was completed. Execution will be removed")
			err = s.executionStorage.Del(ctx, execution.GetCacheKey())
			if err != nil {
				return err
			}

			continue
		}

		s.logger.Debug().Str("execution", execution.GetCacheKey()).Msg("continue abandoned scenario")
		s.scenarioInputChannel <- ExecuteScenariosTask{
			Alarm:                alarm,
			Entity:               execution.Entity,
			AdditionalData:       execution.AdditionalData,
			FifoAckEvent:         execution.FifoAckEvent,
			IsMetaAlarmUpdated:   execution.IsMetaAlarmUpdated,
			IsInstructionMatched: execution.IsInstructionMatched,

			AbandonedExecutionCacheKey: execution.GetCacheKey(),
		}
	}

	return nil
}

func (s *service) sendEventToFifoAck(ctx context.Context, event types.Event) {
	if event.Healthcheck {
		return
	}

	body, err := s.encoder.Encode(event)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to send fifo ack event: failed to encode fifo ack event")
		return
	}

	err = s.fifoChan.PublishWithContext(
		ctx,
		s.fifoExchange,
		s.fifoQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to send fifo ack event: failed to publish message")
	}
}
