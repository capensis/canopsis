package alarmtag

import (
	"context"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"
)

type InternalUpdater interface {
	Update(ctx context.Context, ch <-chan string) error
}

func NewInternalUpdater(client mongo.DbClient, logger zerolog.Logger) InternalUpdater {
	return &internalUpdater{
		client:          client,
		collection:      client.Collection(mongo.AlarmTagCollection),
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
		logger:          logger,
	}
}

type internalUpdater struct {
	client          mongo.DbClient
	collection      mongo.DbCollection
	alarmCollection mongo.DbCollection
	logger          zerolog.Logger

	tagsMx sync.RWMutex
	tags   []AlarmTag
}

func (u *internalUpdater) Update(ctx context.Context, ch <-chan string) error {
	err := u.load(ctx)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		cs, err := u.collection.Watch(ctx, []bson.M{})
		if err != nil {
			panic(err)
		}
		defer cs.Close(ctx)

		for cs.Next(ctx) {
			err = u.load(ctx)
			if err != nil {
				return err
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
				case alarmId, ok := <-ch:
					if !ok {
						return nil
					}

					err := u.updateAlarm(ctx, alarmId)
					if err != nil {
						u.logger.Err(err).Str("alarm", alarmId).Msg("cannot update alarm")
					}
				}
			}
		})
	}

	return g.Wait()
}

func (u *internalUpdater) load(ctx context.Context) error {
	cursor, err := u.collection.Find(ctx, bson.M{"type": TypeInternal})
	if err != nil {
		return err
	}
	var tags []AlarmTag
	err = cursor.All(ctx, &tags)
	if err != nil {
		return err
	}

	u.tagsMx.Lock()
	defer u.tagsMx.Unlock()
	u.tags = tags
	return nil
}

func (u *internalUpdater) updateAlarm(ctx context.Context, alarmId string) error {
	u.tagsMx.RLock()
	defer u.tagsMx.RUnlock()

	if len(u.tags) == 0 {
		return nil
	}

	err := u.client.WithTransaction(ctx, func(ctx context.Context) error {
		cursor, err := u.alarmCollection.Aggregate(ctx, []bson.M{
			{"$match": bson.M{
				"_id":        alarmId,
				"v.resolved": nil,
			}},
			{"$project": bson.M{"v.steps": 0}},
			{"$project": bson.M{
				"alarm": "$$ROOT",
			}},
			{"$lookup": bson.M{
				"from":         mongo.EntityMongoCollection,
				"foreignField": "_id",
				"localField":   "alarm.d",
				"as":           "entity",
			}},
			{"$unwind": "$entity"},
		})
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)
		if !cursor.Next(ctx) {
			return nil
		}

		var alarm types.AlarmWithEntity
		err = cursor.Decode(&alarm)
		if err != nil {
			return err
		}

		tags := make([]string, len(alarm.Alarm.ExternalTags))
		copy(tags, alarm.Alarm.ExternalTags)
		internalTags := make(map[string]types.CpsTime)
		for _, tag := range u.tags {
			if updated, ok := alarm.Alarm.InternalTags[tag.Value]; ok {
				if updated.After(tag.Updated) {
					tags = append(tags, tag.Value)
					internalTags[tag.Value] = tag.Updated
					continue
				}
			}

			if len(tag.EntityPattern) > 0 {
				ok, err := tag.EntityPattern.Match(alarm.Entity)
				if err != nil {
					u.logger.Err(err).Str("tag", tag.ID).Msg("tag has invalid entity pattern")
					continue
				}
				if !ok {
					continue
				}
			}
			if len(tag.AlarmPattern) > 0 {
				ok, err := tag.AlarmPattern.Match(alarm.Alarm)
				if err != nil {
					u.logger.Err(err).Str("tag", tag.ID).Msg("tag has invalid alarm pattern")
					continue
				}
				if !ok {
					continue
				}
			}

			tags = append(tags, tag.Value)
			internalTags[tag.Value] = tag.Updated
		}

		_, err = u.alarmCollection.UpdateOne(ctx, bson.M{"_id": alarm.Alarm.ID}, bson.M{
			"$set": bson.M{
				"tags":          tags,
				"internal_tags": internalTags,
			},
		})
		return err
	})

	return err
}
