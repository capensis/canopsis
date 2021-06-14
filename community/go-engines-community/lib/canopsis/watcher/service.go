package watcher

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"runtime/trace"
	"sync"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/errt"
	"github.com/globalsign/mgo/bson"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

type service struct {
	redisClient 	*redis.Client
	pubChannel      AmqpChannelPublisher
	pubExchangeName string
	pubQueueName    string
	jsonEncoder     encoding.Encoder
	watcherAdapter  Adapter
	alarmAdapter    alarm.Adapter
	countersCache   CountersCache
	logger          zerolog.Logger

	watchersLock sync.RWMutex
	watchers     map[string]Watcher
}

// NewService gives the correct watcher adapter.
func NewService(
	redisClient *redis.Client,
	pubChannel AmqpChannelPublisher,
	pubExchangeName, pubQueueName string,
	jsonEncoder encoding.Encoder,
	watcherAdapter Adapter,
	alarmAdapter alarm.Adapter,
	countersCache CountersCache,
	logger zerolog.Logger,
) Service {
	service := service{
		redisClient:	 redisClient,
		pubChannel:      pubChannel,
		pubExchangeName: pubExchangeName,
		pubQueueName:    pubQueueName,
		jsonEncoder:     jsonEncoder,
		watcherAdapter:  watcherAdapter,
		alarmAdapter:    alarmAdapter,
		countersCache:   countersCache,
		logger:          logger,
	}
	return &service
}

// sendEvent sends a watcher event.
func (s *service) sendEvent(event types.Event) error {
	jevt, err := s.jsonEncoder.Encode(event)
	if err != nil {
		return fmt.Errorf("unable to serialize watcher event: %v", err)
	}

	err = s.pubChannel.Publish(
		s.pubExchangeName,
		s.pubQueueName,
		false,
		false,
		amqp.Publishing{
			Body:        jevt,
			ContentType: "application/json",
		},
	)
	if err != nil {
		return fmt.Errorf("unable to send watcher event: %v", err)
	}

	return nil
}

// updateWatcherState computes the state of a watcher given its AlarmCounters,
// and sends an event to update the corresponding alarm.
func (s *service) updateWatcherState(
	watcher Watcher,
	counters AlarmCounters,
) error {
	if !watcher.Enabled {
		return nil
	}

	output, err := watcher.GetOutput(counters)
	if err != nil {
		return err
	}

	state, err := watcher.GetState(counters)
	if err != nil {
		return err
	}

	err = s.sendEvent(types.Event{
		EventType:     types.EventTypeCheck,
		SourceType:    types.SourceTypeComponent,
		Component:     watcher.ID,
		Connector:     "watcher",
		ConnectorName: "watcher",
		State:         types.CpsNumber(state),
		Output:        output,
		Timestamp:     types.CpsTime{Time: time.Now()},
	})
	if err != nil {
		return err
	}

	return s.watcherAdapter.Update(watcher.ID, bson.M{
		"$set": bson.M{
			"alarms_cumulative_data.watched_count":           counters.All,
			"alarms_cumulative_data.watched_pbehavior_count": counters.PbehaviorCounters,
			"alarms_cumulative_data.watched_not_acked_count": counters.NotAcknowledged,
		},
	})
}

// updateWatchersState computes the states of watchers given their
// AlarmCounters, and sends events to update the corresponding alarms.
func (s *service) updateWatchersState(
	watcherCounters map[string]AlarmCounters,
) {
	s.watchersLock.RLock()
	defer s.watchersLock.RUnlock()

	for watcherID, counters := range watcherCounters {
		watcher, exists := s.watchers[watcherID]
		if !exists {
			log.Printf("no such watcher : %s", watcherID)
		}

		err := s.updateWatcherState(watcher, counters)
		if err != nil {
			log.Printf("Unable to update watcher state %+v", err)
		}
	}
}

// Process processes an event and updates the watchers impacted by the event.
func (s *service) Process(ctx context.Context, event types.Event) error {
	defer trace.StartRegion(ctx, "watcher.Process").End()

	if event.AlarmChange == nil {
		return errt.NewUnknownError(fmt.Errorf("event's alarm_change is nil : %v", event))
	}

	if event.AlarmChange.Type != types.AlarmChangeTypeAck &&
		event.AlarmChange.Type != types.AlarmChangeTypeAckremove &&
		event.AlarmChange.Type != types.AlarmChangeTypeStateIncrease &&
		event.AlarmChange.Type != types.AlarmChangeTypeStateDecrease &&
		event.AlarmChange.Type != types.AlarmChangeTypePbhEnter &&
		event.AlarmChange.Type != types.AlarmChangeTypePbhLeave &&
		event.AlarmChange.Type != types.AlarmChangeTypePbhLeaveAndEnter &&
		event.AlarmChange.Type != types.AlarmChangeTypeCreateAndPbhEnter &&
		event.AlarmChange.Type != types.AlarmChangeTypeCreate {
		return nil
	}

	return s.calculateState(event)
}

func (s *service) ProcessRpc(ctx context.Context, event types.Event) error {
	defer trace.StartRegion(ctx, "watcher.ProcessRpc").End()

	return s.calculateState(event)
}

func (s *service) calculateState(event types.Event) error {
	s.watchersLock.RLock()
	defer s.watchersLock.RUnlock()

	if event.Entity == nil {
		s.logger.Warn().Msgf("event's entity is nil : %+v", event)
		return nil
	}

	// Create a DependencyState (that contains all the information required to
	// update the watchers' states) from the event
	pbehaviorType := event.PbehaviorInfo.TypeID
	dependencyState := NewDependencyState(
		*event.Entity,
		event.Alarm,
		event.PbehaviorInfo.IsActive(),
		pbehaviorType,
		s.watchers,
		time.Now(),
	)
	if len(dependencyState.ImpactedWatchers) == 0 {
		return nil
	}

	// Process this DependencyState, and get the new values of the impacted
	// watchers' AlarmCounters
	impactCounters, err := s.countersCache.ProcessState(dependencyState)
	if err != nil {
		return fmt.Errorf("unable to process state : %+v", err)
	}

	// Compute the new watcher's state, and send an event to update the alarm
	s.updateWatchersState(impactCounters)

	return nil
}

func (s *service) UpdateWatcher(ctx context.Context, watcherID string) error {
	defer trace.StartRegion(ctx, "watcher.UpdateWatcher").End()

	// Update the watchers map
	var watchers []Watcher
	err := s.watcherAdapter.GetAllValidWatchers(&watchers)
	if err != nil {
		return err
	}

	s.watchersLock.Lock()
	s.watchers = map[string]Watcher{}
	for _, watcher := range watchers {
		s.watchers[watcher.ID] = watcher
	}
	s.watchersLock.Unlock()

	// Get the watcher's dependencies, with their alarms and pbehaviors
	iter := s.watcherAdapter.GetAnnotatedDependenciesIter(watcherID)
	if err := iter.Err(); err != nil {
		return err
	}

	now := time.Now()
	watcherCounters := map[string]AlarmCounters{}
	var entity AnnotatedEntity
	for iter.Next(&entity) {
		// Create a DependencyState (that contains all the information required
		// to update the watchers' states) from the AnnotatedEntity
		isActive := entity.Alarm != nil && entity.Alarm.IsInActivePeriod()
		pbehaviorType := ""
		if entity.Alarm != nil {
			pbehaviorType = entity.Alarm.Value.PbehaviorInfo.TypeID
		}
		s.watchersLock.RLock()
		dependencyState := NewDependencyState(
			entity.Entity, entity.Alarm, isActive, pbehaviorType, s.watchers, now)
		s.watchersLock.RUnlock()
		if len(dependencyState.ImpactedWatchers) == 0 {
			continue
		}

		// Process this DependencyState, and get the new values of the impacted
		// watchers' AlarmCounters
		impactCounters, err := s.countersCache.ProcessState(dependencyState)
		if err != nil {
			s.logger.Error().Err(err).Msg("Unable to process state")
		}

		// The watchers may be impacted by multiple entities. Only keep the
		// latest value of the AlarmCounters for each watcher, so that they are
		// only updated once.
		for watcherID, counters := range impactCounters {
			watcherCounters[watcherID] = counters
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}

	if err := iter.Close(); err != nil {
		return err
	}

	// Compute the new watchers' state, and send events to update the alarms
	s.updateWatchersState(watcherCounters)

	return nil
}

func (s *service) ProcessResolvedAlarm(ctx context.Context, alarm types.Alarm, entity types.Entity) error {
	defer trace.StartRegion(ctx, "watcher.ProcessResolvedAlarm").End()
	s.watchersLock.RLock()
	defer s.watchersLock.RUnlock()

	now := time.Now()
	watcherCounters := map[string]AlarmCounters{}

	pbehaviorType := alarm.Value.PbehaviorInfo.TypeID
	dependencyState := NewDependencyState(
		entity, &alarm, alarm.IsInActivePeriod(), pbehaviorType, s.watchers, now)
	if len(dependencyState.ImpactedWatchers) == 0 {
		return nil
	}

	// Process this DependencyState, and get the new values of the impacted
	// watchers' AlarmCounters
	impactCounters, err := s.countersCache.ProcessState(dependencyState)
	if err != nil {
		log.Printf("Unable to process state : %+v", err)
	}

	// The watchers may be impacted by multiple resolved alarms. Only keep
	// the latest value of the AlarmCounters for each watcher, so that they
	// are only updated once.
	for watcherID, counters := range impactCounters {
		watcherCounters[watcherID] = counters
	}

	// Compute the new watchers' state, and send events to update the alarms
	s.updateWatchersState(watcherCounters)

	return nil
}

func (s *service) ComputeAllWatchers(ctx context.Context) error {
	defer trace.StartRegion(ctx, "watcher.ComputeAllWatchers").End()

	// Update the watchers map
	var watchers []Watcher
	err := s.watcherAdapter.GetAllValidWatchers(&watchers)
	if err != nil {
		return err
	}

	s.watchersLock.Lock()
	s.watchers = map[string]Watcher{}
	for _, watcher := range watchers {
		s.watchers[watcher.ID] = watcher
	}
	s.watchersLock.Unlock()

	if len(watchers) == 0 {
		return nil
	}

	// Get all the entities, with their alarms and pbehaviors
	now := time.Now()
	iter := s.watcherAdapter.GetAnnotatedEntitiesIter()
	if err := iter.Err(); err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	watcherCounters := map[string]AlarmCounters{}
	var countersMX sync.Mutex

	var entity AnnotatedEntity
	for iter.Next(&entity) {
		wg.Add(1)

		go func(e AnnotatedEntity) {
			defer wg.Done()
			// Create a DependencyState (that contains all the information required
			// to update the watchers' states) from the AnnotatedEntity
			s.watchersLock.RLock()
			isActive := e.Alarm != nil && e.Alarm.IsInActivePeriod()
			pbehaviorType := ""
			if e.Alarm != nil {
				pbehaviorType = e.Alarm.Value.PbehaviorInfo.TypeID
			}
			dependencyState := NewDependencyState(
				e.Entity, e.Alarm, isActive, pbehaviorType, s.watchers, now)
			s.watchersLock.RUnlock()
			if len(dependencyState.ImpactedWatchers) == 0 {
				return
			}

			// Process this DependencyState, and get the new values of the impacted
			// watchers' AlarmCounters
			impactCounters, err := s.countersCache.ProcessState(dependencyState)
			if err != nil {
				s.logger.Error().Err(err).Msg("Unable to process state")
			}

			// The watchers may be impacted by multiple entities. Only keep the
			// latest value of the AlarmCounters for each watcher, so that they are
			// only updated once.
			for watcherID, counters := range impactCounters {
				countersMX.Lock()
				watcherCounters[watcherID] = counters
				countersMX.Unlock()
			}
		}(entity)
	}

	wg.Wait()

	if err := iter.Err(); err != nil {
		return err
	}

	if err := iter.Close(); err != nil {
		return err
	}

	// Compute the new watchers' state, and send events to update the alarms
	s.updateWatchersState(watcherCounters)

	return nil
}

func (s *service) FlushDB() error {
	st := s.redisClient.FlushDB()
	if st.Err() != nil {
		return st.Err()
	}

	return nil
}
