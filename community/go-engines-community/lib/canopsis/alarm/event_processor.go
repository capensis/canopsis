package alarm

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"runtime/trace"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	liboperation "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
)

const MaxRedisLockRetries = 10

type eventProcessor struct {
	adapter             Adapter
	entityAdapter       entity.Adapter
	ruleAdapter         correlation.RulesAdapter
	redisLockClient     redis.LockClient
	alarmConfigProvider config.AlarmConfigProvider
	executor            liboperation.Executor
	logger              zerolog.Logger
}

func NewEventProcessor(
	adapter Adapter,
	entityAdapter entity.Adapter,
	ruleAdapter correlation.RulesAdapter,
	alarmConfigProvider config.AlarmConfigProvider,
	executor liboperation.Executor,
	redisLockClient redis.LockClient,
	logger zerolog.Logger,
) EventProcessor {
	return &eventProcessor{
		adapter:             adapter,
		entityAdapter:       entityAdapter,
		ruleAdapter:         ruleAdapter,
		alarmConfigProvider: alarmConfigProvider,
		executor:            executor,
		redisLockClient:     redisLockClient,
		logger:              logger,
	}
}

func (s *eventProcessor) Process(ctx context.Context, event *types.Event) (types.AlarmChange, error) {
	defer trace.StartRegion(ctx, "alarm.ProcessAlarmEvent").End()

	alarmChange := types.NewAlarmChange()

	if event.Entity != nil && !event.Entity.Enabled {
		if event.EventType == types.EventTypeEntityToggled ||
			event.EventType == types.EventTypeRecomputeEntityService {
			return s.resolveAlarmForDisabledEntity(ctx, event)
		}

		return alarmChange, nil
	}

	alarm, err := s.adapter.GetOpenedAlarm(ctx, event.Connector, event.ConnectorName, event.GetEID())
	alarmNotFound := false
	if _, ok := err.(errt.NotFound); ok {
		alarmNotFound = true
	} else if err != nil {
		return alarmChange, err
	}

	if !alarmNotFound {
		event.Alarm = &alarm
	}

	if err := s.fillAlarmChange(ctx, event, &alarmChange); err != nil {
		return alarmChange, err
	}

	switch event.EventType {
	case types.EventTypeCheck:
		changeType, err := s.storeAlarm(ctx, event)
		if changeType == types.AlarmChangeTypeStateIncrease || changeType == types.AlarmChangeTypeStateDecrease {
			s.updateMetaChildrenState(ctx, event)
		} else if event.Alarm != nil && event.Alarm.IsMetaChildren() &&
			s.alarmConfigProvider.Get().EnableLastEventDate {
			s.updateMetaLastEventDate(ctx, event)
		}

		alarmChange.Type = changeType
		return alarmChange, err
	case types.EventTypeNoEvents:
		changeType, err := s.processNoEvents(ctx, event)
		alarmChange.Type = changeType
		return alarmChange, err
	case types.EventTypeMetaAlarm:
		changeType, err := s.processMetaAlarmCreateEvent(ctx, event)
		alarmChange.Type = changeType
		return alarmChange, err
	}

	if event.Alarm == nil {
		return alarmChange, nil
	}

	operation := s.createOperationFromEvent(event)
	changeType, err := s.executor.Exec(ctx, operation, event.Alarm, event.Timestamp, event.Role, event.Initiator)
	if err != nil {
		return alarmChange, err
	}

	if changeType == types.AlarmChangeTypeResolve {
		err := s.adapter.CopyAlarmToResolvedCollection(ctx, *event.Alarm)
		if err != nil {
			return alarmChange, err
		}
	}

	if event.IdleRuleApply != "" {
		event.Entity.LastIdleRuleApply = event.IdleRuleApply
		err := s.entityAdapter.UpdateIdleFields(ctx, event.Entity.ID, event.Entity.IdleSince,
			event.Entity.LastIdleRuleApply)
		if err != nil {
			return alarmChange, err
		}
	}

	if err = s.processAckResources(ctx, event, operation); err != nil {
		return alarmChange, err
	}

	if err = s.processMetaAlarmChildren(ctx, event, changeType, operation); err != nil {
		return alarmChange, err
	}

	resolved := false
	if resolved, err = s.handleMetaAlarmChildResolve(ctx, event, changeType); err != nil {
		return alarmChange, err
	}

	if !resolved && event.Alarm.IsMetaChildren() && s.alarmConfigProvider.Get().EnableLastEventDate {
		s.updateMetaLastEventDate(ctx, event)
	}

	alarmChange.Type = changeType

	return alarmChange, nil
}

func (s *eventProcessor) fillAlarmChange(ctx context.Context, event *types.Event, alarmChange *types.AlarmChange) error {
	alarm := event.Alarm
	if alarm == nil {
		lastAlarm, err := s.adapter.GetLastAlarm(ctx, event.Connector, event.ConnectorName, event.GetEID())
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
		alarmChange.PreviousPbehaviorTypeID = alarm.Value.PbehaviorInfo.TypeID
		alarmChange.PreviousPbehaviorCannonicalType = alarm.Value.PbehaviorInfo.CanonicalType
	}

	return nil
}

func (s *eventProcessor) storeAlarm(ctx context.Context, event *types.Event) (types.AlarmChangeType, error) {
	changeType := types.AlarmChangeTypeNone
	if event.Alarm == nil && event.State == types.AlarmStateOK {
		return changeType, nil
	}

	if event.Alarm == nil {
		return s.createAlarm(ctx, event)
	}

	return s.updateAlarm(ctx, event)
}

func (s *eventProcessor) createAlarm(ctx context.Context, event *types.Event) (types.AlarmChangeType, error) {
	changeType := types.AlarmChangeTypeNone

	if event.Entity == nil {
		return changeType, nil
	}

	alarmConfig := s.alarmConfigProvider.Get()
	alarm := newAlarm(*event, alarmConfig)
	err := alarm.PartialUpdateState(event.Timestamp, event.State, event.Output, alarmConfig)
	if err != nil {
		return changeType, err
	}

	alarm.PartialUpdateEventsCount()

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

	err = s.adapter.Insert(ctx, alarm)
	if err != nil {
		return changeType, err
	}

	event.Alarm = &alarm

	return changeType, nil
}

// updateAlarm updates alarm value and crops steps.
// TODO use mongo transactions after migration to mongo v4 because steps crop can override adding step by engine-webhook and engine-correlation.
func (s *eventProcessor) updateAlarm(ctx context.Context, event *types.Event) (types.AlarmChangeType, error) {
	changeType := types.AlarmChangeTypeNone
	alarm, err := s.adapter.GetOpenedAlarmByAlarmId(ctx, event.Alarm.ID)
	if err != nil {
		return changeType, err
	}

	alarmConfig := s.alarmConfigProvider.Get()
	previousState := alarm.CurrentState()
	newState := event.State
	err = alarm.PartialUpdateState(event.Timestamp, event.State, event.Output, alarmConfig)
	if err != nil {
		return changeType, err
	}

	alarm.UpdateOutput(event.Output)
	alarm.UpdateLongOutput(event.LongOutput)
	alarm.PartialUpdateEventsCount()

	if alarmConfig.EnableLastEventDate {
		alarm.PartialUpdateLastEventDate(event.Timestamp)
	}

	err = s.adapter.PartialUpdateOpen(ctx, &alarm)
	if err != nil {
		return changeType, err
	}

	// Update cropped steps if needed
	alarm.PartialUpdateCropSteps()
	err = s.adapter.PartialUpdateOpen(ctx, &alarm)
	if err != nil {
		return changeType, err
	}

	event.Alarm = &alarm

	if newState > previousState {
		changeType = types.AlarmChangeTypeStateIncrease
	} else if newState < previousState {
		changeType = types.AlarmChangeTypeStateDecrease
	}

	if event.Entity.IdleSince != nil || event.Entity.LastIdleRuleApply != "" {
		alarmLastUpdateRule := fmt.Sprintf("%s_%s", idlerule.RuleTypeAlarm, idlerule.RuleAlarmConditionLastUpdate)
		if event.Entity.LastIdleRuleApply == alarmLastUpdateRule {
			if changeType != types.AlarmChangeTypeNone {
				event.Entity.LastIdleRuleApply = ""
			}
		} else {
			event.Entity.LastIdleRuleApply = ""
		}

		event.Entity.IdleSince = nil
		err := s.entityAdapter.UpdateIdleFields(ctx, event.Entity.ID, event.Entity.IdleSince,
			event.Entity.LastIdleRuleApply)
		if err != nil {
			return changeType, err
		}
	}

	return changeType, nil
}

func (s *eventProcessor) processNoEvents(ctx context.Context, event *types.Event) (types.AlarmChangeType, error) {
	changeType := types.AlarmChangeTypeNone
	if event.Entity == nil || event.Alarm == nil && event.State == types.AlarmStateOK {
		return changeType, nil
	}

	alarmConfig := s.alarmConfigProvider.Get()

	if event.Alarm == nil {
		alarm := newAlarm(*event, alarmConfig)
		err := alarm.PartialUpdateNoEvents(event.State, event.Timestamp, event.Author,
			event.Output, event.Role, event.Initiator, alarmConfig)
		if err != nil {
			return changeType, err
		}

		err = s.adapter.Insert(ctx, alarm)
		if err != nil {
			return changeType, err
		}

		event.Alarm = &alarm
		changeType = types.AlarmChangeTypeCreate
	} else {
		alarm := event.Alarm
		previousState := alarm.CurrentState()
		previousStatus := alarm.CurrentStatus(alarmConfig)
		err := alarm.PartialUpdateNoEvents(event.State, event.Timestamp, event.Author,
			event.Output, event.Role, event.Initiator, alarmConfig)
		if err != nil {
			return changeType, err
		}

		err = s.adapter.PartialUpdateOpen(ctx, alarm)
		if err != nil {
			return changeType, err
		}

		newState := alarm.CurrentState()
		newStatus := alarm.CurrentStatus(alarmConfig)
		if newState > previousState {
			changeType = types.AlarmChangeTypeStateIncrease
		} else if newState < previousState {
			changeType = types.AlarmChangeTypeStateDecrease
		} else if newStatus != previousStatus {
			changeType = types.AlarmChangeTypeUpdateStatus
		}
	}

	if event.Alarm.Value.Status.Value == types.AlarmStatusNoEvents {
		event.Entity.IdleSince = &event.Timestamp
		event.Entity.LastIdleRuleApply = event.IdleRuleApply
	} else {
		event.Entity.IdleSince = nil
		event.Entity.LastIdleRuleApply = ""
	}

	err := s.entityAdapter.UpdateIdleFields(ctx, event.Entity.ID, event.Entity.IdleSince,
		event.Entity.LastIdleRuleApply)
	if err != nil {
		return changeType, err
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
	case types.EventTypeInstructionStarted, types.EventTypeInstructionPaused,
		types.EventTypeInstructionResumed, types.EventTypeInstructionCompleted,
		types.EventTypeInstructionFailed, types.EventTypeInstructionAborted,
		types.EventTypeAutoInstructionStarted, types.EventTypeAutoInstructionCompleted,
		types.EventTypeAutoInstructionFailed, types.EventTypeAutoInstructionAlreadyRunning,
		types.EventTypeInstructionJobStarted, types.EventTypeInstructionJobCompleted,
		types.EventTypeInstructionJobAborted, types.EventTypeInstructionJobFailed:
		parameters = types.OperationInstructionParameters{
			Execution: event.Execution,
			Output:    event.Output,
			Author:    event.Author,
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

func (s *eventProcessor) processAckResources(ctx context.Context, event *types.Event, operation types.Operation) error {
	if event.EventType != types.EventTypeAck || !event.AckResources {
		return nil
	}

	alarms, err := s.adapter.GetUnacknowledgedAlarmsByComponent(ctx, event.Component)
	if err != nil {
		return err
	}

	for _, alarm := range alarms {
		_, err := s.executor.Exec(ctx, operation, &alarm, event.Timestamp, event.Role, event.Initiator)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *eventProcessor) processMetaAlarmCreateEvent(ctx context.Context, event *types.Event) (types.AlarmChangeType, error) {
	var childAlarms []types.Alarm
	if event.MetaAlarmChildren != nil {
		err := s.adapter.GetOpenedAlarmsByIDs(ctx, *event.MetaAlarmChildren, &childAlarms)
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
	metaAlarm, err := types.NewAlarm(*event, s.alarmConfigProvider.Get())
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

	err = s.adapter.Insert(ctx, metaAlarm)
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

	err = s.adapter.MassUpdate(ctx, childAlarms, false)
	if err != nil {
		return types.AlarmChangeTypeNone, err
	}

	event.Alarm = &metaAlarm
	return types.AlarmChangeTypeCreate, nil
}

func (s *eventProcessor) processMetaAlarmChildren(ctx context.Context, event *types.Event, changeType types.AlarmChangeType, operation types.Operation) error {
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
	err := s.adapter.GetOpenedAlarmsByIDs(ctx, event.Alarm.Value.Children, &alarms)
	if err != nil {
		s.logger.Error().Err(err).Msg("error getting meta-alarm children")
		return err
	}
	for _, alarm := range alarms {
		_, err := s.executor.Exec(ctx, operation, &alarm, event.Timestamp, event.Role, event.Initiator)
		if err != nil {
			s.logger.Error().Err(err).Msg("error updating meta-alarm child alarm")
			return err
		}
	}

	return nil
}

func (s *eventProcessor) handleMetaAlarmChildResolve(ctx context.Context, event *types.Event, changeType types.AlarmChangeType) (bool, error) {
	if event.Alarm.IsMetaAlarm() || changeType != types.AlarmChangeTypeResolve ||
		len(event.Alarm.Value.Parents) == 0 {

		return false, nil
	}

	var alarms []types.Alarm
	if err := s.adapter.GetOpenedAlarmsByIDs(ctx, event.Alarm.Value.Parents, &alarms); err != nil {
		return false, err
	}

	if len(alarms) == 0 {
		s.logger.Debug().Msg("No opening parent alarms")
		return false, nil
	}

	var wg sync.WaitGroup
	for _, alarm := range alarms {
		wg.Add(1)
		go func(alarm types.Alarm) {
			if err := s.handleAutoResolveMetaAlarm(ctx, alarm); err != nil {
				s.logger.Err(err).
					Str("alarm-id", alarm.ID).
					Msg("handle auto resolve parent alarm had error")
			}
			wg.Done()
		}(alarm)
	}
	wg.Wait()

	return true, nil
}

func (s *eventProcessor) handleAutoResolveMetaAlarm(ctx context.Context, alarm types.Alarm) error {
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

	metaAlarmLock, err := s.redisLockClient.Obtain(ctx, alarm.ID, 100*time.Millisecond, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(11*time.Millisecond), MaxRedisLockRetries),
	})
	if err != nil {
		return err
	}

	defer func() {
		if metaAlarmLock != nil {
			err := metaAlarmLock.Release(ctx)
			if err != nil && err != redislock.ErrLockNotHeld {
				s.logger.Warn().
					Str("alarm_id", alarm.ID)
			}
		}
	}()

	c, err := s.adapter.CountResolvedAlarm(ctx, alarm.Value.Children)
	if err != nil {
		return err
	}

	if c == len(alarm.Value.Children) {
		s.logger.Info().Str("metaalarm", alarm.AlarmID()).
			Msg("All children of metalarm has been resolved.")
		alarm.Resolve(&types.CpsTime{
			Time: time.Now(),
		})

		err := s.adapter.CopyAlarmToResolvedCollection(ctx, alarm)
		if err != nil {
			return err
		}

		return s.adapter.Update(ctx, alarm)
	}

	return nil
}

// updateMetaChildrenState updates alarm's parents according with children state if need
func (s *eventProcessor) updateMetaChildrenState(ctx context.Context, event *types.Event) {
	if !event.Alarm.IsMetaChildren() {
		return
	}
	var parents []types.Alarm
	err := s.adapter.GetOpenedAlarmsByIDs(ctx, event.Alarm.Value.Parents, &parents)
	if err == nil {
		s.logger.Debug().Msgf("change child's %v state of meta-alarms %+v", event.Alarm, parents)
		updatedParents := make([]types.Alarm, 0, len(parents))
		for _, metaAlarm := range parents {
			maCurrentState := metaAlarm.CurrentState()
			if event.Alarm.Value.State.Value > maCurrentState {
				metaAlarm.UpdateState(event.Alarm.Value.State.Value, event.Alarm.Value.LastUpdateDate)
				updatedParents = append(updatedParents, metaAlarm)
			} else if event.Alarm.Value.State.Value < maCurrentState {
				if UpdateToWorstState(ctx, &metaAlarm, []*types.Alarm{event.Alarm}, s.adapter, s.alarmConfigProvider.Get()) {
					updatedParents = append(updatedParents, metaAlarm)
				}
			}
		}
		if len(updatedParents) > 0 {
			err = s.adapter.MassUpdate(ctx, updatedParents, true)
		}
	}
	if err != nil {
		s.logger.Error().Err(err).Msgf("error changestate meta-alarm from children %+v", event.Alarm)
	}
}

func (s *eventProcessor) resolveAlarmForDisabledEntity(ctx context.Context, event *types.Event) (types.AlarmChange, error) {
	alarmChange := types.NewAlarmChange()
	alarm, err := s.adapter.GetOpenedAlarm(ctx, event.Connector, event.ConnectorName, event.GetEID())
	if _, ok := err.(errt.NotFound); ok {
		return alarmChange, nil
	} else if err != nil {
		return alarmChange, err
	}

	if err := s.fillAlarmChange(ctx, event, &alarmChange); err != nil {
		return alarmChange, err
	}

	event.Alarm = &alarm
	operation := types.Operation{
		Type: types.EventTypeEntityToggled,
		Parameters: types.OperationParameters{
			Output: event.Output,
			Author: event.Author,
		},
	}
	changeType, err := s.executor.Exec(ctx, operation, event.Alarm, event.Timestamp, event.Role, event.Initiator)
	if err != nil {
		return alarmChange, err
	}

	alarmChange.Type = changeType

	return alarmChange, nil
}

// updateMetaLastEventDate updates alarm's parents LastEventDate
func (s *eventProcessor) updateMetaLastEventDate(ctx context.Context, event *types.Event) {
	var parents []types.Alarm
	err := s.adapter.GetOpenedAlarmsByIDs(ctx, event.Alarm.Value.Parents, &parents)
	if err == nil {
		updatedParents := make([]string, 0, len(parents))
		var alarm types.Alarm
		for _, metaAlarm := range parents {
			if alarm.ID == "" {
				alarm = metaAlarm
				alarm.PartialUpdateLastEventDate(event.Timestamp)
			}
			updatedParents = append(updatedParents, metaAlarm.ID)
		}
		if len(updatedParents) > 0 {
			err = s.adapter.MassPartialUpdateOpen(ctx, &alarm, updatedParents)
		}
	}
	if err != nil {
		s.logger.Error().Err(err).Msgf("error changestate meta-alarm from children %+v", event.Alarm)
	}
}

func newAlarm(event types.Event, alarmConfig config.AlarmConfig) types.Alarm {
	now := types.CpsTime{Time: time.Now()}

	return types.Alarm{
		EntityID: event.GetEID(),
		ID:       utils.NewID(),
		Time:     now,
		Value: types.AlarmValue{
			Component:         event.Component,
			Connector:         event.Connector,
			ConnectorName:     event.ConnectorName,
			CreationDate:      now,
			DisplayName:       types.GenDisplayName(alarmConfig.DisplayNameScheme),
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
}
