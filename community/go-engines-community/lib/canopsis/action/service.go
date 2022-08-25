package action

import (
	"context"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

const AbandonedDuration = 60
const MaxRetries = 5

type service struct {
	alarmAdapter            libalarm.Adapter
	scenarioInputChannel    chan<- ExecuteScenariosTask
	delayedScenarioManager  DelayedScenarioManager
	executionStorage        ScenarioExecutionStorage
	encoder                 encoding.Encoder
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

				s.logger.Debug().Msgf("scenario for alarm_id = %s finished", result.Alarm.ID)
				// Fetch updated alarm from storage since task manager returns
				// updated alarm after one scenario and not after all scenarios.
				alarm, err := s.alarmAdapter.GetAlarmByAlarmId(ctx, result.Alarm.ID)
				if err != nil {
					s.logger.Error().Err(err).Msg("failed to fetch alarm")
					break
				}

				event := &types.Event{
					Connector:         alarm.Value.Connector,
					ConnectorName:     alarm.Value.ConnectorName,
					Component:         alarm.Value.Component,
					Resource:          alarm.Value.Resource,
					Alarm:             &alarm,
					MetaAlarmParents:  &alarm.Value.Parents,
					MetaAlarmChildren: &alarm.Value.Children,
					// need it for fifo metaalarm lock
					MetaAlarmRelatedParents: result.Alarm.Value.RelatedParents,
				}

				activationSent := false
				if !alarm.IsResolved() && (result.Err == nil ||
					(result.Err != nil && len(result.ActionExecutions) > 0 &&
						result.ActionExecutions[len(result.ActionExecutions)-1].Action.Type == types.ActionTypeWebhook)) {
					// Send activation event
					ok, err = s.activationService.Process(&alarm)
					if err != nil {
						s.logger.Error().Err(err).Msg("failed to send activation")
						break
					}

					if ok {
						activationSent = true
					}
				}

				if !activationSent {
					s.sendEventToFifoAck(event)
				}
			}
		}
	}()
}

func (s *service) Process(ctx context.Context, event *types.Event) error {
	if event.Alarm == nil || event.Entity == nil {
		s.sendEventToFifoAck(event)

		return nil
	}

	alarm := *event.Alarm
	entity := *event.Entity

	// need it for fifo metaalarm lock
	alarm.Value.RelatedParents = event.MetaAlarmRelatedParents

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

	additionalData := AdditionalData{
		AlarmChangeType: event.AlarmChange.Type,
		Author:          event.Author,
		Initiator:       event.Initiator,
	}

	if event.EventType == types.EventTypeRunDelayedScenario {
		s.scenarioInputChannel <- ExecuteScenariosTask{
			Alarm:             alarm,
			Entity:            entity,
			DelayedScenarioID: event.DelayedScenarioID,
			AdditionalData:    additionalData,
		}

		return nil
	}

	triggers := event.AlarmChange.GetTriggers()
	if len(triggers) == 0 {
		s.sendEventToFifoAck(event)
		return nil
	}

	s.scenarioInputChannel <- ExecuteScenariosTask{
		Triggers:       triggers,
		Alarm:          alarm,
		Entity:         entity,
		AckResources:   event.AckResources,
		AdditionalData: additionalData,
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
			if err == mongo.ErrNoDocuments {
				s.logger.Warn().Str("execution_id", execution.ID).Msg("Alarm for scenario execution doesn't exist or resolved. Execution will be removed")
				err = s.executionStorage.Del(ctx, execution.ID)
				if err != nil {
					return err
				}
			}

			continue
		}

		completed := execution.ActionExecutions[len(execution.ActionExecutions)-1].Executed

		if completed {
			s.logger.Debug().Str("execution_id", execution.ID).Msg("Execution was completed. Execution will be removed")
			err = s.executionStorage.Del(ctx, execution.ID)
			if err != nil {
				return err
			}

			continue
		}

		s.scenarioInputChannel <- ExecuteScenariosTask{
			Alarm:                alarm,
			Entity:               execution.Entity,
			AbandonedExecutionID: execution.ID,
			AdditionalData:       execution.AdditionalData,
		}
	}

	return nil
}

func (s *service) sendEventToFifoAck(event *types.Event) {
	body, err := s.encoder.Encode(event)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to send fifo ack event: failed to encode fifo ack event")
		return
	}

	err = s.fifoChan.Publish(
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
