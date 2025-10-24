package che

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewEventFailureCleaner(logger zerolog.Logger) datastorage.Cleaner {
	return &eventFailureCleaner{
		logger: logger,
	}
}

type eventFailureCleaner struct {
	logger zerolog.Logger
}

func (c *eventFailureCleaner) IsEnabled(conf datastorage.Config) bool {
	return datetime.IsDurationEnabledAndValid(conf.EventFilterFailure.DeleteAfter)
}

func (c *eventFailureCleaner) Clean(ctx context.Context, dbClient mongo.DbClient, conf datastorage.Config, t datetime.CpsTime, limit int) (datastorage.CleanResult, error) {
	res := datastorage.CleanResult{}
	if !c.IsEnabled(conf) {
		return res, nil
	}

	defer func() {
		if res.Deleted > 0 {
			c.logger.Info().Int64("count", res.Deleted).Msg("event filter failures were deleted")
		}
	}()

	dbCollection := dbClient.Collection(mongo.EventFilterFailureCollection)
	dbRuleCollection := dbClient.Collection(mongo.EventFilterRuleCollection)
	opts := options.Find().SetProjection(bson.M{
		"_id":    1,
		"rule":   1,
		"unread": 1,
	})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := dbCollection.Find(ctx, bson.M{
		"t": bson.M{"$lte": conf.EventFilterFailure.DeleteAfter.SubFrom(t)},
	}, opts)
	if err != nil {
		return res, fmt.Errorf("failed to find failures: %w", err)
	}

	defer cursor.Close(ctx)
	ids := make([]string, 0, datastorage.BulkSize)
	countsByRule := make(map[string]int64, datastorage.BulkSize)
	unreadCountsByRule := make(map[string]int64, datastorage.BulkSize)
	ruleWriteModels := make([]mongodriver.WriteModel, 0, datastorage.BulkSize)
	for cursor.Next(ctx) {
		var item eventfilter.Failure
		err := cursor.Decode(&item)
		if err != nil {
			return res, fmt.Errorf("failed to decode failure: %w", err)
		}

		ids = append(ids, item.ID)
		countsByRule[item.Rule]++
		if item.Unread {
			unreadCountsByRule[item.Rule]++
		}

		if len(ids) >= datastorage.BulkSize {
			for ruleID, dec := range countsByRule {
				ruleWriteModels = append(ruleWriteModels, c.getRuleUpdateModel(ruleID, dec, unreadCountsByRule[ruleID]))
			}

			_, err = dbRuleCollection.BulkWrite(ctx, ruleWriteModels)
			if err != nil {
				return res, fmt.Errorf("failed to update event filter rules: %w", err)
			}

			d, err := dbCollection.DeleteMany(
				ctx,
				bson.M{"_id": bson.M{"$in": ids}},
			)
			if err != nil {
				return res, fmt.Errorf("failed to delete failures: %w", err)
			}

			res.Deleted += d
			ids = ids[:0]
			clear(countsByRule)
			clear(unreadCountsByRule)
			ruleWriteModels = ruleWriteModels[:0]
		}
	}

	if err = cursor.Err(); err != nil {
		return res, fmt.Errorf("failed to fetch failures: %w", err)
	}

	if len(ids) > 0 {
		for ruleID, dec := range countsByRule {
			ruleWriteModels = append(ruleWriteModels, c.getRuleUpdateModel(ruleID, dec, unreadCountsByRule[ruleID]))
		}

		_, err = dbRuleCollection.BulkWrite(ctx, ruleWriteModels)
		if err != nil {
			return res, fmt.Errorf("failed to update event filter rules: %w", err)
		}

		d, err := dbCollection.DeleteMany(
			ctx,
			bson.M{"_id": bson.M{"$in": ids}},
		)
		if err != nil {
			return res, fmt.Errorf("failed to delete failures: %w", err)
		}

		res.Deleted += d
	}

	return res, nil
}

func (c *eventFailureCleaner) getRuleUpdateModel(ruleID string, dec, unreadDec int64) mongodriver.WriteModel {
	update := bson.M{
		"failures_count": -dec,
	}
	if unreadDec > 0 {
		update["unread_failures_count"] = -unreadDec
	}

	return mongodriver.NewUpdateOneModel().
		SetFilter(bson.M{"_id": ruleID}).
		SetUpdate(bson.M{"$inc": update})
}
