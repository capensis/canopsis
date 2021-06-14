package watcher

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/rs/zerolog"
)

// expiration is the expire time used in the SET operations.
// The expiration is set to 0, so that the values never expire.
const expiration = 0

// maxRetries is the maximum number of times a redis WATCH transaction will be
// retried before aborting.
const maxRetries = 10

// counterName* are the name of the keys used to store each field of an
// AlarmCounter in redis.
const (
	counterNameAll               = "all"
	counterNameAlarms            = "alarms"
	counterNameStateCritical     = "state:critical"
	counterNameStateMajor        = "state:major"
	counterNameStateMinor        = "state:minor"
	counterNameStateInfo         = "state:info"
	counterNameAcknowledged      = "acknowledged"
	counterNameNotAcknowledged   = "not_acknowledged"
	counterNamePbehaviorCounters = "pbehavior_counters"
)

// redisDependencyStateKey returns the name of the key used to store a
// DependencyState in redis.
func redisDependencyStateKey(entityID string) string {
	return fmt.Sprintf("watcher:dependency:%s", entityID)
}

// redisCounterKey returns the name of the key used to store a field of a
// watcher's AlarmCounters in redis.
func redisCounterKey(counterName, watcherID string) string {
	return fmt.Sprintf("watcher:counters:%s:%s", counterName, watcherID)
}

// countersCache is a type that implements the CountersCache interface.
type countersCache struct {
	client  *redis.Client
	encoder encoding.Encoder
	decoder encoding.Decoder
	logger  zerolog.Logger
}

// NewCountersCache creates a new CountersCache
func NewCountersCache(client *redis.Client, logger zerolog.Logger) CountersCache {
	return countersCache{
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

// saveDependencyState writes a dependency state to redis.
func (s countersCache) saveDependencyState(
	pipe redis.Cmdable,
	state DependencyState,
) error {
	encodedState, err := s.encoder.Encode(state)
	if err != nil {
		return err
	}

	return pipe.Set(
		redisDependencyStateKey(state.EntityID),
		encodedState,
		expiration).Err()
}

// getDependencyState returns the dependency state currently stored in redis
// for an entity.
// It returns the zero value of DependencyState if the value does not exist in
// redis.
func (s countersCache) getDependencyState(
	pipe redis.Cmdable,
	entityID string,
) (DependencyState, error) {
	resp := pipe.Get(redisDependencyStateKey(entityID))
	if resp.Err() == redis.Nil {
		return DependencyState{}, nil
	} else if resp.Err() != nil {
		return DependencyState{}, resp.Err()
	}

	encodedState, err := resp.Bytes()
	if err != nil {
		return DependencyState{}, err
	}

	state := DependencyState{}
	err = s.decoder.Decode(encodedState, &state)
	return state, err
}

// incrementCounters increments the counters of a watcher with the values of
// incrementBy. It returns the results of the incr commands in an
// alarmCountersCmds struct, so that the values of the counters can be read
// after the pipeline has been executed.
func (s countersCache) incrementCounters(
	pipe redis.Pipeliner,
	watcherID string,
	incrementBy AlarmCounters,
) alarmCountersCmds {
	pbehaviorCounters := make(map[string]*redis.IntCmd)
	for pbhType, counter := range incrementBy.PbehaviorCounters {
		pbehaviorCounters[pbhType] = pipe.HIncrBy(
			redisCounterKey(counterNamePbehaviorCounters, watcherID),
			pbhType,
			counter,
		)
	}

	return alarmCountersCmds{
		All: pipe.IncrBy(
			redisCounterKey(counterNameAll, watcherID),
			incrementBy.All),
		Alarms: pipe.IncrBy(
			redisCounterKey(counterNameAlarms, watcherID),
			incrementBy.Alarms),
		State: stateCountersCmds{
			Critical: pipe.IncrBy(
				redisCounterKey(counterNameStateCritical, watcherID),
				incrementBy.State.Critical),
			Major: pipe.IncrBy(
				redisCounterKey(counterNameStateMajor, watcherID),
				incrementBy.State.Major),
			Minor: pipe.IncrBy(
				redisCounterKey(counterNameStateMinor, watcherID),
				incrementBy.State.Minor),
			Info: pipe.IncrBy(
				redisCounterKey(counterNameStateInfo, watcherID),
				incrementBy.State.Info),
		},
		Acknowledged: pipe.IncrBy(
			redisCounterKey(counterNameAcknowledged, watcherID),
			incrementBy.Acknowledged),
		NotAcknowledged: pipe.IncrBy(
			redisCounterKey(counterNameNotAcknowledged, watcherID),
			incrementBy.NotAcknowledged),
		PbehaviorCounters: pbehaviorCountersCmd{
			All:  pipe.HGetAll(redisCounterKey(counterNamePbehaviorCounters, watcherID)),
			Incr: pbehaviorCounters,
		},
	}
}

// ProcessState processes an entity's state, and update the counters of the
// watchers that are (or used to be) impacted by it.
// It returns a map containing the impacted watchers and the new values of
// their counters.
// This method can be run safely in multiple goroutines.
//
// This method works as follows :
//
//  1. WATCH the redis key containing the entity's state.
//  2. Get the previous entity's state in redis, and ensure that it is older
//     than the state currently being processed.
//  3. Write the current entity's state to redis.
//  4. Update the counters of the watchers that are impacted by the entity's
//     change.
//     For example, if the entity has an alarm that was previously major and is
//     now critical, the counter of major alarms of each watcher is
//     decremented, and the counter of critical alarms of each watcher is
//     incremented.
//  5. Run the redis transaction. The transaction will fail if the redis key
//     containing the entity's state was changed by another process. In this
//     case, the "previous entity's state" returned in step 2 is outdated, so
//     we go back to step 1.
//  6. Return the counters of the watchers that have been modified.
//
// The transaction and the WATCH command ensure that the entity states are
// consistent with the watchers' alarm counters.
//
// For more informations on transactions and WATCH, read
// https://redis.io/topics/transactions#optimistic-locking-using-check-and-set
func (s countersCache) ProcessState(
	currentState DependencyState,
) (map[string]AlarmCounters, error) {
	increments := map[string]AlarmCounters{}
	countersCmd := map[string]alarmCountersCmds{}

	for retries := maxRetries; retries >= 0; retries-- {
		err := s.client.Watch(func(tx *redis.Tx) error {
			// Get the previous entity's state
			previousState, err := s.getDependencyState(tx, currentState.EntityID)
			if err != nil {
				return err
			}

			if previousState.LastUpdateDate.After(currentState.LastUpdateDate) {
				return nil
			}

			// Create a transaction that will be used to update the watchers'
			// counters and update the entity state.
			// The following redis commands will only be committed if the
			// transaction succeeds and the entity's state has not been
			// modified by another process.
			pipe := tx.TxPipeline()

			// Write the current entity's state
			err = s.saveDependencyState(pipe, currentState)
			if err != nil {
				return err
			}

			// Update the impacted watchers' counters.
			increments = GetCountersIncrementsFromStates(
				previousState, currentState)
			countersCmd = map[string]alarmCountersCmds{}
			for watcherID, increment := range increments {
				countersCmd[watcherID] = s.incrementCounters(
					pipe, watcherID, increment)
			}

			_, err = pipe.Exec()
			return err
		}, redisDependencyStateKey(currentState.EntityID))

		if err == redis.TxFailedErr {
			// The WATCH failed because the entity's state was modified by another
			// process or goroutine. Try again, with the new entity's state.
			s.logger.Warn().
				Str("entity_id", currentState.EntityID).
				Int("remaining_retries", retries).
				Str("help", "The state of this entity will be temporarily out of date in the watchers. This will fix by itself, and should not be a problem has long as it does not happen often.").
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
		for watcherID, cmd := range countersCmd {
			watcherCounters, err := cmd.Result()
			if err != nil {
				s.logger.Error().
					Err(err).
					Str("watcher_id", watcherID).
					Str("help", "The redis cache seems to be corrupted. This may be fixed by flushing the cache and restarting the engines.").
					Msg("unable to read AlarmCounters")
			}
			counters[watcherID] = watcherCounters
		}
		return counters, nil
	}

	return nil, fmt.Errorf(
		"failed to process entity's state : entity_id=%s",
		currentState.EntityID)
}
