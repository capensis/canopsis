package event

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewNoEventsProcessor(
	client mongo.DbClient,
	alarmConfigProvider config.AlarmConfigProvider,
	alarmStatusService alarmstatus.Service,
	pbhTypeResolver pbehavior.EntityTypeResolver,
	autoInstructionMatcher AutoInstructionMatcher,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &noEventsProcessor{
		client:                          client,
		alarmCollection:                 client.Collection(mongo.AlarmMongoCollection),
		entityCollection:                client.Collection(mongo.EntityMongoCollection),
		pbehaviorCollection:             client.Collection(mongo.PbehaviorMongoCollection),
		alarmConfigProvider:             alarmConfigProvider,
		alarmStatusService:              alarmStatusService,
		pbhTypeResolver:                 pbhTypeResolver,
		autoInstructionMatcher:          autoInstructionMatcher,
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		metaAlarmEventProcessor:         metaAlarmEventProcessor,
		metricsSender:                   metricsSender,
		remediationRpcClient:            remediationRpcClient,
		encoder:                         encoder,
		logger:                          logger,
	}
}

type noEventsProcessor struct {
	client                          mongo.DbClient
	alarmCollection                 mongo.DbCollection
	entityCollection                mongo.DbCollection
	pbehaviorCollection             mongo.DbCollection
	alarmConfigProvider             config.AlarmConfigProvider
	alarmStatusService              alarmstatus.Service
	pbhTypeResolver                 pbehavior.EntityTypeResolver
	autoInstructionMatcher          AutoInstructionMatcher
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	eventsSender                    entitycounters.EventsSender
	metaAlarmEventProcessor         libalarm.MetaAlarmEventProcessor
	metricsSender                   metrics.Sender
	remediationRpcClient            engine.RPCClient
	encoder                         encoding.Encoder
	logger                          zerolog.Logger
}

func (p *noEventsProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled || event.Parameters.State == nil {
		return result, nil
	}

	entity := *event.Entity
	var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo

	var componentStateChanged bool
	var newComponentState int

	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil

		alarm := types.Alarm{}
		err := p.alarmCollection.FindOne(ctx, bson.M{
			"d":          entity.ID,
			"v.resolved": nil,
		}).Decode(&alarm)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		if alarm.ID == "" {
			result, err = p.createAlarm(ctx, entity, event.Parameters)
		} else {
			result, err = p.updateAlarm(ctx, alarm, entity, event.Parameters)
		}

		if err != nil {
			return err
		}

		updatedServiceStates, componentStateChanged, newComponentState, err = processComponentAndServiceCounters(
			ctx,
			p.entityServiceCountersCalculator,
			p.componentCountersCalculator,
			&result.Alarm,
			&entity,
			result.AlarmChange,
		)

		return err
	})

	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	result.IsInstructionMatched = isInstructionMatched(event, result, p.autoInstructionMatcher, p.logger)
	go p.postProcess(context.Background(), event, result, updatedServiceStates, componentStateChanged, newComponentState)

	return result, nil
}

func (p *noEventsProcessor) createAlarm(ctx context.Context, entity types.Entity, params rpc.AxeParameters) (Result, error) {
	now := datetime.NewCpsTime()
	result := Result{}
	if *params.State == types.AlarmStateOK {
		return result, nil
	}

	alarmChange := types.NewAlarmChange()
	var pbehaviorInfo types.PbehaviorInfo
	updateEntityPbhInfo := false
	var err error
	if entity.PbehaviorInfo.IsDefaultActive() {
		updateEntityPbhInfo = true
		pbehaviorInfo, err = resolvePbehaviorInfo(ctx, entity, now, p.pbhTypeResolver)
		if err != nil {
			return result, err
		}
	} else {
		pbehaviorInfo = entity.PbehaviorInfo
		pbehaviorInfo.Timestamp = &now
		alarmChange.PreviousPbehaviorTypeID = entity.PbehaviorInfo.TypeID
		alarmChange.PreviousPbehaviorCannonicalType = entity.PbehaviorInfo.CanonicalType
		alarmChange.PreviousEntityPbehaviorTime = entity.PbehaviorInfo.Timestamp
	}

	alarmConfig := p.alarmConfigProvider.Get()
	alarm, err := p.newAlarm(params, entity, now, alarmConfig)
	if err != nil {
		return result, err
	}

	stateStep := types.NewAlarmStep(types.AlarmStepStateIncrease, params.Timestamp, params.Author,
		params.Output, params.User, params.Role, params.Initiator, false)
	stateStep.Value = *params.State
	statusStep := types.NewAlarmStep(types.AlarmStepStatusIncrease, params.Timestamp, params.Author,
		params.Output, params.User, params.Role, params.Initiator, false)
	statusStep.Value = types.AlarmStatusNoEvents
	alarm.Value.State = &stateStep
	err = alarm.Value.Steps.Add(stateStep)
	if err != nil {
		return result, fmt.Errorf("cannot add alarm steps: %w", err)
	}
	alarm.Value.Status = &statusStep
	err = alarm.Value.Steps.Add(statusStep)
	if err != nil {
		return result, fmt.Errorf("cannot add alarm steps: %w", err)
	}

	alarm.Value.TotalStateChanges++

	if pbehaviorInfo.IsDefaultActive() {
		alarmChange.Type = types.AlarmChangeTypeCreate
		alarm.NotAckedSince = &alarm.Time
	} else {
		if pbehaviorInfo.IsActive() {
			alarm.NotAckedSince = &alarm.Time
		} else {
			alarm.Value.InactiveStart = &now
		}

		pbhOutput := fmt.Sprintf(
			"Pbehavior %s. Type: %s. Reason: %s.",
			pbehaviorInfo.Name,
			pbehaviorInfo.TypeName,
			pbehaviorInfo.ReasonName,
		)
		newStep := types.NewAlarmStep(types.AlarmStepPbhEnter, *pbehaviorInfo.Timestamp, canopsis.DefaultEventAuthor,
			pbhOutput, "", "", types.InitiatorSystem, false)
		newStep.PbehaviorCanonicalType = pbehaviorInfo.CanonicalType
		alarm.Value.PbehaviorInfo = pbehaviorInfo
		err := alarm.Value.Steps.Add(newStep)
		if err != nil {
			return result, fmt.Errorf("cannot add alarm steps: %w", err)
		}

		alarmChange.Type = types.AlarmChangeTypeCreateAndPbhEnter
	}

	if p.alarmConfigProvider.Get().ActivateAlarmAfterAutoRemediation {
		matched, err := p.autoInstructionMatcher.Match(alarmChange.GetTriggers(), types.AlarmWithEntity{Alarm: alarm, Entity: entity})
		if err != nil {
			return result, err
		}

		alarm.InactiveAutoInstructionInProgress = matched
	}

	_, err = p.alarmCollection.InsertOne(ctx, alarm)
	if err != nil {
		return result, fmt.Errorf("cannot create alarm: %w", err)
	}

	entityUpdate := bson.M{"$set": bson.M{
		"idle_since":           params.Timestamp,
		"last_idle_rule_apply": params.IdleRuleApply,
	}}
	if alarmChange.Type == types.AlarmChangeTypeCreateAndPbhEnter && updateEntityPbhInfo {
		entityUpdate["pbehavior_info"] = alarm.Value.PbehaviorInfo
		entityUpdate["last_pbehavior_date"] = alarm.Value.PbehaviorInfo.Timestamp
	}

	result.Entity, err = updateEntityByID(ctx, entity.ID, entityUpdate, p.entityCollection)
	if err != nil {
		return result, err
	}

	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}

func (p *noEventsProcessor) updateAlarm(ctx context.Context, alarm types.Alarm, entity types.Entity, params rpc.AxeParameters) (Result, error) {
	result := Result{}
	alarmChange := p.newAlarmChange(alarm)
	previousState := alarm.Value.State.Value
	previousStatus := alarm.Value.Status.Value
	newState := *params.State
	match := bson.M{
		"_id":        alarm.ID,
		"v.resolved": nil,
		"$expr":      bson.M{"$lt": bson.A{bson.M{"$size": "$v.steps"}, types.AlarmStepsHardLimit}},
	}
	set := bson.M{
		"v.output":      params.Output,
		"v.long_output": params.LongOutput,
	}
	push := bson.M{}
	inc := bson.M{}
	unset := bson.M{}

	var stateStep types.AlarmStep
	if newState != previousState {
		stateStep = types.NewAlarmStep(types.AlarmStepStateIncrease, params.Timestamp, params.Author,
			params.Output, params.User, params.Role, params.Initiator, !alarm.Value.PbehaviorInfo.IsDefaultActive())
		stateStep.Value = newState
		alarmChange.Type = types.AlarmChangeTypeStateIncrease
		if newState < previousState {
			alarmChange.Type = types.AlarmChangeTypeStateDecrease
			stateStep.Type = types.AlarmStepStateDecrease
		}

		alarm.Value.State = &stateStep
		err := alarm.Value.Steps.Add(stateStep)
		if err != nil {
			return result, fmt.Errorf("cannot add alarm steps: %w", err)
		}
		set["v.state"] = stateStep
		inc["v.last_update_date"] = params.Timestamp

		if alarm.IsStateLocked() {
			alarm.Value.ChangeState = nil
			unset["v.change_state"] = ""
		}
	}

	newStatus := types.CpsNumber(types.AlarmStatusNoEvents)
	if newState == types.AlarmStateOK {
		newStatus = p.alarmStatusService.ComputeStatus(alarm, entity)
	}

	if newStatus == previousStatus && newState == previousState {
		return result, nil
	}

	if newStatus == previousStatus {
		if stateStep.Type != "" {
			push["v.steps"] = stateStep
			inc["v.total_state_changes"] = 1
			inc["v.state_changes_since_status_update"] = 1
		}
	} else {
		statusStep := types.NewAlarmStep(types.AlarmStepStatusIncrease, params.Timestamp, params.Author,
			params.Output, params.User, params.Role, params.Initiator, !alarm.Value.PbehaviorInfo.IsDefaultActive())
		statusStep.Value = newStatus
		if newStatus < previousStatus {
			statusStep.Type = types.AlarmStepStatusDecrease
		}

		set["v.status"] = statusStep
		set["v.state_changes_since_status_update"] = 0
		if stateStep.Type != "" {
			push["v.steps"] = bson.M{"$each": bson.A{stateStep, statusStep}}
		} else {
			push["v.steps"] = statusStep
		}
	}

	newAlarm := types.Alarm{}
	err := p.alarmCollection.FindOneAndUpdate(ctx, match, bson.M{
		"$set":   set,
		"$push":  push,
		"$inc":   inc,
		"$unset": unset,
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&newAlarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return result, nil
		}
		return result, fmt.Errorf("cannot update alarm: %w", err)
	}

	var entityUpdate bson.M
	if newAlarm.Value.Status.Value == types.AlarmStatusNoEvents {
		entityUpdate = bson.M{"$set": bson.M{
			"idle_since":           params.Timestamp,
			"last_idle_rule_apply": params.IdleRuleApply,
		}}
	} else {
		entityUpdate = bson.M{"$unset": bson.M{
			"idle_since":           "",
			"last_idle_rule_apply": "",
		}}
	}

	result.Entity, err = updateEntityByID(ctx, entity.ID, entityUpdate, p.entityCollection)
	if err != nil {
		return result, err
	}

	result.Forward = true
	result.Alarm = newAlarm
	result.AlarmChange = alarmChange

	return result, nil
}

func (p *noEventsProcessor) newAlarmChange(alarm types.Alarm) types.AlarmChange {
	alarmChange := types.NewAlarmChange()
	alarmChange.PreviousState = alarm.Value.State.Value
	alarmChange.PreviousStateChange = alarm.Value.State.Timestamp
	alarmChange.PreviousStatus = alarm.Value.Status.Value
	return alarmChange
}

func (p *noEventsProcessor) newAlarm(
	params rpc.AxeParameters,
	entity types.Entity,
	timestamp datetime.CpsTime,
	alarmConfig config.AlarmConfig,
) (types.Alarm, error) {
	alarm := types.Alarm{
		EntityID: entity.ID,
		ID:       utils.NewID(),
		Time:     timestamp,
		Value: types.AlarmValue{
			CreationDate:      timestamp,
			DisplayName:       types.GenDisplayName(alarmConfig.DisplayNameScheme),
			InitialOutput:     params.Output,
			Output:            params.Output,
			InitialLongOutput: params.LongOutput,
			LongOutput:        params.LongOutput,
			LongOutputHistory: []string{params.LongOutput},
			LastUpdateDate:    params.Timestamp,
			LastEventDate:     timestamp,
			Parents:           []string{},
			Children:          []string{},
			UnlinkedParents:   []string{},
			Infos:             map[string]map[string]interface{}{},
			RuleVersion:       map[string]string{},
		},
	}

	if params.Initiator != types.InitiatorSystem {
		return types.Alarm{}, fmt.Errorf("unknown initiator %q", params.Initiator)
	}

	connector := ""
	connectorName := ""
	if entity.Connector == "" {
		connector = canopsis.DefaultSystemAlarmConnector
		connectorName = canopsis.DefaultSystemAlarmConnector
	} else {
		connector, connectorName, _ = strings.Cut(entity.Connector, "/")
	}

	switch entity.Type {
	case types.EntityTypeResource:
		alarm.Value.Resource = entity.Name
		alarm.Value.Component = entity.Component
		alarm.Value.Connector = connector
		alarm.Value.ConnectorName = connectorName
	case types.EntityTypeComponent, types.EntityTypeService:
		alarm.Value.Component = entity.Name
		alarm.Value.Connector = connector
		alarm.Value.ConnectorName = connectorName
	case types.EntityTypeConnector:
		alarm.Value.Connector, alarm.Value.ConnectorName, _ = strings.Cut(entity.ID, "/")
	default:
		return types.Alarm{}, fmt.Errorf("unknown entity type %q", entity.Type)
	}

	return alarm, nil
}

func (p *noEventsProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	updatedServiceStates map[string]entitycounters.UpdatedServicesInfo,
	componentStateChanged bool,
	newComponentState int,
) {
	p.metricsSender.SendEventMetrics(
		result.Alarm,
		*event.Entity,
		result.AlarmChange,
		event.Parameters.Timestamp.Time,
		event.Parameters.Initiator,
		event.Parameters.User,
		event.Parameters.Instruction,
		"",
	)

	for servID, servInfo := range updatedServiceStates {
		err := p.eventsSender.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			p.logger.Err(err).Msg("failed to update service state")
		}
	}

	if componentStateChanged {
		err := p.eventsSender.UpdateComponentState(ctx, event.Entity.Component, newComponentState)
		if err != nil {
			p.logger.Err(err).Msg("failed to update component state")
		}
	}

	err := p.metaAlarmEventProcessor.ProcessAxeRpc(ctx, event, rpc.AxeResultEvent{
		Alarm:           &result.Alarm,
		AlarmChangeType: result.AlarmChange.Type,
	})
	if err != nil {
		p.logger.Err(err).Msg("cannot process meta alarm")
	}

	err = sendRemediationEvent(ctx, event, result, p.remediationRpcClient, p.encoder)
	if err != nil {
		p.logger.Err(err).Msg("cannot send event to engine-remediation")
	}

	err = updatePbhLastAlarmDate(ctx, result, p.pbehaviorCollection)
	if err != nil {
		p.logger.Err(err).Msg("cannot update pbehavior")
	}
}
