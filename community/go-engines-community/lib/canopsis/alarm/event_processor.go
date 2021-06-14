package alarm

import (
	"context"
	"fmt"
	"runtime/trace"
	"sync"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm"
	liboperation "git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/errt"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
)

const MaxRedisLockRetries = 10

type eventProcessor struct {
	adapter         Adapter
	ruleAdapter     metaalarm.RulesAdapter
	redisLockClient redis.LockClient
	cfg             config.CanopsisConf
	executor        liboperation.Executor
	logger          zerolog.Logger
}

func NewEventProcessor(
	adapter Adapter,
	ruleAdapter metaalarm.RulesAdapter,
	cfg config.CanopsisConf,
	executor liboperation.Executor,
	redisLockClient redis.LockClient,
	logger zerolog.Logger,
) EventProcessor {
	return &eventProcessor{
		adapter:         adapter,
		ruleAdapter:     ruleAdapter,
		cfg:             cfg,
		executor:        executor,
		redisLockClient: redisLockClient,
		logger:          logger,
	}
}

func (s *eventProcessor) Process(ctx context.Context, event *types.Event) (types.AlarmChange, error) {
	defer trace.StartRegion(ctx, "alarm.ProcessAlarmEvent").End()

	alarmChange := types.NewAlarmChange()

	if event.Entity != nil && !event.Entity.Enabled {
		return alarmChange, nil
	}

	alarm, err := s.adapter.GetOpenedAlarm(event.Connector, event.ConnectorName, event.GetEID())
	alarmNotFound := false
	if _, ok := err.(errt.NotFound); ok {
		alarmNotFound = true
	} else if err != nil {
		return alarmChange, err
	}

	if !alarmNotFound {
		event.Alarm = &alarm
	}

	if err := s.fillAlarmChange(event, &alarmChange); err != nil {
		return alarmChange, err
	}

	switch event.EventType {
	case types.EventTypeCheck, types.EventTypeWatcher:
		changeType, err := s.storeAlarm(event)
		if changeType == types.AlarmChangeTypeStateIncrease || changeType == types.AlarmChangeTypeStateDecrease {
			s.updateMetaChildrenState(event)
		}

		alarmChange.Type = changeType
		return alarmChange, err
	case types.EventTypeMetaAlarm:
		changeType, err := s.processMetaAlarmCreateEvent(event)
		alarmChange.Type = changeType
		return alarmChange, err
	}

	if event.Alarm == nil {
		return alarmChange, nil
	}

	operation := s.createOperationFromEvent(event)
	changeType, err := s.executor.Exec(operation, event.Alarm, event.Timestamp, event.Role, event.Initiator)
	if err != nil {
		return alarmChange, err
	}

	if err := s.processAckResources(event, operation); err != nil {
		return alarmChange, err
	}

	if err := s.processMetaAlarmChildren(event, changeType, operation); err != nil {
		return alarmChange, err
	}

	if err := s.handleMetaAlarmChildResolve(event, changeType); err != nil {
		return alarmChange, err
	}

	alarmChange.Type = changeType

	return alarmChange, nil
}

func (s *eventProcessor) fillAlarmChange(event *types.Event, alarmChange *types.AlarmChange) error {
	alarm := event.Alarm
	if alarm == nil {
		lastAlarm, err := s.adapter.GetLastAlarm(event.Connector, event.ConnectorName, event.GetEID())
		notFound := false
		if _, ok := err.(errt.NotFound); ok {
			notFound = true
		} else if err != nil {
			return err
		}

		if !notFound && lastAlarm.Value.Resolved != nil {
			alarmChange.PreviousStateChange = *lastAlarm.Value.Resolved
			alarmChange.PreviousStatusChange = *lastAlarm.Value.Resolved
		} else {
			alarmChange.PreviousStateChange = event.Timestamp
			alarmChange.PreviousStatusChange = event.Timestamp
		}
	} else {
		alarmChange.PreviousState = alarm.Value.State.Value
		alarmChange.PreviousStateChange = alarm.Value.State.Timestamp
		alarmChange.PreviousStatus = alarm.Value.Status.Value
		alarmChange.PreviousStatusChange = alarm.Value.Status.Timestamp
	}

	return nil
}

func (s *eventProcessor) storeAlarm(event *types.Event) (types.AlarmChangeType, error) {
	changeType := types.AlarmChangeTypeNone
	if event.Alarm == nil && event.State == types.AlarmStateOK {
		return changeType, nil
	}

	if event.Alarm == nil {
		return s.createAlarm(event)
	}

	return s.updateAlarm(event)
}

func (s *eventProcessor) createAlarm(event *types.Event) (types.AlarmChangeType, error) {
	changeType := types.AlarmChangeTypeNone
	now := types.CpsTime{Time: time.Now()}
	alarm := &types.Alarm{
		EntityID: event.GetEID(),
		ID:       utils.NewID(),
		Time:     now,
		Value: types.AlarmValue{
			Component:         event.Component,
			Connector:         event.Connector,
			ConnectorName:     event.ConnectorName,
			CreationDate:      now,
			DisplayName:       types.GenDisplayName(s.cfg.Alarm.DisplayNameScheme),
			InitialOutput:     event.Output,
			Output:            event.Output,
			InitialLongOutput: event.LongOutput,
			LongOutput:        event.LongOutput,
			LongOutputHistory: []string{event.LongOutput},
			LastUpdateDate:    event.Timestamp,
			LastEventDate:     now,
			Resource:          event.Resource,
			Tags:              []string{},
			Parents:           []string{},
			Children:          []string{},
			Extra:             map[string]interface{}{},
			Infos:             map[string]map[string]interface{}{},
			RuleVersion:       map[string]string{},
		},
	}

	err := alarm.PartialUpdateState(event.Timestamp, event.State, event.Output, s.cfg)
	if err != nil {
		return changeType, err
	}

	if event.PbehaviorInfo.IsDefaultActive() {
		changeType = types.AlarmChangeTypeCreate
	} else {
		output := fmt.Sprintf(
			"Pbehavior %s. Type: %s. Reason: %s",
			event.PbehaviorInfo.Name,
			event.PbehaviorInfo.TypeName,
			event.PbehaviorInfo.Reason,
		)

		err := alarm.PartialUpdatePbhEnter(event.Timestamp, event.PbehaviorInfo,
			event.Author, output, event.Role, event.Initiator)
		if err != nil {
			return changeType, err
		}

		changeType = types.AlarmChangeTypeCreateAndPbhEnter
	}

	err = s.adapter.Insert(*alarm)
	if err != nil {
		return changeType, err
	}

	event.Alarm = alarm

	return changeType, nil
}

// updateAlarm updates alarm value and crops steps.
// TODO use mongo transactions after migration to mongo v4 because steps crop can override adding step by engine-webhook and engine-correlation.
func (s *eventProcessor) updateAlarm(event *types.Event) (types.AlarmChangeType, error) {
	changeType := types.AlarmChangeTypeNone
	alarm, err := s.adapter.GetOpenedAlarmByAlarmId(event.Alarm.ID)
	if err != nil {
		return changeType, err
	}

	previousState := alarm.CurrentState()
	newState := event.State
	err = alarm.PartialUpdateState(event.Timestamp, event.State, event.Output, s.cfg)
	if err != nil {
		return changeType, err
	}

	alarm.UpdateOutput(event.Output)
	alarm.UpdateLongOutput(event.LongOutput)

	if s.cfg.Alarm.EnableLastEventDate {
		alarm.PartialUpdateLastEventDate(event.Timestamp)
	}

	err = s.adapter.PartialUpdateOpen(&alarm)
	if err != nil {
		return changeType, err
	}

	// Update cropped steps if needed
	alarm.PartialUpdateCropSteps()
	err = s.adapter.PartialUpdateOpen(&alarm)
	if err != nil {
		return changeType, err
	}

	event.Alarm = &alarm

	if newState > previousState {
		changeType = types.AlarmChangeTypeStateIncrease
	} else if newState < previousState {
		changeType = types.AlarmChangeTypeStateDecrease
	}

	return changeType, nil
}

func (s *eventProcessor) createOperationFromEvent(event *types.Event) types.Operation {
	var parameters interface{}
	switch event.EventType {
	case types.EventTypeAssocTicket:
		parameters = types.OperationAssocTicketParameters{
			Ticket: event.Ticket,
			Output: event.Output,
			Author: event.Author,
		}
	case types.EventTypeSnooze:
		parameters = types.OperationSnoozeParameters{
			Duration: types.DurationWithUnit{
				Seconds: int64(*event.Duration),
			},
			Output: event.Output,
			Author: event.Author,
		}
	case types.EventTypeChangestate, types.EventTypeKeepstate:
		parameters = types.OperationChangeStateParameters{
			State:  event.State,
			Output: event.Output,
			Author: event.Author,
		}
	case types.EventTypePbhEnter, types.EventTypePbhLeave, types.EventTypePbhLeaveAndEnter:
		parameters = types.OperationPbhParameters{
			PbehaviorInfo: event.PbehaviorInfo,
			Output:        event.Output,
			Author:        event.Author,
		}
	default:
		parameters = types.OperationParameters{
			Output: event.Output,
			Author: event.Author,
		}
	}

	return types.Operation{
		Type:       event.EventType,
		Parameters: parameters,
	}
}

func (s *eventProcessor) processAckResources(event *types.Event, operation types.Operation) error {
	if event.EventType != types.EventTypeAck || !event.AckResources {
		return nil
	}

	alarms, err := s.adapter.GetUnacknowledgedAlarmsByComponent(event.Component)
	if err != nil {
		return err
	}

	for _, alarm := range alarms {
		_, err := s.executor.Exec(operation, &alarm, event.Timestamp, event.Role, event.Initiator)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *eventProcessor) processMetaAlarmCreateEvent(event *types.Event) (types.AlarmChangeType, error) {
	var childAlarms []types.Alarm
	if event.MetaAlarmChildren != nil {
		err := s.adapter.GetOpenedAlarmsByIDs(*event.MetaAlarmChildren, &childAlarms)
		if err != nil {
			s.logger.Err(err).Msg("error on geting meta-alarm children")
			return types.AlarmChangeTypeNone, err
		}
		worstState := types.CpsNumber(types.AlarmStateMinor)

		for i := 0; i < len(childAlarms); i++ {
			c := childAlarms[i]
			if c.Value.State != nil && c.Value.State.Value > worstState {
				worstState = c.Value.State.Value
			}
		}
		event.State = worstState
	}
	metaAlarm, err := types.NewAlarm(*event, s.cfg)
	if err != nil {
		s.logger.Err(err).Msg("error on creating new meta-alarm")
		return types.AlarmChangeTypeNone, err
	}

	metaAlarm.Value.Tags = []string{}
	metaAlarm.Value.Extra = map[string]interface{}{}
	metaAlarm.Value.Infos = map[string]map[string]interface{}{}
	metaAlarm.Value.RuleVersion = map[string]string{}
	metaAlarm.Value.Parents = []string{}
	metaAlarm.Value.Children = []string{}

	if event.MetaAlarmChildren != nil {
		for i := 0; i < len(childAlarms); i++ {
			c := childAlarms[i]
			metaAlarm.Value.Children = append(metaAlarm.Value.Children, c.EntityID)
			c.Value.Parents = append(c.Value.Parents, metaAlarm.EntityID)
			childAlarms[i] = c
		}
	}
	metaAlarm.Value.Meta = event.MetaAlarmRuleID
	metaAlarm.Value.MetaValuePath = event.MetaAlarmValuePath

	if event.ExtraInfos != nil {
		if nameInf, ok := event.ExtraInfos["display_name"]; ok {
			if name, isStr := nameInf.(string); isStr {
				metaAlarm.Value.DisplayName = name
			}
		}
	}

	err = s.adapter.Insert(metaAlarm)
	if err != nil {
		s.logger.Err(err).Msg("error on inserting new meta-alarm to db")
		return types.AlarmChangeTypeNone, err
	}

	ruleIdentifier := metaAlarm.Value.Meta
	rule, err := s.ruleAdapter.GetRule(metaAlarm.Value.Meta)
	if err != nil {
		// the rule can be deleted
		if err.Error() != "not found" {
			s.logger.Err(err).Str("rule", metaAlarm.Value.Meta).Msg("Get rule had failed")
		}
	} else {
		ruleIdentifier = rule.Name
	}
	newStep := types.NewMetaAlarmAttachStep(metaAlarm, ruleIdentifier)

	for i := 0; i < len(childAlarms); i++ {
		c := childAlarms[i]
		err := c.Value.Steps.Add(newStep)
		if err != nil {
			s.logger.Err(err).Str("metaalarm", metaAlarm.EntityID).
				Str("child", c.EntityID).
				Msg("Failed to add metaalarmattach step to child")
		}
		childAlarms[i] = c
	}

	err = s.adapter.MassUpdate(childAlarms, false)
	if err != nil {
		return types.AlarmChangeTypeNone, err
	}

	event.Alarm = &metaAlarm
	return types.AlarmChangeTypeCreate, nil
}

func (s *eventProcessor) processMetaAlarmChildren(event *types.Event, changeType types.AlarmChangeType, operation types.Operation) error {
	if !event.Alarm.IsMetaAlarm() ||
		changeType != types.AlarmChangeTypeAck &&
			changeType != types.AlarmChangeTypeAckremove &&
			changeType != types.AlarmChangeTypeAssocTicket &&
			changeType != types.AlarmChangeTypeCancel &&
			changeType != types.AlarmChangeTypeChangeState &&
			changeType != types.AlarmChangeTypeComment &&
			changeType != types.AlarmChangeTypeDone &&
			changeType != types.AlarmChangeTypeSnooze &&
			changeType != types.AlarmChangeTypeUncancel &&
			changeType != types.AlarmChangeTypeUpdateStatus {

		return nil
	}

	var alarms []types.Alarm
	err := s.adapter.GetOpenedAlarmsByIDs(event.Alarm.Value.Children, &alarms)
	if err != nil {
		s.logger.Error().Err(err).Msg("error getting meta-alarm children")
		return err
	}
	for _, alarm := range alarms {
		_, err := s.executor.Exec(operation, &alarm, event.Timestamp, event.Role, event.Initiator)
		if err != nil {
			s.logger.Error().Err(err).Msg("error updating meta-alarm child alarm")
			return err
		}
	}

	return nil
}

func (s *eventProcessor) handleMetaAlarmChildResolve(event *types.Event, changeType types.AlarmChangeType) error {
	if event.Alarm.IsMetaAlarm() || changeType != types.AlarmChangeTypeResolve ||
		len(event.Alarm.Value.Parents) == 0 {

		return nil
	}

	var alarms []types.Alarm
	if err := s.adapter.GetOpenedAlarmsByIDs(event.Alarm.Value.Parents, &alarms); err != nil {
		return err
	}

	if len(alarms) == 0 {
		s.logger.Debug().Msg("No opening parent alarms")
		return nil
	}

	var wg sync.WaitGroup
	for _, alarm := range alarms {
		wg.Add(1)
		go func(alarm types.Alarm) {
			if err := s.handleAutoResolveMetaAlarm(alarm); err != nil {
				s.logger.Err(err).
					Str("alarm-id", alarm.ID).
					Msg("handle auto resolve parent alarm had error")
			}
			wg.Done()
		}(alarm)
	}
	wg.Wait()

	return nil
}

func (s *eventProcessor) handleAutoResolveMetaAlarm(alarm types.Alarm) error {
	rule, err := s.ruleAdapter.GetRule(alarm.Value.Meta)
	if err != nil {
		return err
	}
	if !rule.AutoResolve {
		s.logger.Info().Str("rule", alarm.Value.Meta).
			Msg("Metaalarm rule is no auto resolve")
		return nil
	}

	if alarm.IsResolved() {
		return fmt.Errorf("alarm had already resolved")
	}

	metaAlarmLock, err := s.redisLockClient.Obtain(alarm.ID, 100*time.Millisecond, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(11*time.Millisecond), MaxRedisLockRetries),
	})
	if err != nil {
		return err
	}

	defer func() {
		if metaAlarmLock != nil {
			err := metaAlarmLock.Release()
			if err != nil && err != redislock.ErrLockNotHeld {
				s.logger.Warn().
					Str("alarm_id", alarm.ID)
			}
		}
	}()

	c, err := s.adapter.CountResolvedAlarm(alarm.Value.Children)
	if err != nil {
		return err
	}

	if c == len(alarm.Value.Children) {
		s.logger.Info().Str("metaalarm", alarm.AlarmID()).
			Msg("All children of metalarm has been resolved.")
		alarm.Resolve(&types.CpsTime{
			Time: time.Now(),
		})
		return s.adapter.Update(alarm)
	}

	return nil
}

// updateMetaChildrenState updates alarm's parents according with children state if need
func (s *eventProcessor) updateMetaChildrenState(event *types.Event) {
	if !event.Alarm.IsMetaChildren() {
		return
	}
	var parents []types.Alarm
	err := s.adapter.GetOpenedAlarmsByIDs(event.Alarm.Value.Parents, &parents)
	if err == nil {
		s.logger.Debug().Msgf("change child's %v state of meta-alarms %+v", event.Alarm, parents)
		updatedParents := make([]types.Alarm, 0, len(parents))
		for _, metaAlarm := range parents {
			maCurrentState := metaAlarm.CurrentState()
			if event.Alarm.Value.State.Value > maCurrentState {
				metaAlarm.UpdateState(event.Alarm.Value.State.Value, event.Alarm.Value.LastUpdateDate)
				updatedParents = append(updatedParents, metaAlarm)
			} else if event.Alarm.Value.State.Value < maCurrentState {
				if UpdateToWorstState(&metaAlarm, []*types.Alarm{event.Alarm}, s.adapter, s.cfg) {
					updatedParents = append(updatedParents, metaAlarm)
				}
			}
		}
		if len(updatedParents) > 0 {
			err = s.adapter.MassUpdate(updatedParents, true)
		}
	}
	if err != nil {
		s.logger.Error().Err(err).Msgf("error changestate meta-alarm from children %+v", event.Alarm)
	}
}
