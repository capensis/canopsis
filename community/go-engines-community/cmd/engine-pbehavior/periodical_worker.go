package main

import (
	"context"
	"fmt"
	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type periodicalWorker struct {
	ChannelPub             libamqp.Channel
	PeriodicalInterval     time.Duration
	LockerClient           redis.LockClient
	Store                  redis.Store
	PbhService             pbehavior.Service
	DbClient               mongo.DbClient
	EventManager           pbehavior.EventManager
	FrameDuration          time.Duration
	TimezoneConfigProvider config.TimezoneConfigProvider
	Encoder                encoding.Encoder
	Logger                 zerolog.Logger
}

func (w *periodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *periodicalWorker) Work(ctx context.Context) error {
	w.Logger.Debug().Msg("Periodical process")

	_, err := w.LockerClient.Obtain(ctx, redis.PeriodicalLockKey, w.GetInterval(), nil)
	if err == redislock.ErrNotObtained {
		w.Logger.Debug().Msg("Periodical process: Could not obtain periodical lock! Skip periodical process")
		return nil
	} else if err != nil {
		w.Logger.Error().Err(err).Msg("Periodical process: obtain redis lock - unexpected error! Skip periodical process")
		return nil
	}

	backoff := time.Second
	retries := int(w.GetInterval().Seconds() - 1)
	computeLock, err := w.LockerClient.Obtain(ctx, redis.RecomputeLockKey, redis.RecomputeLockDuration, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(backoff), retries),
	})

	defer func() {
		if computeLock != nil {
			err := computeLock.Release(ctx)
			if err != nil && err != redislock.ErrLockNotHeld {
				w.Logger.Warn().Err(err).Msg("Periodical process: failed to manually release compute-lock, the lock will be released by ttl")
			}
		}
	}()

	if err != nil {
		w.Logger.Err(err).Msg("Periodical process: obtain redlock failed! Skip periodical process")
		return nil
	}

	ok, err := w.Store.Restore(ctx, w.PbhService)
	if err != nil {
		w.Logger.Err(err).Msg("Periodical process: get pbehavior's frames from redis failed! Skip periodical process")
		return nil
	}

	now := time.Now().In(w.TimezoneConfigProvider.Get().Location)

	span := w.PbhService.GetSpan()
	if !ok || span.To().Before(now.Add(w.FrameDuration/2)) {
		err = w.PbhService.Compute(ctx, timespan.New(now, now.Add(w.FrameDuration)))
		if err != nil {
			w.Logger.Err(err).Msg("Periodical process: compute pbehavior's frames failed! Skip periodical process")
			return nil
		}

		err = w.Store.Save(ctx, w.PbhService)
		if err != nil {
			w.Logger.Err(err).Msg("Periodical process: save pbehavior's frames to redis failed! Skip periodical process")
			return nil
		}
	}

	err = computeLock.Release(ctx)
	if err != nil {
		if err == redislock.ErrLockNotHeld {
			return nil
		} else {
			w.Logger.Warn().Msg("Periodical process: failed to manually release compute-lock, the lock will be released by ttl")
		}
	}

	alarmCollection := w.DbClient.Collection(mongo.AlarmMongoCollection)

	cursor, err := alarmCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"v.resolved": bson.M{"$in": bson.A{false, nil}},
			"t":          bson.M{"$lt": types.CpsTime{Time: now}},
		}},
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
		{"$match": bson.M{"entity.enabled": true}},
	})

	if err != nil {
		w.Logger.Err(err).Msg("Periodical process: get alarms from mongo failed! Skip periodical process")
		return nil
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var alarmWithEntity types.AlarmWithEntity

		err = cursor.Decode(&alarmWithEntity)
		if err != nil {
			w.Logger.Err(err).Msg("Periodical process: decode alarm with entity failed! Skip periodical process")
			return nil
		}

		alarm := alarmWithEntity.Alarm
		entity := alarmWithEntity.Entity

		resolveResult, err := w.PbhService.Resolve(ctx, &entity, now)
		if err != nil {
			w.Logger.Err(err).Str("entity_id", entity.ID).Msg("Periodical process: resolve an entity failed! Skip periodical process")
			return nil
		}

		event := w.EventManager.GetEvent(resolveResult, alarm, now)
		if event.EventType != "" {
			err := w.publishToEngineFIFO(event)
			if err != nil {
				w.Logger.Err(err).Str("alarm_id", alarm.ID).Msgf("failed to send %s event", event.EventType)
				return nil
			}

			w.Logger.Debug().Str("resolve pbehavior", resolveResult.ResolvedPbhID).Str("resolve type", fmt.Sprintf("%+v", resolveResult.ResolvedType)).Str("alarm", alarm.ID).Msgf("Periodical process: send %s event", event.EventType)
		}
	}

	return nil
}

func (w *periodicalWorker) publishToEngineFIFO(event types.Event) error {
	return w.publishTo(event, canopsis.FIFOQueueName)
}

func (w *periodicalWorker) publishTo(event types.Event, queue string) error {
	bevent, err := w.Encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("publishTo(): error while encoding event %+v", err)
	}

	return errt.NewIOError(w.ChannelPub.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json", // this type is mandatory to avoid bad conversions into Python.
			Body:         bevent,
			DeliveryMode: amqp.Persistent,
		},
	))
}
