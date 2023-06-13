package alarmtag

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

const (
	WatchMessageTypeInsert = iota
	WatchMessageTypeUpdate
	WatchMessageTypeDelete
)

const workers = 1000

type WatchMessage struct {
	Tags []AlarmTag
	Type int
}

type InternalWatcher interface {
	Update(ctx context.Context, ch <-chan WatchMessage) error
}

func NewInternalWatcher(client mongo.DbClient, logger zerolog.Logger) InternalWatcher {
	return &internalWatcher{
		client:          client,
		collection:      client.Collection(mongo.AlarmTagCollection),
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
		logger:          logger,
	}
}

type internalWatcher struct {
	client          mongo.DbClient
	collection      mongo.DbCollection
	alarmCollection mongo.DbCollection
	logger          zerolog.Logger
}

func (u *internalWatcher) Update(ctx context.Context, ch <-chan WatchMessage) error {
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < workers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case m, ok := <-ch:
					if !ok {
						return nil
					}

					err := u.updateAlarms(ctx, m)
					if err != nil {
						u.logger.Err(err).Msg("cannot update alarms on tag change")
					}
				}
			}
		})
	}

	return g.Wait()
}

func (u *internalWatcher) updateAlarms(ctx context.Context, m WatchMessage) error {
	for _, tag := range m.Tags {
		switch m.Type {
		case WatchMessageTypeInsert:
			err := u.addTag(ctx, tag)
			if err != nil {
				return err
			}
		case WatchMessageTypeUpdate:
			err := u.addTag(ctx, tag)
			if err != nil {
				return err
			}

			err = u.removeTag(ctx, tag, bson.M{"$lt": tag.Updated})
			if err != nil {
				return err
			}
		case WatchMessageTypeDelete:
			err := u.removeTag(ctx, tag, bson.M{"$ne": nil})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (u *internalWatcher) addTag(ctx context.Context, tag AlarmTag) error {
	pipeline := []bson.M{
		{"$match": bson.M{"v.resolved": nil}},
	}
	if len(tag.AlarmPattern) > 0 {
		q, err := tag.AlarmPattern.ToMongoQuery("")
		if err != nil {
			return err
		}
		pipeline = append(pipeline, bson.M{"$match": q})
	}
	pipeline = append(pipeline,
		bson.M{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"foreignField": "_id",
			"localField":   "d",
			"as":           "entity",
		}},
		bson.M{"$unwind": "$entity"},
	)
	if len(tag.EntityPattern) > 0 {
		q, err := tag.EntityPattern.ToMongoQuery("entity")
		if err != nil {
			return err
		}
		pipeline = append(pipeline, bson.M{"$match": q})
	}
	pipeline = append(pipeline, bson.M{"$project": bson.M{"_id": 1}})
	cursor, err := u.alarmCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	writeModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
	field := "internal_tags." + tag.Value
	for cursor.Next(ctx) {
		var alarm types.Alarm
		err = cursor.Decode(&alarm)
		if err != nil {
			return err
		}

		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{
				"_id":        alarm.ID,
				"v.resolved": nil,
				"$or": []bson.M{
					{field: nil},
					{field: bson.M{"$lt": tag.Updated}},
				},
			}).
			SetUpdate(bson.M{
				"$set": bson.M{
					field: tag.Updated,
				},
				"$addToSet": bson.M{
					"tags": tag.Value,
				},
			}))

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err = u.alarmCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return err
			}
			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err = u.alarmCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *internalWatcher) removeTag(
	ctx context.Context,
	tag AlarmTag,
	tagMatch bson.M,
) error {
	field := "internal_tags." + tag.Value
	pipeline := []bson.M{
		{"$match": bson.M{
			"v.resolved": nil,
			field:        tagMatch,
		}},
		{"$project": bson.M{"_id": 1}},
	}

	cursor, err := u.alarmCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	writeModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
	for cursor.Next(ctx) {
		var alarm types.Alarm
		err = cursor.Decode(&alarm)
		if err != nil {
			return err
		}

		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{
				"_id":        alarm.ID,
				"v.resolved": nil,
				field:        tagMatch,
			}).
			SetUpdate(bson.M{
				"$unset": bson.M{
					field: "",
				},
				"$pull": bson.M{
					"tags": tag.Value,
				},
			}))

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err = u.alarmCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return err
			}
			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err = u.alarmCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}
	}

	return nil
}
