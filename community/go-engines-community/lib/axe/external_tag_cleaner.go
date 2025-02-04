package axe

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewExternalTagCleaner(logger zerolog.Logger) datastorage.Cleaner {
	return &externalTagCleaner{
		logger: logger,
	}
}

type externalTagCleaner struct {
	logger zerolog.Logger
}

func (c *externalTagCleaner) IsEnabled(conf datastorage.Config) bool {
	return datetime.IsDurationEnabledAndValid(conf.AlarmExternalTag.DeleteAfter)
}

func (c *externalTagCleaner) Clean(ctx context.Context, dbClient mongo.DbClient, conf datastorage.Config, t datetime.CpsTime, limit int) (datastorage.CleanResult, error) {
	res := datastorage.CleanResult{}
	if !c.IsEnabled(conf) {
		return res, nil
	}

	var deletedColors int64
	defer func() {
		if res.Deleted > 0 {
			c.logger.Info().Int64("count", res.Deleted).Msg("alarm external tags are deleted")
		}

		if deletedColors > 0 {
			c.logger.Info().Int64("count", deletedColors).Msg("alarm tag colors are deleted")
		}
	}()

	dbCollection := dbClient.Collection(mongo.AlarmTagCollection)
	colorDbCollection := dbClient.Collection(mongo.AlarmTagColorCollection)
	var err error
	res.Deleted, err = c.delete(ctx, conf.AlarmExternalTag.DeleteAfter.SubFrom(t), limit, datastorage.BulkSize, dbCollection)
	if err != nil {
		return res, fmt.Errorf("cannot delete tags: %w", err)
	}

	deletedColors, err = c.deleteColors(ctx, limit, datastorage.BulkSize, colorDbCollection)
	if err != nil {
		return res, fmt.Errorf("cannot delete colors: %w", err)
	}

	return res, nil
}

func (c *externalTagCleaner) delete(
	ctx context.Context,
	before datetime.CpsTime,
	limit int,
	bulkSize int,
	dbCollection mongo.DbCollection,
) (int64, error) {
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := dbCollection.Find(ctx, bson.M{
		"type":            alarmtag.TypeExternal,
		"last_event_date": bson.M{"$lte": before},
	}, opts)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)
	ids := make([]string, 0, bulkSize)
	var deleted int64
	for cursor.Next(ctx) {
		var tag alarmtag.AlarmTag
		err := cursor.Decode(&tag)
		if err != nil {
			return deleted, err
		}

		ids = append(ids, tag.ID)
		if len(ids) >= bulkSize {
			res, err := dbCollection.DeleteMany(
				ctx,
				bson.M{"_id": bson.M{"$in": ids}},
			)
			if err != nil {
				return deleted, err
			}

			deleted += res
			ids = ids[:0]
		}
	}

	if err = cursor.Err(); err != nil {
		return deleted, err
	}

	if len(ids) > 0 {
		res, err := dbCollection.DeleteMany(
			ctx,
			bson.M{"_id": bson.M{"$in": ids}},
		)
		if err != nil {
			return deleted, err
		}

		deleted += res
	}

	return deleted, nil
}
func (c *externalTagCleaner) deleteColors(
	ctx context.Context,
	limit int,
	bulkSize int,
	dbCollection mongo.DbCollection,
) (int64, error) {
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.AlarmTagCollection,
			"localField":   "_id",
			"foreignField": "label",
			"as":           "tags",
			"pipeline": []bson.M{
				{"$limit": 1},
				{"$project": bson.M{
					"_id": 1,
				}},
			},
		}},
		{"$match": bson.M{
			"tags": bson.A{},
		}},
		{"$project": bson.M{
			"_id": 1,
		}},
	}

	if limit > 0 {
		pipeline = append(pipeline, bson.M{
			"$limit": limit,
		})
	}

	cursor, err := dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)
	ids := make([]string, 0, bulkSize)
	var deleted int64
	for cursor.Next(ctx) {
		var color struct {
			ID string `bson:"_id"`
		}
		err := cursor.Decode(&color)
		if err != nil {
			return deleted, err
		}

		ids = append(ids, color.ID)
		if len(ids) >= bulkSize {
			res, err := dbCollection.DeleteMany(
				ctx,
				bson.M{"_id": bson.M{"$in": ids}},
			)
			if err != nil {
				return deleted, err
			}

			deleted += res
			ids = ids[:0]
		}
	}

	if err = cursor.Err(); err != nil {
		return deleted, err
	}

	if len(ids) > 0 {
		res, err := dbCollection.DeleteMany(
			ctx,
			bson.M{"_id": bson.M{"$in": ids}},
		)
		if err != nil {
			return deleted, err
		}

		deleted += res
	}

	return deleted, nil
}
