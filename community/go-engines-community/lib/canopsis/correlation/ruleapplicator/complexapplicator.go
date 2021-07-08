package ruleapplicator

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation/storage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

const MaxRedisRetries = 10
const MaxMongoRetries = 10
const MaxRedisLockRetries = 10

// ComplexApplicator implements RuleApplicator interface
type ComplexApplicator struct {
	alarmAdapter      alarm.Adapter
	metaAlarmService  service.MetaAlarmService
	storage           storage.GroupingStorage
	redisClient       *redis.Client
	ruleEntityCounter correlation.RuleEntityCounter
	logger            zerolog.Logger
}

// Apply called by RulesService.ProcessEvent
func (a ComplexApplicator) Apply(ctx context.Context, event types.Event, rule correlation.Rule) ([]types.Event, error) {
	var metaAlarmEvent types.Event
	var watchErr error

	if rule.Type == correlation.RuleTypeTimeBased {
		rule.Config.ThresholdCount = new(int64)
		*rule.Config.ThresholdCount = 2
	}

	if rule.Config.ThresholdRate == nil && rule.Config.ThresholdCount == nil {
		return nil, nil
	}

	if rule.Config.ThresholdRate != nil && rule.Config.ThresholdCount != nil {
		return nil, nil
	}

	if (!rule.Config.AlarmPatterns.Matches(event.Alarm)) ||
		(!rule.Config.EntityPatterns.Matches(event.Entity)) ||
		(!rule.Config.EventPatterns.Matches(event)) {
		return nil, nil
	}

	timeInterval := int64(rule.Config.TimeInterval)
	if rule.Config.TimeInterval == 0 {
		timeInterval = DefaultConfigTimeInterval
	}

	belongs, err := AlreadyBelongsToMetaalarm(a.alarmAdapter, event.GetEID(), rule.ID, "")
	if err != nil {
		return nil, err
	}
	if belongs {
		return nil, nil
	}

	if rule.Config.ThresholdCount != nil {
		for redisRetries := MaxRedisRetries; redisRetries >= 0; redisRetries-- {
			watchErr = a.redisClient.Watch(ctx, func(tx *redis.Tx) error {
				alarmGroup, openedAlarms, err := a.getGroupWithOpenedAlarmsWithEntity(ctx, tx, rule.ID, timeInterval)
				if err != nil {
					return err
				}

				if alarmGroup.GetGroupLength() >= *rule.Config.ThresholdCount {
					if event.Alarm.Value.LastUpdateDate.Unix() > alarmGroup.GetOpenTime() + timeInterval {
						return a.storage.Set(ctx, tx, storage.NewAlarmGroup(rule.ID), DefaultConfigTimeInterval)
					}

					maxRetries := 0
					if watchErr == redis.TxFailedErr {
						maxRetries = MaxMongoRetries
					}

					for mongoRetries := maxRetries; mongoRetries >= 0; mongoRetries-- {
						metaAlarm, err := a.alarmAdapter.GetOpenedMetaAlarm(rule.ID, "")
						switch err.(type) {
						case errt.NotFound:
							a.logger.Warn().
								Str("rule_id", rule.ID).
								Str("alarm_id", event.Alarm.ID).
								Msgf("Another instance has created meta-alarm, but couldn't find an opened meta-alarm. Retry mongo query. Remaining retries: %d", mongoRetries)

							time.Sleep(10 * time.Millisecond)

							continue
						case nil:
							metaAlarmEvent, err = a.metaAlarmService.AddChildToMetaAlarm(
								ctx,
								event,
								metaAlarm,
								types.AlarmWithEntity{Alarm: *event.Alarm, Entity: *event.Entity},
								rule,
							)

							return err
						default:
							return err
						}
					}

					metaAlarmEvent, err = a.metaAlarmService.CreateMetaAlarm(event, append(openedAlarms, types.AlarmWithEntity{
						Alarm:  *event.Alarm,
						Entity: *event.Entity,
					}), rule)

					return err
				}

				alarmGroup.Push(*event.Alarm, timeInterval)
				err = a.storage.Set(ctx, tx, alarmGroup, timeInterval)
				if err != nil {
					return err
				}

				if alarmGroup.GetGroupLength() < *rule.Config.ThresholdCount {
					return nil
				}

				metaAlarmEvent, err = a.metaAlarmService.CreateMetaAlarm(event, append(openedAlarms, types.AlarmWithEntity{
					Alarm:  *event.Alarm,
					Entity: *event.Entity,
				}), rule)

				return err
			}, rule.ID)

			if watchErr == redis.TxFailedErr {
				a.logger.Warn().
					Str("rule_id", rule.ID).
					Str("alarm_id", event.Alarm.ID).
					Msgf("Redis transaction failed because of instances concurrency. Retry the alarm process. Remaining retries: %d", redisRetries)

				continue
			}

			break
		}

		if watchErr == redis.TxFailedErr {
			a.logger.Err(watchErr).
				Str("rule_id", rule.ID).
				Str("alarm_id", event.Alarm.ID).
				Msgf("Failed to process meta-alarm group after %d retries, alarm will be skipped", MaxRedisRetries)
		}
	}

	if rule.Config.ThresholdRate != nil {
		for redisRetries := MaxRedisRetries; redisRetries >= 0; redisRetries-- {
			watchErr = a.redisClient.Watch(ctx, func(tx *redis.Tx) error {
				alarmGroup, openedAlarms, err := a.getGroupWithOpenedAlarmsWithEntity(ctx, tx, rule.ID, timeInterval)
				if err != nil {
					return err
				}

				ratioReached, err := a.isRatioReached(ctx, alarmGroup, rule, true)
				if err != nil {
					return err
				}

				if ratioReached {
					if event.Alarm.Value.LastUpdateDate.Unix() > alarmGroup.GetOpenTime() + timeInterval {
						return a.storage.Set(ctx, tx, storage.NewAlarmGroup(rule.ID), DefaultConfigTimeInterval)
					}

					maxRetries := 0
					if watchErr == redis.TxFailedErr {
						maxRetries = MaxMongoRetries
					}

					for mongoRetries := maxRetries; mongoRetries >= 0; mongoRetries-- {
						metaAlarm, err := a.alarmAdapter.GetOpenedMetaAlarm(rule.ID, "")
						switch err.(type) {
						case errt.NotFound:
							a.logger.Warn().
								Str("rule_id", rule.ID).
								Str("alarm_id", event.Alarm.ID).
								Msgf("Another instance has created meta-alarm, but couldn't find an opened meta-alarm. Retry mongo query. Remaining retries: %d", mongoRetries)

							time.Sleep(10 * time.Millisecond)

							continue
						case nil:
							metaAlarmEvent, err = a.metaAlarmService.AddChildToMetaAlarm(
								ctx,
								event,
								metaAlarm,
								types.AlarmWithEntity{Alarm: *event.Alarm, Entity: *event.Entity},
								rule,
							)

							return err
						default:
							return err
						}
					}

					metaAlarmEvent, err = a.metaAlarmService.CreateMetaAlarm(event, append(openedAlarms, types.AlarmWithEntity{
						Alarm:  *event.Alarm,
						Entity: *event.Entity,
					}), rule)

					return err
				}

				alarmGroup.Push(*event.Alarm, timeInterval)
				err = a.storage.Set(ctx, tx, alarmGroup, timeInterval)
				if err != nil {
					return err
				}

				ratioReached, err = a.isRatioReached(ctx, alarmGroup, rule, false)
				if err != nil {
					return err
				}

				if !ratioReached {
					return nil
				}

				metaAlarmEvent, err = a.metaAlarmService.CreateMetaAlarm(event, append(openedAlarms, types.AlarmWithEntity{
					Alarm:  *event.Alarm,
					Entity: *event.Entity,
				}), rule)

				return err
			}, rule.ID)

			if watchErr == redis.TxFailedErr {
				a.logger.Warn().
					Str("rule_id", rule.ID).
					Str("alarm_id", event.Alarm.ID).
					Msgf("Redis transaction failed because of instances concurrency. Retry the alarm process. Remaining retries: %d", redisRetries)

				continue
			}

			break
		}

		if watchErr == redis.TxFailedErr {
			a.logger.Err(watchErr).
				Str("rule_id", rule.ID).
				Str("alarm_id", event.Alarm.ID).
				Msgf("Failed to process meta-alarm group after %d retries, alarm will be skipped", MaxRedisRetries)
		}
	}

	if watchErr != nil {
		return nil, watchErr
	}

	if metaAlarmEvent.EventType != "" {
		return []types.Event{metaAlarmEvent}, nil
	}

	return nil, nil
}

func (a ComplexApplicator) isRatioReached(ctx context.Context, alarmGroup storage.TimeBasedAlarmGroup, rule correlation.Rule, includeNewAlarmOnRecompute bool) (bool, error) {
	groupLen := alarmGroup.GetGroupLength()

	total, err := a.ruleEntityCounter.GetTotalEntitiesAmount(ctx, rule)
	if err != nil {
		return false, err
	}

	if total != 0 && float64(groupLen)/float64(total) >= *rule.Config.ThresholdRate {
		total, err = a.recomputeTotal(ctx, rule)
		if err != nil {
			return false, err
		}

		if includeNewAlarmOnRecompute {
			groupLen++
		}
	}

	return total != 0 && float64(groupLen)/float64(total) >= *rule.Config.ThresholdRate, nil
}

func (a ComplexApplicator) recomputeTotal(ctx context.Context, rule correlation.Rule) (int, error) {
	err := a.ruleEntityCounter.CountTotalEntitiesAmount(ctx, rule)
	if err != nil {
		return 0, err
	}

	total, err := a.ruleEntityCounter.GetTotalEntitiesAmount(ctx, rule)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (a ComplexApplicator) getGroupWithOpenedAlarmsWithEntity(ctx context.Context, tx *redis.Tx, key string, timeInterval int64) (storage.TimeBasedAlarmGroup, []types.AlarmWithEntity, error) {
	var alarms []types.AlarmWithEntity

	alarmGroup, err := a.storage.Get(ctx, tx, key)
	if err != nil {
		return nil, nil, err
	}

	err = a.alarmAdapter.GetOpenedAlarmsWithEntityByAlarmIDs(alarmGroup.GetAlarmIds(), &alarms)
	if err != nil {
		return nil, nil, err
	}

	alarmGroup = storage.NewAlarmGroup(key)
	for _, v := range alarms {
		alarmGroup.Push(v.Alarm, timeInterval)
	}

	return alarmGroup, alarms, err
}

// NewComplexApplicator instantiates ComplexApplicator with MetaAlarmService
func NewComplexApplicator(alarmAdapter alarm.Adapter, metaAlarmService service.MetaAlarmService, storage storage.GroupingStorage, redisClient *redis.Client, ruleEntityCounter correlation.RuleEntityCounter, logger zerolog.Logger) ComplexApplicator {
	return ComplexApplicator{
		alarmAdapter:      alarmAdapter,
		metaAlarmService:  metaAlarmService,
		storage:           storage,
		redisClient:       redisClient,
		ruleEntityCounter: ruleEntityCounter,
		logger:            logger,
	}
}
