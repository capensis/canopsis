package alarm

import (
	"context"
	"fmt"
	"runtime/trace"
	"time"

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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

type eventProcessor struct {
	dbClient            mongo.DbClient
	adapter             Adapter
	entityAdapter       entity.Adapter
	ruleAdapter         correlation.RulesAdapter
	alarmConfigProvider config.AlarmConfigProvider
	executor            liboperation.Executor
	alarmStatusService  alarmstatus.Service
	logger              zerolog.Logger
	metricsSender       metrics.Sender

	metaAlarmEventProcessor MetaAlarmEventProcessor

	statisticsSender statistics.EventStatisticsSender
}

func NewEventProcessor(
	dbClient mongo.DbClient,
	adapter Adapter,
	entityAdapter entity.Adapter,
	ruleAdapter correlation.RulesAdapter,
	alarmConfigProvider config.AlarmConfigProvider,
	executor liboperation.Executor,
	alarmStatusService alarmstatus.Service,
	metricsSender metrics.Sender,
	metaAlarmEventProcessor MetaAlarmEventProcessor,
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
		metricsSender:       metricsSender,
		statisticsSender:    statisticsSender,
		logger:              logger,

		metaAlarmEventProcessor: metaAlarmEventProcessor,
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
		return alarmChange, fmt.Errorf("cannot fetch alarm: %w", err)
	}

	if !alarmNotFound {
		event.Alarm = &alarm
	}

	if err := s.fillAlarmChange(event.Alarm, &alarmChange); err != nil {
		return alarmChange, err
	}

	switch event.EventType {
	case types.EventTypeCheck:
		changeType, err := s.storeAlarm(ctx, event)
		if err == nil {
			go s.sendEventStatistics(ctx, *event)
		}

		go func() {
			err := s.metaAlarmEventProcessor.Process(context.Background(), *event)
			if err != nil {
				s.logger.Err(err).Msg("cannot process meta alarm")
			}
		}()

		alarmChange.Type = changeType
		return alarmChange, err
	case types.EventTypeNoEvents:
		changeType, err := s.processNoEvents(ctx, event)
		alarmChange.Type = changeType
		return alarmChange, err
	case types.EventTypeMetaAlarm:
		alarm, err := s.metaAlarmEventProcessor.CreateMetaAlarm(ctx, *event)
		if err != nil {
			return alarmChange, err
		}

		event.Alarm = alarm
		alarmChange.Type = types.AlarmChangeTypeCreate
		return alarmChange, err
	}

	if event.Alarm == nil {
		if event.Entity == nil {
			return alarmChange, nil
		}

		err = s.processPbhEventsForEntity(ctx, event, &alarmChange)

		return alarmChange, err
	}

	if event.Entity == nil {
		return alarmChange, nil
	}

	operation := s.createOperationFromEvent(*event)
	changeType, err := s.executor.Exec(ctx, operation, event.Alarm, event.Entity, event.Timestamp, event.UserID, event.Role, event.Initiator)
	if err != nil {
		return alarmChange, fmt.Errorf("cannot update alarm: %w", err)
	}

	if changeType == types.AlarmChangeTypeResolve {
		go func() {
			err := s.adapter.CopyAlarmToResolvedCollection(context.Background(), *event.Alarm)
			if err != nil {
				s.logger.Err(err).Msg("cannot update resolved alarm")
			}
		}()
	}

	if event.IdleRuleApply != "" {
		event.Entity.LastIdleRuleApply = event.IdleRuleApply
		err := s.entityAdapter.UpdateIdleFields(ctx, event.Entity.ID, event.Entity.IdleSince,
			event.Entity.LastIdleRuleApply)
		if err != nil {
			return alarmChange, fmt.Errorf("cannot update alarm: %w", err)
		}
	}

	go func() {
		err := s.metaAlarmEventProcessor.Process(context.Background(), *event)
		if err != nil {
			s.logger.Err(err).Msg("cannot process meta alarm")
		}
	}()

	go func() {
		if err = s.processAckResources(context.Background(), *event); err != nil {
			s.logger.Err(err).Msg("cannot ack resources")
		}
	}()

	alarmChange.Type = changeType

	return alarmChange, nil
}

func (s *eventProcessor) fillAlarmChange(alarm *types.Alarm, alarmChange *types.AlarmChange) error {
	if alarm != nil {
		alarmChange.PreviousState = alarm.Value.State.Value
		alarmChange.PreviousStateChange = alarm.Value.State.Timestamp
		alarmChange.PreviousStatus = alarm.Value.Status.Value
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
			event.Author, output, event.UserID, event.Role, event.Initiator)
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
		go s.metricsSender.SendCreate(context.Background(), alarm, alarm.Value.CreationDate.Time)
	}

	if changeType == types.AlarmChangeTypeCreateAndPbhEnter {
		event.Entity.PbehaviorInfo = alarm.Value.PbehaviorInfo
		err := s.entityAdapter.UpdatePbehaviorInfo(ctx, event.Entity.ID, event.Entity.PbehaviorInfo)
		if err != nil {
			return changeType, fmt.Errorf("cannot update entity: %w", err)
		}

		go s.metricsSender.SendCreateAndPbhEnter(context.Background(), alarm, alarm.Value.CreationDate.Time)
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

	go s.metricsSender.SendUpdateState(context.Background(), *event.Alarm, *event.Entity, previousState)

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
				event.Author, output, event.UserID, event.Role, event.Initiator)
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
		go s.metricsSender.SendCreate(context.Background(), alarm, alarm.Value.CreationDate.Time)
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
			go s.metricsSender.SendUpdateState(context.Background(), *alarm, *event.Entity, previousState)
		} else if newState < previousState {
			changeType = types.AlarmChangeTypeStateDecrease
			go s.metricsSender.SendUpdateState(context.Background(), *alarm, *event.Entity, previousState)
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

		go s.metricsSender.SendCreateAndPbhEnter(context.Background(), *event.Alarm, event.Alarm.Value.CreationDate.Time)
	}

	return changeType, nil
}

func (s *eventProcessor) createOperationFromEvent(event types.Event) types.Operation {
	var parameters interface{}
	switch event.EventType {
	case types.EventTypeAssocTicket:
		parameters = types.OperationAssocTicketParameters{
			Ticket: event.Ticket,
			Output: event.Output,
			Author: event.Author,
		}
	case types.EventTypeDeclareTicketWebhook:
		parameters = types.OperationDeclareTicketParameters{
			Ticket: event.Ticket,
			Data:   event.TicketData,
			Output: event.Output,
			Author: event.Author,
		}
	case types.EventTypeSnooze:
		parameters = types.OperationSnoozeParameters{
			Duration: types.DurationWithUnit{
				Value: int64(event.Duration),
				Unit:  "s",
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

func (s *eventProcessor) processAckResources(ctx context.Context, event types.Event) error {
	if event.AckResources {
		return s.metaAlarmEventProcessor.ProcessAckResources(ctx, event)
	}

	return nil
}

func (s *eventProcessor) resolveAlarmForDisabledEntity(ctx context.Context, event *types.Event) (types.AlarmChange, error) {
	alarmChange := types.NewAlarmChange()
	alarm, err := s.adapter.GetOpenedAlarm(ctx, event.Connector, event.ConnectorName, event.GetEID())
	if _, ok := err.(errt.NotFound); ok {
		return alarmChange, nil
	} else if err != nil {
		return alarmChange, fmt.Errorf("cannot fetch alarm: %w", err)
	}

	if err := s.fillAlarmChange(event.Alarm, &alarmChange); err != nil {
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
				go s.metricsSender.SendPbhEnter(context.Background(), nil, *event.Entity)
			} else if event.PbehaviorInfo.TypeID == "" {
				alarmChange.Type = types.AlarmChangeTypePbhLeave
				go s.metricsSender.SendPbhLeave(context.Background(), *event.Entity, event.Timestamp.Time, curPbehaviorInfo.CanonicalType, curPbehaviorInfo.Timestamp.Time)
			} else {
				alarmChange.Type = types.AlarmChangeTypePbhLeaveAndEnter
				go s.metricsSender.SendPbhLeaveAndEnter(context.Background(), nil, *event.Entity, curPbehaviorInfo.CanonicalType, curPbehaviorInfo.Timestamp.Time)
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

	stats := statistics.EventStatistics{Timestamp: &event.Timestamp}
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
