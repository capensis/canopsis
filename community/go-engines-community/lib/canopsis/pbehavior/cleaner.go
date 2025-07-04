package pbehavior

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewCleaner(logger zerolog.Logger) datastorage.Cleaner {
	return &cleaner{
		logger:   logger,
		bulkSize: datastorage.BulkSize,
	}
}

type cleaner struct {
	logger   zerolog.Logger
	bulkSize int
}

func (c *cleaner) IsEnabled(conf datastorage.Config) bool {
	return datetime.IsDurationEnabledAndValid(conf.Pbehavior.DeleteAfter)
}

func (c *cleaner) Clean(ctx context.Context, dbClient mongo.DbClient, conf datastorage.Config, t datetime.CpsTime, limit int) (datastorage.CleanResult, error) {
	res := datastorage.CleanResult{}
	if !c.IsEnabled(conf) {
		return res, nil
	}

	defer func() {
		if res.Deleted > 0 {
			c.logger.Info().Int64("count", res.Deleted).Msg("pbehaviors were deleted")
		}
	}()

	opts := options.Find().SetProjection(bson.M{"_id": 1})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	before := conf.Pbehavior.DeleteAfter.SubFrom(t)
	collection := dbClient.Collection(mongo.PbehaviorMongoCollection)
	cursor, err := collection.Find(ctx, bson.M{"$or": []bson.M{
		{
			"rrule": bson.M{"$in": bson.A{"", nil}},
			"tstop": bson.M{
				"$ne":  nil,
				"$lte": before,
			},
		},
		{
			"rrule": bson.M{"$nin": bson.A{"", nil}},
			"rrule_end": bson.M{
				"$ne":  nil,
				"$lte": before,
			},
		},
	}}, opts)
	if err != nil {
		return res, fmt.Errorf("failed to find pbehaviors: %w", err)
	}

	defer cursor.Close(ctx)
	ids := make([]string, 0, c.bulkSize)
	for cursor.Next(ctx) {
		var pbh PBehavior
		err := cursor.Decode(&pbh)
		if err != nil {
			return res, fmt.Errorf("failed to decode pbehavior: %w", err)
		}

		ids = append(ids, pbh.ID)

		if len(ids) >= c.bulkSize {
			d, err := collection.DeleteMany(
				ctx,
				bson.M{"_id": bson.M{"$in": ids}},
			)
			if err != nil {
				return res, fmt.Errorf("failed to delete pbehaviors: %w", err)
			}

			res.Deleted += d
			ids = ids[:0]
		}
	}

	if err = cursor.Err(); err != nil {
		return res, fmt.Errorf("failed to fetch pbehaviors: %w", err)
	}

	if len(ids) > 0 {
		d, err := collection.DeleteMany(
			ctx,
			bson.M{"_id": bson.M{"$in": ids}},
		)
		if err != nil {
			return res, fmt.Errorf("failed to delete pbehaviors: %w", err)
		}

		res.Deleted += d
	}

	return res, nil
}
