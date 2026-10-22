package axe

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pausejob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	amqp091 "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/sync/errgroup"
)

const (
	pauseMaxWorkers               = 10
	pauseRetryInterval            = time.Minute
	pauseRedisLockTTLDuration     = 10 * time.Second
	pauseRedisLockRefreshDuration = pauseRedisLockTTLDuration / 2
	pauseRedisLockAcquireRetries  = 5
	pauseRedisLockAcquireInterval = pauseRedisLockRefreshDuration / pauseRedisLockAcquireRetries
)

type pauseWorker struct {
	client              mongo.DbClient
	collection          mongo.DbCollection
	noTimeoutCollection mongo.DbCollection
	entityCollection    mongo.DbCollection
	pauseJobService     pausejob.Service
	redisLockClient     libredis.LockClient
	publisher           amqp.Publisher
	encoder             encoding.Encoder
	eventGenerator      libevent.Generator
	logger              zerolog.Logger
}

func newPauseWorker(dbClient, noTimeoutClient mongo.DbClient, redisLockClient libredis.LockClient,
	publisher amqp.Publisher, encoder encoding.Encoder, eventGenerator libevent.Generator, logger zerolog.Logger,
) *pauseWorker {
	return &pauseWorker{
		client:              dbClient,
		collection:          dbClient.Collection(mongo.PauseJobCollection),
		noTimeoutCollection: noTimeoutClient.Collection(mongo.PauseJobCollection),
		entityCollection:    dbClient.Collection(mongo.EntityMongoCollection),
		pauseJobService:     pausejob.NewService(dbClient, eventGenerator, logger),
		redisLockClient:     redisLockClient,
		publisher:           publisher,
		encoder:             encoder,
		eventGenerator:      eventGenerator,
		logger:              logger,
	}
}

func (w *pauseWorker) Run(ctx context.Context) (err error) {
	var lock libredis.Lock

	defer func() {
		if lock == nil {
			return
		}

		lockErr := lock.Release(context.WithoutCancel(ctx))
		if lockErr != nil && !errors.Is(lockErr, redislock.ErrLockNotHeld) {
			w.logger.Err(lockErr).Msg("failed to release pause worker lock")
		}

		w.logger.Debug().Msg("pause worker: redis lock is released")
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		lock, err = w.obtainLock(ctx)
		if err != nil {
			w.logger.Warn().Err(err).Msg("failed to obtain lock for pause worker, waiting for next attempt")
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(pauseRedisLockAcquireInterval):
				w.logger.Debug().Msg("pause worker: retrying to obtain lock")
				continue
			}
		}

		err = w.runWorkers(ctx, lock)
		if err != nil {
			if errors.Is(err, libredis.ErrFailedToRefreshLock) {
				continue
			}

			return fmt.Errorf("failed to run pause worker: %w", err)
		}
	}
}

func (w *pauseWorker) runWorkers(ctx context.Context, lock libredis.Lock) error {
	ticker := time.NewTicker(pauseRetryInterval)
	defer ticker.Stop()

	jobChan := make(chan pausejob.Job)
	defer close(jobChan)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return w.startLockRefresher(gCtx, lock)
	})

	for i := 0; i < pauseMaxWorkers; i++ {
		g.Go(func() error {
			w.runJobWorker(gCtx, jobChan)

			return nil
		})
	}

	g.Go(func() error {
		return w.runWatcher(gCtx, jobChan)
	})

	g.Go(func() error {
		for {
			cursor, err := w.collection.Find(gCtx, bson.M{"retry_at": bson.M{"$lte": time.Now()}})
			if err != nil {
				w.logger.Err(err).Msg("failed to find pause propagation jobs")
			} else {
				err = w.processCursor(gCtx, cursor, jobChan)
				if err != nil {
					w.logger.Err(err).Msg("failed to process pause propagation jobs")
				}
			}

			select {
			case <-gCtx.Done():
				return nil
			case <-ticker.C:
			}
		}
	})

	return g.Wait()
}

func (w *pauseWorker) processCursor(ctx context.Context, cursor mongo.Cursor, ch chan<- pausejob.Job) error {
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var job pausejob.Job
		if err := cursor.Decode(&job); err != nil {
			return fmt.Errorf("failed to decode pause propagation job: %w", err)
		}

		select {
		case ch <- job:
		case <-ctx.Done():
			return nil
		}
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("failed to process pause propagation cursor: %w", err)
	}

	return nil
}

func (w *pauseWorker) runWatcher(ctx context.Context, ch chan<- pausejob.Job) error {
	// Only inserts wake the worker. Receipts and retry updates must not republish
	// every outstanding recipient. The timer also covers stream gaps and restarts.
	stream, err := w.noTimeoutCollection.Watch(ctx, []bson.M{
		{"$match": bson.M{"operationType": mongo.ChangeStreamTypeInsert}},
	})
	if err != nil {
		return fmt.Errorf("failed to watch pause propagation jobs: %w", err)
	}

	defer stream.Close(ctx)

	for stream.Next(ctx) {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		var event struct {
			FullDocument pausejob.Job `bson:"fullDocument"`
		}

		if err := stream.Decode(&event); err != nil {
			return fmt.Errorf("failed to decode pause propagation change: %w", err)
		}

		select {
		case ch <- event.FullDocument:
		case <-ctx.Done():
			return nil
		}
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf("failed to watch pause propagation jobs: %w", err)
	}

	return nil
}

func (w *pauseWorker) runJobWorker(ctx context.Context, jobChan <-chan pausejob.Job) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobChan:
			if !ok {
				return
			}

			err := w.processJob(ctx, job.ID)
			if err != nil {
				w.logger.Err(err).Str("job", job.ID).Msg("failed to process pause propagation job")
			}
		}
	}
}

func (w *pauseWorker) processJob(ctx context.Context, id string) error {
	// Claim a due attempt atomically: watcher and timer can queue the same job.
	// Delivery can overlap again after retry_at; recipient transactions make
	// duplicate delivery safe without keeping a separate in-memory job registry.
	var job pausejob.Job
	err := w.collection.FindOneAndUpdate(ctx,
		bson.M{"_id": id, "retry_at": bson.M{"$lte": time.Now()}},
		bson.M{"$set": bson.M{"retry_at": time.Now().Add(pauseRetryInterval)}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&job)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil
		}

		return fmt.Errorf("failed to find pause propagation job: %w", err)
	}

	if !job.Initialized {
		err = w.initializeJob(ctx, &job)
		if err != nil {
			return err
		}
	}

	for _, id := range job.Pending {
		var entity types.Entity
		err := w.entityCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return fmt.Errorf("failed to find pause recipient %s: %w", id, err)
		}

		if entity.ID == "" || !entity.Enabled {
			err := w.completeUnavailable(ctx, job, id)
			if err != nil {
				return err
			}

			// If it was enabled concurrently, retry from a fresh entity snapshot.
			continue
		}

		event, err := w.eventGenerator.Generate(entity)
		if err != nil {
			return fmt.Errorf("failed to generate pause for %s: %w", id, err)
		}

		if job.Event.Connector != "" {
			event.Connector = job.Event.Connector
			event.ConnectorName = job.Event.ConnectorName
		}

		event.EventType = types.EventTypePause
		event.Timestamp = job.Event.Timestamp
		event.Output = job.Event.Output
		event.Author = job.Event.Author
		event.UserID = job.Event.UserID
		event.Role = job.Event.Role
		event.Initiator = job.Event.Initiator
		event.Inherited = true
		event.PauseSource = job.Event.PauseSource
		event.PauseSourceType = job.Event.PauseSourceType
		event.PausePath = job.Event.PausePath
		event.PauseType = job.Event.PauseType
		event.PauseJob = job.ID

		body, err := w.encoder.Encode(event)
		if err != nil {
			return fmt.Errorf("failed to encode pause for %s: %w", id, err)
		}

		err = w.publisher.PublishWithContext(ctx, "", canopsis.FIFOQueueName, false, false, amqp091.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		})
		if err != nil {
			return fmt.Errorf("failed to publish pause for %s: %w", id, err)
		}
	}

	_, err = w.collection.DeleteOne(ctx, bson.M{"_id": job.ID, "initialized": true, "pending": bson.M{"$size": 0}})
	if err != nil {
		return fmt.Errorf("failed to delete completed pause propagation job: %w", err)
	}

	return nil
}

func (w *pauseWorker) initializeJob(ctx context.Context, job *pausejob.Job) error {
	var relationship bson.M
	switch job.ParentType {
	case types.EntityTypeService:
		relationship = bson.M{"services": job.Parent}
	case types.EntityTypeComponent:
		relationship = bson.M{"component": job.Parent, "type": types.EntityTypeResource}
	default:
		relationship = bson.M{"$or": bson.A{
			bson.M{"services": job.Parent},
			bson.M{"component": job.Parent, "type": types.EntityTypeResource},
		}}
	}

	cursor, err := w.entityCollection.Find(ctx,
		bson.M{"_id": bson.M{"$ne": job.Parent}, "enabled": true, "$and": bson.A{relationship}},
		options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return fmt.Errorf("failed to find pause dependencies: %w", err)
	}

	defer cursor.Close(ctx)

	pending := make([]string, 0)
	for cursor.Next(ctx) {
		var entity types.Entity
		if err := cursor.Decode(&entity); err != nil {
			return fmt.Errorf("failed to decode pause dependency: %w", err)
		}

		pending = append(pending, entity.ID)
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("failed to fetch pause dependencies: %w", err)
	}

	_, err = w.collection.UpdateOne(ctx, bson.M{"_id": job.ID, "initialized": false}, bson.M{
		"$set": bson.M{"pending": pending, "initialized": true},
	})
	if err != nil {
		return fmt.Errorf("failed to initialize pause propagation job: %w", err)
	}

	job.Pending = pending
	job.Initialized = true

	return nil
}

func (w *pauseWorker) completeUnavailable(ctx context.Context, job pausejob.Job, id string) error {
	return w.client.WithTransaction(ctx, func(ctx context.Context) error {
		pending, err := w.pauseJobService.IsPending(ctx, job.ID, id)
		if err != nil {
			return fmt.Errorf("failed to find pause propagation receipt: %w", err)
		}

		if !pending {
			return nil
		}

		var entity types.Entity
		err = w.entityCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return fmt.Errorf("failed to find unavailable pause recipient: %w", err)
		}

		if entity.Enabled {
			return nil
		}

		if err := w.pauseJobService.Complete(ctx, job.ID, id); err != nil {
			return err
		}

		return nil
	})
}

func (w *pauseWorker) obtainLock(ctx context.Context) (libredis.Lock, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
		}

		lock, err := w.redisLockClient.Obtain(ctx, libredis.AxePauseLockKey, pauseRedisLockTTLDuration, &redislock.Options{
			RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(pauseRedisLockAcquireInterval), pauseRedisLockAcquireRetries),
		})
		if err != nil {
			if errors.Is(err, redislock.ErrNotObtained) {
				w.logger.Debug().Msg("pause worker redis lock is not obtained, retry")
				continue
			}

			return nil, fmt.Errorf("cannot obtain pause worker lock: %w", err)
		}

		w.logger.Debug().Msg("pause worker redis lock is obtained")

		return lock, nil
	}
}

func (w *pauseWorker) startLockRefresher(ctx context.Context, lock libredis.Lock) error {
	ticker := time.NewTicker(pauseRedisLockRefreshDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := lock.Refresh(ctx, pauseRedisLockTTLDuration, &redislock.Options{})
			if err != nil {
				w.logger.Err(err).Msg("failed to refresh pause worker lock")

				return libredis.ErrFailedToRefreshLock
			}

			w.logger.Debug().Msg("pause worker redis lock is refreshed")
		}
	}
}
