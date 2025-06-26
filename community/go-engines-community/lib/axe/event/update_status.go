package event

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewUpdateStatusProcessor(
	dbClient mongo.DbClient,
	alarmConfigProvider config.AlarmConfigProvider,
	alarmStatusService alarmstatus.Service,
	pbhTypeResolver pbehavior.EntityTypeResolver,
	autoInstructionMatcher AutoInstructionMatcher,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	internalTagAlarmMatcher alarmtag.InternalTagAlarmMatcher,
	eventGenerator libevent.Generator,
	amqpPublisher amqp.Publisher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &updateStatusProcessor{
		dbClient:                        dbClient,
		alarmCollection:                 dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:                dbClient.Collection(mongo.EntityMongoCollection),
		pbehaviorCollection:             dbClient.Collection(mongo.PbehaviorMongoCollection),
		alarmConfigProvider:             alarmConfigProvider,
		alarmStatusService:              alarmStatusService,
		pbhTypeResolver:                 pbhTypeResolver,
		autoInstructionMatcher:          autoInstructionMatcher,
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		metaAlarmPostProcessor:          metaAlarmPostProcessor,
		metricsSender:                   metricsSender,
		remediationRpcClient:            remediationRpcClient,
		internalTagAlarmMatcher:         internalTagAlarmMatcher,
		eventGenerator:                  eventGenerator,
		amqpPublisher:                   amqpPublisher,
		encoder:                         encoder,
		logger:                          logger,
	}
}

type updateStatusProcessor struct {
	dbClient                        mongo.DbClient
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
	metaAlarmPostProcessor          MetaAlarmPostProcessor
	metricsSender                   metrics.Sender
	remediationRpcClient            engine.RPCClient
	internalTagAlarmMatcher         alarmtag.InternalTagAlarmMatcher
	eventGenerator                  libevent.Generator
	amqpPublisher                   amqp.Publisher
	encoder                         encoding.Encoder
	logger                          zerolog.Logger
}

func (p *updateStatusProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || event.Entity.ID == "" || !event.Entity.Enabled {
		return result, nil
	}

	entity := *event.Entity
	var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo
	var componentStateChanged bool
	var newComponentState int
	match := getOpenAlarmMatch(event)
	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		entity = *event.Entity
		updatedServiceStates = nil
		componentStateChanged = false
		newComponentState = 0

		alarm := types.Alarm{}
		err := p.alarmCollection.FindOne(ctx, match).Decode(&alarm)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		currentStatus := types.CpsNumber(types.AlarmStatusOff)
		if alarm.ID != "" {
			currentStatus = alarm.Value.Status.Value
		}

		newStatus, statusRuleName, err := p.alarmStatusService.ComputeStatusOnStatusChange(ctx, alarm, *event.Entity)
		if err != nil {
			return fmt.Errorf("cannot compute alarm status: %w", err)
		}

		if newStatus == currentStatus {
			return nil
		}

		if alarm.ID == "" {
			if newStatus != types.AlarmStatusUnknown {
				return nil
			}

			var v types.Entity
			err = p.entityCollection.FindOne(ctx, bson.M{"_id": entity.ID}).Decode(&v)
			if err != nil {
				return err
			}

			entity = v
			result, err = p.createAlarm(ctx, entity, event.Parameters, newStatus, statusRuleName)
			if err != nil {
				return err
			}

			result.IsCountersUpdated, updatedServiceStates, componentStateChanged, newComponentState, err = processComponentAndServiceCounters(
				ctx,
				p.entityServiceCountersCalculator,
				p.componentCountersCalculator,
				&result.Alarm,
				&entity,
				result.AlarmChange,
			)

			return err
		}

		prevStatus := alarm.Value.Status.Value
		if prevStatus == types.AlarmStatusUnknown && alarm.Value.InitialStatus == types.AlarmStatusUnknown &&
			alarm.Value.LastEventDate == nil && !alarm.IsStateLocked() {
			newState := types.CpsNumber(types.AlarmStateOK)
			if newStatus == types.AlarmStatusOngoing {
				newStatus = types.AlarmStatusOff
			}

			result, err = p.updateAlarmStateAndStatus(ctx, alarm, entity, event.Parameters, newState, newStatus, statusRuleName)
			if err != nil {
				return err
			}

			result.IsCountersUpdated, updatedServiceStates, componentStateChanged, newComponentState, err = processComponentAndServiceCounters(
				ctx,
				p.entityServiceCountersCalculator,
				p.componentCountersCalculator,
				&result.Alarm,
				&entity,
				result.AlarmChange,
			)

			return err
		}

		result, err = p.updateAlarmStatus(ctx, alarm, event, newStatus, statusRuleName)

		return err
	})
	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	if result.Alarm.ID != "" {
		go p.postProcess(context.WithoutCancel(ctx), event, result, updatedServiceStates, componentStateChanged, newComponentState)
	}

	return result, nil
}

func (p *updateStatusProcessor) updateAlarmStatus(
	ctx context.Context,
	alarm types.Alarm,
	event rpc.AxeEvent,
	newStatus types.CpsNumber,
	statusRuleName string,
) (Result, error) {
	result := Result{}
	alarmStepType := types.AlarmStepStatusIncrease
	prevStatus := alarm.Value.Status.Value
	if prevStatus > newStatus {
		alarmStepType = types.AlarmStepStatusDecrease
	}

	statusStepMessage := ConcatOutputAndRuleName(event.Parameters.Output, statusRuleName)
	newStepStatusQuery := valStepUpdateQueryWithInPbhInterval(alarmStepType, newStatus, statusStepMessage, event.Parameters)
	matchUpdate := getOpenAlarmMatchWithStepsLimit(event)
	update := []bson.M{
		{"$set": bson.M{
			"v.status":                            newStepStatusQuery,
			"v.state_changes_since_status_update": 0,
			"v.last_update_date":                  event.Parameters.Timestamp,
			"v.last_st_upd_dt":                    event.Parameters.Timestamp,
			"v.steps":                             addStepUpdateQuery(newStepStatusQuery),
		}},
	}
	if newStatus == types.AlarmStatusUnknown {
		update = append(update, bson.M{"$unset": bson.A{"v.canceled"}})
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	updatedAlarm := types.Alarm{}
	err := p.alarmCollection.FindOneAndUpdate(ctx, matchUpdate, update, opts).Decode(&updatedAlarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return result, nil
		}

		return result, err
	}

	alarmChange := types.NewAlarmChange()
	alarmChange.PreviousStatus = prevStatus
	alarmChange.Type = types.AlarmChangeTypeUpdateStatus
	result.Forward = true
	result.Alarm = updatedAlarm
	result.AlarmChange = alarmChange

	return result, nil
}

func (p *updateStatusProcessor) createAlarm(
	ctx context.Context,
	entity types.Entity,
	params rpc.AxeParameters,
	newStatus types.CpsNumber,
	statusRuleName string,
) (Result, error) {
	now := datetime.NewCpsTime()
	result := Result{}
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

	stateStep := NewAlarmStep(types.AlarmStepStateIncrease, params, false)
	stateStep.Value = types.AlarmStateMinor
	stateStep.Message = ConcatOutputAndRuleName(params.Output, statusRuleName)
	alarm.Value.State = &stateStep
	alarm.Value.MaxState = stateStep.Value
	alarm.Value.InitialState = stateStep.Value
	err = alarm.Value.Steps.Add(stateStep)
	if err != nil {
		return result, fmt.Errorf("cannot add alarm steps: %w", err)
	}

	statusStep := NewAlarmStep(types.AlarmStepStatusIncrease, params, false)
	statusStep.Value = newStatus
	statusStep.Message = ConcatOutputAndRuleName(params.Output, statusRuleName)
	alarm.Value.Status = &statusStep
	alarm.Value.InitialStatus = statusStep.Value
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

		newStep := types.NewPbhAlarmStep(types.AlarmStepPbhEnter, *pbehaviorInfo.Timestamp, pbehaviorInfo.Author,
			pbehaviorInfo.GetStepMessage(), "", "", types.InitiatorSystem, pbehaviorInfo.CanonicalType,
			pbehaviorInfo.IconName, pbehaviorInfo.Color)
		alarm.Value.PbehaviorInfo = pbehaviorInfo
		err := alarm.Value.Steps.Add(newStep)
		if err != nil {
			return result, fmt.Errorf("cannot add alarm steps: %w", err)
		}

		alarmChange.Type = types.AlarmChangeTypeCreateAndPbhEnter
	}

	result.IsInstructionMatched, err = p.autoInstructionMatcher.Match(alarmChange.GetTriggers(), types.AlarmWithEntity{Alarm: alarm, Entity: entity})
	if err != nil {
		return result, err
	}

	if p.alarmConfigProvider.Get().ActivateAlarmAfterAutoRemediation {
		alarm.InactiveAutoInstructionInProgress = result.IsInstructionMatched
	}

	alarm.InternalTags = p.internalTagAlarmMatcher.Match(entity, alarm)
	alarm.InternalTagsUpdated = datetime.NewMicroTime()
	alarm.Tags = slices.Clone(alarm.InternalTags)

	_, err = p.alarmCollection.InsertOne(ctx, types.AlarmWithEntityField{
		Alarm:  alarm,
		Entity: entity,
	})
	if err != nil {
		return result, fmt.Errorf("cannot create alarm: %w", err)
	}

	entityUpdate := bson.M{}
	if alarmChange.Type == types.AlarmChangeTypeCreateAndPbhEnter && updateEntityPbhInfo {
		entityUpdate["pbehavior_info"] = alarm.Value.PbehaviorInfo
		entityUpdate["last_pbehavior_date"] = alarm.Value.PbehaviorInfo.Timestamp
	}

	if len(entityUpdate) > 0 {
		result.Entity, err = updateEntityByID(ctx, entity.ID, bson.M{"$set": entityUpdate}, p.entityCollection)
		if err != nil {
			return result, err
		}
	}

	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}

func (p *updateStatusProcessor) updateAlarmStateAndStatus(
	ctx context.Context,
	alarm types.Alarm,
	entity types.Entity,
	params rpc.AxeParameters,
	newState, newStatus types.CpsNumber,
	statusRuleName string,
) (Result, error) {
	result := Result{}
	alarmChange := types.NewAlarmChange()
	alarmChange.PreviousState = alarm.Value.State.Value
	alarmChange.PreviousStatusChange = alarm.Value.State.Timestamp
	alarmChange.PreviousStatus = alarm.Value.Status.Value
	stateStep := NewAlarmStep(types.AlarmStepStateIncrease, params, !alarm.Value.PbehaviorInfo.IsDefaultActive())
	stateStep.Value = newState
	alarmChange.Type = types.AlarmChangeTypeStateIncrease
	set := bson.M{}
	if newState < alarmChange.PreviousState {
		alarmChange.Type = types.AlarmChangeTypeStateDecrease
		stateStep.Type = types.AlarmStepStateDecrease
	} else if alarm.Value.MaxState < newState {
		alarm.Value.MaxState = newState
		set["v.max_state"] = newState
	}

	alarm.Value.State = &stateStep
	err := alarm.Value.Steps.Add(stateStep)
	if err != nil {
		return result, fmt.Errorf("cannot add alarm steps: %w", err)
	}

	set["v.state"] = stateStep
	set["v.last_update_date"] = params.Timestamp
	set["v.last_st_upd_dt"] = params.Timestamp

	statusStep := NewAlarmStep(types.AlarmStepStatusIncrease, params, !alarm.Value.PbehaviorInfo.IsDefaultActive())
	statusStep.Value = newStatus
	statusStep.Message = statusRuleName
	if newStatus < alarmChange.PreviousStatus {
		statusStep.Type = types.AlarmStepStatusDecrease
	}

	set["v.status"] = statusStep
	set["v.state_changes_since_status_update"] = 0
	push := bson.M{}
	if stateStep.Type != "" {
		push["v.steps"] = bson.M{"$each": bson.A{stateStep, statusStep}}
	} else {
		push["v.steps"] = statusStep
	}

	match := bson.M{
		"_id":        alarm.ID,
		"v.resolved": nil,
		"$expr":      bson.M{"$lt": bson.A{bson.M{"$size": "$v.steps"}, types.AlarmStepsHardLimit}},
	}
	newAlarm := types.Alarm{}
	err = p.alarmCollection.FindOneAndUpdate(ctx, match, bson.M{
		"$set":  set,
		"$push": push,
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&newAlarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return result, nil
		}

		return result, fmt.Errorf("cannot update alarm: %w", err)
	}

	result.Forward = true
	result.Alarm = newAlarm
	result.AlarmChange = alarmChange
	result.IsInstructionMatched, err = p.autoInstructionMatcher.Match(alarmChange.GetTriggers(), types.AlarmWithEntity{Alarm: newAlarm, Entity: entity})

	return result, err
}

func (p *updateStatusProcessor) newAlarm(
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
			CreationDate:                timestamp,
			DisplayName:                 types.GenDisplayName(alarmConfig.DisplayNameScheme),
			InitialOutput:               params.Output,
			Output:                      params.Output,
			InitialLongOutput:           params.LongOutput,
			LongOutput:                  params.LongOutput,
			LongOutputHistory:           []string{params.LongOutput},
			LastUpdateDate:              params.Timestamp,
			LastStateOrStatusUpdateDate: params.Timestamp,
			Parents:                     []string{},
			Children:                    []string{},
			UnlinkedParents:             []string{},
			Infos:                       map[string]map[string]interface{}{},
		},
	}

	if params.Initiator != types.InitiatorSystem {
		return types.Alarm{}, fmt.Errorf("unknown initiator %q", params.Initiator)
	}

	connector := canopsis.DefaultSystemAlarmConnector
	connectorName := canopsis.DefaultSystemAlarmConnector
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
	default:
		return types.Alarm{}, fmt.Errorf("unknown entity type %q", entity.Type)
	}

	return alarm, nil
}

func (p *updateStatusProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	updatedServiceStates map[string]entitycounters.UpdatedServicesInfo,
	componentStateChanged bool,
	newComponentState int,
) {
	entity := *event.Entity
	if result.Entity.ID != "" {
		entity = result.Entity
	}

	p.metricsSender.SendEventMetrics(
		result.Alarm,
		entity,
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

	err := p.metaAlarmPostProcessor.Process(ctx, event, rpc.AxeResultEvent{
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

	if result.AlarmChange.Type == types.AlarmChangeTypeCreateAndPbhEnter {
		err = updatePbehaviorLastAlarmDate(ctx, p.pbehaviorCollection, result.Alarm.Value.PbehaviorInfo.ID, result.Alarm.Value.PbehaviorInfo.Timestamp)
		if err != nil {
			p.logger.Err(err).Msg("cannot update pbehavior")
		}

		if !result.Alarm.Value.PbehaviorInfo.IsDefaultActive() {
			err = updatePbehaviorAlarmCount(ctx, p.pbehaviorCollection, result.Alarm.Value.PbehaviorInfo.ID, "")
			if err != nil {
				p.logger.Err(err).Msg("cannot update pbehavior")
			}
		}
	}

	switch result.AlarmChange.Type {
	case types.AlarmChangeTypeUpdateStatus:
		if result.Alarm.Value.Status.Value == types.AlarmStatusOff ||
			result.AlarmChange.PreviousStatus == types.AlarmStatusUnknown {
			err = sendDownstreamEventsOnOK(ctx, entity, p.entityCollection, p.eventGenerator, p.encoder, p.amqpPublisher)
			if err != nil {
				p.logger.Err(err).Msg("cannot send downstream events")
			}
		} else if result.Alarm.Value.Status.Value == types.AlarmStatusUnknown {
			err = sendDownstreamEventsOnKO(ctx, entity, p.entityCollection, p.eventGenerator, p.encoder, p.amqpPublisher)
			if err != nil {
				p.logger.Err(err).Msg("cannot send downstream events")
			}
		}
	case types.AlarmChangeTypeCreate, types.AlarmChangeTypeCreateAndPbhEnter:
		err = sendDownstreamEventsOnKO(ctx, entity, p.entityCollection, p.eventGenerator, p.encoder, p.amqpPublisher)
		if err != nil {
			p.logger.Err(err).Msg("cannot send downstream events")
		}
	case types.AlarmChangeTypeStateDecrease:
		alarmStatus := result.Alarm.Value.Status.Value
		alarmState := result.Alarm.Value.State.Value
		if result.AlarmChange.PreviousStatus != alarmStatus && alarmStatus == types.AlarmStatusOff ||
			result.AlarmChange.PreviousStatus == alarmStatus && alarmStatus == types.AlarmStatusCancelled && alarmState == types.AlarmStateOK {
			err = sendDownstreamEventsOnOK(ctx, entity, p.entityCollection, p.eventGenerator, p.encoder, p.amqpPublisher)
			if err != nil {
				p.logger.Err(err).Msg("cannot send downstream events")
			}
		}
	}
}
