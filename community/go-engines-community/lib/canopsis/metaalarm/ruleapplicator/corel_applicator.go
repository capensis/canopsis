package ruleapplicator

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/storage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"html/template"
	"strings"
	"time"
)

const (
	CorelTypeParent = iota
	CorelTypeChild
)

const DayInSeconds = 86400

// CorelApplicator implements RuleApplicator interface
type CorelApplicator struct {
	alarmAdapter      alarm.Adapter
	metaAlarmService  service.MetaAlarmService
	storage           storage.GroupingStorageNew
	redisClient       *redis.Client
	redisLockClient   *redislock.Client
	ruleEntityCounter metaalarm.RuleEntityCounter
	logger            zerolog.Logger
}

// Apply called by RulesService.ProcessEvent
func (a CorelApplicator) Apply(ctx context.Context, event *types.Event, rule metaalarm.Rule) ([]types.Event, error) {
	var metaAlarmEvents []types.Event
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

	if (!rule.Config.AlarmPatterns.Matches(event.Alarm)) ||
		(!rule.Config.EntityPatterns.Matches(event.Entity)) ||
		(!rule.Config.EventPatterns.Matches(*event)) {
		return nil, nil
	}

	if rule.Config.ThresholdCount == nil {
		rule.Config.ThresholdCount = new(int64)
		*rule.Config.ThresholdCount = 2
	}

	if rule.Config.CorelStatus == "" || rule.Config.CorelID == "" || rule.Config.CorelParent == "" || rule.Config.CorelChild == "" {
		return nil, nil
	}

	if event.Alarm.IsMetaAlarm() {
		return nil, nil
	}

	corelID, err := a.renderTemplate(
		rule.Config.CorelID,
		types.AlarmWithEntity{
			Alarm:  *event.Alarm,
			Entity: *event.Entity,
		}, nil)
	if err != nil {
		return nil, err
	}

	if corelID == "" {
		return nil, nil
	}

	corelStatus, err := a.renderTemplate(
		rule.Config.CorelStatus,
		types.AlarmWithEntity{
			Alarm:  *event.Alarm,
			Entity: *event.Entity,
		}, nil)
	if err != nil {
		return nil, err
	}

	var corelType int

	switch corelStatus {
	case rule.Config.CorelParent:
		corelType = CorelTypeParent
	case rule.Config.CorelChild:
		corelType = CorelTypeChild
	default:
		return nil, nil
	}

	belongs, err := AlreadyBelongsToMetaalarm(a.alarmAdapter, event.GetEID(), rule.ID, corelID)
	if err != nil {
		return nil, err
	}
	if belongs {
		return nil, nil
	}

	// we gather two separate groups: one for children, one for parents
	childGroupID := fmt.Sprintf("%s&&%s", rule.ID, corelID)
	parentGroupId := fmt.Sprintf("%s$$parent", childGroupID)

	//we take threshold - 1, because 1 place is resolved by the parent.
	childrenThreshold := int(*rule.Config.ThresholdCount - 1)
	timeInterval := int64(rule.Config.TimeInterval)
	if rule.Config.TimeInterval == 0 {
		timeInterval = DayInSeconds
	}

	for redisRetries := MaxRedisRetries; redisRetries >= 0; redisRetries-- {
		watchErr = a.redisClient.Watch(ctx, func(tx *redis.Tx) error {
			parentGroup, parentOpenedAlarms, err := a.getGroupWithOpenedAlarmsWithEntity(ctx, tx, parentGroupId, timeInterval)
			if err != nil {
				return err
			}

			childrenGroup, childrenOpenedAlarms, err := a.getGroupWithOpenedAlarmsWithEntity(ctx, tx, childGroupID, timeInterval)
			if err != nil {
				return err
			}

			// get minimal from both groups
			openTime := parentGroup.GetOpenTime()
			if childrenGroup.GetOpenTime() < openTime {
				openTime = childrenGroup.GetOpenTime()
			}

			if childrenGroup.GetGroupLength() >= childrenThreshold {
				if event.Alarm.Value.LastUpdateDate.Unix() > openTime+timeInterval {
					//if groups time is over and we had at least one parent,
					//that means that metaAlarms were created and we can reset
					//both groups and start to gather again.
					if len(parentOpenedAlarms) != 0 {
						parentGroup = storage.NewAlarmGroup(parentGroupId)
						childrenGroup = storage.NewAlarmGroup(childGroupID)

						if corelType == CorelTypeParent {
							parentGroup.Push(*event.Alarm, timeInterval)
						} else {
							childrenGroup.Push(*event.Alarm, timeInterval)
						}

						// there is 100% group is not gathered, because were reset, so we can return
						return a.storage.SetMany(ctx, tx, timeInterval, parentGroup, childrenGroup)
					}

					//If there were no parent, then basically we can assume, that group wasn't completed,
					//we need to shift intervals and check if threshold is reached, so process is always.
				}
			}

			//if it's parent, then push to parent group and shift child group and vice versa
			if corelType == CorelTypeParent {
				parentGroup.Push(*event.Alarm, timeInterval)
				childrenGroup.RemoveBefore(event.Alarm.Value.LastUpdateDate.Unix() - timeInterval)
			} else {
				childrenGroup.Push(*event.Alarm, timeInterval)
				parentGroup.RemoveBefore(event.Alarm.Value.LastUpdateDate.Unix() - timeInterval)
			}

			metaAlarmLock, err = a.redisLockClient.Obtain(ctx, childGroupID, 100*time.Millisecond, &redislock.Options{
				RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(11*time.Millisecond), MaxRedisLockRetries),
			})

			if err != nil {
				a.logger.Err(err).
					Str("rule_id", rule.ID).
					Str("alarm_id", event.Alarm.ID).
					Msg("Update meta-alarm: obtain redlock failed, alarm will be skipped")

				return err
			}

			// get updated alarms for alarm groups
			err = a.alarmAdapter.GetOpenedAlarmsWithEntityByAlarmIDs(childrenGroup.GetAlarmIds(), &childrenOpenedAlarms)
			if err != nil {
				return err
			}

			parentId := ""
			parentIds := parentGroup.GetAlarmIds()
			if len(parentIds) != 0 {
				// we are interested only in first parent
				parentId = parentIds[0]
			}

			err = a.alarmAdapter.GetOpenedAlarmsWithEntityByAlarmIDs([]string{parentId}, &parentOpenedAlarms)
			if err != nil {
				return err
			}

			// update groups
			err = a.storage.SetMany(ctx, tx, timeInterval, parentGroup, childrenGroup)

			if childrenGroup.GetGroupLength() >= childrenThreshold {
				if corelType == CorelTypeParent {
					//parent should be the first one
					if len(parentOpenedAlarms) != 0 && parentOpenedAlarms[0].Alarm.ID != event.Alarm.ID {
						return nil
					}

					event.Alarm.Value.Meta = rule.ID
					event.Alarm.Value.MetaValuePath = childGroupID

					metaAlarmEvent, err := a.metaAlarmService.AddMultipleChildsToMetaAlarm(
						event,
						*event.Alarm,
						childrenOpenedAlarms,
						rule,
					)
					if err != nil {
						return err
					}

					metaAlarmEvents = append(metaAlarmEvents, metaAlarmEvent)
					return nil
				}

				if len(parentOpenedAlarms) != 0 {
					parentOpenedAlarms[0].Alarm.Value.Meta = rule.ID
					parentOpenedAlarms[0].Alarm.Value.MetaValuePath = childGroupID

					metaAlarmEvent, err := a.metaAlarmService.AddChildToMetaAlarm(
						event,
						parentOpenedAlarms[0].Alarm,
						types.AlarmWithEntity{
							Alarm:  *event.Alarm,
							Entity: *event.Entity,
						},
						rule,
					)
					if err != nil {
						return err
					}

					metaAlarmEvents = append(metaAlarmEvents, metaAlarmEvent)
				}
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

			return err
		}, childGroupID, parentGroupId)

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

	if watchErr != nil {
		return nil, watchErr
	}

	return metaAlarmEvents, nil
}

func (a CorelApplicator) getGroupWithOpenedAlarmsWithEntity(ctx context.Context, tx *redis.Tx, key string, timeInterval int64) (storage.TimeBasedAlarmGroup, []types.AlarmWithEntity, error) {
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

func (a CorelApplicator) renderTemplate(templateStr string, data interface{}, f template.FuncMap) (string, error) {
	t, err := template.New("template").Funcs(f).Parse(templateStr)
	if err != nil {
		return "", err
	}
	b := strings.Builder{}
	err = t.Execute(&b, data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// NewCorelApplicator instantiates CorelApplicator with MetaAlarmService
func NewCorelApplicator(alarmAdapter alarm.Adapter, metaAlarmService service.MetaAlarmService, storage storage.GroupingStorageNew, redisClient *redis.Client, redisLockClient *redislock.Client, logger zerolog.Logger) CorelApplicator {
	return CorelApplicator{
		alarmAdapter:     alarmAdapter,
		metaAlarmService: metaAlarmService,
		storage:          storage,
		redisClient:      redisClient,
		redisLockClient:  redisLockClient,
		logger:           logger,
	}
}
