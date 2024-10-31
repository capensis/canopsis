package event

import (
	"context"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMetaAlarmProcessor(
	autoInstructionMatcher AutoInstructionMatcher,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	dbClient mongo.DbClient,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	adapter libalarm.Adapter,
	ruleAdapter correlation.RulesAdapter,
	pbhTypeResolver pbehavior.EntityTypeResolver,
	alarmStatusService alarmstatus.Service,
	alarmConfigProvider config.AlarmConfigProvider,
	templateExecutor template.Executor,
	encoder encoding.Encoder,
	eventGenerator libevent.Generator,
	amqpPublisher libamqp.Publisher,
	logger zerolog.Logger,
) Processor {
	return &metaAlarmProcessor{
		autoInstructionMatcher:    autoInstructionMatcher,
		metricsSender:             metricsSender,
		remediationRpcClient:      remediationRpcClient,
		dbClient:                  dbClient,
		alarmCollection:           dbClient.Collection(mongo.AlarmMongoCollection),
		metaAlarmStatesCollection: dbClient.Collection(mongo.MetaAlarmStatesCollection),
		entityCollection:          dbClient.Collection(mongo.EntityMongoCollection),
		metaAlarmStatesService:    metaAlarmStatesService,
		adapter:                   adapter,
		ruleAdapter:               ruleAdapter,
		pbhTypeResolver:           pbhTypeResolver,
		alarmStatusService:        alarmStatusService,
		alarmConfigProvider:       alarmConfigProvider,
		templateExecutor:          templateExecutor,
		encoder:                   encoder,
		eventGenerator:            eventGenerator,
		amqpPublisher:             amqpPublisher,
		logger:                    logger,
	}
}

type metaAlarmProcessor struct {
	dbClient                  mongo.DbClient
	alarmCollection           mongo.DbCollection
	metaAlarmStatesCollection mongo.DbCollection
	entityCollection          mongo.DbCollection
	metaAlarmStatesService    correlation.MetaAlarmStateService
	adapter                   libalarm.Adapter
	ruleAdapter               correlation.RulesAdapter
	pbhTypeResolver           pbehavior.EntityTypeResolver
	alarmStatusService        alarmstatus.Service
	alarmConfigProvider       config.AlarmConfigProvider
	autoInstructionMatcher    AutoInstructionMatcher
	metricsSender             metrics.Sender
	remediationRpcClient      engine.RPCClient
	encoder                   encoding.Encoder
	eventGenerator            libevent.Generator
	amqpPublisher             libamqp.Publisher
	templateExecutor          template.Executor
	logger                    zerolog.Logger
}

func (p *metaAlarmProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	if event.Entity == nil {
		return Result{}, nil
	}

	result, updatedChildrenAlarms, err := p.createMetaAlarm(ctx, event)
	if err != nil {
		return result, err
	}

	for _, child := range updatedChildrenAlarms {
		p.metricsSender.SendCorrelation(event.Parameters.Timestamp.Time, child)
	}

	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *metaAlarmProcessor) createMetaAlarm(ctx context.Context, event rpc.AxeEvent) (Result, []types.Alarm, error) {
	if event.Entity == nil {
		return Result{}, nil, nil
	}

	var updatedChildrenAlarms []types.Alarm
	var metaAlarm types.Alarm
	var alarmChange types.AlarmChange
	var result Result
	var activateChildEvents []types.Event

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedChildrenAlarms = updatedChildrenAlarms[:0]
		activateChildEvents = activateChildEvents[:0]
		metaAlarm = types.Alarm{}
		alarmChange = types.NewAlarmChange()
		result = Result{Forward: true}

		rule, err := p.ruleAdapter.GetRule(ctx, event.Parameters.MetaAlarmRuleID)
		if err != nil {
			return fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", event.Parameters.MetaAlarmRuleID, err)
		} else if rule.ID == "" {
			return fmt.Errorf("meta alarm rule id=%q not found", event.Parameters.MetaAlarmRuleID)
		}

		entity := types.Entity{}
		err = p.entityCollection.FindOne(ctx, bson.M{"_id": event.Entity.ID}).Decode(&entity)
		if err != nil {
			return err
		}

		metaAlarm = p.newMetaAlarm(event.Parameters, entity, p.alarmConfigProvider.Get())
		metaAlarm.Value.Meta = event.Parameters.MetaAlarmRuleID
		metaAlarm.Value.MetaValuePath = event.Parameters.MetaAlarmValuePath
		metaAlarm.Value.LastEventDate = datetime.CpsTime{} // should be empty

		if event.Parameters.DisplayName != "" {
			metaAlarm.Value.DisplayName = event.Parameters.DisplayName
		}

		stateID := rule.GetStateID(event.Parameters.MetaAlarmValuePath)
		var childEntityIDs []string
		var archived bool

		if rule.IsManual() {
			childEntityIDs = event.Parameters.MetaAlarmChildren
		} else {
			metaAlarmState, err := p.metaAlarmStatesService.GetMetaAlarmState(ctx, stateID)
			if err != nil {
				return err
			}

			if metaAlarmState.GetEntityID() != entity.ID {
				// try to get archived state
				metaAlarmState, err = p.metaAlarmStatesService.GetMetaAlarmState(ctx, stateID+"-"+entity.ID)
				if err != nil {
					return err
				}

				if metaAlarmState.ID == "" {
					return fmt.Errorf("meta alarm state for rule id=%q and meta alarm name=%q not found",
						event.Parameters.MetaAlarmRuleID, entity.Name)
				}

				archived = true
			}

			childEntityIDs = metaAlarmState.ChildrenEntityIDs
		}

		var lastChild types.AlarmWithEntity
		worstState := types.CpsNumber(types.AlarmStateMinor)
		eventsCount := types.CpsNumber(0)
		var writeModels []mongodriver.WriteModel
		var childAlarms []types.AlarmWithEntity
		if len(childEntityIDs) > 0 {
			childAlarms, err = getAlarmsWithEntityByMatch(ctx, p.alarmCollection, bson.M{
				"d":          bson.M{"$in": childEntityIDs},
				"v.resolved": nil,
			})
			if err != nil {
				return fmt.Errorf("cannot fetch children alarms: %w", err)
			}

			writeModels = make([]mongodriver.WriteModel, 0, len(childAlarms))
			for _, childAlarm := range childAlarms {
				if childAlarm.Alarm.Value.State != nil {
					childState := childAlarm.Alarm.Value.State.Value
					if childState > worstState {
						worstState = childState
					}
				}

				if !childAlarm.Alarm.HasParentByEID(metaAlarm.EntityID) {
					metaAlarm.AddChild(childAlarm.Alarm.EntityID)
					childAlarm.Alarm.AddParent(metaAlarm.EntityID)
					newStep := NewAlarmStep(types.AlarmStepMetaAlarmAttach, event.Parameters, !childAlarm.Alarm.Value.PbehaviorInfo.IsDefaultActive())
					newStep.Message = getMetaAlarmChildStepMsg(rule, metaAlarm, event)
					err := childAlarm.Alarm.Value.Steps.Add(newStep)
					if err != nil {
						return err
					}

					if childAlarm.Alarm.InactiveDelayMetaAlarmInProgress {
						activateChildEvent, err := p.eventGenerator.Generate(childAlarm.Entity)
						if err != nil {
							return err
						}

						activateChildEvent.EventType = types.EventTypeMetaAlarmChildActivate
						activateChildEvent.Timestamp = datetime.NewCpsTime()
						activateChildEvents = append(activateChildEvents, activateChildEvent)
					}

					writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
						SetFilter(bson.M{
							"_id":       childAlarm.Alarm.ID,
							"v.parents": bson.M{"$ne": metaAlarm.EntityID},
						}).
						SetUpdate(bson.M{
							"$addToSet": bson.M{"v.parents": metaAlarm.EntityID},
							"$push":     bson.M{"v.steps": newStep},
						}))

					updatedChildrenAlarms = append(updatedChildrenAlarms, childAlarm.Alarm)
					eventsCount += childAlarm.Alarm.Value.EventsCount
					lastChild = childAlarm
					if metaAlarm.Value.LastEventDate.Before(childAlarm.Alarm.Value.LastEventDate) {
						metaAlarm.Value.LastEventDate = childAlarm.Alarm.Value.LastEventDate
					}
				}
			}
		}

		output := ""
		entityUpdate := bson.M{}
		var ruleInfos []types.CorrelationRuleInfo
		if rule.IsManual() {
			output = event.Parameters.Output
			ruleInfos = event.Parameters.MetaAlarmInfos
			if event.Parameters.MetaAlarmTags != nil {
				metaAlarm.CopyTagsFromChildren = event.Parameters.MetaAlarmTags.CopyFromChildren
				metaAlarm.FilterChildrenTagsByLabel = event.Parameters.MetaAlarmTags.FilterByLabel
			}
		} else {
			output, err = executeMetaAlarmOutputTpl(p.templateExecutor, correlation.EventExtraInfosMeta{
				Rule:      rule,
				Count:     int64(len(updatedChildrenAlarms)),
				LastChild: lastChild,
				Children:  lastChild,
			})
			if err != nil {
				return err
			}

			ruleInfos = rule.Infos
			metaAlarm.CopyTagsFromChildren = rule.Tags.CopyFromChildren
			metaAlarm.FilterChildrenTagsByLabel = rule.Tags.FilterByLabel
		}

		if metaAlarm.CopyTagsFromChildren {
			metaAlarm.ExternalTags = getMetaAlarmExternalTags(metaAlarm.FilterChildrenTagsByLabel, childAlarms, nil)
			metaAlarm.Tags = metaAlarm.ExternalTags
		}

		copyInfoNames := make([]string, 0)
		for _, info := range ruleInfos {
			if info.CopyFromChildren {
				copyInfoNames = append(copyInfoNames, info.Name)
				continue
			}

			entityUpdate["infos."+info.Name] = types.Info{
				Name:        info.Name,
				Description: info.Description,
				Value:       info.Value,
			}
		}

		if len(copyInfoNames) > 0 {
			metaAlarm.EntityInfosFromChildren = copyInfoNames
			copiedEntityInfos := getMetaAlarmEntityInfos(copyInfoNames, childAlarms, nil)
			for k, v := range copiedEntityInfos {
				entityUpdate["infos."+k] = v
			}
		}

		metaAlarm.Value.Output = output
		_, _, err = updateMetaAlarmState(&metaAlarm, entity, event.Parameters.Timestamp, worstState,
			output, p.alarmStatusService)
		if err != nil {
			return err
		}

		metaAlarm.Value.EventsCount = eventsCount

		pbehaviorInfo, err := resolvePbehaviorInfo(ctx, entity, metaAlarm.Time, p.pbhTypeResolver)
		if err != nil {
			return fmt.Errorf("failed to resolve pbehavior info for metaalarm: %w", err)
		}

		if pbehaviorInfo.IsDefaultActive() {
			alarmChange.Type = types.AlarmChangeTypeCreate
			metaAlarm.NotAckedSince = &metaAlarm.Time
		} else {
			if pbehaviorInfo.IsActive() {
				metaAlarm.NotAckedSince = &metaAlarm.Time
			} else {
				metaAlarm.Value.InactiveStart = &metaAlarm.Time
			}

			newStep := types.NewPbhAlarmStep(types.AlarmStepPbhEnter, *pbehaviorInfo.Timestamp, pbehaviorInfo.Author,
				pbehaviorInfo.GetStepMessage(), "", "", types.InitiatorSystem, pbehaviorInfo.CanonicalType,
				pbehaviorInfo.IconName, pbehaviorInfo.Color)
			metaAlarm.Value.PbehaviorInfo = pbehaviorInfo
			err := metaAlarm.Value.Steps.Add(newStep)
			if err != nil {
				return fmt.Errorf("cannot add pbhenter step on metaalarm create: %w", err)
			}

			alarmChange.Type = types.AlarmChangeTypeCreateAndPbhEnter
			entityUpdate["pbehavior_info"] = metaAlarm.Value.PbehaviorInfo
			entityUpdate["last_pbehavior_date"] = metaAlarm.Value.PbehaviorInfo.Timestamp
		}

		if len(entityUpdate) > 0 {
			entity = types.Entity{}
			err := p.entityCollection.FindOneAndUpdate(ctx, bson.M{"_id": metaAlarm.EntityID},
				bson.M{"$set": entityUpdate}, options.FindOneAndUpdate().SetReturnDocument(options.After)).
				Decode(&entity)
			if err != nil {
				return fmt.Errorf("cannot update meta alarm entity: %w", err)
			}

			result.Entity = entity
		}

		result.IsInstructionMatched, err = p.autoInstructionMatcher.Match(alarmChange.GetTriggers(), types.AlarmWithEntity{Alarm: metaAlarm, Entity: entity})
		if err != nil {
			return err
		}

		if p.alarmConfigProvider.Get().ActivateAlarmAfterAutoRemediation {
			metaAlarm.InactiveAutoInstructionInProgress = result.IsInstructionMatched
		}

		writeModels = append(writeModels, mongodriver.NewInsertOneModel().SetDocument(types.AlarmWithEntityField{
			Alarm:  metaAlarm,
			Entity: entity,
		}))
		_, err = p.alarmCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}

		if !rule.IsManual() && !archived {
			ok, err := p.metaAlarmStatesService.SwitchStateToCreated(ctx, stateID)
			if err != nil || !ok {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return Result{}, nil, err
	}

	for _, e := range activateChildEvents {
		err = p.sendToFifo(ctx, e)
		if err != nil {
			return Result{}, nil, err
		}
	}

	result.Alarm = metaAlarm
	result.AlarmChange = alarmChange

	return result, updatedChildrenAlarms, nil
}

func (p *metaAlarmProcessor) newMetaAlarm(
	params rpc.AxeParameters,
	entity types.Entity,
	alarmConfig config.AlarmConfig,
) types.Alarm {
	now := datetime.NewCpsTime()
	alarm := types.Alarm{
		EntityID: entity.ID,
		ID:       utils.NewID(),
		Time:     now,
		Value: types.AlarmValue{
			Connector:         canopsis.DefaultSystemAlarmConnector,
			ConnectorName:     canopsis.DefaultSystemAlarmConnector,
			Component:         entity.Component,
			Resource:          entity.Name,
			CreationDate:      now,
			DisplayName:       types.GenDisplayName(alarmConfig.DisplayNameScheme),
			InitialOutput:     params.Output,
			Output:            params.Output,
			InitialLongOutput: params.LongOutput,
			LongOutput:        params.LongOutput,
			LongOutputHistory: []string{params.LongOutput},
			LastUpdateDate:    params.Timestamp,
			LastEventDate:     now,
			Parents:           []string{},
			Children:          []string{},
			UnlinkedParents:   []string{},
			Infos:             map[string]map[string]interface{}{},
			RuleVersion:       map[string]string{},
		},
	}

	return alarm
}

func (p *metaAlarmProcessor) sendToFifo(ctx context.Context, event types.Event) error {
	body, err := p.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("cannot encode event: %w", err)
	}

	err = p.amqpPublisher.PublishWithContext(
		ctx,
		canopsis.FIFOExchangeName,
		canopsis.FIFOQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot send child event: %w", err)
	}

	return nil
}

func (p *metaAlarmProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
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

	err := sendRemediationEvent(ctx, event, result, p.remediationRpcClient, p.encoder)
	if err != nil {
		p.logger.Err(err).Msg("cannot send event to engine-remediation")
	}
}
