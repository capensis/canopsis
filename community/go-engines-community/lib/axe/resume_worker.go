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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/resumejob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
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
	resumeMaxWorkers               = 10
	resumeRetryInterval            = time.Minute
	resumeRedisLockTTLDuration     = 10 * time.Second
	resumeRedisLockRefreshDuration = resumeRedisLockTTLDuration / 2
	resumeRedisLockAcquireRetries  = 5
	resumeRedisLockAcquireInterval = resumeRedisLockRefreshDuration / resumeRedisLockAcquireRetries
)

type resumeWorker struct {
	client              mongo.DbClient
	collection          mongo.DbCollection
	noTimeoutCollection mongo.DbCollection
	entityCollection    mongo.DbCollection
	resumeJobService    resumejob.Service
	redisLockClient     libredis.LockClient
	publisher           amqp.Publisher
	encoder             encoding.Encoder
	eventGenerator      libevent.Generator
	logger              zerolog.Logger
}

func newResumeWorker(dbClient, noTimeoutClient mongo.DbClient, redisLockClient libredis.LockClient,
	publisher amqp.Publisher, encoder encoding.Encoder, eventGenerator libevent.Generator, logger zerolog.Logger,
) *resumeWorker {
	return &resumeWorker{
		client:              dbClient,
		collection:          dbClient.Collection(mongo.ResumeJobCollection),
		noTimeoutCollection: noTimeoutClient.Collection(mongo.ResumeJobCollection),
		entityCollection:    dbClient.Collection(mongo.EntityMongoCollection),
		resumeJobService:    resumejob.NewService(dbClient, eventGenerator, logger),
		redisLockClient:     redisLockClient,
		publisher:           publisher,
		encoder:             encoder,
		eventGenerator:      eventGenerator,
		logger:              logger,
	}
}

func (w *resumeWorker) Run(ctx context.Context) (err error) {
	var lock libredis.Lock

	defer func() {
		if lock == nil {
			return
		}

		lockErr := lock.Release(context.WithoutCancel(ctx))
		if lockErr != nil && !errors.Is(lockErr, redislock.ErrLockNotHeld) {
			w.logger.Err(lockErr).Msg("failed to release resume worker lock")
		}

		w.logger.Debug().Msg("resume worker: redis lock is released")
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		lock, err = w.obtainLock(ctx)
		if err != nil {
			w.logger.Warn().Err(err).Msg("failed to obtain lock for resume worker, waiting for next attempt")
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(resumeRedisLockAcquireInterval):
				w.logger.Debug().Msg("resume worker: retrying to obtain lock")
				continue
			}
		}

		err = w.runWorkers(ctx, lock)
		if err != nil {
			if errors.Is(err, libredis.ErrFailedToRefreshLock) {
				continue
			}

			return fmt.Errorf("failed to run resume worker: %w", err)
		}
	}
}

func (w *resumeWorker) runWorkers(ctx context.Context, lock libredis.Lock) error {
	ticker := time.NewTicker(resumeRetryInterval)
	defer ticker.Stop()

	jobChan := make(chan resumejob.Job)
	defer close(jobChan)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return w.startLockRefresher(gCtx, lock)
	})

	for i := 0; i < resumeMaxWorkers; i++ {
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
				w.logger.Err(err).Msg("failed to find resume propagation jobs")
			} else {
				err = w.processCursor(gCtx, cursor, jobChan)
				if err != nil {
					w.logger.Err(err).Msg("failed to process resume propagation jobs")
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

func (w *resumeWorker) processCursor(ctx context.Context, cursor mongo.Cursor, ch chan<- resumejob.Job) error {
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var job resumejob.Job
		if err := cursor.Decode(&job); err != nil {
			return fmt.Errorf("failed to decode resume propagation job: %w", err)
		}

		select {
		case ch <- job:
		case <-ctx.Done():
			return nil
		}
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("failed to process resume propagation cursor: %w", err)
	}

	return nil
}

func (w *resumeWorker) runWatcher(ctx context.Context, ch chan<- resumejob.Job) error {
	// Only inserts wake the worker. Receipts and retry updates must not republish
	// every outstanding recipient. The timer also covers stream gaps and restarts.
	stream, err := w.noTimeoutCollection.Watch(ctx, []bson.M{
		{"$match": bson.M{"operationType": mongo.ChangeStreamTypeInsert}},
	})
	if err != nil {
		return fmt.Errorf("failed to watch resume propagation jobs: %w", err)
	}

	defer stream.Close(ctx)

	for stream.Next(ctx) {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		var event struct {
			FullDocument resumejob.Job `bson:"fullDocument"`
		}

		if err := stream.Decode(&event); err != nil {
			return fmt.Errorf("failed to decode resume propagation change: %w", err)
		}

		select {
		case ch <- event.FullDocument:
		case <-ctx.Done():
			return nil
		}
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf("failed to watch resume propagation jobs: %w", err)
	}

	return nil
}

func (w *resumeWorker) runJobWorker(ctx context.Context, jobChan <-chan resumejob.Job) {
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
				w.logger.Err(err).Str("job", job.ID).Msg("failed to process resume propagation job")
			}
		}
	}
}

func (w *resumeWorker) processJob(ctx context.Context, id string) error {
	// Claim a due attempt atomically: watcher and timer can queue the same job.
	// Delivery can overlap again after retry_at; recipient transactions make
	// duplicate delivery safe without keeping a separate in-memory job registry.
	var job resumejob.Job
	err := w.collection.FindOneAndUpdate(ctx,
		bson.M{"_id": id, "retry_at": bson.M{"$lte": time.Now()}},
		bson.M{"$set": bson.M{"retry_at": time.Now().Add(resumeRetryInterval)}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&job)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil
		}

		return fmt.Errorf("failed to find resume propagation job: %w", err)
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
			return fmt.Errorf("failed to find resume recipient %s: %w", id, err)
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
			return fmt.Errorf("failed to generate resume for %s: %w", id, err)
		}

		if job.Event.Connector != "" {
			event.Connector = job.Event.Connector
			event.ConnectorName = job.Event.ConnectorName
		}

		event.EventType = types.EventTypeResume
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
		event.PauseStart = job.Event.PauseStart
		event.ResumeJob = job.ID

		body, err := w.encoder.Encode(event)
		if err != nil {
			return fmt.Errorf("failed to encode resume for %s: %w", id, err)
		}

		err = w.publisher.PublishWithContext(ctx, "", canopsis.FIFOQueueName, false, false, amqp091.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		})
		if err != nil {
			return fmt.Errorf("failed to publish resume for %s: %w", id, err)
		}
	}

	_, err = w.collection.DeleteOne(ctx, bson.M{"_id": job.ID, "initialized": true, "pending": bson.M{"$size": 0}})
	if err != nil {
		return fmt.Errorf("failed to delete completed resume propagation job: %w", err)
	}

	return nil
}

func (w *resumeWorker) initializeJob(ctx context.Context, job *resumejob.Job) error {
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

	// A removed member can still hold this exact pause. Keep recovery tied to
	// both its source/start and its immediate pause parent. Ancestors must not
	// also enumerate the descendants owned by an intermediary's child job.
	retained := bson.M{
		"pause":        bson.M{"$exists": true},
		"pause.source": job.Event.PauseSource, "pause.timestamp": job.Event.PauseStart,
	}
	if job.Parent == job.Event.PauseSource {
		retained["pause.path.0"] = bson.M{"$exists": false}
	} else {
		retained["pause.path"] = job.Parent
		retained["$expr"] = bson.M{"$eq": bson.A{
			bson.M{"$arrayElemAt": bson.A{bson.M{"$ifNull": bson.A{"$pause.path", bson.A{}}}, -1}},
			job.Parent,
		}}
	}

	cursor, err := w.entityCollection.Find(ctx,
		bson.M{"_id": bson.M{"$ne": job.Parent}, "$or": bson.A{relationship, retained}},
		options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return fmt.Errorf("failed to find resume dependencies: %w", err)
	}

	defer cursor.Close(ctx)

	pending := make([]string, 0)
	for cursor.Next(ctx) {
		var entity types.Entity
		if err := cursor.Decode(&entity); err != nil {
			return fmt.Errorf("failed to decode resume dependency: %w", err)
		}

		pending = append(pending, entity.ID)
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("failed to fetch resume dependencies: %w", err)
	}

	_, err = w.collection.UpdateOne(ctx, bson.M{"_id": job.ID, "initialized": false}, bson.M{
		"$set": bson.M{"pending": pending, "initialized": true},
	})
	if err != nil {
		return fmt.Errorf("failed to initialize resume propagation job: %w", err)
	}

	job.Pending = pending
	job.Initialized = true

	return nil
}

func (w *resumeWorker) completeUnavailable(ctx context.Context, job resumejob.Job, id string) error {
	return w.client.WithTransaction(ctx, func(ctx context.Context) error {
		pending, err := w.resumeJobService.IsPending(ctx, job.ID, id)
		if err != nil {
			return fmt.Errorf("failed to find resume propagation receipt: %w", err)
		}

		if !pending {
			return nil
		}

		var entity types.Entity
		err = w.entityCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return fmt.Errorf("failed to find unavailable resume recipient: %w", err)
		}

		if entity.Enabled {
			return nil
		}

		entity.ID = id
		// Disabled or deleted intermediaries cannot consume a normal event, but
		// their descendants still need a durable continuation before this receipt.
		err = w.resumeJobService.Create(ctx, entity, types.LatestPause{
			Timestamp:  *job.Event.PauseStart,
			Source:     job.Event.PauseSource,
			SourceType: job.Event.PauseSourceType,
			Path:       job.Event.PausePath,
			Inherited:  true,
		}, rpc.AxeParameters{
			Timestamp: job.Event.Timestamp,
			Output:    job.Event.Output,
			Author:    job.Event.Author,
			User:      job.Event.UserID,
			Role:      job.Event.Role,
			Initiator: job.Event.Initiator,
		})
		if err != nil {
			return err
		}

		// TODO: a rapid disable/enable can strand the recipient's pause if this
		// receipt is completed before AXE handles disable, and enable then makes
		// that handler skip cleanup. The pause is retained without a pending retry.
		if err := w.resumeJobService.Complete(ctx, job.ID, id); err != nil {
			return err
		}

		return nil
	})
}

func (w *resumeWorker) obtainLock(ctx context.Context) (libredis.Lock, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
		}

		lock, err := w.redisLockClient.Obtain(ctx, libredis.AxeResumeLockKey, resumeRedisLockTTLDuration, &redislock.Options{
			RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(resumeRedisLockAcquireInterval), resumeRedisLockAcquireRetries),
		})
		if err != nil {
			if errors.Is(err, redislock.ErrNotObtained) {
				w.logger.Debug().Msg("resume worker redis lock is not obtained, retry")
				continue
			}

			return nil, fmt.Errorf("cannot obtain resume worker lock: %w", err)
		}

		w.logger.Debug().Msg("resume worker redis lock is obtained")

		return lock, nil
	}
}

func (w *resumeWorker) startLockRefresher(ctx context.Context, lock libredis.Lock) error {
	ticker := time.NewTicker(resumeRedisLockRefreshDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := lock.Refresh(ctx, resumeRedisLockTTLDuration, &redislock.Options{})
			if err != nil {
				w.logger.Err(err).Msg("failed to refresh resume worker lock")

				return libredis.ErrFailedToRefreshLock
			}

			w.logger.Debug().Msg("resume worker redis lock is refreshed")
		}
	}
}
