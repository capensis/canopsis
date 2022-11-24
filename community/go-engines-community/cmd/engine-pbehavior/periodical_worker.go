package main

import (
	"context"
	"fmt"
	"time"

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
)

type periodicalWorker struct {
	ChannelPub             libamqp.Channel
	PeriodicalInterval     time.Duration
	LockerClient           redis.LockClient
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
	_, err := w.LockerClient.Obtain(ctx, redis.PeriodicalLockKey, w.GetInterval()-100*time.Millisecond, nil)
	if err == redislock.ErrNotObtained {
		w.Logger.Debug().Msg("Could not obtain periodical lock! Skip periodical process")
		return nil
	} else if err != nil {
		w.Logger.Error().Err(err).Msg("obtain redis lock - unexpected error")
		return nil
	}
	now := time.Now().In(w.TimezoneConfigProvider.Get().Location)
	w.compute(ctx, now)
	w.processAlarms(ctx, now)

	return nil
}

func (w *periodicalWorker) compute(ctx context.Context, now time.Time) {
	newSpan := timespan.New(now, now.Add(w.FrameDuration))
	count, err := w.PbhService.Compute(ctx, newSpan)
	if err != nil {
		w.Logger.Err(err).Msg("compute pbehavior's frames failed")
		return
	}

	if count >= 0 {
		w.Logger.Info().
			Time("interval_from", newSpan.From()).
			Time("interval_to", newSpan.To()).
			Int("count", count).
			Msg("pbehaviors are recomputed")
	}
}

func (w *periodicalWorker) processAlarms(ctx context.Context, now time.Time) {
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
		w.Logger.Err(err).Msg("get alarms from mongo failed")
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var alarmWithEntity types.AlarmWithEntity

		err = cursor.Decode(&alarmWithEntity)
		if err != nil {
			w.Logger.Err(err).Msg("decode alarm with entity failed")
			return
		}

		alarm := alarmWithEntity.Alarm
		entity := alarmWithEntity.Entity

		resolveResult, err := w.PbhService.Resolve(ctx, entity.ID, now)
		if err != nil {
			w.Logger.Err(err).Str("entity_id", entity.ID).Msg("resolve an entity failed")
			return
		}

		event := w.EventManager.GetEvent(resolveResult, alarm, now)
		if event.EventType != "" {
			err := w.publishToEngineFIFO(event)
			if err != nil {
				w.Logger.Err(err).Str("alarm_id", alarm.ID).Msgf("failed to send %s event", event.EventType)
				return
			}

			w.Logger.Debug().
				Str("resolve_pbehavior", resolveResult.ResolvedPbhID).
				Str("resolve_type", fmt.Sprintf("%+v", resolveResult.ResolvedType)).
				Str("alarm", alarm.ID).
				Msgf("send %s event", event.EventType)
		}
	}
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
