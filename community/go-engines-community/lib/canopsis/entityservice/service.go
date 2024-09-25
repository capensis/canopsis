package entityservice

import (
	"context"
	"fmt"
	"math"
	"runtime/trace"
	"sync"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

const (
	lockServiceKey       = "lock-service-"
	lockServiceTtl       = time.Second
	lockServiceUpdateKey = "lock-service-update-"
	lockServiceUpdateTtl = time.Minute
	skippedQueueKey      = "skipped-entities-"
)

const BulkMaxSize = 10000

type service struct {
	pubChannel      libamqp.Publisher
	pubExchangeName string
	pubQueueName    string
	encoder         encoding.Encoder
	adapter         Adapter
	entityAdapter   entity.Adapter
	alarmAdapter    libalarm.Adapter
	countersCache   CountersCache
	storage         Storage
	lockClient      libredis.LockClient
	redisClient     redis.Cmdable
	logger          zerolog.Logger

	templateExecutor template.Executor
}

// NewService gives the correct service adapter.
func NewService(
	pubChannel libamqp.Publisher,
	pubExchangeName, pubQueueName string,
	encoder encoding.Encoder,
	adapter Adapter,
	entityAdapter entity.Adapter,
	alarmAdapter libalarm.Adapter,
	countersCache CountersCache,
	storage Storage,
	lockClient libredis.LockClient,
	redisClient redis.Cmdable,
	templateExecutor template.Executor,
	logger zerolog.Logger,
) Service {
	service := service{
		pubChannel:      pubChannel,
		pubExchangeName: pubExchangeName,
		pubQueueName:    pubQueueName,
		encoder:         encoder,
		adapter:         adapter,
		entityAdapter:   entityAdapter,
		alarmAdapter:    alarmAdapter,
		countersCache:   countersCache,
		storage:         storage,
		lockClient:      lockClient,
		redisClient:     redisClient,
		logger:          logger,

		templateExecutor: templateExecutor,
	}
	return &service
}

func (s *service) sendEvent(ctx context.Context, event types.Event) error {
	body, err := s.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("unable to serialize service event: %v", err)
	}

	err = s.pubChannel.PublishWithContext(
		ctx,
		s.pubExchangeName,
		s.pubQueueName,
		false,
		false,
		amqp.Publishing{
			Body:        body,
			ContentType: "application/json",
		},
	)
	if err != nil {
		return fmt.Errorf("unable to send service event: %v", err)
	}

	return nil
}

// updateServiceState computes the state of a service given its AlarmCounters,
// and sends an event to update the corresponding alarm.
func (s *service) updateServiceState(
	ctx context.Context,
	serviceID, serviceOutput string,
	counters AlarmCounters,
) error {
	count, err := s.adapter.GetDependenciesCount(ctx, serviceID)
	if err != nil {
		return err
	}

	counters.Depends = count
	counters.UnderPbehavior = 0
	for _, c := range counters.PbehaviorCounters {
		counters.UnderPbehavior += c
	}
	output, err := s.templateExecutor.Execute(serviceOutput, counters)
	if err != nil {
		return err
	}

	err = s.sendEvent(ctx, types.Event{
		EventType:     types.EventTypeCheck,
		SourceType:    types.SourceTypeService,
		Component:     serviceID,
		Connector:     types.ConnectorEngineService,
		ConnectorName: types.ConnectorEngineService,
		State:         types.CpsNumber(GetServiceState(counters)),
		Output:        output,
		Timestamp:     types.CpsTime{Time: time.Now()},
	})
	if err != nil {
		return err
	}

	err = s.adapter.UpdateCounters(ctx, serviceID, counters)
	if err != nil {
		s.logger.Error().Err(err).Msg("unable to update service counters")
		return err
	}

	return nil
}

// Process processes an event and updates the services impacted by the event.
func (s *service) Process(ctx context.Context, event types.Event) error {
	defer trace.StartRegion(ctx, "service.Process").End()
	if event.Entity == nil {
		s.logger.Warn().Msgf("event's entity is nil : %+v", event)
		return nil
	}
	if event.AlarmChange == nil {
		return errt.NewUnknownError(fmt.Errorf("event's alarm_change is nil : %v", event))
	}

	err := s.calculateState(ctx, event)
	if err != nil {
		return err
	}

	if event.EventType == types.EventTypeResolveDeleted && event.Entity.SoftDeleted != nil {
		err = s.entityAdapter.SetResolveDeletedEventProcessed(ctx, event.Entity.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) markServices(parentCtx context.Context, idleSinceMap *ServicesIdleSinceMap, services []EntityService, impacts []string, timestamp int64) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	wg := sync.WaitGroup{}
	workerCh := make(chan string)
	defer close(workerCh)

	wg.Add(len(impacts))
	go func() {
		for _, impact := range impacts {
			workerCh <- impact
		}
	}()

	for i := 0; i < maxWorkersCount; i++ {
		go func() {
			for impact := range workerCh {
				func() {
					defer wg.Done()

					select {
					case <-ctx.Done():
						return
					default:
					}

					if !idleSinceMap.Mark(impact, timestamp) {
						return
					}

					for _, service := range services {
						if service.ID == impact && len(service.Services) > 0 {
							wg.Add(len(service.Services))
							go func(service EntityService) {
								for _, impact := range service.Services {
									select {
									case <-ctx.Done():
										return
									case workerCh <- impact:
									}
								}
							}(service)

							return
						}
					}
				}()
			}
		}()
	}

	wg.Wait()
}

func (s *service) RecomputeIdleSince(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	defer trace.StartRegion(ctx, "service.RecomputeIdleSince").End()

	services, err := s.adapter.GetEnabled(ctx)
	if err != nil {
		return err
	}

	if len(services) == 0 {
		return nil
	}

	idleSinceMap := NewServicesIdleSinceMap()
	cursor, err := s.entityAdapter.GetWithIdleSince(ctx)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	wg := sync.WaitGroup{}
	workerCh := make(chan types.Entity)
	go func() {
		defer close(workerCh)
		for cursor.Next(ctx) {
			var ent types.Entity
			err := cursor.Decode(&ent)
			if err != nil {
				s.logger.Err(err).Msg("Can't decode entity")
			}

			select {
			case <-ctx.Done():
				return
			case workerCh <- ent:
			}
		}
	}()

	errCh := make(chan error, maxWorkersCount)
	defer close(errCh)

	for i := 0; i < maxWorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case ent, ok := <-workerCh:
					if !ok {
						return
					}

					if ent.IdleSince == nil {
						continue
					}

					if ent.Type == types.EntityTypeResource || ent.Type == types.EntityTypeComponent {
						s.markServices(ctx, &idleSinceMap, services, ent.Services, ent.IdleSince.Unix())
					}

					if ent.Type == types.EntityTypeConnector {
						s.markServices(ctx, &idleSinceMap, services, ent.ImpactedServices, ent.IdleSince.Unix())
					}
				}
			}
		}()
	}

	wg.Wait()

	writeModels := make([]mongodriver.WriteModel, 0, BulkMaxSize)
	for _, service := range services {
		idleSince := idleSinceMap.idleMap[service.ID]
		if idleSince > 0 {
			writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": service.ID}).
				SetUpdate(bson.M{"$set": bson.M{"idle_since": idleSince}}))
		} else {
			writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": service.ID}).
				SetUpdate(bson.M{"$unset": bson.M{"idle_since": ""}}))
		}

		if len(writeModels) == BulkMaxSize {
			err := s.adapter.UpdateBulk(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		err = s.adapter.UpdateBulk(ctx, writeModels)
	}

	return err
}

func (s *service) ProcessRpc(ctx context.Context, event types.Event) error {
	defer trace.StartRegion(ctx, "service.ProcessRpc").End()

	if event.Entity == nil {
		s.logger.Warn().Msgf("event's entity is nil : %+v", event)
		return nil
	}
	if event.AlarmChange == nil {
		return errt.NewUnknownError(fmt.Errorf("event's alarm_change is nil : %v", event))
	}

	return s.calculateState(ctx, event)
}

// calculateState adds or removes alarm from services counters.
// If lockClient is defined service is locked to not conflict with service update. If
// service update in progress event is ignored. If event occurs after service update
// method checks alarm cache to detect which alarm state was used previously.
// If lockClient is not defined service is not locked.
func (s *service) calculateState(ctx context.Context, event types.Event) error {
	services, err := s.storage.GetAll(ctx)
	if err != nil {
		return err
	}
	serviceIDs := make([]string, len(services))
	servicesByID := make(map[string]ServiceData, len(services))
	i := 0
	for _, service := range services {
		serviceIDs[i] = service.ID
		servicesByID[service.ID] = service
		i++
	}

	oldCounters, newCounters, isAlarmChanged := GetAlarmCountersFromEvent(event)
	addedToServices, removedFromServices, unchangedServices := GetServiceIDsFromEvent(event, serviceIDs)

	g, ctx := errgroup.WithContext(ctx)
	workers := int(math.Min(float64(len(services)), float64(maxWorkersCount)))
	workerCh := make(chan workerMsg)
	g.Go(func() error {
		defer close(workerCh)

		for _, serviceID := range unchangedServices {
			data, ok := servicesByID[serviceID]
			if !ok {
				s.logger.Debug().Str("service", serviceID).Msgf("service is missing")
				continue
			}

			select {
			case <-ctx.Done():
				return nil
			case workerCh <- workerMsg{
				Service:   data,
				Unchanged: true,
			}:
			}
		}

		for _, serviceID := range addedToServices {
			data, ok := servicesByID[serviceID]
			if !ok {
				s.logger.Debug().Str("service", serviceID).Msgf("service is missing")
				continue
			}

			select {
			case <-ctx.Done():
				return nil
			case workerCh <- workerMsg{
				Service: data,
				Added:   true,
			}:
			}
		}

		for _, serviceID := range removedFromServices {
			data, ok := servicesByID[serviceID]
			if !ok {
				s.logger.Debug().Str("service", serviceID).Msgf("service is missing")
				continue
			}

			select {
			case <-ctx.Done():
				return nil
			case workerCh <- workerMsg{
				Service: data,
				Removed: true,
			}:
			}
		}

		return nil
	})

	for i := 0; i < workers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case msg, ok := <-workerCh:
					if !ok {
						return nil
					}

					err := s.calculateServiceState(ctx, msg, event, oldCounters,
						newCounters, isAlarmChanged)
					if err != nil {
						return err
					}
				}
			}
		})
	}

	return g.Wait()
}

type workerMsg struct {
	Service                   ServiceData
	Added, Removed, Unchanged bool
}

// calculateServiceState updates service counters based on alarm and update service state.
func (s *service) calculateServiceState(
	ctx context.Context,
	msg workerMsg,
	event types.Event,
	oldCounters, newCounters *AlarmCounters,
	isAlarmChanged bool,
) error {
	var cacheAlarmCounters *AlarmCounters

	if s.lockClient != nil {
		lock, err := s.lockService(ctx, msg.Service.ID, event)
		if err != nil {
			return err
		}

		if lock == nil {
			alarmID := ""
			if event.Alarm != nil {
				alarmID = event.Alarm.ID
			}
			s.logger.Debug().
				Str("service", msg.Service.ID).
				Str("entity", event.GetEID()).
				Str("alarm", alarmID).
				Msg("service update in progress, skip event")

			return nil
		}

		defer func() {
			err := lock.Release(ctx)
			if err != nil {
				s.logger.Err(err).Msg("fail to release lock")
			}
		}()

		key := fmt.Sprintf("%s&&%s", msg.Service.ID, event.GetEID())
		cacheAlarmCounters, err = s.countersCache.RemoveAndGet(ctx, key)
		if err != nil {
			return err
		}
	}

	increments := AlarmCounters{}

	if msg.Added {
		if newCounters == nil {
			return nil
		}

		increments = increments.Add(*newCounters)
		if cacheAlarmCounters != nil {
			increments = increments.Add(cacheAlarmCounters.Negate())
		}
	} else if msg.Removed {
		if oldCounters == nil && cacheAlarmCounters == nil {
			return nil
		}

		if cacheAlarmCounters != nil {
			increments = increments.Add(cacheAlarmCounters.Negate())
		} else if oldCounters != nil {
			increments = increments.Add(oldCounters.Negate())
		}
	} else if msg.Unchanged {
		if !isAlarmChanged && cacheAlarmCounters == nil {
			return nil
		}

		if cacheAlarmCounters != nil {
			increments = increments.Add(cacheAlarmCounters.Negate())
		} else if oldCounters != nil {
			increments = increments.Add(oldCounters.Negate())
		}

		if newCounters != nil {
			increments = increments.Add(*newCounters)
		}
	}
	// Do nothing if service counters don't change.
	if increments.IsZero() {
		return nil
	}

	counters, err := s.countersCache.Update(ctx, map[string]AlarmCounters{msg.Service.ID: increments})
	if err != nil {
		return err
	}

	return s.updateServiceState(ctx, msg.Service.ID, msg.Service.OutputTemplate, counters[msg.Service.ID])
}

// UpdateService recomputes service counters and alarm state.
// If lockClient is defined service is locked during recomputing and all events
// are ignored for service. In this case counters for each alarm are stored in cache to
// determinate which alarm state was used during service recomputing.
// If lockClient is not defined service is not locked and cache for alarms is not used.
func (s *service) UpdateService(ctx context.Context, event types.Event) error {
	serviceID := event.GetEID()
	if s.lockClient != nil {
		lock, err := s.lockServiceUpdate(ctx, serviceID)
		if err != nil {
			return err
		}
		defer func() {
			err := lock.Release(ctx)
			if err != nil {
				s.logger.Err(err).Msg("fail to release lock")
			}
		}()
	}

	serviceData, _, _, isSoftDeleted, err := s.storage.Reload(ctx, serviceID)
	if err != nil {
		return err
	}

	if serviceData == nil {
		err = s.countersCache.Remove(ctx, serviceID)
		if err != nil {
			return err
		}

		err = s.calculateState(ctx, event)
		if err != nil {
			return err
		}

		if isSoftDeleted {
			err = s.entityAdapter.SetResolveDeletedEventProcessed(ctx, serviceID)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err = s.computeServiceCounters(ctx, serviceID, serviceData.OutputTemplate, s.lockClient != nil)
	if err != nil {
		return err
	}

	return s.processSkippedQueue(ctx, serviceID)
}

func (s *service) ReloadService(ctx context.Context, serviceID string) error {
	data, _, _, _, err := s.storage.Reload(ctx, serviceID)
	if err != nil {
		return err
	}

	if data == nil {
		return s.countersCache.Remove(ctx, serviceID)
	}

	return nil
}

// ComputeAllServices recomputes all services counters and alarm state.
// It doesn't lock services during recompute.
func (s *service) ComputeAllServices(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	defer trace.StartRegion(ctx, "service.ComputeAllServices").End()

	services, err := s.loadServices(ctx)
	if err != nil {
		return err
	}

	if len(services) == 0 {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)
	workers := int(math.Min(float64(len(services)), float64(maxWorkersCount)))
	workerCh := make(chan ServiceData)
	g.Go(func() error {
		defer close(workerCh)
		for _, data := range services {
			select {
			case <-ctx.Done():
				return nil
			case workerCh <- data:
			}
		}

		return nil
	})

	for i := 0; i < workers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case data, ok := <-workerCh:
					if !ok {
						return nil
					}

					err := s.computeServiceCounters(ctx, data.ID, data.OutputTemplate, false)
					if err != nil {
						return err
					}
				}
			}
		})
	}

	return g.Wait()
}

func (s *service) ClearCache(ctx context.Context) error {
	return s.countersCache.ClearAll(ctx)
}

func (s *service) loadServices(ctx context.Context) ([]ServiceData, error) {
	data, err := s.storage.ReloadAll(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(data))
	for i := range data {
		ids[i] = data[i].ID
	}

	err = s.countersCache.KeepOnly(ctx, ids)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// lockService locks service during alarm event processing or adds event to queue if
// service update is in progress.
func (s *service) lockService(ctx context.Context, serviceID string, event types.Event) (libredis.Lock, error) {
	key := fmt.Sprintf("%s%s", lockServiceKey, serviceID)
	updateKey := fmt.Sprintf("%s%s", lockServiceUpdateKey, serviceID)
	// Lock service
	lock, err := s.lockClient.Obtain(ctx, key, lockServiceTtl, &redislock.Options{
		RetryStrategy: redislock.LinearBackoff(2 * time.Millisecond),
	})
	if err != nil {
		return nil, err
	}
	// Try to lock update lock to check if service update in progress
	updateLock, err := s.lockClient.Obtain(ctx, updateKey, lockServiceUpdateTtl,
		&redislock.Options{})
	if err != nil {
		// Release lock if update in progress and add entity to queue.
		if err == redislock.ErrNotObtained {
			var encoded []byte
			encoded, err = s.encoder.Encode(types.Event{
				EventType:     types.EventTypeAlarmSkipped,
				Connector:     event.Connector,
				ConnectorName: event.ConnectorName,
				Component:     event.Component,
				Resource:      event.Resource,
				SourceType:    event.SourceType,
				Output:        "recompute service counters",
			})
			if err == nil {
				queueKey := fmt.Sprintf("%s%s", skippedQueueKey, serviceID)
				err = s.redisClient.HSetNX(ctx, queueKey, event.GetEID(), encoded).Err()
			}
		}

		releaseErr := lock.Release(ctx)
		if releaseErr != nil {
			s.logger.Err(releaseErr).Msg("fail to release lock")
		}

		return nil, err
	}
	// Release update lock
	err = updateLock.Release(ctx)
	if err != nil {
		return nil, err
	}

	return lock, nil
}

// lockServiceUpdate locks service during service update event processing.
func (s *service) lockServiceUpdate(ctx context.Context, id string) (libredis.Lock, error) {
	key := fmt.Sprintf("%s%s", lockServiceKey, id)
	updateKey := fmt.Sprintf("%s%s", lockServiceUpdateKey, id)
	// Lock update
	updateLock, err := s.lockClient.Obtain(ctx, updateKey, lockServiceUpdateTtl, &redislock.Options{
		RetryStrategy: redislock.LinearBackoff(30 * time.Millisecond),
	})
	if err != nil {
		return nil, err
	}
	// Try to lock service to wait if the end of service processing
	lock, err := s.lockClient.Obtain(ctx, key, lockServiceUpdateTtl, &redislock.Options{
		RetryStrategy: redislock.LinearBackoff(2 * time.Millisecond),
	})
	if err != nil {
		releaseErr := updateLock.Release(ctx)
		if releaseErr != nil {
			s.logger.Err(releaseErr).Msg("fail to release lock")
		}

		return nil, err
	}
	// Release lock
	err = lock.Release(ctx)
	if err != nil {
		return nil, err
	}

	return updateLock, nil
}

// processSkippedQueue locks service and sends events from queue which were skipped during
// service update
func (s *service) processSkippedQueue(ctx context.Context, serviceID string) error {
	if s.lockClient == nil {
		return nil
	}

	key := fmt.Sprintf("%s%s", lockServiceKey, serviceID)
	lock, err := s.lockClient.Obtain(ctx, key, lockServiceTtl, &redislock.Options{
		RetryStrategy: redislock.LinearBackoff(2 * time.Millisecond),
	})
	if err != nil {
		return err
	}

	defer func() {
		err := lock.Release(ctx)
		if err != nil {
			s.logger.Err(err).Msg("fail to release lock")
		}
	}()

	queueKey := fmt.Sprintf("%s%s", skippedQueueKey, serviceID)
	res := s.redisClient.HGetAll(ctx, queueKey)
	if err := res.Err(); err != nil {
		return err
	}

	events := res.Val()
	if len(events) == 0 {
		return nil
	}

	err = s.redisClient.Del(ctx, queueKey).Err()
	if err != nil {
		return err
	}

	for entityID, body := range events {
		err := s.pubChannel.PublishWithContext(
			ctx,
			s.pubExchangeName,
			s.pubQueueName,
			false,
			false,
			amqp.Publishing{
				Body:        []byte(body),
				ContentType: "application/json",
			},
		)
		if err != nil {
			return err
		}

		s.logger.Debug().Str("entity", entityID).Msgf("send skipped event")
	}

	return nil
}

func (s *service) computeServiceCounters(ctx context.Context, serviceID, outputTemplate string,
	saveAlarmCounters bool) error {

	processedEntities := make(map[string]bool)
	counters := AlarmCounters{}
	entityCursor, err := s.adapter.GetServiceDependencies(ctx, serviceID)
	if err != nil {
		return err
	}

	defer entityCursor.Close(ctx)

	alarms := make([]types.Alarm, 0)
	for entityCursor.Next(ctx) {
		e := types.Entity{}
		err := entityCursor.Decode(&e)
		if err != nil {
			return err
		}

		if processedEntities[e.ID] {
			continue
		}

		alarms = alarms[:0]
		err = s.alarmAdapter.GetOpenedAlarmsByIDs(ctx, []string{e.ID}, &alarms)
		if err != nil {
			return err
		}
		for _, alarm := range alarms {
			alarmCounters := GetAlarmCountersFromAlarm(alarm)
			counters = counters.Add(alarmCounters)
			processedEntities[alarm.EntityID] = true
			if saveAlarmCounters {
				key := fmt.Sprintf("%s&&%s", serviceID, alarm.EntityID)
				err := s.countersCache.Replace(ctx, key, alarmCounters)
				if err != nil {
					return err
				}
			}
		}

		if processedEntities[e.ID] {
			continue
		}

		alarmCounters := GetAlarmCountersFromEntity(e)
		counters = counters.Add(alarmCounters)
		if saveAlarmCounters && !alarmCounters.IsZero() {
			key := fmt.Sprintf("%s&&%s", serviceID, e.ID)
			err := s.countersCache.Replace(ctx, key, alarmCounters)
			if err != nil {
				return err
			}
		}
	}

	err = s.countersCache.Replace(ctx, serviceID, counters)
	if err != nil {
		s.logger.Error().Err(err).Msg("Unable to process state")
	}

	err = s.updateServiceState(ctx, serviceID, outputTemplate, counters)
	if err != nil {
		return err
	}

	return nil
}
