package ruleapplicator

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
)

type ParentChildApplicator struct {
	alarmAdapter     alarm.Adapter
	metaAlarmService service.MetaAlarmService
	redisLockClient  *redislock.Client
	logger           zerolog.Logger
}

func (a ParentChildApplicator) Apply(ctx context.Context, event *types.Event, rule metaalarm.Rule) ([]types.Event, error) {
	var metaAlarmEvent types.Event
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

	if event.SourceType == types.SourceTypeComponent {
		//skip is component alarm is already a meta-alarm
		if !event.Alarm.IsMetaAlarm() {
			resourceAlarms, err := a.alarmAdapter.GetAllOpenedResourceAlarmsByComponent(event.Component)
			if err != nil {
				return nil, err
			}

			if len(resourceAlarms) != 0 {
				componentAlarm := *event.Alarm
				componentAlarm.Value.Meta = rule.ID

				metaAlarmLock, err = a.redisLockClient.Obtain(ctx, componentAlarm.ID, 100*time.Millisecond, &redislock.Options{
					RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(11*time.Millisecond), MaxRedisLockRetries),
				})

				if err != nil {
					a.logger.Err(err).
						Str("rule_id", rule.ID).
						Str("alarm_id", event.Alarm.ID).
						Msg("Update meta-alarm: obtain redlock failed, alarm will be skipped")

					return nil, err
				}

				metaAlarmEvent, err = a.metaAlarmService.AddMultipleChildsToMetaAlarm(event, componentAlarm, resourceAlarms, rule)
				if err != nil {
					return nil, err
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
			}
		}
	}

	if event.SourceType == types.SourceTypeResource {
		// Check if component alarm exists and if it's already a meta-alarm
		componentAlarm, err := a.alarmAdapter.GetLastAlarm(event.Connector, event.ConnectorName, event.Component)
		if err != nil {
			if !a.isNotFound(err) {
				return nil, err
			}

			return nil, nil
		}

		if !componentAlarm.IsMalfunctioning() {
			return nil, nil
		}

		if !componentAlarm.IsMetaAlarm() {
			// transform to meta-alarm
			componentAlarm.Value.Meta = rule.ID
		}

		metaAlarmLock, err = a.redisLockClient.Obtain(ctx, componentAlarm.ID, 100*time.Millisecond, &redislock.Options{
			RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(11*time.Millisecond), MaxRedisLockRetries),
		})

		if err != nil {
			a.logger.Err(err).
				Str("rule_id", rule.ID).
				Str("alarm_id", event.Alarm.ID).
				Msg("Update meta-alarm: obtain redlock failed, alarm will be skipped")

			return nil, err
		}

		metaAlarmEvent, err = a.metaAlarmService.AddChildToMetaAlarm(
			event,
			componentAlarm,
			types.AlarmWithEntity{Alarm: *event.Alarm, Entity: *event.Entity},
			rule,
		)
		if err != nil {
			return nil, err
		}

		err = metaAlarmLock.Release(ctx)
		if err != nil {
			if err == redislock.ErrLockNotHeld {
				a.logger.Err(err).
					Str("rule_id", rule.ID).
					Str("alarm_id", event.Alarm.ID).
					Msg("Update meta-alarm: the update process took more time than redlock ttl, data might be inconsistent")
			}

			a.logger.Warn().
				Str("rule_id", rule.ID).
				Str("alarm_id", event.Alarm.ID).
				Msg("Update meta-alarm: failed to manually release redlock, the lock will be released by ttl")
		}

		metaAlarmLock = nil
	}

	if metaAlarmEvent.EventType != "" {
		return []types.Event{metaAlarmEvent}, nil
	}

	return nil, nil
}

func (a ParentChildApplicator) isNotFound(err error) bool {
	_, ok := err.(errt.NotFound)

	return ok
}

func NewParentChildApplicator(alarmAdapter alarm.Adapter, metaAlarmService service.MetaAlarmService, redisLockClient *redislock.Client, logger zerolog.Logger) ParentChildApplicator {
	return ParentChildApplicator{
		alarmAdapter:     alarmAdapter,
		metaAlarmService: metaAlarmService,
		redisLockClient:  redisLockClient,
		logger:           logger,
	}
}
