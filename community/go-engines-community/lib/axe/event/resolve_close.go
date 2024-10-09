package event

import (
	"context"
	"errors"
	"fmt"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice/statecounters"
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
	stateCountersService statecounters.StateCountersService,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &resolveCloseProcessor{
		dbClient:                dbClient,
		alarmCollection:         dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:        dbClient.Collection(mongo.EntityMongoCollection),
		resolvedAlarmCollection: dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		metaAlarmRuleCollection: dbClient.Collection(mongo.MetaAlarmRulesMongoCollection),
		stateCountersService:    stateCountersService,
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		metaAlarmStatesService:  metaAlarmStatesService,
		metricsSender:           metricsSender,
		remediationRpcClient:    remediationRpcClient,
		encoder:                 encoder,
		logger:                  logger,
	}
}

type resolveCloseProcessor struct {
	dbClient                mongo.DbClient
	alarmCollection         mongo.DbCollection
	entityCollection        mongo.DbCollection
	resolvedAlarmCollection mongo.DbCollection
	metaAlarmRuleCollection mongo.DbCollection
	stateCountersService    statecounters.StateCountersService
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	metaAlarmStatesService  correlation.MetaAlarmStateService
	metricsSender           metrics.Sender
	remediationRpcClient    engine.RPCClient
	encoder                 encoding.Encoder
	logger                  zerolog.Logger
}

func (p *resolveCloseProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	match := getOpenAlarmMatch(event)
	match["v.state.val"] = types.AlarmStateOK
	result, updatedServiceStates, notAckedMetricType, err := processResolve(ctx, match, event, p.stateCountersService, p.metaAlarmStatesService, p.dbClient, p.alarmCollection, p.entityCollection, p.resolvedAlarmCollection, p.metaAlarmRuleCollection)
	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	go postProcessResolve(context.Background(), event, result, updatedServiceStates, notAckedMetricType, p.stateCountersService, p.metaAlarmEventProcessor, p.metricsSender, p.remediationRpcClient, p.encoder, p.logger)

	return result, nil
}

func processResolve(
	ctx context.Context,
	match bson.M,
	event rpc.AxeEvent,
	stateCountersService statecounters.StateCountersService,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	dbClient mongo.DbClient,
	alarmCollection, entityCollection, resolvedCollection, metaAlarmRuleCollection mongo.DbCollection,
) (Result, map[string]statecounters.UpdatedServicesInfo, string, error) {
	result := Result{}
	var updatedServiceStates map[string]statecounters.UpdatedServicesInfo
	notAckedMetricType := ""
	err := dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil
		notAckedMetricType = ""

		beforeAlarm, err := updateAlarmToResolve(ctx, alarmCollection, match)
		if err != nil || beforeAlarm.ID == "" {
			return err
		}

		if beforeAlarm.NotAckedMetricSendTime != nil {
			notAckedMetricType = beforeAlarm.NotAckedMetricType
		}

		entity, err := updateEntityOfResolvedAlarm(ctx, entityCollection, event.Entity.ID)
		if err != nil || entity.ID == "" {
			return err
		}

		alarm, err := copyAlarmToResolvedCollection(ctx, alarmCollection, resolvedCollection, beforeAlarm.ID)
		if err != nil || alarm.ID == "" {
			return err
		}

		alarmChange := types.NewAlarmChange()
		alarmChange.Type = types.AlarmChangeTypeResolve
		result.Forward = true
		result.Alarm = alarm
		result.Entity = entity
		result.AlarmChange = alarmChange

		updatedServiceStates, err = stateCountersService.UpdateServiceCounters(ctx, entity, &result.Alarm, result.AlarmChange)
		if err != nil {
		    return err
		}

		if !result.Alarm.IsMetaAlarm() {
			return nil
		}

		var rule correlation.Rule
		err = metaAlarmRuleCollection.FindOne(ctx, bson.M{"_id": result.Alarm.Value.Meta}).Decode(&rule)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return fmt.Errorf("meta alarm rule %s not found", result.Alarm.Value.Meta)
			}

			return fmt.Errorf("cannot fetch meta alarm rule: %w", err)
		}

		return RemoveMetaAlarmState(ctx, result.Alarm, rule, metaAlarmStatesService)
	})
	if err != nil || result.Alarm.ID == "" {
		return result, nil, "", err
	}

	return result, updatedServiceStates, notAckedMetricType, nil
}

func updateAlarmToResolve(ctx context.Context, alarmCollection mongo.DbCollection, match bson.M) (types.Alarm, error) {
	update := getResolveAlarmUpdate(types.NewCpsTime())
	beforeAlarm := types.Alarm{}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.Before).
		SetProjection(bson.M{
			"not_acked_metric_type":      1,
			"not_acked_metric_send_time": 1,
		})
	err := alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&beforeAlarm)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return beforeAlarm, err
	}

	return beforeAlarm, nil
}

func copyAlarmToResolvedCollection(
	ctx context.Context,
	alarmCollection, resolvedCollection mongo.DbCollection,
	alarmID string,
) (types.Alarm, error) {
	// extend alarm struct with bookmarks to copy user's bookmarks to a resolved alarm document
	var alarm struct {
		types.Alarm `bson:"inline"`
		Bookmarks   []string `bson:"bookmarks"`
	}

	err := alarmCollection.FindOne(ctx, bson.M{"_id": alarmID}).Decode(&alarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return alarm.Alarm, nil
		}

		return alarm.Alarm, err
	}

	_, err = resolvedCollection.UpdateOne(
		ctx,
		bson.M{"_id": alarm.ID},
		bson.M{"$set": alarm},
		options.Update().SetUpsert(true),
	)

	return alarm.Alarm, err
}

func updateEntityOfResolvedAlarm(ctx context.Context, entityCollection mongo.DbCollection, entityID string) (types.Entity, error) {
	entity := types.Entity{}
	entityUpdate := getResolveEntityUpdate()
	err := entityCollection.FindOneAndUpdate(ctx, bson.M{"_id": entityID}, entityUpdate,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&entity)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return entity, err
	}

	return entity, nil
}

func postProcessResolve(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	updatedServiceStates map[string]statecounters.UpdatedServicesInfo,
	notAckedMetricType string,
	stateCountersService statecounters.StateCountersService,
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
		err := stateCountersService.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			logger.Err(err).Msg("failed to update service state")
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

func getResolveAlarmUpdate(t types.CpsTime) []bson.M {
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
						{"$gt": bson.A{"$v.inactive_start", 0}},
						{"$or": []bson.M{
							{"$ne": bson.A{
								bson.M{"$cond": bson.M{
									"if":   "$v.snooze",
									"then": "$v.snooze",
									"else": nil,
								}},
								nil,
							}},
							{"$not": bson.M{"$in": bson.A{
								bson.M{"$cond": bson.M{
									"if":   "$v.pbehavior_info",
									"then": "$v.pbehavior_info.canonical_type",
									"else": nil,
								}},
								bson.A{nil, "", pbehavior.TypeActive},
							}}},
							{"$eq": bson.A{"$auto_instruction_in_progress", true}},
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
					"if": bson.M{"$ne": bson.A{
						bson.M{"$cond": bson.M{
							"if":   "$v.snooze",
							"then": "$v.snooze",
							"else": nil,
						}},
						nil,
					}},
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
					"if": bson.M{"$not": bson.M{"$in": bson.A{
						bson.M{"$cond": bson.M{
							"if":   "$v.pbehavior_info",
							"then": "$v.pbehavior_info.canonical_type",
							"else": nil,
						}},
						bson.A{nil, "", pbehavior.TypeActive},
					}}},
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
