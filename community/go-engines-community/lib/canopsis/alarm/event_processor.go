package alarm

import (
	"context"
	"fmt"
	"runtime/trace"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	liboperation "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

const MaxRedisLockRetries = 10

type eventProcessor struct {
	dbClient            mongo.DbClient
	adapter             Adapter
	entityAdapter       entity.Adapter
	ruleAdapter         correlation.RulesAdapter
	redisLockClient     redis.LockClient
	alarmConfigProvider config.AlarmConfigProvider
	executor            liboperation.Executor
	alarmStatusService  alarmstatus.Service
	logger              zerolog.Logger
	metricsSender       metrics.Sender
	statisticsSender    statistics.EventStatisticsSender
}

func NewEventProcessor(
	dbClient mongo.DbClient,
	adapter Adapter,
	entityAdapter entity.Adapter,
	ruleAdapter correlation.RulesAdapter,
	alarmConfigProvider config.AlarmConfigProvider,
	executor liboperation.Executor,
	alarmStatusService alarmstatus.Service,
	redisLockClient redis.LockClient,
	metricsSender metrics.Sender,
	statisticsSender statistics.EventStatisticsSender,
	logger zerolog.Logger,
) EventProcessor {
	return &eventProcessor{
		dbClient:            dbClient,
		adapter:             adapter,
		entityAdapter:       entityAdapter,
		ruleAdapter:         ruleAdapter,
		alarmConfigProvider: alarmConfigProvider,
		executor:            executor,
		alarmStatusService:  alarmStatusService,
		redisLockClient:     redisLockClient,
		metricsSender:       metricsSender,
		statisticsSender:    statisticsSender,
		logger:              logger,
	}
}

func (s *eventProcessor) Process(ctx context.Context, event *types.Event) (types.AlarmChange, error) {
	defer trace.StartRegion(ctx, "alarm.ProcessAlarmEvent").End()

	alarmChange := types.NewAlarmChange()

	if event.Entity == nil {
		return alarmChange, nil
	}

	if !event.Entity.Enabled {
		var err error

		if event.EventType == types.EventTypeEntityToggled ||
			event.EventType == types.EventTypeRecomputeEntityService {
			alarmChange, err = s.resolveAlarmForDisabledEntity(ctx, event)

			if err == nil && alarmChange.Type == types.AlarmChangeTypeNone {
				alarmChange.Type = types.AlarmChangeTypeEntityToggled
			}
		}

		return alarmChange, err
	}

	alarm, err := s.adapter.GetOpenedAlarm(ctx, event.GetEID())
	alarmNotFound := false
	if _, ok := err.(errt.NotFound); ok {
		alarmNotFound = true
	} else if err != nil {
		return alarmChange, fmt.Errorf("cannot fetch alarm: %w", err)
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
		if err == nil {
			go s.sendEventStatistics(ctx, *event)
		}

		if changeType == types.AlarmChangeTypeStateIncrease || changeType == types.AlarmChangeTypeStateDecrease {
			s.updateMetaChildrenState(ctx, event)
		}
		if event.Alarm != nil && event.Alarm.IsMetaChildren() &&
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
	case types.EventTypeEntityToggled:
		alarmChange.Type = types.AlarmChangeTypeEntityToggled
		return alarmChange, nil
	case types.EventTypeTrigger:
		if event.AlarmChange == nil || event.Alarm == nil {
			return types.NewAlarmChange(), nil
		}

		return *event.AlarmChange, nil
	}

	if event.Alarm == nil {
		err = s.processPbhEventsForEntity(ctx, event, &alarmChange)

		return alarmChange, err
	}

	entityOldIdleSince, entityOldLastIdleRuleApply := event.Entity.IdleSince, event.Entity.LastIdleRuleApply

	operation := s.createOperationFromEvent(event)
	changeType, err := s.executor.Exec(ctx, operation, event.Alarm, event.Entity, event.Timestamp, event.UserID, event.Role, event.Initiator)
	if err != nil {
		return alarmChange, fmt.Errorf("cannot update alarm: %w", err)
	}

	mustUpdateIdleFields := entityOldIdleSince != event.Entity.IdleSince ||
		entityOldLastIdleRuleApply != event.Entity.LastIdleRuleApply

	if changeType == types.AlarmChangeTypeResolve {
		err := s.adapter.CopyAlarmToResolvedCollection(ctx, *event.Alarm)
		if err != nil {
			return alarmChange, fmt.Errorf("cannot update resolved alarm: %w", err)
		}
	}

	if event.IdleRuleApply != "" {
		event.Entity.LastIdleRuleApply = event.IdleRuleApply
		mustUpdateIdleFields = true
	}
	if mustUpdateIdleFields {
		err = s.entityAdapter.UpdateIdleFields(ctx, event.Entity.ID, event.Entity.IdleSince,
			event.Entity.LastIdleRuleApply)
		if err != nil {
			return alarmChange, fmt.Errorf("cannot update entity: %w", err)
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
			return fmt.Errorf("cannot fetch alarm: %w", err)
		}

		if !notFound && lastAlarm.Value.Resolved != nil {
			alarmChange.PreviousStateChange = *lastAlarm.Value.Resolved
			alarmChange.PreviousStatusChange = *lastAlarm.Value.Resolved
		} else {
			alarmChange.PreviousStateChange = event.Timestamp
			alarmChange.PreviousStatusChange = event.Timestamp
		}

		alarmChange.PreviousPbehaviorTypeID = event.Entity.PbehaviorInfo.TypeID
		alarmChange.PreviousPbehaviorCannonicalType = event.Entity.PbehaviorInfo.CanonicalType
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

	alarmConfig := s.alarmConfigProvider.Get()
	alarm := newAlarm(*event, alarmConfig)
	err := UpdateAlarmState(&alarm, *event.Entity, event.Timestamp, event.State, event.Output, s.alarmStatusService)
	if err != nil {
		return changeType, err
	}

	alarm.PartialUpdateEventsCount()

	if event.PbehaviorInfo.IsDefaultActive() {
		changeType = types.AlarmChangeTypeCreate
	} else {
		output := fmt.Sprintf(
			"Pbehavior %s. Type: %s. Reason: %s.",
			event.PbehaviorInfo.Name,
			event.PbehaviorInfo.TypeName,
			event.PbehaviorInfo.Reason,
		)

		err := alarm.PartialUpdatePbhEnter(event.Timestamp, event.PbehaviorInfo,
			canopsis.DefaultEventAuthor, output, "", event.Role, event.Initiator)
		if err != nil {
			return changeType, fmt.Errorf("cannot add alarm steps: %w", err)
		}

		changeType = types.AlarmChangeTypeCreateAndPbhEnter
	}

	err = s.adapter.Insert(ctx, alarm)
	if err != nil {
		return changeType, fmt.Errorf("cannot create alarm: %w", err)
	}

	if changeType == types.AlarmChangeTypeCreate {
		s.metricsSender.SendCreate(alarm, alarm.Value.CreationDate.Time)
	}

	if changeType == types.AlarmChangeTypeCreateAndPbhEnter {
		event.Entity.PbehaviorInfo = alarm.Value.PbehaviorInfo
		err := s.entityAdapter.UpdatePbehaviorInfo(ctx, event.Entity.ID, event.Entity.PbehaviorInfo)
		if err != nil {
			return changeType, fmt.Errorf("cannot update entity: %w", err)
		}

		s.metricsSender.SendCreateAndPbhEnter(alarm, alarm.Value.CreationDate.Time)
	}

	event.Alarm = &alarm

	return changeType, nil
}

// updateAlarm updates alarm value and crops steps.
func (s *eventProcessor) updateAlarm(ctx context.Context, event *types.Event) (types.AlarmChangeType, error) {
	changeType := types.AlarmChangeTypeNone
	alarm := event.Alarm
	alarmConfig := s.alarmConfigProvider.Get()
	previousState := alarm.CurrentState()
	newState := event.State
	err := UpdateAlarmState(alarm, *event.Entity, event.Timestamp, event.State, event.Output, s.alarmStatusService)
	if err != nil {
		return changeType, err
	}

	alarm.UpdateOutput(event.Output)
	alarm.UpdateLongOutput(event.LongOutput)
	alarm.PartialUpdateEventsCount()

	if alarmConfig.EnableLastEventDate {
		alarm.PartialUpdateLastEventDate(event.Timestamp)
	}

	err = s.adapter.PartialUpdateOpen(ctx, alarm)
	if err != nil {
		return changeType, fmt.Errorf("cannot update alarm: %w", err)
	}

	// Update cropped steps if needed
	err = s.dbClient.WithTransaction(ctx, func(tranCtx context.Context) error {
		alarm, err := s.adapter.GetOpenedAlarmByAlarmId(tranCtx, event.Alarm.ID)
		if err != nil {
			return fmt.Errorf("cannot fetch alarm: %w", err)
		}
		if alarm.CropSteps() {
			alarm.AddUpdate("$set", bson.M{"v.steps": alarm.Value.Steps})
			err = s.adapter.PartialUpdateOpen(tranCtx, &alarm)
			if err != nil {
				return fmt.Errorf("cannot update alarm: %w", err)
			}
			event.Alarm = &alarm
		}

		return nil
	})
	if err != nil {
		return changeType, err
	}

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
			return changeType, fmt.Errorf("cannot update entity: %w", err)
		}
	}

	s.metricsSender.SendUpdateState(*event.Alarm, *event.Entity, previousState)

	return changeType, nil
}

func (s *eventProcessor) processNoEvents(ctx context.Context, event *types.Event) (types.AlarmChangeType, error) {
	changeType := types.AlarmChangeTypeNone
	if event.Alarm == nil && event.State == types.AlarmStateOK {
		return changeType, nil
	}

	alarmConfig := s.alarmConfigProvider.Get()

	if event.Alarm == nil {
		alarm := newAlarm(*event, alarmConfig)
		err := s.updateAlarmOnNoEventsEvent(&alarm, *event.Entity, *event)
		if err != nil {
			return changeType, err
		}

		changeType = types.AlarmChangeTypeCreate

		if !event.PbehaviorInfo.IsDefaultActive() {
			output := fmt.Sprintf(
				"Pbehavior %s. Type: %s. Reason: %s.",
				event.PbehaviorInfo.Name,
				event.PbehaviorInfo.TypeName,
				event.PbehaviorInfo.Reason,
			)

			err := alarm.PartialUpdatePbhEnter(event.Timestamp, event.PbehaviorInfo,
				canopsis.DefaultEventAuthor, output, "", event.Role, event.Initiator)
			if err != nil {
				return changeType, fmt.Errorf("cannot add alarm steps: %w", err)
			}

			changeType = types.AlarmChangeTypeCreateAndPbhEnter
		}

		err = s.adapter.Insert(ctx, alarm)
		if err != nil {
			return changeType, fmt.Errorf("cannot create alarm: %w", err)
		}

		event.Alarm = &alarm
		s.metricsSender.SendCreate(alarm, alarm.Value.CreationDate.Time)
	} else {
		alarm := event.Alarm
		previousState := alarm.Value.State.Value
		previousStatus := alarm.Value.Status.Value
		err := s.updateAlarmOnNoEventsEvent(alarm, *event.Entity, *event)
		if err != nil {
			return changeType, err
		}

		err = s.adapter.PartialUpdateOpen(ctx, alarm)
		if err != nil {
			return changeType, fmt.Errorf("cannot update alarm: %w", err)
		}

		newState := alarm.Value.State.Value
		newStatus := alarm.Value.Status.Value
		if newState > previousState {
			changeType = types.AlarmChangeTypeStateIncrease
			s.metricsSender.SendUpdateState(*alarm, *event.Entity, previousState)
		} else if newState < previousState {
			changeType = types.AlarmChangeTypeStateDecrease
			s.metricsSender.SendUpdateState(*alarm, *event.Entity, previousState)
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
		return changeType, fmt.Errorf("cannot update entity: %w", err)
	}

	if changeType == types.AlarmChangeTypeCreateAndPbhEnter {
		event.Entity.PbehaviorInfo = event.Alarm.Value.PbehaviorInfo
		err := s.entityAdapter.UpdatePbehaviorInfo(ctx, event.Entity.ID, event.Entity.PbehaviorInfo)
		if err != nil {
			return changeType, fmt.Errorf("cannot update entity: %w", err)
		}

		s.metricsSender.SendCreateAndPbhEnter(*event.Alarm, event.Alarm.Value.CreationDate.Time)
	}

	return changeType, nil
}

func (s *eventProcessor) createOperationFromEvent(event *types.Event) types.Operation {
	parameters := types.OperationParameters{
		Ticket:    event.Ticket,
		Output:    event.Output,
		Author:    event.Author,
		Execution: event.Execution,
	}
	switch event.EventType {
	case types.EventTypeSnooze:
		parameters.Duration = &types.DurationWithUnit{
			Value: int64(event.Duration),
			Unit:  "s",
		}
	case types.EventTypeChangestate, types.EventTypeKeepstate:
		parameters.State = &event.State
	case types.EventTypePbhEnter, types.EventTypePbhLeave, types.EventTypePbhLeaveAndEnter:
		parameters.PbehaviorInfo = &event.PbehaviorInfo
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
		return fmt.Errorf("cannot fetch alarms: %w", err)
	}

	for _, alarm := range alarms {
		_, err := s.executor.Exec(ctx, operation, &alarm.Alarm, &alarm.Entity, event.Timestamp, event.UserID, event.Role, event.Initiator)
		if err != nil {
			return fmt.Errorf("cannot update alarm: %w", err)
		}
	}

	return nil
}

func (s *eventProcessor) processMetaAlarmCreateEvent(ctx context.Context, event *types.Event) (types.AlarmChangeType, error) {
	var childAlarms []types.Alarm
	if event.MetaAlarmChildren != nil {
		err := s.adapter.GetOpenedAlarmsByIDs(ctx, *event.MetaAlarmChildren, &childAlarms)
		if err != nil {
			return types.AlarmChangeTypeNone, fmt.Errorf("cannot fetch children alarms: %w", err)
		}
		worstState := types.CpsNumber(types.AlarmStateMinor)

		for i := 0; i < len(childAlarms); i++ {
			c := childAlarms[i]
			if c.Value.State != nil {
				childState := c.Value.State.Value
				if childState > worstState {
					worstState = childState
				}
			}
		}
		event.State = worstState
	}
	metaAlarm, err := types.NewAlarm(*event, s.alarmConfigProvider.Get())
	if err != nil {
		return types.AlarmChangeTypeNone, fmt.Errorf("cannot create alarm: %w", err)
	}

	metaAlarm.Value.Tags = []string{}
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
		return types.AlarmChangeTypeNone, fmt.Errorf("cannot create alarm: %w", err)
	}

	ruleIdentifier := metaAlarm.Value.Meta
	rule, err := s.ruleAdapter.GetRule(metaAlarm.Value.Meta)
	if err != nil {
		// the rule can be deleted
		if err.Error() != "not found" {
			return types.AlarmChangeTypeNone, fmt.Errorf("cannot fetch rule id=%q: %w", metaAlarm.Value.Meta, err)
		}
	} else {
		ruleIdentifier = rule.Name
	}
	newStep := types.NewMetaAlarmAttachStep(metaAlarm, ruleIdentifier)

	for i := 0; i < len(childAlarms); i++ {
		c := childAlarms[i]
		err := c.Value.Steps.Add(newStep)
		if err != nil {
			s.logger.Err(err).
				Str("metaalarm", metaAlarm.EntityID).
				Str("child entity", c.EntityID).
				Str("child alarm", c.ID).
				Msg("cannot add metaalarmattach step to child")
		}
		childAlarms[i] = c
	}

	err = s.adapter.MassUpdate(ctx, childAlarms, false)
	if err != nil {
		return types.AlarmChangeTypeNone, fmt.Errorf("cannot update children alarms: %w", err)
	}

	go func() {
		timestamp := time.Now()
		if !event.Timestamp.IsZero() {
			timestamp = event.Timestamp.Time
		}

		for _, child := range childAlarms {
			s.metricsSender.SendCorrelation(timestamp, child)
		}
	}()

	event.Alarm = &metaAlarm
	return types.AlarmChangeTypeCreate, nil
}

func (s *eventProcessor) processMetaAlarmChildren(ctx context.Context, event *types.Event, changeType types.AlarmChangeType, operation types.Operation) error {
	if !event.Alarm.IsMetaAlarm() ||
		changeType != types.AlarmChangeTypeAck &&
			changeType != types.AlarmChangeTypeDoubleAck &&
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

	var alarms []types.AlarmWithEntity
	err := s.adapter.GetOpenedAlarmsWithEntityByIDs(ctx, event.Alarm.Value.Children, &alarms)
	if err != nil {
		return fmt.Errorf("cannot fetch children alarms: %w", err)
	}
	for _, alarm := range alarms {
		_, err := s.executor.Exec(ctx, operation, &alarm.Alarm, &alarm.Entity, event.Timestamp, event.UserID, event.Role, event.Initiator)
		if err != nil {
			return fmt.Errorf("cannot update children alarms: %w", err)
		}
	}

	return nil
}

func (s *eventProcessor) handleMetaAlarmChildResolve(ctx context.Context, event *types.Event, changeType types.AlarmChangeType) (bool, error) {
	if event.Alarm.IsMetaAlarm() || changeType != types.AlarmChangeTypeResolve || len(event.Alarm.Value.Parents) == 0 {
		return false, nil
	}

	var alarms []types.Alarm
	err := s.adapter.GetOpenedAlarmsByIDs(ctx, event.Alarm.Value.Parents, &alarms)
	if err != nil {
		return false, fmt.Errorf("cannot fetch parent alarms: %w", err)
	}

	if len(alarms) == 0 {
		return false, nil
	}

	var wg sync.WaitGroup
	for _, alarm := range alarms {
		wg.Add(1)
		go func(alarm types.Alarm) {
			if err := s.handleAutoResolveMetaAlarm(ctx, alarm); err != nil {
				s.logger.Err(err).
					Str("alarm", alarm.ID).
					Msg("cannot auto resolve parent alarm")
			}
			wg.Done()
		}(alarm)
	}
	wg.Wait()

	s.updateMetaChildrenState(ctx, event)

	return true, nil
}

func (s *eventProcessor) handleAutoResolveMetaAlarm(ctx context.Context, alarm types.Alarm) error {
	rule, err := s.ruleAdapter.GetRule(alarm.Value.Meta)
	if err != nil {
		return err
	}
	if !rule.AutoResolve || alarm.IsResolved() {
		return nil
	}

	metaAlarmLock, err := s.redisLockClient.Obtain(ctx, alarm.ID, 100*time.Millisecond, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(11*time.Millisecond), MaxRedisLockRetries),
	})
	if err != nil {
		return fmt.Errorf("cannot obtain meta alarm lock: %w", err)
	}

	defer func() {
		if metaAlarmLock != nil {
			err := metaAlarmLock.Release(ctx)
			if err != nil && err != redislock.ErrLockNotHeld {
				s.logger.Err(err).Msg("cannot release meta alarm lock")
			}
		}
	}()

	c, err := s.adapter.CountResolvedAlarm(ctx, alarm.Value.Children)
	if err != nil {
		return fmt.Errorf("cannot fetch alarms: %w", err)
	}

	if c == len(alarm.Value.Children) {
		err := alarm.PartialUpdateResolve(types.CpsTime{
			Time: time.Now(),
		})
		if err != nil {
			return fmt.Errorf("cannot update alarm: %w", err)
		}

		err = s.adapter.CopyAlarmToResolvedCollection(ctx, alarm)
		if err != nil {
			return fmt.Errorf("cannot update alarm: %w", err)
		}

		err = s.adapter.PartialUpdateOpen(ctx, &alarm)
		if err != nil {
			return fmt.Errorf("cannot update alarm: %w", err)
		}
	}

	return nil
}

// updateMetaChildrenState updates alarm's parents according with children state if need
func (s *eventProcessor) updateMetaChildrenState(ctx context.Context, event *types.Event) {
	if !event.Alarm.IsMetaChildren() {
		return
	}

	alarmState := event.Alarm.Value.State.Value
	lastUpdateDate := event.Alarm.Value.LastUpdateDate
	if event.Alarm.IsResolved() {
		if alarmState == types.AlarmStateOK {
			return
		}
		alarmState = types.AlarmStateOK
		lastUpdateDate = *event.Alarm.Value.Resolved
	}

	var parents []types.AlarmWithEntity
	err := s.adapter.GetOpenedAlarmsWithEntityByIDs(ctx, event.Alarm.Value.Parents, &parents)
	if err != nil {
		s.logger.Error().Err(err).Msgf("cannot fetch parent alarms")
		return
	}

	updatedParents := make([]types.Alarm, 0, len(parents))
	for _, metaAlarm := range parents {
		maCurrentState := metaAlarm.Alarm.Value.State.Value
		if alarmState > maCurrentState {
			err := UpdateAlarmState(&metaAlarm.Alarm, metaAlarm.Entity, lastUpdateDate, alarmState, metaAlarm.Alarm.Value.Output, s.alarmStatusService)
			if err != nil {
				s.logger.Error().Err(err).Str("alarm", metaAlarm.Alarm.ID).Msgf("cannot update alarm")
				return
			}

			updatedParents = append(updatedParents, metaAlarm.Alarm)
		} else if alarmState < maCurrentState {
			err := s.updateMetaAlarmToWorstState(ctx, &metaAlarm.Alarm, metaAlarm.Entity, []*types.Alarm{event.Alarm})
			if err != nil {
				s.logger.Error().Err(err).Str("alarm", metaAlarm.Alarm.ID).Msgf("cannot update alarm")
				return
			}

			updatedParents = append(updatedParents, metaAlarm.Alarm)
		}
	}

	if len(updatedParents) > 0 {
		err = s.adapter.PartialMassUpdateOpen(ctx, updatedParents)
		if err != nil {
			s.logger.Error().Err(err).Msgf("cannot update parent alarms")
		}
	}
}

func (s *eventProcessor) resolveAlarmForDisabledEntity(ctx context.Context, event *types.Event) (types.AlarmChange, error) {
	alarmChange := types.NewAlarmChange()
	alarm, err := s.adapter.GetOpenedAlarm(ctx, event.GetEID())
	event.Entity.IdleSince = nil
	event.Entity.LastIdleRuleApply = ""

	if _, ok := err.(errt.NotFound); ok {
		return alarmChange, nil
	} else if err != nil {
		return alarmChange, fmt.Errorf("cannot fetch alarm: %w", err)
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
	changeType, err := s.executor.Exec(ctx, operation, event.Alarm, event.Entity, event.Timestamp, event.UserID, event.Role, event.Initiator)
	if err != nil {
		return alarmChange, fmt.Errorf("cannot update alarm: %w", err)
	}

	alarmChange.Type = changeType

	return alarmChange, nil
}

// updateMetaLastEventDate updates alarm's parents LastEventDate
func (s *eventProcessor) updateMetaLastEventDate(ctx context.Context, event *types.Event) {
	var parents []types.Alarm
	err := s.adapter.GetOpenedAlarmsByIDs(ctx, event.Alarm.Value.Parents, &parents)
	if err != nil {
		s.logger.Err(err).Msg("cannot fetch parent alarms")
		return
	}
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
		err := s.adapter.MassPartialUpdateOpen(ctx, &alarm, updatedParents)
		if err != nil {
			s.logger.Err(err).Msg("cannot update parent alarms")
			return
		}
	}
}

func (s *eventProcessor) updateAlarmOnNoEventsEvent(alarm *types.Alarm, entity types.Entity, event types.Event) error {
	var currentState, currentStatus types.CpsNumber
	if alarm.Value.State != nil {
		currentState = alarm.Value.State.Value
		currentStatus = alarm.Value.Status.Value
	}
	stateUpdated := false
	state := event.State
	if currentState != state {
		// Create new Step to keep track of the alarm history
		newStep := types.NewAlarmStep(types.AlarmStepStateIncrease, event.Timestamp, event.Author, event.Output, event.UserID, event.Role, event.Initiator)
		newStep.Value = state

		if state < currentState {
			newStep.Type = types.AlarmStepStateDecrease
		}

		alarm.Value.State = &newStep
		err := alarm.Value.Steps.Add(newStep)
		if err != nil {
			return err
		}

		stateUpdated = true
	}

	newStatus := types.CpsNumber(types.AlarmStatusNoEvents)
	if state == types.AlarmStateOK {
		newStatus = s.alarmStatusService.ComputeStatus(*alarm, entity)
	}

	if newStatus == currentStatus {
		if stateUpdated {
			alarm.AddUpdate("$set", bson.M{"v.state": alarm.Value.State})
			alarm.AddUpdate("$push", bson.M{"v.steps": alarm.Value.State})
		}
		return nil
	}

	// Create new Step to keep track of the alarm history
	newStepStatus := types.NewAlarmStep(types.AlarmStepStatusIncrease, event.Timestamp, event.Author, event.Output, event.UserID, event.Role, event.Initiator)
	newStepStatus.Value = newStatus

	if newStatus < currentStatus {
		newStepStatus.Type = types.AlarmStepStatusDecrease
	}

	alarm.Value.Status = &newStepStatus
	err := alarm.Value.Steps.Add(newStepStatus)
	if err != nil {
		return err
	}

	alarm.Value.StateChangesSinceStatusUpdate = 0
	alarm.Value.LastUpdateDate = event.Timestamp

	set := bson.M{
		"v.status":                            alarm.Value.Status,
		"v.state_changes_since_status_update": alarm.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  alarm.Value.LastUpdateDate,
	}
	newSteps := bson.A{}
	if stateUpdated {
		set["v.state"] = alarm.Value.State
		newSteps = append(newSteps, alarm.Value.State)
	}
	newSteps = append(newSteps, alarm.Value.Status)
	alarm.AddUpdate("$set", set)
	alarm.AddUpdate("$push", bson.M{"v.steps": bson.M{"$each": newSteps}})

	return nil
}

// updateMetaAlarmToWorstState updates meta-alarm's state from its worst children, return true when meta-alarm has updated
func (s *eventProcessor) updateMetaAlarmToWorstState(ctx context.Context, metaAlarm *types.Alarm, metaAlarmEntity types.Entity, updateChildren []*types.Alarm) error {
	var (
		alarms     []types.Alarm
		stepTs, ts types.CpsTime
		worstState = types.CpsNumber(types.AlarmStateOK)
	)

	if metaAlarm.Value.Children != nil && len(metaAlarm.Value.Children) > 0 {
		if len(metaAlarm.Value.Children) == len(updateChildren) {
			// all children to update
			for _, child := range updateChildren {
				childState := child.CurrentState()
				if child.IsResolved() {
					childState = types.AlarmStateOK
				}
				if childState >= worstState {
					worstState = childState
					stepTs = child.Value.LastUpdateDate
				}
			}
		} else {
			err := s.adapter.GetOpenedAlarmsByIDs(ctx, metaAlarm.Value.Children, &alarms)
			if err != nil {
				return fmt.Errorf("cannot fetch children alarms: %w", err)
			}
			for _, child := range alarms {
				childState := child.CurrentState()
				for _, uc := range updateChildren {
					ts = uc.Value.LastUpdateDate
					if uc.ID == child.ID {
						childState = uc.CurrentState()
						if uc.IsResolved() {
							childState = types.AlarmStateOK
						}
						break
					}
				}
				if childState > worstState {
					worstState = childState
					stepTs = ts
				}
			}
			if stepTs.IsZero() {
				stepTs = ts
			}
		}
	}

	return UpdateAlarmState(metaAlarm, metaAlarmEntity, stepTs, worstState, metaAlarm.Value.Output, s.alarmStatusService)
}

func (s *eventProcessor) processPbhEventsForEntity(ctx context.Context, event *types.Event, alarmChange *types.AlarmChange) error {
	switch event.EventType {
	case types.EventTypePbhEnter, types.EventTypePbhLeave, types.EventTypePbhLeaveAndEnter:
		curPbehaviorInfo := event.Entity.PbehaviorInfo
		if curPbehaviorInfo != event.PbehaviorInfo {
			alarmChange.PreviousPbehaviorCannonicalType = event.Entity.PbehaviorInfo.CanonicalType
			alarmChange.PreviousPbehaviorTypeID = event.Entity.PbehaviorInfo.TypeID
			event.Entity.PbehaviorInfo = event.PbehaviorInfo
			err := s.entityAdapter.UpdatePbehaviorInfo(ctx, event.Entity.ID, event.Entity.PbehaviorInfo)
			if err != nil {
				return fmt.Errorf("cannot update entity: %w", err)
			}

			if alarmChange.PreviousPbehaviorTypeID == "" {
				alarmChange.Type = types.AlarmChangeTypePbhEnter
				s.metricsSender.SendPbhEnter(nil, *event.Entity)
			} else if event.PbehaviorInfo.TypeID == "" {
				alarmChange.Type = types.AlarmChangeTypePbhLeave
				s.metricsSender.SendPbhLeave(*event.Entity, event.Timestamp.Time, curPbehaviorInfo.CanonicalType, curPbehaviorInfo.Timestamp.Time)
			} else {
				alarmChange.Type = types.AlarmChangeTypePbhLeaveAndEnter
				s.metricsSender.SendPbhLeaveAndEnter(nil, *event.Entity, curPbehaviorInfo.CanonicalType, curPbehaviorInfo.Timestamp.Time)
			}
		}
	}

	return nil
}

func (s *eventProcessor) sendEventStatistics(ctx context.Context, event types.Event) {
	if event.Entity == nil {
		return
	}

	if event.Entity.PbehaviorInfo.Is(pbehavior.TypeInactive) {
		return
	}

	stats := statistics.EventStatistics{LastEvent: &event.Timestamp}
	if event.State == types.AlarmStateOK {
		stats.OK = 1
	} else {
		stats.KO = 1
	}

	s.statisticsSender.Send(ctx, event.GetEID(), stats)
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
			Infos:             map[string]map[string]interface{}{},
			RuleVersion:       map[string]string{},
		},
	}
}

func UpdateAlarmState(alarm *types.Alarm, entity types.Entity, timestamp types.CpsTime, state types.CpsNumber, output string,
	service alarmstatus.Service) error {
	var currentState, currentStatus types.CpsNumber
	if alarm.Value.State != nil {
		currentState = alarm.Value.State.Value
		currentStatus = alarm.Value.Status.Value
	}

	if state != currentState {
		// Event is an OK, so the alarm should be resolved anyway
		if alarm.IsStateLocked() && state != types.AlarmStateOK {
			return nil
		}

		// Create new Step to keep track of the alarm history
		newStep := types.NewAlarmStep(types.AlarmStepStateIncrease, timestamp, alarm.Value.Connector+"."+alarm.Value.ConnectorName, output, "", "", "")
		newStep.Value = state

		if state < currentState {
			newStep.Type = types.AlarmStepStateDecrease
		}

		alarm.Value.State = &newStep
		err := alarm.Value.Steps.Add(newStep)
		if err != nil {
			return err
		}

		alarm.Value.TotalStateChanges++
		alarm.Value.LastUpdateDate = timestamp
	}

	newStatus := service.ComputeStatus(*alarm, entity)

	if newStatus == currentStatus {
		if state != currentState {
			alarm.Value.StateChangesSinceStatusUpdate++

			alarm.AddUpdate("$set", bson.M{
				"v.state":                             alarm.Value.State,
				"v.state_changes_since_status_update": alarm.Value.StateChangesSinceStatusUpdate,
				"v.total_state_changes":               alarm.Value.TotalStateChanges,
				"v.last_update_date":                  alarm.Value.LastUpdateDate,
			})
			alarm.AddUpdate("$push", bson.M{"v.steps": alarm.Value.State})
		}

		return nil
	}

	// Create new Step to keep track of the alarm history
	newStepStatus := types.NewAlarmStep(types.AlarmStepStatusIncrease, timestamp, alarm.Value.Connector+"."+alarm.Value.ConnectorName, output, "", "", "")
	newStepStatus.Value = newStatus

	if newStatus < currentStatus {
		newStepStatus.Type = types.AlarmStepStatusDecrease
	}

	alarm.Value.Status = &newStepStatus
	err := alarm.Value.Steps.Add(newStepStatus)
	if err != nil {
		return err
	}

	alarm.Value.StateChangesSinceStatusUpdate = 0
	alarm.Value.LastUpdateDate = timestamp

	set := bson.M{
		"v.status":                            alarm.Value.Status,
		"v.state_changes_since_status_update": alarm.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  alarm.Value.LastUpdateDate,
	}
	newSteps := bson.A{}
	if state != currentState {
		set["v.total_state_changes"] = alarm.Value.TotalStateChanges
		set["v.state"] = alarm.Value.State
		newSteps = append(newSteps, alarm.Value.State)
	}
	newSteps = append(newSteps, alarm.Value.Status)

	alarm.AddUpdate("$set", set)
	alarm.AddUpdate("$push", bson.M{"v.steps": bson.M{"$each": newSteps}})

	return nil
}
