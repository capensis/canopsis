package pbehavior

import (
	"context"
	"encoding/json"
	"fmt"
	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

const (
	OperationCreate = iota
	OperationUpdate
	OperationDelete
)

const (
	CalculateAll  = ""
	CalculateNone = "none"
)

// Computer is used to implement pbehavior intervals recompute on singnal from channel.
type Computer interface {
	// Compute recomputes all pbehaviors on empty signal
	// or recomputes only one pbehavior on signal which conatins pbehavior id.
	Compute(ctx context.Context, ch <-chan ComputeTask)
}

type ComputeTask struct {
	// PbehaviorID defines for which pbehavior intervals should be recomputed.
	// If empty all intervals are recomputed.
	PbehaviorID string
	// UpdateAlarms prescribes to update alarms by pbehavior. PbehaviorID shouldn't be empty.
	OperationType int
}

// cancelableComputer recomputes pbehaviors.
type cancelableComputer struct {
	logger       zerolog.Logger
	lockClient   redis.LockClient
	redisStore   redis.Store
	service      Service
	dbClient     mongo.DbClient
	eventManager EventManager
	encoder      encoding.Encoder
	publisher    libamqp.Publisher
	queue        string
}

// NewCancelableComputer creates new computer.
func NewCancelableComputer(
	lockClient redis.LockClient,
	redisStore redis.Store,
	service Service,
	dbClient mongo.DbClient,
	publisher libamqp.Publisher,
	eventManager EventManager,
	encoder encoding.Encoder,
	queue string,
	logger zerolog.Logger,
) Computer {
	return &cancelableComputer{
		logger:       logger,
		lockClient:   lockClient,
		redisStore:   redisStore,
		service:      service,
		dbClient:     dbClient,
		eventManager: eventManager,
		encoder:      encoder,
		publisher:    publisher,
		queue:        queue,
	}
}

func (c *cancelableComputer) Compute(parentCtx context.Context, ch <-chan ComputeTask) {
	var ctx context.Context
	var cancel context.CancelFunc
	c.logger.Debug().Msg("compute started")
	defer func() {
		if cancel != nil {
			cancel()
		}
		c.logger.Debug().Msg("compute ended")
	}()

	ongoingID := CalculateNone
	ongoingMx := sync.Mutex{}

	wg := sync.WaitGroup{}

	for {
		select {
		case <-parentCtx.Done():
			return
		case task, ok := <-ch:
			if !ok {
				return
			}

			pbehaviorID := task.PbehaviorID

			ongoingMx.Lock()
			currentOngoing := ongoingID
			ongoingMx.Unlock()

			// Cancel ongoing calculation if
			// - ongoing calculated pbehavior is the same as in the signal
			// - ongoing calculation calculates all pbehavior
			// - new calculation calculates all pbehavior
			if currentOngoing != CalculateNone && (pbehaviorID == CalculateAll || currentOngoing == CalculateAll || currentOngoing == pbehaviorID) {
				c.logger.Info().Str("pbehaviorID", pbehaviorID).Str("ongoingID", currentOngoing).Msg("Test cancel log")
				c.logger.Debug().Msg("API pbehavior recompute: STOP SIGNAL!!!")

				cancel()
				wg.Wait()
				// If canceled calculation calculates all pbehavior keep it
				if currentOngoing != CalculateAll {
					ongoingID = pbehaviorID
				} else {
					ongoingID = CalculateAll
				}
			} else {
				// Wait ongoing calculation to avoid conflicts
				wg.Wait()
				ongoingID = pbehaviorID
			}

			ctx, cancel = context.WithCancel(parentCtx)
			wg.Add(1)

			go func(ctx context.Context, pbhId string, opType int) {
				defer func() {
					ongoingMx.Lock()
					ongoingID = CalculateNone
					ongoingMx.Unlock()

					wg.Done()
				}()

				c.computePbehavior(ctx, pbhId, opType)
			}(ctx, ongoingID, task.OperationType)
		}
	}
}

// computePbehavior obtains lock and calls computer.
func (c *cancelableComputer) computePbehavior(
	ctx context.Context,
	pbehaviorID string,
	operationType int,
) {
	computeLock, err := c.lockClient.Obtain(
		ctx,
		redis.RecomputeLockKey,
		redis.RecomputeLockDuration,
		&redislock.Options{
			RetryStrategy: redislock.LinearBackoff(time.Second),
		},
	)
	if err != nil {
		c.logger.Err(err).Msg("API pbehavior recompute: obtain redlock failed! Skip recompute")
		return
	}
	c.logger.Debug().Msg("API pbehavior recompute: obtain redlock")
	defer func() {
		if computeLock != nil {
			err := computeLock.Release(ctx)
			if err != nil {
				if err == redislock.ErrLockNotHeld {
					c.logger.Err(err).Msg("API pbehavior recompute: the pbehavior's frames computing took more time than redlock ttl, the data might be inconsistent")
				} else {
					c.logger.Warn().Msg("API pbehavior recompute: failed to manually release compute-lock, the lock will be released by ttl")
				}
			} else {
				c.logger.Debug().Msg("API pbehavior recompute: release redlock")
			}
		}
	}()

	ok, err := c.redisStore.Restore(ctx, c.service)
	if err != nil || !ok {
		if err == nil {
			err = fmt.Errorf("pbehavior intervals are not computed, cache is empty")
		}

		c.logger.Err(err).Msgf("API pbehavior recompute failed")
		return
	}

	if pbehaviorID == "" {
		err = c.service.Compute(ctx, c.service.GetSpan())
	} else {
		err = c.service.Recompute(ctx, pbehaviorID)
	}

	if err != nil {
		c.logger.Err(err).Msgf("API pbehavior recompute failed")
		return
	}

	err = c.redisStore.Save(ctx, c.service)
	if err != nil {
		c.logger.Err(err).Msg("API pbehavior recompute: save pbehavior's frames to redis failed! The data might be inconsistent")
		return
	}

	c.logger.Info().Str("pbehavior", pbehaviorID).Msg("API pbehavior recompute: pbehavior recomputed")

	if pbehaviorID != "" {
		err := c.updateAlarms(ctx, pbehaviorID, operationType)
		if err != nil {
			return
		}
	}
}

func (c *cancelableComputer) updateAlarms(
	ctx context.Context,
	pbehaviorID string,
	operationType int,
) error {
	switch operationType {
	case OperationCreate:
		cursor, err := c.findAlarmsMatchPbhFilter(ctx, pbehaviorID)
		if err != nil {
			c.logger.Err(err).Msg("API pbehavior recompute: cannot find alarms")
			return err
		}

		err = c.sendAlarmEvents(ctx, cursor, pbehaviorID)
		if err != nil {
			return err
		}
	case OperationUpdate:
		cursor, err := c.findAlarmsMatchPbhFilter(ctx, pbehaviorID)
		if err != nil {
			return err
		}

		err = c.sendAlarmEvents(ctx, cursor, pbehaviorID)
		if err != nil {
			return err
		}

		cursor, err = c.findAlarmsMatchPbhID(ctx, pbehaviorID)
		if err != nil {
			return err
		}

		err = c.sendAlarmEvents(ctx, cursor, pbehaviorID)
		if err != nil {
			return err
		}
	case OperationDelete:
		cursor, err := c.findAlarmsMatchPbhID(ctx, pbehaviorID)
		if err != nil {
			c.logger.Err(err).Msg("API pbehavior recompute: cannot find alarms")
			return err
		}

		err = c.sendAlarmEvents(ctx, cursor, pbehaviorID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cancelableComputer) sendAlarmEvents(ctx context.Context, cursor mongo.Cursor, pbehaviorID string) error {
	now := time.Now()
	for cursor.Next(ctx) {
		alarm := types.AlarmWithEntity{}
		err := cursor.Decode(&alarm)
		if err != nil {
			c.logger.Err(err).Msg("API pbehavior recompute: cannot decode alarm")
			return err
		}

		resolveResult, err := c.service.Resolve(ctx, &alarm.Entity, now)
		if err != nil {
			c.logger.Err(err).Msg("API pbehavior recompute: cannot resolve pbehavior for entity")
			return err
		}

		event := c.eventManager.GetEvent(resolveResult, alarm.Alarm, now)
		if event.EventType == "" {
			continue
		}

		body, err := c.encoder.Encode(event)
		if err != nil {
			c.logger.Err(err).Msg("API pbehavior recompute: cannot encode event")
			return err
		}

		err = c.publisher.Publish(
			"",
			c.queue,
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: amqp.Persistent,
			},
		)

		if err != nil {
			c.logger.Err(err).Str("pbehavior", pbehaviorID).Str("alarm", alarm.Alarm.ID).Msgf("API pbehavior recompute: cannot send %s event", event.EventType)
			return err
		}

		c.logger.Info().Str("pbehavior", pbehaviorID).Str("alarm", alarm.Alarm.ID).Msgf("API pbehavior recompute: send %s event", event.EventType)
	}

	return nil
}

func (c *cancelableComputer) findAlarmsMatchPbhFilter(
	parentCtx context.Context,
	pbehaviorID string,
) (mongo.Cursor, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	pbehaviorCollection := c.dbClient.Collection(PBehaviorCollectionName)
	pbehavior := PBehavior{}
	err := pbehaviorCollection.FindOne(ctx, bson.M{"_id": pbehaviorID},
		options.FindOne().SetProjection(bson.M{
			"filter": 1,
		})).Decode(&pbehavior)

	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, fmt.Errorf("cannot find pbehavior id=%s", pbehaviorID)
		} else {
			return nil, err
		}
	}

	var filter interface{}
	err = json.Unmarshal([]byte(pbehavior.Filter), &filter)
	if err != nil {
		return nil, err
	}

	entityCollection := c.dbClient.Collection(mongo.EntityMongoCollection)
	return entityCollection.Aggregate(ctx, []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"entity": "$$ROOT",
		}},
		{"$lookup": bson.M{
			"from": mongo.AlarmMongoCollection,
			"let":  bson.M{"eid": "$entity._id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$$eid", "$d"}}}},
				{"$match": bson.M{"$or": []bson.M{
					{"v.resolved": bson.M{"$in": bson.A{false, nil}}},
					{"v.resolved": bson.M{"$exists": false}},
				}}},
			},
			"as": "alarm",
		}},
		{"$unwind": "$alarm"},
	})
}

func (c *cancelableComputer) findAlarmsMatchPbhID(
	parentCtx context.Context,
	pbehaviorID string,
) (mongo.Cursor, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	return c.dbClient.Collection(mongo.AlarmMongoCollection).Aggregate(ctx, []bson.M{
		{"$match": bson.M{"v.pbehavior_info.id": pbehaviorID}},
		{"$match": bson.M{"$or": []bson.M{
			{"v.resolved": bson.M{"$in": bson.A{false, nil}}},
			{"v.resolved": bson.M{"$exists": false}},
		}}},
		{"$project": bson.M{
			"alarm": "$$ROOT",
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
	})
}
