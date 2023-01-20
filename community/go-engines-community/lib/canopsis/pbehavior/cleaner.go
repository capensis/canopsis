package pbehavior

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"github.com/teambition/rrule-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Cleaner interface {
	Clean(ctx context.Context, before types.CpsTime, limit int64) (int64, error)
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

func (c *cleaner) Clean(ctx context.Context, before types.CpsTime, limit int64) (int64, error) {
	deleted, err := c.cleanPbhWithoutRrule(ctx, before, limit)
	if err != nil {
		return 0, err
	}

	deletedWithRrule, err := c.cleanPbhWithRrule(ctx, before, limit)
	if err != nil {
		return 0, err
	}

	deleted += deletedWithRrule
	return deleted, nil
}

func (c *cleaner) cleanPbhWithoutRrule(ctx context.Context, before types.CpsTime, limit int64) (int64, error) {
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cursor, err := c.collection.Find(ctx, bson.M{
		"rrule": bson.M{"$in": bson.A{"", nil}},
		"tstop": bson.M{
			"$gt": 0,
			"$lt": before,
		},
	}, opts)
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

func (c *cleaner) cleanPbhWithRrule(ctx context.Context, before types.CpsTime, limit int64) (int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cursor, err := c.collection.Find(ctx, bson.M{
		"rrule": bson.M{"$nin": bson.A{"", nil}},
		"tstop": bson.M{"$gt": 0},
	}, opts)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	ids := make([]string, 0, c.bulkSize)
	var deleted int64

	for cursor.Next(ctx) {
		var pbehavior PBehavior
		err := cursor.Decode(&pbehavior)
		if err != nil {
			return 0, err
		}

		rOption, err := rrule.StrToROption(pbehavior.RRule)
		if err != nil {
			c.logger.Err(err).
				Str("pbehavior", pbehavior.ID).
				Str("rrule", pbehavior.RRule).
				Msg("invalid rrule in pbehavior")
			continue
		}

		rOption.Dtstart = pbehavior.Start.Time
		rRule, err := rrule.NewRRule(*rOption)
		if err != nil {
			c.logger.Err(err).
				Str("pbehavior", pbehavior.ID).
				Str("rrule", pbehavior.RRule).
				Msg("invalid rrule in pbehavior")
			continue
		}

		if !isRruleFinished(*rRule, before.Time) {
			continue
		}

		ids = append(ids, pbehavior.ID)

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

func isRruleFinished(rRule rrule.RRule, before time.Time) bool {
	if !rRule.Options.Until.IsZero() && rRule.Options.Until.Before(before) {
		return true
	}

	if rRule.Options.Count > 0 {
		after := rRule.After(before, false)
		if after.IsZero() {
			return true
		}
	}

	return false
}
