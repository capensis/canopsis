package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewChangeStateProcessor(
	client mongo.DbClient,
	userInterfaceConfigProvider config.UserInterfaceConfigProvider,
	alarmStatusService alarmstatus.Service,
	autoInstructionMatcher AutoInstructionMatcher,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &changeStateProcessor{
		client:                          client,
		alarmCollection:                 client.Collection(mongo.AlarmMongoCollection),
		entityCollection:                client.Collection(mongo.EntityMongoCollection),
		userInterfaceConfigProvider:     userInterfaceConfigProvider,
		alarmStatusService:              alarmStatusService,
		autoInstructionMatcher:          autoInstructionMatcher,
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		metaAlarmPostProcessor:          metaAlarmPostProcessor,
		metricsSender:                   metricsSender,
		remediationRpcClient:            remediationRpcClient,
		encoder:                         encoder,
		logger:                          logger,
	}
}

type changeStateProcessor struct {
	client                          mongo.DbClient
	alarmCollection                 mongo.DbCollection
	entityCollection                mongo.DbCollection
	userInterfaceConfigProvider     config.UserInterfaceConfigProvider
	alarmStatusService              alarmstatus.Service
	autoInstructionMatcher          AutoInstructionMatcher
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	eventsSender                    entitycounters.EventsSender
	metaAlarmPostProcessor          MetaAlarmPostProcessor
	metricsSender                   metrics.Sender
	remediationRpcClient            engine.RPCClient
	encoder                         encoding.Encoder
	logger                          zerolog.Logger
}

func (p *changeStateProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled || event.Parameters.State == nil {
		return Result{}, nil
	}

	if *event.Parameters.State == types.AlarmStateOK && !p.userInterfaceConfigProvider.Get().IsAllowChangeSeverityToInfo {
		return result, errors.New("cannot change to ok state")
	}

	entity := *event.Entity
	match := getOpenAlarmMatchWithStepsLimit(event)
	match["$and"] = []bson.M{
		{"v.state.val": bson.M{"$ne": types.AlarmStateOK}},
		{"$or": []bson.M{
			{"v.state.val": bson.M{"$ne": event.Parameters.State}},
			{"v.change_state": nil},
		}},
	}
	matchUpdate := getOpenAlarmMatch(event)
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo

	var componentStateChanged bool
	var newComponentState int

	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil

		alarm := types.Alarm{}
		err := p.alarmCollection.FindOne(ctx, match).Decode(&alarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		newStepState := NewAlarmStep(types.AlarmStepChangeState, event.Parameters, !alarm.Value.PbehaviorInfo.IsDefaultActive())
		newStepState.Value = *event.Parameters.State
		alarmChange := types.NewAlarmChange()
		alarmChange.PreviousState = alarm.Value.State.Value
		alarmChange.PreviousStateChange = alarm.Value.State.Timestamp
		alarmChange.PreviousStatus = alarm.Value.Status.Value
		alarm.Value.ChangeState = &newStepState
		alarm.Value.State = &newStepState
		err = alarm.Value.Steps.Add(newStepState)
		if err != nil {
			return err
		}

		currentStatus := alarm.Value.Status.Value
		newStatus, statusRuleName := p.alarmStatusService.ComputeStatus(alarm, *event.Entity)
		var update bson.M
		if newStatus == currentStatus {
			update = bson.M{
				"$set": bson.M{
					"v.state":        newStepState,
					"v.change_state": newStepState,
				},
				"$push": bson.M{"v.steps": newStepState},
			}
		} else {
			newStepStatus := NewAlarmStep(types.AlarmStepStatusIncrease, event.Parameters, !alarm.Value.PbehaviorInfo.IsDefaultActive())
			newStepStatus.Message = ConcatOutputAndRuleName(event.Parameters.Output, statusRuleName)
			newStepStatus.Value = newStatus
			if alarm.Value.Status.Value > newStatus {
				newStepStatus.Type = types.AlarmStepStatusDecrease
			}

			update = bson.M{
				"$set": bson.M{
					"v.state":                             newStepState,
					"v.change_state":                      newStepState,
					"v.status":                            newStepStatus,
					"v.state_changes_since_status_update": 0,
					"v.last_update_date":                  event.Parameters.Timestamp,
				},
				"$push": bson.M{"v.steps": bson.M{"$each": bson.A{newStepState, newStepStatus}}},
			}
		}

		updatedAlarm := types.Alarm{}
		err = p.alarmCollection.FindOneAndUpdate(ctx, matchUpdate, update, opts).Decode(&updatedAlarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		alarmChange.Type = types.AlarmChangeTypeChangeState
		result.Forward = true
		result.Alarm = updatedAlarm
		result.AlarmChange = alarmChange

		if event.Parameters.IdleRuleApply != "" {
			result.Entity, err = updateEntityByID(ctx, entity.ID, bson.M{"$set": bson.M{
				"last_idle_rule_apply": event.Parameters.IdleRuleApply,
			}}, p.entityCollection)
			if err != nil {
				return err
			}
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
	})
	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	result.IsInstructionMatched = isInstructionMatched(event, result, p.autoInstructionMatcher, p.logger)
	go p.postProcess(context.Background(), event, result, updatedServiceStates, componentStateChanged, newComponentState)

	return result, nil
}

func (p *changeStateProcessor) postProcess(
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
}
