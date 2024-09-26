package event

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewCheckProcessor(
	client mongo.DbClient,
	alarmConfigProvider config.AlarmConfigProvider,
	alarmStatusService alarmstatus.Service,
	pbhTypeResolver pbehavior.EntityTypeResolver,
	autoInstructionMatcher AutoInstructionMatcher,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metricsSender metrics.Sender,
	eventStatisticsSender statistics.EventStatisticsSender,
	remediationRpcClient engine.RPCClient,
	externalTagUpdater alarmtag.ExternalUpdater,
	internalTagAlarmMatcher alarmtag.InternalTagAlarmMatcher,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &checkProcessor{
		client:                          client,
		alarmCollection:                 client.Collection(mongo.AlarmMongoCollection),
		entityCollection:                client.Collection(mongo.EntityMongoCollection),
		pbehaviorCollection:             client.Collection(mongo.PbehaviorMongoCollection),
		alarmConfigProvider:             alarmConfigProvider,
		alarmStatusService:              alarmStatusService,
		pbhTypeResolver:                 pbhTypeResolver,
		autoInstructionMatcher:          autoInstructionMatcher,
		metaAlarmPostProcessor:          metaAlarmPostProcessor,
		metricsSender:                   metricsSender,
		eventStatisticsSender:           eventStatisticsSender,
		remediationRpcClient:            remediationRpcClient,
		externalTagUpdater:              externalTagUpdater,
		internalTagAlarmMatcher:         internalTagAlarmMatcher,
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		encoder:                         encoder,
		logger:                          logger,
	}
}

type checkProcessor struct {
	client                          mongo.DbClient
	alarmCollection                 mongo.DbCollection
	entityCollection                mongo.DbCollection
	pbehaviorCollection             mongo.DbCollection
	alarmConfigProvider             config.AlarmConfigProvider
	alarmStatusService              alarmstatus.Service
	pbhTypeResolver                 pbehavior.EntityTypeResolver
	autoInstructionMatcher          AutoInstructionMatcher
	metaAlarmPostProcessor          MetaAlarmPostProcessor
	metricsSender                   metrics.Sender
	eventStatisticsSender           statistics.EventStatisticsSender
	remediationRpcClient            engine.RPCClient
	externalTagUpdater              alarmtag.ExternalUpdater
	internalTagAlarmMatcher         alarmtag.InternalTagAlarmMatcher
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	eventsSender                    entitycounters.EventsSender
	encoder                         encoding.Encoder
	logger                          zerolog.Logger
}

func (p *checkProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil ||
		!event.Entity.Enabled ||
		event.Parameters.State == nil ||
		event.Entity.StateInfo != nil && event.Parameters.Initiator != types.InitiatorSystem && !event.Parameters.StateSettingUpdated {
		return result, nil
	}

	entity := *event.Entity
	var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo

	var componentStateChanged bool
	var newComponentState int

	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		entity = *event.Entity
		updatedServiceStates = nil
		componentStateChanged = false
		newComponentState = 0

		alarm := types.Alarm{}
		err := p.alarmCollection.FindOne(ctx, bson.M{
			"d":          entity.ID,
			"v.resolved": nil,
		}).Decode(&alarm)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		if alarm.ID == "" {
			var v types.Entity
			err = p.entityCollection.FindOne(ctx, bson.M{"_id": entity.ID}).Decode(&v)
			if err != nil {
				return err
			}

			entity = v
			result, err = p.createAlarm(ctx, entity, event)
		} else {
			result, err = p.updateAlarm(ctx, alarm, entity, event.Parameters)
		}

		if err != nil {
			return err
		}

		if result.Entity.ID != "" {
			entity = result.Entity
		}

		if !event.Healthcheck {
			result.IsCountersUpdated, updatedServiceStates, componentStateChanged, newComponentState, err = processComponentAndServiceCounters(
				ctx,
				p.entityServiceCountersCalculator,
				p.componentCountersCalculator,
				&result.Alarm,
				&entity,
				result.AlarmChange,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return result, err
	}

	if result.Alarm.ID != "" {
		result.AlarmChange.EventsCount = int(result.Alarm.Value.EventsCount)
	}

	if !event.Healthcheck {
		go p.postProcess(context.Background(), event, result, updatedServiceStates, componentStateChanged, newComponentState)
	}

	return result, nil
}

func (p *checkProcessor) createAlarm(ctx context.Context, entity types.Entity, event rpc.AxeEvent) (Result, error) {
	params := event.Parameters
	now := datetime.NewCpsTime()
	result := Result{
		Forward: true,
	}

	if event.Parameters.StateSettingUpdated {
		componentState, err := p.componentCountersCalculator.RecomputeCounters(ctx, &entity)
		if err != nil {
			return Result{}, err
		}

		*params.State = types.CpsNumber(componentState)
	}

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

	author := ""
	if params.Initiator == types.InitiatorExternal {
		author = params.Connector + "." + params.ConnectorName
	} else {
		author = params.Author
	}

	alarmConfig := p.alarmConfigProvider.Get()
	alarm, err := p.newAlarm(params, entity, now, alarmConfig)
	if err != nil {
		return result, err
	}

	stateStep := NewAlarmStep(types.AlarmStepStateIncrease, params, false)
	stateStep.Author = author
	stateStep.Value = *params.State
	statusStep := NewAlarmStep(types.AlarmStepStatusIncrease, params, false)
	statusStep.Author = author
	statusStep.Value = types.AlarmStatusOngoing
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

	alarm.Value.EventsCount++
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

	if alarmConfig.ActivateAlarmAfterAutoRemediation {
		alarm.InactiveAutoInstructionInProgress = result.IsInstructionMatched
	}

	alarm.InternalTags = p.internalTagAlarmMatcher.Match(entity, alarm)
	alarm.InternalTagsUpdated = datetime.NewMicroTime()
	alarm.Tags = append(alarm.Tags, alarm.InternalTags...)
	alarm.Healthcheck = event.Healthcheck
	_, err = p.alarmCollection.InsertOne(ctx, types.AlarmWithEntityField{
		Alarm:  alarm,
		Entity: entity,
	})
	if err != nil {
		return result, fmt.Errorf("cannot create alarm: %w", err)
	}

	if alarmChange.Type == types.AlarmChangeTypeCreateAndPbhEnter && updateEntityPbhInfo {
		updateRes, err := p.entityCollection.UpdateOne(ctx, bson.M{"_id": entity.ID},
			bson.M{"$set": bson.M{
				"pbehavior_info":      alarm.Value.PbehaviorInfo,
				"last_pbehavior_date": alarm.Value.PbehaviorInfo.Timestamp,
			}},
		)
		if err != nil {
			return result, fmt.Errorf("cannot update entity: %w", err)
		}

		if updateRes.ModifiedCount > 0 {
			entity.PbehaviorInfo = alarm.Value.PbehaviorInfo
			result.Entity = entity
		}
	}

	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}

func (p *checkProcessor) updateAlarm(ctx context.Context, alarm types.Alarm, entity types.Entity, params rpc.AxeParameters) (Result, error) {
	result := Result{
		Forward: true,
	}

	var newState types.CpsNumber

	if params.StateSettingUpdated {
		componentState, err := p.componentCountersCalculator.RecomputeCounters(ctx, &entity)
		if err != nil {
			return Result{}, err
		}

		newState = types.CpsNumber(componentState)
	} else {
		newState = *params.State
	}

	alarmChange := p.newAlarmChange(alarm)
	previousState := alarm.Value.State.Value
	previousStatus := alarm.Value.Status.Value
	match := bson.M{"_id": alarm.ID, "v.resolved": nil}
	set := bson.M{
		"v.output":          params.Output,
		"v.last_event_date": params.Timestamp,
		"v.long_output":     params.LongOutput,
	}
	push := bson.M{}
	inc := bson.M{
		"v.events_count": 1,
	}
	unset := bson.M{}
	author := ""
	if params.Initiator == types.InitiatorExternal {
		author = params.Connector + "." + params.ConnectorName
	} else {
		author = params.Author
	}

	if alarm.Value.LongOutputHistory[len(alarm.Value.LongOutputHistory)-1] != params.LongOutput {
		if len(alarm.Value.LongOutputHistory) < types.AlarmLongOutputHistoryLimit {
			push["v.long_output_history"] = params.LongOutput
		} else {
			history := append(alarm.Value.LongOutputHistory, params.LongOutput)
			set["v.long_output_history"] = history[len(history)-types.AlarmLongOutputHistoryLimit:]
		}
	}

	var stateStep types.AlarmStep
	if newState != previousState && (!alarm.IsStateLocked() || newState == types.AlarmStateOK) {
		stateStep = NewAlarmStep(types.AlarmStepStateIncrease, params, !alarm.Value.PbehaviorInfo.IsDefaultActive())
		stateStep.Author = author
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
		set["v.last_update_date"] = params.Timestamp
		inc["v.total_state_changes"] = 1

		if alarm.IsStateLocked() {
			alarm.Value.ChangeState = nil
			unset["v.change_state"] = ""
		}
	}

	newStatus, statusRuleName := p.alarmStatusService.ComputeStatus(alarm, entity)
	if newStatus == previousStatus {
		if stateStep.Type != "" {
			match["$expr"] = bson.M{"$lt": bson.A{bson.M{"$size": "$v.steps"}, types.AlarmStepsHardLimit}}
			push["v.steps"] = stateStep
			inc["v.state_changes_since_status_update"] = 1
		}
	} else {
		if newStatus != types.AlarmStatusOff {
			match["$expr"] = bson.M{"$lt": bson.A{bson.M{"$size": "$v.steps"}, types.AlarmStepsHardLimit}}
		}

		statusStep := NewAlarmStep(types.AlarmStepStatusIncrease, params, !alarm.Value.PbehaviorInfo.IsDefaultActive())
		statusStep.Message = ConcatOutputAndRuleName(params.Output, statusRuleName)
		statusStep.Author = author
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

	addToSet := bson.M{}
	newExternalTags := p.getNewExternalTags(alarm.ExternalTags, params.Tags)
	if len(newExternalTags) > 0 {
		addToSet["tags"] = bson.M{"$each": newExternalTags}
		addToSet["etags"] = bson.M{"$each": newExternalTags}
	}
	newAlarm := types.Alarm{}
	err := p.alarmCollection.FindOneAndUpdate(ctx, match, bson.M{
		"$set":      set,
		"$push":     push,
		"$inc":      inc,
		"$addToSet": addToSet,
		"$unset":    unset,
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&newAlarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return result, nil
		}
		return result, fmt.Errorf("cannot update alarm: %w", err)
	}

	// Update cropped steps if needed
	cropStepsNumber := p.alarmConfigProvider.Get().CropStepsNumber
	if cropStepsNumber > 0 && newAlarm.CropSteps(cropStepsNumber) {
		_, err = p.alarmCollection.UpdateOne(ctx, bson.M{"_id": newAlarm.ID, "v.resolved": nil}, bson.M{
			"$set": bson.M{
				"v.steps": newAlarm.Value.Steps,
			},
		})
		if err != nil {
			return result, fmt.Errorf("cannot update alarm: %w", err)
		}
	}

	if entity.IdleSince != nil || entity.LastIdleRuleApply != "" {
		unsetIdleFields := bson.M{"idle_since": ""}
		alarmLastUpdateRule := fmt.Sprintf("%s_%s", idlerule.RuleTypeAlarm, idlerule.RuleAlarmConditionLastUpdate)
		if entity.LastIdleRuleApply != alarmLastUpdateRule ||
			entity.LastIdleRuleApply == alarmLastUpdateRule && alarmChange.Type != types.AlarmChangeTypeNone {
			unsetIdleFields["last_idle_rule_apply"] = ""
		}

		result.Entity, err = updateEntityByID(ctx, entity.ID, bson.M{"$unset": unsetIdleFields}, p.entityCollection)
		if err != nil {
			return result, err
		}
	}

	result.Alarm = newAlarm
	result.AlarmChange = alarmChange

	return result, nil
}

func (p *checkProcessor) newAlarmChange(alarm types.Alarm) types.AlarmChange {
	alarmChange := types.NewAlarmChange()
	alarmChange.PreviousState = alarm.Value.State.Value
	alarmChange.PreviousStateChange = alarm.Value.State.Timestamp
	alarmChange.PreviousStatus = alarm.Value.Status.Value
	return alarmChange
}

func (p *checkProcessor) newAlarm(
	params rpc.AxeParameters,
	entity types.Entity,
	timestamp datetime.CpsTime,
	alarmConfig config.AlarmConfig,
) (types.Alarm, error) {
	tags := types.TransformEventTags(params.Tags)
	alarm := types.Alarm{
		EntityID:     entity.ID,
		ID:           utils.NewID(),
		Time:         timestamp,
		Tags:         tags,
		ExternalTags: tags,
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

	connector := ""
	connectorName := ""
	switch params.Initiator {
	case types.InitiatorExternal, types.InitiatorUser:
		connector = params.Connector
		connectorName = params.ConnectorName
	case types.InitiatorSystem:
		if entity.Connector == "" {
			connector = canopsis.DefaultSystemAlarmConnector
			connectorName = canopsis.DefaultSystemAlarmConnector
		} else {
			connector, connectorName, _ = strings.Cut(entity.Connector, "/")
		}
	default:
		return types.Alarm{}, fmt.Errorf("unknown initiator %q", params.Initiator)
	}

	switch entity.Type {
	case types.EntityTypeResource:
		alarm.Value.Resource = entity.Name
		alarm.Value.Component = entity.Component
		alarm.Value.Connector = connector
		alarm.Value.ConnectorName = connectorName
	case types.EntityTypeComponent, types.EntityTypeService:
		alarm.Value.Component = entity.ID
		alarm.Value.Connector = connector
		alarm.Value.ConnectorName = connectorName
	case types.EntityTypeConnector:
		alarm.Value.Connector, alarm.Value.ConnectorName, _ = strings.Cut(entity.ID, "/")
	default:
		return types.Alarm{}, fmt.Errorf("unknown entity type %q", entity.Type)
	}

	return alarm, nil
}

func (p *checkProcessor) getNewExternalTags(alarmTags []string, eventTags map[string]string) []string {
	exists := make(map[string]struct{}, len(alarmTags))
	for _, tag := range alarmTags {
		exists[tag] = struct{}{}
	}

	tags := types.TransformEventTags(eventTags)
	var k = 0
	for _, tag := range tags {
		if _, ok := exists[tag]; !ok {
			tags[k] = tag
			k++
		}
	}
	if k == 0 {
		return nil
	}
	tags = tags[:k]
	return tags
}

func (p *checkProcessor) postProcess(
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

	p.externalTagUpdater.Add(event.Parameters.Tags)

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

	p.sendEventStatistics(ctx, event)

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

		if event.Entity.PbehaviorInfo.IsDefaultActive() {
			err = updatePbehaviorAlarmCount(ctx, p.pbehaviorCollection, result.Alarm.Value.PbehaviorInfo.ID, "")
			if err != nil {
				p.logger.Err(err).Msg("cannot update pbehavior")
			}
		}
	}
}

func (p *checkProcessor) sendEventStatistics(ctx context.Context, event rpc.AxeEvent) {
	if event.Entity.PbehaviorInfo.Is(pbehavior.TypeInactive) {
		return
	}

	stats := statistics.EventStatistics{LastEvent: &event.Parameters.Timestamp}
	if *event.Parameters.State == types.AlarmStateOK {
		stats.OK = 1
	} else {
		stats.KO = 1
	}

	p.eventStatisticsSender.Send(ctx, event.Entity.ID, stats)
}
