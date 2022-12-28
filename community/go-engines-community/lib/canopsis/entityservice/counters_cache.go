package entityservice

import (
	"context"
	"fmt"
	"regexp"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

// maxRetries is the maximum number of times a redis WATCH transaction will be
// retried before aborting.
const maxRetries = 10

// counterName* are the name of the keys used to store each field of an
// AlarmCounter in redis.
const (
	cacheKeyTpl                     = "service:counters:%s:%s"
	counterNameAll                  = "all"
	counterNameActive               = "active"
	counterNameStateCritical        = "state:critical"
	counterNameStateMajor           = "state:major"
	counterNameStateMinor           = "state:minor"
	counterNameStateOk              = "state:ok"
	counterNameAcknowledged         = "acknowledged"
	counterNameNotAcknowledged      = "not_acknowledged"
	counterNameAcknowledgedUnderPbh = "acknowledged_under_pbh"
	counterNamePbehaviorCounters    = "pbehavior_counters"
)

// redisCounterKey returns the name of the key used to store a field of a
// service's AlarmCounters in redis.
func redisCounterKey(counterName, serviceID string) string {
	return fmt.Sprintf(cacheKeyTpl, counterName, serviceID)
}

// countersCache is a type that implements the CountersCache interface.
type countersCache struct {
	client  redis.UniversalClient
	encoder encoding.Encoder
	decoder encoding.Decoder
	logger  zerolog.Logger
}

// NewCountersCache creates a new CountersCache
func NewCountersCache(client redis.UniversalClient, logger zerolog.Logger) CountersCache {
	return &countersCache{
		client: client,

		// The DependencyState and AlarmCounters values are stored as JSON into
		// the cache because :
		//  - if fields are added to or removed from these structs, we will
		//    still be able to decode old values stored in the cache (this
		//    might not be the case with the gob format).
		//  - having human-readable data in cache will make debugging easier.
		decoder: json.NewDecoder(),
		encoder: json.NewEncoder(),

		logger: logger,
	}
}

// incrementCounters increments the counters of a service with the values of
// incrementBy. It returns the results of the incr commands in an
// alarmCountersCmds struct, so that the values of the counters can be read
// after the pipeline has been executed.
func (s *countersCache) incrementCounters(
	ctx context.Context,
	pipe redis.Pipeliner,
	serviceID string,
	incrementBy AlarmCounters,
) alarmCountersCmds {
	pbehaviorCounters := make(map[string]*redis.IntCmd)
	for pbhType, counter := range incrementBy.PbehaviorCounters {
		pbehaviorCounters[pbhType] = pipe.HIncrBy(
			ctx,
			redisCounterKey(counterNamePbehaviorCounters, serviceID),
			pbhType,
			counter,
		)
	}

	return alarmCountersCmds{
		All: pipe.IncrBy(
			ctx,
			redisCounterKey(counterNameAll, serviceID),
			incrementBy.All),
		Active: pipe.IncrBy(
			ctx,
			redisCounterKey(counterNameActive, serviceID),
			incrementBy.Active),
		State: stateCountersCmds{
			Critical: pipe.IncrBy(
				ctx,
				redisCounterKey(counterNameStateCritical, serviceID),
				incrementBy.State.Critical),
			Major: pipe.IncrBy(
				ctx,
				redisCounterKey(counterNameStateMajor, serviceID),
				incrementBy.State.Major),
			Minor: pipe.IncrBy(
				ctx,
				redisCounterKey(counterNameStateMinor, serviceID),
				incrementBy.State.Minor),
			Ok: pipe.IncrBy(
				ctx,
				redisCounterKey(counterNameStateOk, serviceID),
				incrementBy.State.Ok),
		},
		Acknowledged: pipe.IncrBy(
			ctx,
			redisCounterKey(counterNameAcknowledged, serviceID),
			incrementBy.Acknowledged),
		NotAcknowledged: pipe.IncrBy(
			ctx,
			redisCounterKey(counterNameNotAcknowledged, serviceID),
			incrementBy.NotAcknowledged),
		AcknowledgedUnderPbh: pipe.IncrBy(
			ctx,
			redisCounterKey(counterNameAcknowledgedUnderPbh, serviceID),
			incrementBy.AcknowledgedUnderPbh),
		PbehaviorCounters: pbehaviorCountersCmd{
			All:  pipe.HGetAll(ctx, redisCounterKey(counterNamePbehaviorCounters, serviceID)),
			Incr: pbehaviorCounters,
		},
	}
}

func (s *countersCache) removeCounters(
	ctx context.Context,
	pipe redis.Pipeliner,
	serviceID string,
) *redis.IntCmd {
	return pipe.Del(
		ctx,
		redisCounterKey(counterNameAll, serviceID),
		redisCounterKey(counterNameActive, serviceID),
		redisCounterKey(counterNameStateCritical, serviceID),
		redisCounterKey(counterNameStateMajor, serviceID),
		redisCounterKey(counterNameStateMinor, serviceID),
		redisCounterKey(counterNameStateOk, serviceID),
		redisCounterKey(counterNameAcknowledged, serviceID),
		redisCounterKey(counterNameNotAcknowledged, serviceID),
		redisCounterKey(counterNameAcknowledgedUnderPbh, serviceID),
		redisCounterKey(counterNamePbehaviorCounters, serviceID),
	)
}

func (s countersCache) getCounters(
	ctx context.Context,
	pipe redis.Cmdable,
	serviceID string,
) getAlarmCountersCmds {
	return getAlarmCountersCmds{
		All:    pipe.Get(ctx, redisCounterKey(counterNameAll, serviceID)),
		Active: pipe.Get(ctx, redisCounterKey(counterNameActive, serviceID)),
		State: getStateCountersCmds{
			Critical: pipe.Get(ctx, redisCounterKey(counterNameStateCritical, serviceID)),
			Major:    pipe.Get(ctx, redisCounterKey(counterNameStateMajor, serviceID)),
			Minor:    pipe.Get(ctx, redisCounterKey(counterNameStateMinor, serviceID)),
			Ok:       pipe.Get(ctx, redisCounterKey(counterNameStateOk, serviceID)),
		},
		Acknowledged:         pipe.Get(ctx, redisCounterKey(counterNameAcknowledged, serviceID)),
		NotAcknowledged:      pipe.Get(ctx, redisCounterKey(counterNameNotAcknowledged, serviceID)),
		AcknowledgedUnderPbh: pipe.Get(ctx, redisCounterKey(counterNameAcknowledgedUnderPbh, serviceID)),
		PbehaviorCounters: getPbehaviorCountersCmd{
			All: pipe.HGetAll(ctx, redisCounterKey(counterNamePbehaviorCounters, serviceID)),
		},
	}
}

func (s *countersCache) Update(ctx context.Context, update map[string]AlarmCounters) (map[string]AlarmCounters, error) {
	countersCmd := map[string]alarmCountersCmds{}

	for retries := maxRetries; retries >= 0; retries-- {
		err := s.client.Watch(ctx, func(tx *redis.Tx) error {
			pipe := tx.TxPipeline()

			countersCmd = map[string]alarmCountersCmds{}
			for serviceID, increment := range update {
				countersCmd[serviceID] = s.incrementCounters(
					ctx, pipe, serviceID, increment)
			}

			_, err := pipe.Exec(ctx)
			return err
		})

		if err == redis.TxFailedErr {
			// The WATCH failed because the entity's state was modified by another
			// process or goroutine. Try again, with the new entity's state.
			s.logger.Warn().
				Int("remaining_retries", retries).
				Str("help", "The state of this entity will be temporarily out of date in the services. This will fix by itself, and should not be a problem has long as it does not happen often.").
				Msg("failed to update entity's state")
			continue
		} else if err != nil {
			// The function called in Watch returned an error other than a
			// transaction failure. We do not need to retry in this case.
			return nil, err
		}

		// The results of the increment commands need to be read here (after
		// the pipe has been executed).
		counters := map[string]AlarmCounters{}
		for serviceID, cmd := range countersCmd {
			serviceCounters, err := cmd.Result()
			if err != nil {
				s.logger.Error().
					Err(err).
					Str("service_id", serviceID).
					Str("help", "The redis cache seems to be corrupted. This may be fixed by flushing the cache and restarting the engines.").
					Msg("unable to read AlarmCounters")
			}
			counters[serviceID] = serviceCounters
		}
		return counters, nil
	}

	return nil, fmt.Errorf("failed to process entity's state ")
}

func (s *countersCache) Replace(ctx context.Context, serviceID string, counters AlarmCounters) error {
	cmd := alarmCountersCmds{}

	for retries := maxRetries; retries >= 0; retries-- {
		err := s.client.Watch(ctx, func(tx *redis.Tx) error {
			pipe := tx.TxPipeline()

			s.removeCounters(ctx, pipe, serviceID)
			cmd = s.incrementCounters(ctx, pipe, serviceID, counters)

			_, err := pipe.Exec(ctx)
			return err
		})

		if err == redis.TxFailedErr {
			s.logger.Warn().
				Int("remaining_retries", retries).
				Str("help", "The state of this entity will be temporarily out of date in the services. This will fix by itself, and should not be a problem has long as it does not happen often.").
				Msg("failed to update entity's state")
			continue
		} else if err != nil {
			return err
		}

		// The results of the increment commands need to be read here (after
		// the pipe has been executed).
		_, err = cmd.Result()
		if err != nil {
			s.logger.Error().
				Err(err).
				Str("service_id", serviceID).
				Str("help", "The redis cache seems to be corrupted. This may be fixed by flushing the cache and restarting the engines.").
				Msg("unable to read AlarmCounters")
		}
		return nil
	}

	return fmt.Errorf("failed to process entity's state ")
}

func (s *countersCache) RemoveAndGet(ctx context.Context, serviceID string) (*AlarmCounters, error) {
	cmd := getAlarmCountersCmds{}

	for retries := maxRetries; retries >= 0; retries-- {
		err := s.client.Watch(ctx, func(tx *redis.Tx) error {
			pipe := tx.TxPipeline()

			cmd = s.getCounters(ctx, pipe, serviceID)
			s.removeCounters(ctx, pipe, serviceID)

			_, err := pipe.Exec(ctx)
			return err
		})

		if err == redis.TxFailedErr {
			s.logger.Warn().
				Int("remaining_retries", retries).
				Str("help", "The state of this entity will be temporarily out of date in the services. This will fix by itself, and should not be a problem has long as it does not happen often.").
				Msg("failed to update entity's state")
			continue
		} else if err != nil {
			if err == redis.Nil {
				return nil, nil
			}

			return nil, err
		}

		// The results of the increment commands need to be read here (after
		// the pipe has been executed).
		res, err := cmd.Result()
		if err != nil {
			if err == redis.Nil {
				return nil, nil
			}

			s.logger.Error().
				Err(err).
				Str("service_id", serviceID).
				Str("help", "The redis cache seems to be corrupted. This may be fixed by flushing the cache and restarting the engines.").
				Msg("unable to read AlarmCounters")
			return nil, err
		}

		return &res, nil
	}

	return nil, fmt.Errorf("failed to process entity's state ")
}

func (s *countersCache) Remove(ctx context.Context, serviceID string) error {
	for retries := maxRetries; retries >= 0; retries-- {
		err := s.client.Watch(ctx, func(tx *redis.Tx) error {
			pipe := tx.TxPipeline()

			s.removeCounters(ctx, pipe, serviceID)

			_, err := pipe.Exec(ctx)
			return err
		})

		if err == redis.TxFailedErr {
			s.logger.Warn().
				Int("remaining_retries", retries).
				Str("help", "The state of this entity will be temporarily out of date in the services. This will fix by itself, and should not be a problem has long as it does not happen often.").
				Msg("failed to update entity's state")
			continue
		} else if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("failed to process entity's state ")
}

func (s *countersCache) ClearAll(ctx context.Context) error {
	match := fmt.Sprintf(cacheKeyTpl, "*", "*")
	var cursor uint64

	for {
		res := s.client.Scan(ctx, cursor, match, 50)
		if err := res.Err(); err != nil {
			return err
		}

		var keys []string
		keys, cursor = res.Val()
		if len(keys) == 0 {
			break
		}

		err := s.client.Del(ctx, keys...).Err()
		if err != nil {
			return err
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (s *countersCache) KeepOnly(ctx context.Context, ids []string) error {
	match := fmt.Sprintf(cacheKeyTpl, "*", "*")
	var cursor uint64
	r, err := regexp.Compile("^" + fmt.Sprintf(cacheKeyTpl, "(.+)", "(.+)"))
	if err != nil {
		return err
	}
	remove := make([]string, 0)
	for {
		res := s.client.Scan(ctx, cursor, match, 50)
		if err := res.Err(); err != nil {
			return err
		}

		var keys []string
		keys, cursor = res.Val()
		if len(keys) == 0 {
			break
		}

		for _, k := range keys {
			m := r.FindStringSubmatch(k)
			if len(m) >= 3 {
				serviceID := m[2]
				found := false
				for _, v := range ids {
					if v == serviceID {
						found = true
						break
					}
				}

				if !found {
					remove = append(remove, k)
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	if len(remove) > 0 {
		err := s.client.Del(ctx, remove...).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
