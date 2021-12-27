package pbehavior

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"github.com/teambition/rrule-go"
	"go.mongodb.org/mongo-driver/bson"
)

type Cleaner interface {
	Clean(ctx context.Context, before types.CpsTime) (int64, error)
}

func NewCleaner(client mongo.DbClient, logger zerolog.Logger) Cleaner {
	return &cleaner{
		collection: client.Collection(mongo.PbehaviorMongoCollection),
		logger:     logger,
	}
}

type cleaner struct {
	collection mongo.DbCollection
	logger     zerolog.Logger
}

func (c *cleaner) Clean(ctx context.Context, before types.CpsTime) (int64, error) {
	// Delete pbehaviors without rrule.
	deleted, err := c.collection.DeleteMany(ctx, bson.M{
		"rrule": bson.M{"$in": bson.A{"", nil}},
		"tstop": bson.M{
			"$gt": 0,
			"$lt": before,
		},
	})
	if err != nil {
		return deleted, err
	}

	// Delete pbehaviors with rrule.
	cursor, err := c.collection.Find(ctx, bson.M{
		"rrule": bson.M{"$nin": bson.A{"", nil}},
		"tstop": bson.M{"$gt": 0},
	})
	if err != nil {
		return deleted, err
	}

	defer cursor.Close(ctx)

	ids := make([]string, 0)
	for cursor.Next(ctx) {
		pbehavior := PBehavior{}
		err := cursor.Decode(&pbehavior)
		if err != nil {
			c.logger.Err(err).Msg("cannot decode pbehavior")
			continue
		}

		rOption, err := rrule.StrToROption(pbehavior.RRule)
		if err != nil {
			c.logger.Err(err).Str("pbehavior", pbehavior.ID).Str("rrule", pbehavior.RRule).
				Msg("invalid rrule in pbehavior")
			continue
		}

		rOption.Dtstart = pbehavior.Start.Time
		rRule, err := rrule.NewRRule(*rOption)
		if err != nil {
			c.logger.Err(err).Str("pbehavior", pbehavior.ID).Str("rrule", pbehavior.RRule).
				Msg("invalid rrule in pbehavior")
			continue
		}

		if isRruleFinished(*rRule, before.Time) {
			ids = append(ids, pbehavior.ID)
		}
	}

	if len(ids) > 0 {
		rruleDeleted, err := c.collection.DeleteMany(ctx, bson.M{
			"_id": bson.M{"$in": ids},
		})
		if err != nil {
			return deleted, err
		}

		deleted += rruleDeleted
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
