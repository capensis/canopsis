package ruleapplicator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/storage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

// ValueApplicator implements RuleApplicator interface
type ValueApplicator struct {
	alarmAdapter      alarm.Adapter
	metaAlarmService  service.MetaAlarmService
	storage           storage.GroupingStorage
	redisClient       *redis.Client
	ruleEntityCounter metaalarm.ValueGroupEntityCounter
	logger            zerolog.Logger
}

// Apply called by RulesService.ProcessEvent
func (a ValueApplicator) Apply(ctx context.Context, event types.Event, rule metaalarm.Rule) ([]types.Event, error) {
	var metaAlarmEvent types.Event
	var watchErr error

	if len(rule.Config.ValuePaths) == 0 || rule.Config.TimeInterval == 0 {
		return nil, nil
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

	valuePath, valuePathMap := a.extractValue(event, &rule)
	if valuePath == "" {
		return nil, nil
	}

	belongs, err := AlreadyBelongsToMetaalarm(a.alarmAdapter, event.GetEID(), rule.ID, valuePath)
	if err != nil {
		return nil, err
	}
	if belongs {
		return nil, nil
	}

	ruleKey := fmt.Sprintf("%s&&%s", rule.ID, valuePath)

	if rule.Config.ThresholdCount != nil {
		for redisRetries := MaxRedisRetries; redisRetries >= 0; redisRetries-- {
			watchErr = a.redisClient.Watch(ctx, func(tx *redis.Tx) error {
				alarmGroup, openedAlarms, err := a.getGroupWithOpenedAlarmsWithEntity(ctx, tx, ruleKey, timeInterval)
				if err != nil {
					return err
				}

				if alarmGroup.GetGroupLength() >= *rule.Config.ThresholdCount {
					if event.Alarm.Value.LastUpdateDate.Unix() > alarmGroup.GetOpenTime() + timeInterval {
						return a.storage.Set(ctx, tx, storage.NewAlarmGroup(ruleKey), DefaultConfigTimeInterval)
					}

					maxRetries := 0
					if watchErr == redis.TxFailedErr {
						maxRetries = MaxMongoRetries
					}

					for mongoRetries := maxRetries; mongoRetries >= 0; mongoRetries-- {
						metaAlarm, err := a.alarmAdapter.GetOpenedMetaAlarm(rule.ID, valuePath)
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
					if err != nil {
						return err
					}

					metaAlarmEvent.MetaAlarmValuePath = valuePath
					return nil
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
				if err != nil {
					return err
				}

				metaAlarmEvent.MetaAlarmValuePath = valuePath
				return nil
			}, ruleKey)

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
				alarmGroup, openedAlarms, err := a.getGroupWithOpenedAlarmsWithEntity(ctx, tx, ruleKey, timeInterval)
				if err != nil {
					return err
				}

				ratioReached, err := a.isRatioReached(ctx, alarmGroup, rule, valuePath, valuePathMap)
				if err != nil {
					return err
				}

				if ratioReached {
					alarmGroup, err := a.storage.Get(ctx, tx, ruleKey)
					if err != nil {
						return err
					}

					if event.Alarm.Value.LastUpdateDate.Unix() > alarmGroup.GetOpenTime() + timeInterval {
						return a.storage.Set(ctx, tx, storage.NewAlarmGroup(ruleKey), DefaultConfigTimeInterval)
					}

					maxRetries := 0
					if watchErr == redis.TxFailedErr {
						maxRetries = MaxMongoRetries
					}

					for mongoRetries := maxRetries; mongoRetries >= 0; mongoRetries-- {
						metaAlarm, err := a.alarmAdapter.GetOpenedMetaAlarm(rule.ID, valuePath)
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
							if err != nil {
								return err
							}
						default:
							return err
						}
					}

					metaAlarmEvent, err = a.metaAlarmService.CreateMetaAlarm(event, append(openedAlarms, types.AlarmWithEntity{
						Alarm:  *event.Alarm,
						Entity: *event.Entity,
					}), rule)
					if err != nil {
						return err
					}

					metaAlarmEvent.MetaAlarmValuePath = valuePath
					return nil
				}

				alarmGroup.Push(*event.Alarm, timeInterval)
				err = a.storage.Set(ctx, tx, alarmGroup, timeInterval)
				if err != nil {
					return err
				}

				ratioReached, err = a.isRatioReached(ctx, alarmGroup, rule, valuePath, valuePathMap)
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
				if err != nil {
					return err
				}

				metaAlarmEvent.MetaAlarmValuePath = valuePath
				return nil
			}, ruleKey)

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

func (a ValueApplicator) getGroupWithOpenedAlarmsWithEntity(ctx context.Context, tx *redis.Tx, key string, timeInterval int64) (storage.TimeBasedAlarmGroup, []types.AlarmWithEntity, error) {
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

func (a ValueApplicator) extractValue(event types.Event, rule *metaalarm.Rule) (string, map[string]string) {
	bytes, err := json.Marshal(types.AlarmWithEntity{
		Alarm:  *event.Alarm,
		Entity: *event.Entity,
	})
	if err != nil {
		a.logger.Err(err).Str("rule-id", rule.ID).Msg("")
		return "", nil
	}

	var paths []string
	valuePathsMap := make(map[string]string)
	for _, valuePath := range rule.Config.ValuePaths {
		value := gjson.GetBytes(bytes, valuePath).String()
		if value == "" {
			return "", nil
		}

		paths = append(paths, value)
		valuePathsMap[valuePath] = value
	}

	return strings.Join(paths, "."), valuePathsMap
}

func (a ValueApplicator) isRatioReached(ctx context.Context, alarmGroup storage.TimeBasedAlarmGroup, rule metaalarm.Rule, valuePath string, valuePathsMap map[string]string) (bool, error) {
	groupLen := alarmGroup.GetGroupLength()

	total, err := a.ruleEntityCounter.GetTotalEntitiesAmount(ctx, rule.ID, valuePath)
	if err != nil {
		return false, err
	}

	if total != 0 && float64(groupLen) / float64(total) >= *rule.Config.ThresholdRate {
		err = a.ruleEntityCounter.CountTotalEntitiesAmountForValuePaths(ctx, rule, valuePathsMap)
		if err != nil {
			return false, err
		}

		total, err = a.ruleEntityCounter.GetTotalEntitiesAmount(ctx, rule.ID, valuePath)
		if err != nil {
			return false, err
		}
	}

	return total != 0 && float64(groupLen) / float64(total) >= *rule.Config.ThresholdRate, nil
}

// NewValueApplicator instantiates ValueApplicator with MetaAlarmService
func NewValueGroupApplicator(alarmAdapter alarm.Adapter, metaAlarmService service.MetaAlarmService, storage storage.GroupingStorage, redisClient *redis.Client, ruleEntityCounter metaalarm.ValueGroupEntityCounter, logger zerolog.Logger) ValueApplicator {
	return ValueApplicator{
		alarmAdapter:      alarmAdapter,
		metaAlarmService:  metaAlarmService,
		storage:           storage,
		redisClient:       redisClient,
		ruleEntityCounter: ruleEntityCounter,
		logger:            logger,
	}
}
