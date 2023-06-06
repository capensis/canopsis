package axe

import (
	"context"
	"fmt"
	"runtime/trace"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice/statecounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	liboperation "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

type eventProcessor struct {
	dbClient            mongo.DbClient
	adapter             alarm.Adapter
	entityAdapter       libentity.Adapter
	ruleAdapter         correlation.RulesAdapter
	alarmConfigProvider config.AlarmConfigProvider
	executor            liboperation.Executor
	alarmStatusService  alarmstatus.Service
	logger              zerolog.Logger
	metricsSender       metrics.Sender

	metaAlarmEventProcessor alarm.MetaAlarmEventProcessor

	statisticsSender statistics.EventStatisticsSender

	stateCountersService statecounters.StateCountersService

	pbhTypeResolver pbehavior.EntityTypeResolver

	autoInstructionMatcher AutoInstructionMatcher
}

func NewEventProcessor(
	dbClient mongo.DbClient,
	adapter alarm.Adapter,
	entityAdapter libentity.Adapter,
	ruleAdapter correlation.RulesAdapter,
	alarmConfigProvider config.AlarmConfigProvider,
	executor liboperation.Executor,
	alarmStatusService alarmstatus.Service,
	metricsSender metrics.Sender,
	metaAlarmEventProcessor alarm.MetaAlarmEventProcessor,
	statisticsSender statistics.EventStatisticsSender,
	stateCountersService statecounters.StateCountersService,
	pbhTypeResolver pbehavior.EntityTypeResolver,
	autoInstructionMatcher AutoInstructionMatcher,
	logger zerolog.Logger,
) alarm.EventProcessor {
	return &eventProcessor{
		dbClient:            dbClient,
		adapter:             adapter,
		entityAdapter:       entityAdapter,
		ruleAdapter:         ruleAdapter,
		alarmConfigProvider: alarmConfigProvider,
		executor:            executor,
		alarmStatusService:  alarmStatusService,
		metricsSender:       metricsSender,
		statisticsSender:    statisticsSender,
		pbhTypeResolver:     pbhTypeResolver,
		logger:              logger,

		metaAlarmEventProcessor: metaAlarmEventProcessor,

		stateCountersService:   stateCountersService,
		autoInstructionMatcher: autoInstructionMatcher,
	}
}

func (s *eventProcessor) Process(ctx context.Context, event *types.Event) (types.AlarmChange, error) {
	defer trace.StartRegion(ctx, "alarm.ProcessAlarmEvent").End()

	alarmChange := types.NewAlarmChange()

	if event.Entity == nil {
		return alarmChange, nil
	}

	var entityOldIdleSince *types.CpsTime
	var entityOldLastIdleRuleApply string
	var updatedServiceStates map[string]statecounters.UpdatedServicesInfo
	firstTimeTran := true

	err := s.dbClient.WithTransaction(ctx, func(tCtx context.Context) error {
		if !firstTimeTran {
			entity, exist := s.entityAdapter.Get(tCtx, event.GetEID())
			if !exist {
				return fmt.Errorf("entity with id = %s is not found after transaction rollback", event.GetEID())
			}

			event.Entity = &entity
			event.Alarm = nil
		}

		firstTimeTran = false
		alarm, err := s.adapter.GetOpenedAlarm(tCtx, event.GetEID())
		alarmNotFound := false
		if _, ok := err.(errt.NotFound); ok {
			alarmNotFound = true
			err = nil
		} else if err != nil {
			return fmt.Errorf("cannot fetch alarm: %w", err)
		}

		if !alarmNotFound {
			event.Alarm = &alarm
		}

		s.fillAlarmChange(event.Alarm, *event.Entity, &alarmChange)

		switch event.EventType {
		case types.EventTypeUpdateCounters:
			alarmChange.Type = types.AlarmChangeTypeNone
		case types.EventTypeEntityToggled:
			if !event.Entity.Enabled {
				alarmChange, err = s.resolveAlarmForDisabledEntity(tCtx, event)
			} else {
				alarmChange.Type = types.AlarmChangeTypeEnabled
			}
		case types.EventTypeRecomputeEntityService:
			if !event.Entity.Enabled {
				alarmChange, err = s.resolveAlarmForDisabledEntity(tCtx, event)
				break
			}

			updatedServiceStates, err = s.stateCountersService.RecomputeEntityServiceCounters(tCtx, *event)
			return err
		case types.EventTypeCheck:
			alarmChange.Type, err = s.storeAlarm(tCtx, event)
		case types.EventTypeNoEvents:
			alarmChange.Type, err = s.processNoEvents(tCtx, event)
		case types.EventTypeMetaAlarm:
			event.Alarm, err = s.metaAlarmEventProcessor.CreateMetaAlarm(tCtx, *event)
			alarmChange.Type = types.AlarmChangeTypeCreate
		case types.EventTypeTrigger:
			if event.AlarmChange != nil && event.Alarm != nil {
				alarmChange = *event.AlarmChange
			}
		default:
			if event.Entity == nil {
				return nil
			}

			entityOldIdleSince, entityOldLastIdleRuleApply = event.Entity.IdleSince, event.Entity.LastIdleRuleApply

			if event.Alarm == nil {
				err = s.processPbhEventsForEntity(tCtx, event, &alarmChange)
			} else {
				operation := s.createOperationFromEvent(event)
				alarmChange.Type, err = s.executor.Exec(tCtx, operation, event.Alarm, event.Entity, event.Timestamp, event.UserID, event.Role, event.Initiator)
				if err != nil {
					err = fmt.Errorf("cannot update alarm: %w", err)
				}
			}
		}

		if err != nil {
			return err
		}

		updatedServiceStates, err = s.stateCountersService.UpdateServiceCounters(tCtx, *event.Entity, event.Alarm, alarmChange)
		return err
	})

	if err != nil {
		return alarmChange, err
	}

	// services alarms
	go func() {
		for servID, servInfo := range updatedServiceStates {
			err := s.stateCountersService.UpdateServiceState(context.Background(), servID, servInfo)
			if err != nil {
				s.logger.Err(err).Msg("failed to update service state")
			}
		}
	}()

	// metrics
	go func() {
		alarm := types.Alarm{}
		if event.Alarm != nil {
			alarm = *event.Alarm
		}

		s.metricsSender.SendEventMetrics(
			context.Background(),
			alarm,
			*event.Entity,
			alarmChange,
			event.Timestamp.Time,
			event.Initiator,
			event.UserID,
		)
	}()

	if event.EventType == types.EventTypeCheck {
		go s.sendEventStatistics(ctx, *event)
	}

	if event.Alarm == nil {
		return alarmChange, nil
	}

	mustUpdateIdleFields := entityOldIdleSince != event.Entity.IdleSince ||
		entityOldLastIdleRuleApply != event.Entity.LastIdleRuleApply

	if alarmChange.Type == types.AlarmChangeTypeResolve {
		go func() {
			err := s.adapter.CopyAlarmToResolvedCollection(context.Background(), *event.Alarm)
			if err != nil {
				s.logger.Err(err).Msg("cannot update resolved alarm")
			}
		}()
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

	go func() {
		err := s.metaAlarmEventProcessor.Process(context.Background(), *event)
		if err != nil {
			s.logger.Err(err).Msg("cannot process meta alarm")
		}
	}()

	return alarmChange, nil
}

func (s *eventProcessor) fillAlarmChange(alarm *types.Alarm, entity types.Entity, alarmChange *types.AlarmChange) {
	if alarm == nil {
		alarmChange.PreviousPbehaviorTypeID = entity.PbehaviorInfo.TypeID
		alarmChange.PreviousPbehaviorCannonicalType = entity.PbehaviorInfo.CanonicalType
		alarmChange.PreviousPbehaviorTime = entity.PbehaviorInfo.Timestamp
	} else {
		alarmChange.PreviousState = alarm.Value.State.Value
		alarmChange.PreviousStateChange = alarm.Value.State.Timestamp
		alarmChange.PreviousStatus = alarm.Value.Status.Value
		alarmChange.PreviousPbehaviorTypeID = alarm.Value.PbehaviorInfo.TypeID
		alarmChange.PreviousPbehaviorCannonicalType = alarm.Value.PbehaviorInfo.CanonicalType
		alarmChange.PreviousPbehaviorTime = alarm.Value.PbehaviorInfo.Timestamp
	}
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
	entity := *event.Entity

	pbehaviorInfo, err := s.resolvePbehaviorInfo(ctx, entity)
	if err != nil {
		return changeType, err
	}

	alarmConfig := s.alarmConfigProvider.Get()
	alarm := newAlarm(*event, alarmConfig)
	err = UpdateAlarmState(&alarm, entity, event.Timestamp, event.State, event.Output, s.alarmStatusService)
	if err != nil {
		return changeType, err
	}

	alarm.PartialUpdateEventsCount()

	if pbehaviorInfo.IsDefaultActive() {
		changeType = types.AlarmChangeTypeCreate
		alarm.NotAckedSince = &alarm.Time
	} else {
		if pbehaviorInfo.IsActive() {
			alarm.NotAckedSince = &alarm.Time
		}

		output := fmt.Sprintf(
			"Pbehavior %s. Type: %s. Reason: %s.",
			pbehaviorInfo.Name,
			pbehaviorInfo.TypeName,
			pbehaviorInfo.ReasonName,
		)

		err := alarm.PartialUpdatePbhEnter(event.Timestamp, pbehaviorInfo,
			canopsis.DefaultEventAuthor, output, "", event.Role, event.Initiator)
		if err != nil {
			return changeType, fmt.Errorf("cannot add alarm steps: %w", err)
		}

		changeType = types.AlarmChangeTypeCreateAndPbhEnter
	}

	if s.alarmConfigProvider.Get().ActivateAlarmAfterAutoRemediation {
		matched, err := s.autoInstructionMatcher.Match(types.AlarmWithEntity{Alarm: alarm, Entity: entity})
		if err != nil {
			return changeType, err
		}

		alarm.InactiveAutoInstructionInProgress = matched
	}

	err = s.adapter.Insert(ctx, alarm)
	if err != nil {
		return changeType, fmt.Errorf("cannot create alarm: %w", err)
	}

	if changeType == types.AlarmChangeTypeCreateAndPbhEnter {
		event.Entity.PbehaviorInfo = alarm.Value.PbehaviorInfo
		err := s.entityAdapter.UpdatePbehaviorInfo(ctx, event.Entity.ID, event.Entity.PbehaviorInfo)
		if err != nil {
			return changeType, fmt.Errorf("cannot update entity: %w", err)
		}
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
	alarm.PartialUpdateTags(event.Tags)

	if alarmConfig.EnableLastEventDate {
		alarm.PartialUpdateLastEventDate(event.Timestamp)
	}

	err = s.adapter.PartialUpdateOpen(ctx, alarm)
	if err != nil {
		return changeType, fmt.Errorf("cannot update alarm: %w", err)
	}

	// Update cropped steps if needed
	if alarm.CropSteps() {
		alarm.AddUpdate("$set", bson.M{"v.steps": alarm.Value.Steps})
		err = s.adapter.PartialUpdateOpen(ctx, alarm)
		if err != nil {
			return changeType, fmt.Errorf("cannot update alarm: %w", err)
		}
		event.Alarm = alarm
	}
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

	return changeType, nil
}

func (s *eventProcessor) processNoEvents(ctx context.Context, event *types.Event) (types.AlarmChangeType, error) {
	changeType := types.AlarmChangeTypeNone
	if event.Alarm == nil && event.State == types.AlarmStateOK {
		return changeType, nil
	}

	pbehaviorInfo, err := s.resolvePbehaviorInfo(ctx, *event.Entity)
	if err != nil {
		return changeType, err
	}

	if event.Alarm == nil {
		alarmConfig := s.alarmConfigProvider.Get()
		alarm := newAlarm(*event, alarmConfig)
		err := s.updateAlarmOnNoEventsEvent(&alarm, *event.Entity, *event)
		if err != nil {
			return changeType, err
		}

		changeType = types.AlarmChangeTypeCreate

		if !pbehaviorInfo.IsDefaultActive() {
			output := fmt.Sprintf(
				"Pbehavior %s. Type: %s. Reason: %s.",
				pbehaviorInfo.Name,
				pbehaviorInfo.TypeName,
				pbehaviorInfo.ReasonName,
			)

			err := alarm.PartialUpdatePbhEnter(event.Timestamp, pbehaviorInfo,
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

	err = s.entityAdapter.UpdateIdleFields(ctx, event.Entity.ID, event.Entity.IdleSince,
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
	}

	return changeType, nil
}

func (s *eventProcessor) createOperationFromEvent(event *types.Event) types.Operation {
	t := event.EventType
	parameters := types.OperationParameters{
		TicketInfo:  event.TicketInfo,
		Output:      event.Output,
		Author:      event.Author,
		Execution:   event.Execution,
		Instruction: event.Instruction,
	}

	switch event.EventType {
	case types.EventTypeSnooze:
		parameters.Duration = &types.DurationWithUnit{
			Value: int64(event.Duration),
			Unit:  "s",
		}
	case types.EventTypeChangestate:
		parameters.State = &event.State
	case types.EventTypePbhEnter, types.EventTypePbhLeave, types.EventTypePbhLeaveAndEnter:
		parameters.PbehaviorInfo = &event.PbehaviorInfo
	case types.EventTypeRecomputeEntityService:
		if event.Entity.SoftDeleted == nil {
			return types.Operation{}
		}

		t = types.EventTypeResolveDeleted
	}

	return types.Operation{
		Type:       t,
		Parameters: parameters,
	}
}

func (s *eventProcessor) resolveAlarmForDisabledEntity(ctx context.Context, event *types.Event) (types.AlarmChange, error) {
	alarmChange := types.NewAlarmChange()
	event.Entity.IdleSince = nil
	event.Entity.LastIdleRuleApply = ""
	if event.Alarm == nil {
		return alarmChange, nil
	}

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

func (s *eventProcessor) processPbhEventsForEntity(ctx context.Context, event *types.Event, alarmChange *types.AlarmChange) error {
	switch event.EventType {
	case types.EventTypePbhEnter, types.EventTypePbhLeave, types.EventTypePbhLeaveAndEnter:
		curPbehaviorInfo := event.Entity.PbehaviorInfo
		if !curPbehaviorInfo.Same(event.PbehaviorInfo) {
			alarmChange.PreviousPbehaviorCannonicalType = event.Entity.PbehaviorInfo.CanonicalType
			alarmChange.PreviousPbehaviorTypeID = event.Entity.PbehaviorInfo.TypeID
			event.Entity.PbehaviorInfo = event.PbehaviorInfo
			err := s.entityAdapter.UpdatePbehaviorInfo(ctx, event.Entity.ID, event.Entity.PbehaviorInfo)
			if err != nil {
				return fmt.Errorf("cannot update entity: %w", err)
			}

			if alarmChange.PreviousPbehaviorTypeID == "" {
				alarmChange.Type = types.AlarmChangeTypePbhEnter
			} else if event.PbehaviorInfo.TypeID == "" {
				alarmChange.Type = types.AlarmChangeTypePbhLeave
			} else {
				alarmChange.Type = types.AlarmChangeTypePbhLeaveAndEnter
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

func (s *eventProcessor) resolvePbehaviorInfo(ctx context.Context, entity types.Entity) (types.PbehaviorInfo, error) {
	if !entity.PbehaviorInfo.IsDefaultActive() {
		return entity.PbehaviorInfo, nil
	}

	now := time.Now()
	result, err := s.pbhTypeResolver.Resolve(ctx, entity, now)
	if err != nil {
		return types.PbehaviorInfo{}, err
	}

	return pbehavior.NewPBehaviorInfo(types.CpsTime{Time: now}, result), nil
}

func newAlarm(event types.Event, alarmConfig config.AlarmConfig) types.Alarm {
	now := types.CpsTime{Time: time.Now()}

	return types.Alarm{
		EntityID: event.GetEID(),
		ID:       utils.NewID(),
		Time:     now,
		Tags:     types.TransformEventTags(event.Tags),
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
