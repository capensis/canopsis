package event

import (
	"context"
	"errors"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func newResolveHelper(
	dbClient mongo.DbClient,
	alarmConfigProvider config.AlarmConfigProvider,
	alarmStatusService alarmstatus.Service,
	pbhTypeResolver pbehavior.EntityTypeResolver,
	autoInstructionMatcher AutoInstructionMatcher,
	internalTagAlarmMatcher alarmtag.InternalTagAlarmMatcher,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	eventsSender entitycounters.EventsSender,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	eventGenerator event.Generator,
	amqpPublisher libamqp.Publisher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) *resolveHelper {
	return &resolveHelper{
		dbClient:                          dbClient,
		metaAlarmStatesService:            metaAlarmStatesService,
		metaAlarmPostProcessor:            metaAlarmPostProcessor,
		metricsSender:                     metricsSender,
		remediationRpcClient:              remediationRpcClient,
		alarmCollection:                   dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:                  dbClient.Collection(mongo.EntityMongoCollection),
		resolvedCollection:                dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		metaAlarmRuleCollection:           dbClient.Collection(mongo.MetaAlarmRulesMongoCollection),
		closeDelayJobCollection:           dbClient.Collection(mongo.CloseDelayJobCollection),
		pbehaviorCollection:               dbClient.Collection(mongo.PbehaviorMongoCollection),
		encoder:                           encoder,
		logger:                            logger,
		componentAndServiceCountersHelper: newComponentAndServiceCountersHelper(entityServiceCountersCalculator, componentCountersCalculator, eventsSender, logger),
		upstreamHelper: newUpstreamHelper(
			dbClient,
			alarmConfigProvider,
			alarmStatusService,
			pbhTypeResolver,
			autoInstructionMatcher,
			metaAlarmPostProcessor,
			metricsSender,
			remediationRpcClient,
			internalTagAlarmMatcher,
			entityServiceCountersCalculator,
			componentCountersCalculator,
			eventsSender,
			eventGenerator,
			amqpPublisher,
			encoder,
			logger,
		),
	}
}

type resolveHelper struct {
	dbClient                          mongo.DbClient
	metaAlarmStatesService            correlation.MetaAlarmStateService
	metaAlarmPostProcessor            MetaAlarmPostProcessor
	metricsSender                     metrics.Sender
	remediationRpcClient              engine.RPCClient
	alarmCollection                   mongo.DbCollection
	entityCollection                  mongo.DbCollection
	resolvedCollection                mongo.DbCollection
	metaAlarmRuleCollection           mongo.DbCollection
	closeDelayJobCollection           mongo.DbCollection
	pbehaviorCollection               mongo.DbCollection
	encoder                           encoding.Encoder
	logger                            zerolog.Logger
	componentAndServiceCountersHelper *componentAndServiceCountersHelper
	upstreamHelper                    *upstreamHelper
}

func (h *resolveHelper) Process(ctx context.Context, match bson.M, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	countersRes := componentAndServiceCountersResult{}
	notAckedMetricType := ""
	err := h.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		countersRes = componentAndServiceCountersResult{}
		notAckedMetricType = ""

		beforeAlarm, err := h.UpdateAlarmToResolve(ctx, match, event.Parameters)
		if err != nil || beforeAlarm.ID == "" {
			return err
		}

		if beforeAlarm.NotAckedMetricSendTime != nil {
			notAckedMetricType = beforeAlarm.NotAckedMetricType
		}

		entity, err := h.UpdateEntityOfResolvedAlarm(ctx, event.Entity.ID)
		if err != nil || entity.ID == "" {
			return err
		}

		alarm, err := h.CopyAlarmToResolvedCollection(ctx, beforeAlarm.ID)
		if err != nil || alarm.ID == "" {
			return err
		}

		alarmChange := types.NewAlarmChange()
		alarmChange.Type = types.AlarmChangeTypeResolve
		result.Forward = true
		result.Alarm = alarm
		result.Entity = entity
		result.AlarmChange = alarmChange

		result.IsCountersUpdated, countersRes, err = h.componentAndServiceCountersHelper.Process(
			ctx,
			&result.Alarm,
			&entity,
			result.AlarmChange,
		)
		if err != nil {
			return err
		}

		_, err = h.closeDelayJobCollection.DeleteOne(ctx, bson.M{"_id": alarm.ID})
		if err != nil {
			return fmt.Errorf("failed to delete close_delay job on resolve: %w", err)
		}

		return h.RemoveMetaAlarmStateOnResolve(ctx, result.Alarm)
	})
	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	if result.AlarmChange.Type == types.AlarmChangeTypeResolve {
		go h.PostProcess(context.WithoutCancel(ctx), event, result, countersRes, notAckedMetricType)
	}

	return result, nil
}

func (h *resolveHelper) PostProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	countersRes componentAndServiceCountersResult,
	notAckedMetricType string,
) {
	h.metricsSender.SendEventMetrics(
		result.Alarm,
		*event.Entity,
		result.AlarmChange,
		event.Parameters.Timestamp.Time,
		event.Parameters.Initiator,
		event.Parameters.User,
		event.Parameters.Instruction,
		notAckedMetricType,
	)

	h.componentAndServiceCountersHelper.PostProcess(ctx, countersRes)

	err := h.metaAlarmPostProcessor.Process(ctx, event, rpc.AxeResultEvent{
		Alarm:           &result.Alarm,
		AlarmChangeType: result.AlarmChange.Type,
	})
	if err != nil {
		h.logger.Err(err).Msg("cannot process meta alarm")
	}

	err = sendRemediationEvent(ctx, event, result, h.remediationRpcClient, h.encoder)
	if err != nil {
		h.logger.Err(err).Msg("cannot send event to engine-remediation")
	}

	if !result.Alarm.Value.PbehaviorInfo.IsDefaultActive() {
		err = updatePbehaviorAlarmCount(ctx, h.pbehaviorCollection, "", result.Alarm.Value.PbehaviorInfo.ID)
		if err != nil {
			h.logger.Err(err).Msg("cannot update pbehavior")
		}
	}

	if result.Alarm.Value.State.Value != types.AlarmStateOK {
		err = h.upstreamHelper.SendDownstreamEventsOnOK(ctx, *event.Entity)
		if err != nil {
			h.logger.Err(err).Msg("cannot send downstream events")
		}
	}
}

func (h *resolveHelper) UpdateAlarmToResolve(ctx context.Context, match bson.M, params rpc.AxeParameters) (types.Alarm, error) {
	beforeAlarm := types.Alarm{}
	update := h.GetResolveAlarmUpdate(datetime.NewCpsTime(), params)
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.Before).
		SetProjection(bson.M{
			"not_acked_metric_type":      1,
			"not_acked_metric_send_time": 1,
		})
	err := h.alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&beforeAlarm)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return beforeAlarm, err
	}

	return beforeAlarm, nil
}

func (h *resolveHelper) CopyAlarmToResolvedCollection(ctx context.Context, alarmID string) (types.Alarm, error) {
	// extend alarm struct with bookmarks to copy user's bookmarks to a resolved alarm document
	var alarm struct {
		types.Alarm `bson:"inline"`
		Bookmarks   []string `bson:"bookmarks"`
	}
	err := h.alarmCollection.FindOne(ctx, bson.M{"_id": alarmID}).Decode(&alarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return alarm.Alarm, nil
		}

		return alarm.Alarm, err
	}

	_, err = h.resolvedCollection.UpdateOne(
		ctx,
		bson.M{"_id": alarm.ID},
		bson.M{"$set": alarm},
		options.UpdateOne().SetUpsert(true),
	)

	return alarm.Alarm, err
}

func (h *resolveHelper) UpdateEntityOfResolvedAlarm(ctx context.Context, entityID string) (types.Entity, error) {
	entity := types.Entity{}
	entityUpdate := h.GetResolveEntityUpdate()
	err := h.entityCollection.FindOneAndUpdate(ctx, bson.M{"_id": entityID}, entityUpdate,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&entity)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return entity, err
	}

	return entity, nil
}

func (h *resolveHelper) RemoveMetaAlarmStateOnResolve(ctx context.Context, alarm types.Alarm) error {
	if !alarm.IsMetaAlarm() {
		return nil
	}

	var rule correlation.Rule
	err := h.metaAlarmRuleCollection.FindOne(ctx, bson.M{"_id": alarm.Value.Meta}).Decode(&rule)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return fmt.Errorf("meta alarm rule %s not found", alarm.Value.Meta)
		}

		return fmt.Errorf("cannot fetch meta alarm rule: %w", err)
	}

	return removeMetaAlarmState(ctx, alarm, rule, h.metaAlarmStatesService)
}

func (h *resolveHelper) GetResolveAlarmUpdate(t datetime.CpsTime, params rpc.AxeParameters) []bson.M {
	newStep := NewAlarmStep(types.AlarmStepResolve, params, false)
	newStep.Timestamp = t

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
								bson.M{"$ifNull": bson.A{"$v.snooze", nil}},
								nil,
							}},
							{"$not": bson.M{"$in": bson.A{
								bson.M{"$ifNull": bson.A{"$v.pbehavior_info.canonical_type", nil}},
								bson.A{nil, "", pbehavior.TypeActive},
							}}},
							{"$eq": bson.A{"$auto_instruction_in_progress", true}},
							{"$eq": bson.A{"$inactive_delay_meta_alarm_in_progress", true}},
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
			"v.steps":    bson.M{"$concatArrays": bson.A{"$v.steps", bson.A{bson.M{"$literal": newStep}}}},
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
						bson.M{"$ifNull": bson.A{"$v.snooze", nil}},
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
						bson.M{"$ifNull": bson.A{"$v.pbehavior_info.canonical_type", nil}},
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
			"v.close_delay_value",
		}},
	}
}

func (h *resolveHelper) GetResolveEntityUpdate() bson.M {
	return bson.M{"$unset": bson.M{
		"idle_since":           "",
		"last_idle_rule_apply": "",
	}}
}
