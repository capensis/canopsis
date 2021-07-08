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

// AttributeApplicator implements RuleApplicator interface
type AttributeApplicator struct {
	alarmAdapter alarm.Adapter
	service      service.MetaAlarmService
	storage      storage.GroupingStorage
	redisClient  *redis.Client
	logger       zerolog.Logger
}

// Apply called by RulesService.ProcessEvent
func (a AttributeApplicator) Apply(ctx context.Context, event types.Event, rule correlation.Rule) ([]types.Event, error) {
	var metaAlarmEvent types.Event
	var watchErr error

	// Check alarm attributes matched with rule
	if !a.AlarmMatched(event, rule) || !a.EntityMatched(event, rule) || !a.EventMatched(event, rule) {
		return nil, nil
	}

	belongs, err := AlreadyBelongsToMetaalarm(a.alarmAdapter, event.GetEID(), rule.ID, "")
	if err != nil {
		return nil, err
	}
	if belongs {
		return nil, nil
	}

	for redisRetries := MaxRedisRetries; redisRetries >= 0; redisRetries-- {
		watchErr = a.redisClient.Watch(ctx, func(tx *redis.Tx) error {
			maxRetries := 0
			if watchErr == redis.TxFailedErr {
				maxRetries = MaxMongoRetries
			}

			for mongoRetries := maxRetries; maxRetries >= 0; maxRetries-- {
				// Check if meta-alarm already exists
				metaAlarm, err := a.alarmAdapter.GetOpenedMetaAlarm(rule.ID, "")
				switch err.(type) {
				case errt.NotFound:
					a.logger.Warn().
						Str("rule_id", rule.ID).
						Str("alarm_id", event.Alarm.ID).
						Msgf("Another instance has created meta-alarm, but couldn't find an opened meta-alarm. Retry mongo query. Remaining retries: %d", mongoRetries)

					time.Sleep(50 * time.Millisecond)

					continue
				case nil:
					metaAlarmEvent, err = a.service.AddChildToMetaAlarm(
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

			err = a.storage.Set(ctx, tx, storage.NewAlarmGroup(rule.ID), DefaultConfigTimeInterval)
			if err != nil {
				return err
			}

			metaAlarmEvent, err = a.service.CreateMetaAlarm(
				event,
				[]types.AlarmWithEntity{{
					Alarm:  *event.Alarm,
					Entity: *event.Entity,
				}},
				rule,
			)

			return err
		}, rule.ID)

		if watchErr == redis.TxFailedErr {
			a.logger.Warn().
				Str("rule_id", rule.ID).
				Str("alarm_id", event.Alarm.ID).
				Msgf("Redis transaction failed because of instances concurrency. Retry the alarm process. Remaining retries: %d", MaxRedisRetries)

			continue
		}

		break
	}

	if watchErr != nil {
		return nil, watchErr
	}

	if metaAlarmEvent.EventType != "" {
		return []types.Event{metaAlarmEvent}, nil
	}

	return nil, nil
}

// AlarmMatched checks alarm attributes agiainst the AttributePatterns in rule configuration
func (a AttributeApplicator) AlarmMatched(event types.Event, rule correlation.Rule) bool {
	patternsMatch := rule.Config.AlarmPatterns.Matches(event.Alarm)
	a.logger.Debug().Msgf("Alarm matched event %+v with rule %v %t", event, rule.Config.AlarmPatterns.AsMongoDriverQuery(), patternsMatch)
	return patternsMatch
}

// EntityMatched checks entity attributes agiainst the EntityPatterns in rule configuration
func (a AttributeApplicator) EntityMatched(event types.Event, rule correlation.Rule) bool {
	patternsMatch := rule.Config.EntityPatterns.Matches(event.Entity)
	a.logger.Debug().Msgf("Entity matched event %+v with rule %v %t", event, rule.Config.EntityPatterns.AsMongoDriverQuery(), patternsMatch)
	return patternsMatch
}

// EventMatched checks event attributes agiainst the EventPatterns in rule configuration
func (a AttributeApplicator) EventMatched(event types.Event, rule correlation.Rule) bool {
	patternsMatch := rule.Config.EventPatterns.Matches(event)
	a.logger.Debug().Msgf("Event matched event %+v with rule %v", event, rule.Config.EventPatterns)
	return patternsMatch
}

// NewAttributeApplicator instantiates AttributeApplicator with MetaAlarmService
func NewAttributeApplicator(alarmAdapter alarm.Adapter, logger zerolog.Logger, metaAlarmService service.MetaAlarmService, redisClient *redis.Client) AttributeApplicator {
	return AttributeApplicator{
		alarmAdapter: alarmAdapter,
		service:      metaAlarmService,
		storage:      storage.NewRedisGroupingStorage(),
		redisClient:  redisClient,
		logger:       logger,
	}
}
