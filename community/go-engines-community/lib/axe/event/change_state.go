package event

import (
	"context"
	"errors"
	"fmt"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice/statecounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewChangeStateProcessor(
	client mongo.DbClient,
	alarmConfigProvider config.AlarmConfigProvider,
	userInterfaceConfigProvider config.UserInterfaceConfigProvider,
	alarmStatusService alarmstatus.Service,
	autoInstructionMatcher AutoInstructionMatcher,
	stateCountersService statecounters.StateCountersService,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &changeStateProcessor{
		client:                      client,
		alarmCollection:             client.Collection(mongo.AlarmMongoCollection),
		entityCollection:            client.Collection(mongo.EntityMongoCollection),
		alarmConfigProvider:         alarmConfigProvider,
		userInterfaceConfigProvider: userInterfaceConfigProvider,
		alarmStatusService:          alarmStatusService,
		autoInstructionMatcher:      autoInstructionMatcher,
		stateCountersService:        stateCountersService,
		metaAlarmEventProcessor:     metaAlarmEventProcessor,
		metricsSender:               metricsSender,
		remediationRpcClient:        remediationRpcClient,
		encoder:                     encoder,
		logger:                      logger,
	}
}

type changeStateProcessor struct {
	client                      mongo.DbClient
	alarmCollection             mongo.DbCollection
	entityCollection            mongo.DbCollection
	alarmConfigProvider         config.AlarmConfigProvider
	userInterfaceConfigProvider config.UserInterfaceConfigProvider
	alarmStatusService          alarmstatus.Service
	autoInstructionMatcher      AutoInstructionMatcher
	stateCountersService        statecounters.StateCountersService
	metaAlarmEventProcessor     libalarm.MetaAlarmEventProcessor
	metricsSender               metrics.Sender
	remediationRpcClient        engine.RPCClient
	encoder                     encoding.Encoder
	logger                      zerolog.Logger
}

func (p *changeStateProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled || event.Parameters.State == nil {
		return Result{}, nil
	}

	if *event.Parameters.State == types.AlarmStateOK && !p.userInterfaceConfigProvider.Get().IsAllowChangeSeverityToInfo {
		return result, fmt.Errorf("cannot change to ok state")
	}

	entity := *event.Entity
	var match bson.M
	if *event.Parameters.State == types.AlarmStateOK {
		match = getOpenAlarmMatch(event)
	} else {
		match = getOpenAlarmMatchWithStepsLimit(event)
	}
	match["$and"] = []bson.M{
		{"v.state.val": bson.M{"$ne": types.AlarmStateOK}},
		{"$or": []bson.M{
			{"v.state.val": bson.M{"$ne": event.Parameters.State}},
			{"v.change_state": nil},
		}},
	}
	matchUpdate := getOpenAlarmMatch(event)
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	conf := p.alarmConfigProvider.Get()
	output := utils.TruncateString(event.Parameters.Output, conf.OutputLength)
	newStepState := types.NewAlarmStep(types.AlarmStepChangeState, event.Parameters.Timestamp, event.Parameters.Author, output,
		event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator)
	newStepState.Value = *event.Parameters.State
	var updatedServiceStates map[string]statecounters.UpdatedServicesInfo

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
		newStatus := p.alarmStatusService.ComputeStatus(alarm, *event.Entity)
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
			newStepStatus := types.NewAlarmStep(types.AlarmStepStatusIncrease, event.Parameters.Timestamp, event.Parameters.Author, output,
				event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator)
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

		updatedServiceStates, err = p.stateCountersService.UpdateServiceCounters(ctx, entity, &result.Alarm, result.AlarmChange)
		return err
	})
	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	result.IsInstructionMatched = isInstructionMatched(event, result, p.autoInstructionMatcher, p.logger)
	go p.postProcess(context.Background(), event, result, updatedServiceStates)

	return result, nil
}

func (p *changeStateProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	updatedServiceStates map[string]statecounters.UpdatedServicesInfo,
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
		err := p.stateCountersService.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			p.logger.Err(err).Msg("failed to update service state")
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
}
