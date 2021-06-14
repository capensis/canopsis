package ruleapplicator

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/storage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"github.com/bsm/redislock"
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
	redisLockClient   *redislock.Client
	ruleEntityCounter metaalarm.RuleEntityCounter
	logger            zerolog.Logger
}

// Apply called by RulesService.ProcessEvent
func (a ComplexApplicator) Apply(ctx context.Context, event *types.Event, rule metaalarm.Rule) ([]types.Event, error) {
	var metaAlarmEvent types.Event
	var watchErr error
	var metaAlarmLock *redislock.Lock

	defer func() {
		if metaAlarmLock != nil {
			err := metaAlarmLock.Release(ctx)
			if err != nil && err != redislock.ErrLockNotHeld {
				a.logger.Warn().
					Str("rule_id", rule.ID).
					Str("alarm_id", event.Alarm.ID).
					Msg("Update meta-alarm: failed to manually release redlock, the lock will be released by ttl")
			}
		}
	}()

	if rule.Type == metaalarm.RuleTypeTimeBased {
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
		(!rule.Config.EventPatterns.Matches(*event)) {
		return nil, nil
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
				groupLen, err := a.getGroupLen(ctx, tx, rule.ID)
				if err != nil {
					return err
				}

				if groupLen >= *rule.Config.ThresholdCount {
					alarmGroup, err := a.storage.Get(ctx, tx, rule.ID)
					if err != nil {
						return err
					}

					if event.Alarm.Value.LastUpdateDate.After(alarmGroup.OpenTime.Add(time.Duration(rule.Config.TimeInterval) * time.Second)) {
						err := a.storage.CleanPush(ctx, tx, rule, *event.Alarm, "")
						if err != nil {
							return err
						}
					} else {
						maxRetries := 0

						if watchErr == redis.TxFailedErr {
							maxRetries = MaxMongoRetries
						}

						updated := false

						for mongoRetries := maxRetries; mongoRetries >= 0 && !updated; mongoRetries-- {
							metaAlarm, err := a.alarmAdapter.GetOpenedMetaAlarm(rule.ID, "")
							switch err.(type) {
							case errt.NotFound:
								if mongoRetries == maxRetries {
									metaAlarmEvent, err = a.createMetaAlarm(ctx, tx, event, rule)
									if err != nil {
										return err
									}

									updated = true

									break
								}

								a.logger.Warn().
									Str("rule_id", rule.ID).
									Str("alarm_id", event.Alarm.ID).
									Msgf("Another instance has created meta-alarm, but couldn't find an opened meta-alarm. Retry mongo query. Remaining retries: %d", mongoRetries)

								time.Sleep(10 * time.Millisecond)

								continue
							case nil:
								metaAlarmLock, err = a.redisLockClient.Obtain(ctx, metaAlarm.ID, 100*time.Millisecond, &redislock.Options{
									RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(11*time.Millisecond), MaxRedisLockRetries),
								})

								if err != nil {
									a.logger.Err(err).
										Str("rule_id", rule.ID).
										Str("alarm_id", event.Alarm.ID).
										Msg("Update meta-alarm: obtain redlock failed, alarm will be skipped")

									return err
								}

								metaAlarmEvent, err = a.metaAlarmService.AddChildToMetaAlarm(
									event,
									metaAlarm,
									types.AlarmWithEntity{Alarm: *event.Alarm, Entity: *event.Entity},
									rule,
								)
								if err != nil {
									return err
								}

								err = metaAlarmLock.Release(ctx)
								if err != nil {
									if err == redislock.ErrLockNotHeld {
										a.logger.Err(err).
											Str("rule_id", rule.ID).
											Str("alarm_id", event.Alarm.ID).
											Msg("Update meta-alarm: the update process took more time than redlock ttl, data might be inconsistent")
									} else {
										a.logger.Warn().
											Str("rule_id", rule.ID).
											Str("alarm_id", event.Alarm.ID).
											Msg("Update meta-alarm: failed to manually release redlock, the lock will be released by ttl")
									}
								}

								metaAlarmLock = nil

								updated = true
							default:
								return err
							}
						}
					}
				} else {
					err = a.storage.Push(ctx, tx, rule, *event.Alarm, "")
					if err != nil {
						return err
					}

					groupLen, err = a.getGroupLen(ctx, tx, rule.ID)
					if err != nil {
						return err
					}

					if groupLen >= *rule.Config.ThresholdCount {
						metaAlarmEvent, err = a.createMetaAlarm(ctx, tx, event, rule)
					}
				}

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
				ratioReached, err := a.isRatioReached(ctx, tx, rule, true)
				if err != nil {
					return err
				}

				if ratioReached {
					alarmGroup, err := a.storage.Get(ctx, tx, rule.ID)
					if err != nil {
						return err
					}

					if event.Alarm.Value.LastUpdateDate.After(alarmGroup.OpenTime.Add(time.Duration(rule.Config.TimeInterval) * time.Second)) {
						err := a.storage.CleanPush(ctx, tx, rule, *event.Alarm, "")
						if err != nil {
							return err
						}
					} else {
						maxRetries := 0

						if watchErr == redis.TxFailedErr {
							maxRetries = MaxMongoRetries
						}

						updated := false

						for mongoRetries := maxRetries; mongoRetries >= 0 && !updated; mongoRetries-- {
							metaAlarm, err := a.alarmAdapter.GetOpenedMetaAlarm(rule.ID, "")
							switch err.(type) {
							case errt.NotFound:
								if mongoRetries == maxRetries {
									metaAlarmEvent, err = a.createMetaAlarm(ctx, tx, event, rule)
									if err != nil {
										return err
									}

									updated = true

									break
								}

								a.logger.Warn().
									Str("rule_id", rule.ID).
									Str("alarm_id", event.Alarm.ID).
									Msgf("Another instance has created meta-alarm, but couldn't find an opened meta-alarm. Retry mongo query. Remaining retries: %d", mongoRetries)

								time.Sleep(10 * time.Millisecond)

								continue
							case nil:
								metaAlarmLock, err = a.redisLockClient.Obtain(ctx, metaAlarm.ID, 100*time.Millisecond, &redislock.Options{
									RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(11*time.Millisecond), MaxRedisLockRetries),
								})

								if err != nil {
									a.logger.Err(err).
										Str("rule_id", rule.ID).
										Str("alarm_id", event.Alarm.ID).
										Msg("Update meta-alarm: obtain redlock failed, alarm will be skipped")

									return err
								}

								err = a.storage.Push(ctx, tx, rule, *event.Alarm, "")
								if err != nil {
									return err
								}

								metaAlarmEvent, err = a.metaAlarmService.AddChildToMetaAlarm(
									event,
									metaAlarm,
									types.AlarmWithEntity{Alarm: *event.Alarm, Entity: *event.Entity},
									rule,
								)
								if err != nil {
									return err
								}

								err = metaAlarmLock.Release(ctx)
								if err != nil {
									if err == redislock.ErrLockNotHeld {
										a.logger.Err(err).
											Str("rule_id", rule.ID).
											Str("alarm_id", event.Alarm.ID).
											Msg("Update meta-alarm: the update process took more time than redlock ttl, data might be inconsistent")
									} else {
										a.logger.Warn().
											Str("rule_id", rule.ID).
											Str("alarm_id", event.Alarm.ID).
											Msg("Update meta-alarm: failed to manually release redlock, the lock will be released by ttl")
									}
								}

								metaAlarmLock = nil

								updated = true
							default:
								return err
							}
						}
					}
				} else {
					err = a.storage.Push(ctx, tx, rule, *event.Alarm, "")
					if err != nil {
						return err
					}

					ratioReached, err = a.isRatioReached(ctx, tx, rule, false)
					if err != nil {
						return err
					}

					if ratioReached {
						metaAlarmEvent, err = a.createMetaAlarm(ctx, tx, event, rule)
					}
				}

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

func (a ComplexApplicator) isRatioReached(ctx context.Context, tx *redis.Tx, rule metaalarm.Rule, includeNewAlarmOnRecompute bool) (bool, error) {
	groupLen, err := a.getGroupLen(ctx, tx, rule.ID)
	if err != nil {
		return false, err
	}

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

func (a ComplexApplicator) recomputeTotal(ctx context.Context, rule metaalarm.Rule) (int, error) {
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

func (a ComplexApplicator) createMetaAlarm(ctx context.Context, tx *redis.Tx, event *types.Event, rule metaalarm.Rule) (types.Event, error) {
	var children []types.AlarmWithEntity

	alarmGroup, err := a.storage.Get(ctx, tx, rule.ID)
	if err != nil {
		return types.Event{}, err
	}

	err = a.alarmAdapter.GetOpenedAlarmsWithEntityByAlarmIDs(alarmGroup.GetAlarmIds(), &children)
	if err != nil {
		return types.Event{}, err
	}

	return a.metaAlarmService.CreateMetaAlarm(event, children, rule)
}

func (a ComplexApplicator) getGroupLen(ctx context.Context, tx *redis.Tx, ruleId string) (int64, error) {
	var children []types.AlarmWithEntity

	alarmGroup, err := a.storage.Get(ctx, tx, ruleId)
	if err != nil {
		return 0, err
	}

	// We need to check for resolved alarms here
	err = a.alarmAdapter.GetOpenedAlarmsWithEntityByAlarmIDs(alarmGroup.GetAlarmIds(), &children)

	return int64(len(children)), err
}

// NewComplexApplicator instantiates ComplexApplicator with MetaAlarmService
func NewComplexApplicator(alarmAdapter alarm.Adapter, metaAlarmService service.MetaAlarmService, redisClient *redis.Client, redisLockClient *redislock.Client, ruleEntityCounter metaalarm.RuleEntityCounter, logger zerolog.Logger) ComplexApplicator {
	return ComplexApplicator{
		alarmAdapter:      alarmAdapter,
		metaAlarmService:  metaAlarmService,
		storage:           storage.NewRedisGroupingStorage(),
		redisClient:       redisClient,
		redisLockClient:   redisLockClient,
		ruleEntityCounter: ruleEntityCounter,
		logger:            logger,
	}
}
