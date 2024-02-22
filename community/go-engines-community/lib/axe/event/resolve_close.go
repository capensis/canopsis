package event

import (
	"context"
	"errors"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
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
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewResolveCloseProcessor(
	dbClient mongo.DbClient,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &resolveCloseProcessor{
		dbClient:                        dbClient,
		alarmCollection:                 dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:                dbClient.Collection(mongo.EntityMongoCollection),
		resolvedAlarmCollection:         dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		metaAlarmEventProcessor:         metaAlarmEventProcessor,
		metricsSender:                   metricsSender,
		remediationRpcClient:            remediationRpcClient,
		eventsSender:                    eventsSender,
		encoder:                         encoder,
		logger:                          logger,
	}
}

type resolveCloseProcessor struct {
	dbClient                        mongo.DbClient
	alarmCollection                 mongo.DbCollection
	entityCollection                mongo.DbCollection
	resolvedAlarmCollection         mongo.DbCollection
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	metaAlarmEventProcessor         libalarm.MetaAlarmEventProcessor
	metricsSender                   metrics.Sender
	remediationRpcClient            engine.RPCClient
	eventsSender                    entitycounters.EventsSender
	encoder                         encoding.Encoder
	logger                          zerolog.Logger
}

func (p *resolveCloseProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	match := getOpenAlarmMatch(event)
	match["v.state.val"] = types.AlarmStateOK
	result, updatedServiceStates, notAckedMetricType, componentStateChanged, newComponentState, err := processResolve(
		ctx,
		match,
		event,
		p.entityServiceCountersCalculator,
		p.componentCountersCalculator,
		p.dbClient,
		p.alarmCollection,
		p.entityCollection,
		p.resolvedAlarmCollection,
	)
	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	go postProcessResolve(
		context.Background(),
		event,
		result,
		updatedServiceStates,
		componentStateChanged,
		newComponentState,
		notAckedMetricType,
		p.eventsSender,
		p.metaAlarmEventProcessor,
		p.metricsSender,
		p.remediationRpcClient,
		p.encoder,
		p.logger,
	)

	return result, nil
}

func processResolve(
	ctx context.Context,
	match bson.M,
	event rpc.AxeEvent,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	dbClient mongo.DbClient,
	alarmCollection, entityCollection, resolvedCollection mongo.DbCollection,
) (Result, map[string]entitycounters.UpdatedServicesInfo, string, bool, int, error) {
	result := Result{}
	update := getResolveAlarmUpdate(datetime.NewCpsTime())
	var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo
	notAckedMetricType := ""

	var componentStateChanged bool
	var newComponentState int

	err := dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil
		notAckedMetricType = ""

		beforeAlarm := types.Alarm{}
		opts := options.FindOneAndUpdate().
			SetReturnDocument(options.Before).
			SetProjection(bson.M{
				"not_acked_metric_type":      1,
				"not_acked_metric_send_time": 1,
			})
		err := alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&beforeAlarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		if beforeAlarm.NotAckedMetricSendTime != nil {
			notAckedMetricType = beforeAlarm.NotAckedMetricType
		}

		// extend alarm struct with bookmarks to copy user's bookmarks to a resolved alarm document
		var alarm struct {
			types.Alarm `bson:"inline"`
			Bookmarks   []string `bson:"bookmarks"`
		}
		err = alarmCollection.FindOne(ctx, bson.M{"_id": beforeAlarm.ID}).Decode(&alarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}
			return err
		}

		entity := types.Entity{}
		entityUpdate := getResolveEntityUpdate()
		err = entityCollection.FindOneAndUpdate(ctx, bson.M{"_id": event.Entity.ID}, entityUpdate,
			options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&entity)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		alarmChange := types.NewAlarmChange()
		alarmChange.Type = types.AlarmChangeTypeResolve
		result.Forward = true
		result.Alarm = alarm.Alarm
		result.Entity = entity
		result.AlarmChange = alarmChange

		_, err = resolvedCollection.UpdateOne(
			ctx,
			bson.M{"_id": alarm.ID},
			bson.M{"$set": alarm},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return err
		}

		updatedServiceStates, componentStateChanged, newComponentState, err = processComponentAndServiceCounters(
			ctx,
			entityServiceCountersCalculator,
			componentCountersCalculator,
			&result.Alarm,
			&entity,
			result.AlarmChange,
		)

		return err
	})
	if err != nil || result.Alarm.ID == "" {
		return result, nil, "", false, 0, err
	}

	return result, updatedServiceStates, notAckedMetricType, componentStateChanged, newComponentState, nil
}

func postProcessResolve(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	updatedServiceStates map[string]entitycounters.UpdatedServicesInfo,
	componentChanged bool,
	newComponentState int,
	notAckedMetricType string,
	eventsSender entitycounters.EventsSender,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) {
	metricsSender.SendEventMetrics(
		result.Alarm,
		*event.Entity,
		result.AlarmChange,
		event.Parameters.Timestamp.Time,
		event.Parameters.Initiator,
		event.Parameters.User,
		event.Parameters.Instruction,
		notAckedMetricType,
	)

	for servID, servInfo := range updatedServiceStates {
		err := eventsSender.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			logger.Err(err).Msg("failed to update service state")
		}
	}

	if componentChanged {
		err := eventsSender.UpdateComponentState(ctx, event.Entity.Component, event.Entity.Connector, newComponentState)
		if err != nil {
			logger.Err(err).Msg("failed to update component state")
		}
	}

	err := metaAlarmEventProcessor.ProcessAxeRpc(ctx, event, rpc.AxeResultEvent{
		Alarm:           &result.Alarm,
		AlarmChangeType: result.AlarmChange.Type,
	})
	if err != nil {
		logger.Err(err).Msg("cannot process meta alarm")
	}

	err = sendRemediationEvent(ctx, event, result, remediationRpcClient, encoder)
	if err != nil {
		logger.Err(err).Msg("cannot send event to engine-remediation")
	}
}

func getResolveAlarmUpdate(t datetime.CpsTime) []bson.M {
	return []bson.M{
		{"$set": bson.M{
			"v.duration": bson.M{"$subtract": bson.A{
				t,
				"$t",
			}},
			"v.inactive_duration": bson.M{"$sum": bson.A{
				"$v.inactive_duration",
				bson.M{"$cond": bson.M{
					"if": bson.M{"$and": []bson.M{
						{"$ne": bson.A{"$v.inactive_start", nil}},
						{"$or": []bson.M{
							{"$ne": bson.A{"$v.snooze", nil}},
							{"$not": bson.M{"$in": bson.A{"$v.pbehavior_info.canonical_type", bson.A{nil, "", pbehavior.TypeActive}}}},
						}},
					}},
					"then": bson.M{"$subtract": bson.A{
						t,
						"$v.inactive_start",
					}},
					"else": 0,
				}},
			}},
		}},
		{"$set": bson.M{
			"v.resolved": t,
			"v.current_state_duration": bson.M{"$subtract": bson.A{
				t,
				"$v.state.t",
			}},
			"v.active_duration": bson.M{"$subtract": bson.A{
				"$v.duration",
				"$v.inactive_duration",
			}},
			"v.snooze_duration": bson.M{"$sum": bson.A{
				"$v.snooze_duration",
				bson.M{"$cond": bson.M{
					"if": bson.M{"$ne": bson.A{"$v.snooze", nil}},
					"then": bson.M{"$subtract": bson.A{
						t,
						"$v.snooze.t",
					}},
					"else": 0,
				}},
			}},
			"v.pbh_inactive_duration": bson.M{"$sum": bson.A{
				"$v.pbh_inactive_duration",
				bson.M{"$cond": bson.M{
					"if": bson.M{"$not": bson.M{"$in": bson.A{"$v.pbehavior_info.canonical_type", bson.A{nil, "", pbehavior.TypeActive}}}},
					"then": bson.M{"$subtract": bson.A{
						t,
						"$v.pbehavior_info.timestamp",
					}},
					"else": 0,
				}},
			}},
		}},
		{"$unset": bson.A{
			"not_acked_metric_type",
			"not_acked_metric_send_time",
			"not_acked_since",
		}},
	}
}

func getResolveEntityUpdate() bson.M {
	return bson.M{"$unset": bson.M{
		"idle_since":           "",
		"last_idle_rule_apply": "",
	}}
}
