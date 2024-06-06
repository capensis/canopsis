package event

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice/statecounters"
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
	stateCountersService statecounters.StateCountersService,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	eventStatisticsSender statistics.EventStatisticsSender,
	remediationRpcClient engine.RPCClient,
	externalTagUpdater alarmtag.ExternalUpdater,
	internalTagAlarmMatcher alarmtag.InternalTagAlarmMatcher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &checkProcessor{
		client:                  client,
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		entityCollection:        client.Collection(mongo.EntityMongoCollection),
		pbehaviorCollection:     client.Collection(mongo.PbehaviorMongoCollection),
		alarmConfigProvider:     alarmConfigProvider,
		alarmStatusService:      alarmStatusService,
		pbhTypeResolver:         pbhTypeResolver,
		autoInstructionMatcher:  autoInstructionMatcher,
		stateCountersService:    stateCountersService,
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		metricsSender:           metricsSender,
		eventStatisticsSender:   eventStatisticsSender,
		remediationRpcClient:    remediationRpcClient,
		externalTagUpdater:      externalTagUpdater,
		internalTagAlarmMatcher: internalTagAlarmMatcher,
		encoder:                 encoder,
		logger:                  logger,
	}
}

type checkProcessor struct {
	client                  mongo.DbClient
	alarmCollection         mongo.DbCollection
	entityCollection        mongo.DbCollection
	pbehaviorCollection     mongo.DbCollection
	alarmConfigProvider     config.AlarmConfigProvider
	alarmStatusService      alarmstatus.Service
	pbhTypeResolver         pbehavior.EntityTypeResolver
	autoInstructionMatcher  AutoInstructionMatcher
	stateCountersService    statecounters.StateCountersService
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	metricsSender           metrics.Sender
	eventStatisticsSender   statistics.EventStatisticsSender
	remediationRpcClient    engine.RPCClient
	externalTagUpdater      alarmtag.ExternalUpdater
	internalTagAlarmMatcher alarmtag.InternalTagAlarmMatcher
	encoder                 encoding.Encoder
	logger                  zerolog.Logger
}

func (p *checkProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled || event.Parameters.State == nil {
		return result, nil
	}

	entity := *event.Entity
	var updatedServiceStates map[string]statecounters.UpdatedServicesInfo

	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		entity = *event.Entity
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
			if result.Alarm.ID == "" {
				updatedServiceStates, err = p.stateCountersService.UpdateServiceCounters(ctx, entity, nil, result.AlarmChange)
			} else {
				updatedServiceStates, err = p.stateCountersService.UpdateServiceCounters(ctx, entity, &result.Alarm, result.AlarmChange)
			}
		}

		return err
	})

	if err != nil {
		return result, err
	}

	if result.Alarm.ID != "" {
		result.IsInstructionMatched = isInstructionMatched(event, result, p.autoInstructionMatcher, p.logger)
		result.AlarmChange.EventsCount = int(result.Alarm.Value.EventsCount)
	}

	if !event.Healthcheck {
		go p.postProcess(context.Background(), event, result, updatedServiceStates)
	}

	return result, nil
}

func (p *checkProcessor) createAlarm(ctx context.Context, entity types.Entity, event rpc.AxeEvent) (Result, error) {
	params := event.Parameters
	now := types.NewCpsTime()
	result := Result{
		Forward: true,
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
	if entity.Type != types.EntityTypeService {
		author = strings.Replace(entity.Connector, "/", ".", 1)
	} else {
		author = params.Connector + "." + params.ConnectorName
	}

	alarmConfig := p.alarmConfigProvider.Get()
	alarm := p.newAlarm(params, entity, now, alarmConfig)
	stateStep := types.NewAlarmStep(types.AlarmStepStateIncrease, params.Timestamp, author,
		params.Output, params.User, params.Role, params.Initiator)
	stateStep.Value = *params.State
	statusStep := types.NewAlarmStep(types.AlarmStepStatusIncrease, params.Timestamp, author,
		params.Output, params.User, params.Role, params.Initiator)
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

		pbhOutput := fmt.Sprintf(
			"Pbehavior %s. Type: %s. Reason: %s.",
			pbehaviorInfo.Name,
			pbehaviorInfo.TypeName,
			pbehaviorInfo.ReasonName,
		)
		newStep := types.NewAlarmStep(types.AlarmStepPbhEnter, *pbehaviorInfo.Timestamp, canopsis.DefaultEventAuthor,
			pbhOutput, "", "", types.InitiatorSystem)
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

	alarm.InternalTags = p.internalTagAlarmMatcher.Match(entity, alarm)
	alarm.InternalTagsUpdated = types.NewMicroTime()
	alarm.Tags = append(alarm.Tags, alarm.InternalTags...)
	alarm.Healthcheck = event.Healthcheck
	_, err = p.alarmCollection.InsertOne(ctx, alarm)
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
	newState := *params.State
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
	if entity.Type != types.EntityTypeService {
		author = strings.Replace(entity.Connector, "/", ".", 1)
	} else {
		author = params.Connector + "." + params.ConnectorName
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
		stateStep = types.NewAlarmStep(types.AlarmStepStateIncrease, params.Timestamp, author,
			params.Output, params.User, params.Role, params.Initiator)
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

	newStatus := p.alarmStatusService.ComputeStatus(alarm, entity)
	if newStatus == previousStatus {
		if stateStep.Type != "" {
			match["$expr"] = bson.M{"$lt": bson.A{bson.M{"$size": "$v.steps"}, types.AlarmStepsHardLimit}}
			push["v.steps"] = stateStep
			inc["v.state_changes_since_status_update"] = 1
		}
	} else {
		if len(alarm.Value.Steps) <= types.AlarmStepsHardLimit {
			match["$expr"] = bson.M{"$lt": bson.A{bson.M{"$size": "$v.steps"}, types.AlarmStepsHardLimit}}
		}
		statusStep := types.NewAlarmStep(types.AlarmStepStatusIncrease, params.Timestamp, author,
			params.Output, params.User, params.Role, params.Initiator)
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
	if newAlarm.CropSteps() {
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
		unset := bson.M{"idle_since": ""}
		alarmLastUpdateRule := fmt.Sprintf("%s_%s", idlerule.RuleTypeAlarm, idlerule.RuleAlarmConditionLastUpdate)
		if entity.LastIdleRuleApply != alarmLastUpdateRule ||
			entity.LastIdleRuleApply == alarmLastUpdateRule && alarmChange.Type != types.AlarmChangeTypeNone {
			unset["last_idle_rule_apply"] = ""
		}

		result.Entity, err = updateEntityByID(ctx, entity.ID, bson.M{"$unset": unset}, p.entityCollection)
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
	timestamp types.CpsTime,
	alarmConfig config.AlarmConfig,
) types.Alarm {
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

	switch entity.Type {
	case types.EntityTypeResource:
		alarm.Value.Resource = entity.Name
		alarm.Value.Component = entity.Component
		alarm.Value.Connector = params.Connector
		alarm.Value.ConnectorName = params.ConnectorName
	case types.EntityTypeComponent, types.EntityTypeService:
		alarm.Value.Component = entity.ID
		alarm.Value.Connector = params.Connector
		alarm.Value.ConnectorName = params.ConnectorName
	case types.EntityTypeConnector:
		alarm.Value.Connector, alarm.Value.ConnectorName, _ = strings.Cut(entity.ID, "/")
	}

	return alarm
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
	updatedServiceStates map[string]statecounters.UpdatedServicesInfo,
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
		err := p.stateCountersService.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			p.logger.Err(err).Msg("failed to update service state")
		}
	}

	p.sendEventStatistics(ctx, event)

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

func resolvePbehaviorInfo(ctx context.Context, entity types.Entity, now types.CpsTime, pbhTypeResolver pbehavior.EntityTypeResolver) (types.PbehaviorInfo, error) {
	result, err := pbhTypeResolver.Resolve(ctx, entity, now.Time)
	if err != nil {
		return types.PbehaviorInfo{}, err
	}

	return pbehavior.NewPBehaviorInfo(now, result), nil
}

func sendRemediationEvent(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
) error {
	if remediationRpcClient == nil {
		return nil
	}

	switch result.AlarmChange.Type {
	case types.AlarmChangeTypeNone:
		if result.AlarmChange.EventsCount < types.MinimalEventsCountThreshold {
			return nil
		}
	case
		types.AlarmChangeTypeCreate,
		types.AlarmChangeTypeCreateAndPbhEnter,
		types.AlarmChangeTypeStateIncrease,
		types.AlarmChangeTypeStateDecrease,
		types.AlarmChangeTypeChangeState,
		types.AlarmChangeTypeUnsnooze,
		types.AlarmChangeTypeActivate,
		types.AlarmChangeTypePbhEnter,
		types.AlarmChangeTypePbhLeave,
		types.AlarmChangeTypePbhLeaveAndEnter,
		types.AlarmChangeTypeResolve:
	default:
		return nil
	}

	entity := event.Entity
	if result.Entity.ID != "" {
		entity = &result.Entity
	}

	body, err := encoder.Encode(types.RPCRemediationEvent{
		Alarm:       &result.Alarm,
		Entity:      entity,
		AlarmChange: result.AlarmChange,
	})
	if err != nil {
		return fmt.Errorf("cannot encode remediation event: %w", err)
	}

	err = remediationRpcClient.Call(ctx, engine.RPCMessage{
		CorrelationID: result.Alarm.ID,
		Body:          body,
	})
	if err != nil {
		return fmt.Errorf("cannot send rpc call to remediation: %w", err)
	}

	return nil
}

func updatePbhLastAlarmDate(ctx context.Context, result Result, pbehaviorCollection mongo.DbCollection) error {
	switch result.AlarmChange.Type {
	case types.AlarmChangeTypeCreateAndPbhEnter,
		types.AlarmChangeTypePbhEnter,
		types.AlarmChangeTypePbhLeaveAndEnter:
	default:
		return nil
	}

	pbhId := ""
	var lastAlarmDate *types.CpsTime
	if result.Alarm.ID == "" {
		pbhId = result.Entity.PbehaviorInfo.ID
		lastAlarmDate = result.Entity.PbehaviorInfo.Timestamp
	} else {
		pbhId = result.Alarm.Value.PbehaviorInfo.ID
		lastAlarmDate = result.Alarm.Value.PbehaviorInfo.Timestamp
	}

	_, err := pbehaviorCollection.UpdateOne(ctx,
		bson.M{"_id": pbhId},
		bson.M{"$set": bson.M{"last_alarm_date": lastAlarmDate}},
	)
	if err != nil {
		return fmt.Errorf("cannot update pbehavior: %w", err)
	}

	return nil
}

func isInstructionMatched(event rpc.AxeEvent, result Result, autoInstructionMatcher AutoInstructionMatcher, logger zerolog.Logger) bool {
	triggers := result.AlarmChange.GetTriggers()
	if len(triggers) == 0 {
		return false
	}

	entity := *event.Entity
	if result.Entity.ID != "" {
		entity = result.Entity
	}

	matched, err := autoInstructionMatcher.Match(triggers, types.AlarmWithEntity{Alarm: result.Alarm, Entity: entity})
	if err != nil {
		logger.Err(err).Str("alarm", result.Alarm.ID).Msg("cannot match auto instructions")
		return false
	}

	return matched
}

func updateEntityByID(ctx context.Context, entityID string, update bson.M, entityCollection mongo.DbCollection) (types.Entity, error) {
	newEntity := types.Entity{}
	err := entityCollection.FindOneAndUpdate(ctx, bson.M{"_id": entityID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).
		Decode(&newEntity)
	if err != nil {
		return newEntity, fmt.Errorf("cannot update entity: %w", err)
	}

	return newEntity, nil
}

func updateEntity(ctx context.Context, match, update bson.M, entityCollection mongo.DbCollection) (types.Entity, error) {
	newEntity := types.Entity{}
	err := entityCollection.FindOneAndUpdate(ctx, match, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).
		Decode(&newEntity)
	if err != nil {
		return newEntity, fmt.Errorf("cannot update entity: %w", err)
	}

	return newEntity, nil
}
