package pbehavior

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Cleaner interface {
	Clean(ctx context.Context, before datetime.CpsTime, limit int64) (int64, error)
}

func NewCleaner(client mongo.DbClient, bulkSize int, logger zerolog.Logger) Cleaner {
	return &cleaner{
		collection: client.Collection(mongo.PbehaviorMongoCollection),
		bulkSize:   bulkSize,
		logger:     logger,
	}
}

type cleaner struct {
	collection mongo.DbCollection
	bulkSize   int
	logger     zerolog.Logger
}

func (c *cleaner) Clean(ctx context.Context, before datetime.CpsTime, limit int64) (int64, error) {
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cursor, err := c.collection.Find(ctx, bson.M{"$or": []bson.M{
		{
			"rrule": bson.M{"$in": bson.A{"", nil}},
			"tstop": bson.M{
				"$ne": nil,
				"$lt": before,
			},
		},
		{
			"rrule": bson.M{"$nin": bson.A{"", nil}},
			"rrule_end": bson.M{
				"$ne": nil,
				"$lt": before,
			},
		},
	}}, opts)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	ids := make([]string, 0, c.bulkSize)
	var deleted int64

	for cursor.Next(ctx) {
		var pbh PBehavior
		err := cursor.Decode(&pbh)
		if err != nil {
			return 0, err
		}

		ids = append(ids, pbh.ID)

		if len(ids) >= c.bulkSize {
			res, err := c.collection.DeleteMany(
				ctx,
				bson.M{"_id": bson.M{"$in": ids}},
			)
			if err != nil {
				return 0, err
			}

			deleted += res
			ids = ids[:0]
		}
	}

	if len(ids) > 0 {
		res, err := c.collection.DeleteMany(
			ctx,
			bson.M{"_id": bson.M{"$in": ids}},
		)
		if err != nil {
			return 0, err
		}

		deleted += res
	}

	return deleted, nil
}
