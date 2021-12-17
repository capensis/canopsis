package ratelimit

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type adapter struct {
	dbCollection mongo.DbCollection
}

func NewAdapter(dbClient mongo.DbClient) Adapter {
	return &adapter{dbCollection: dbClient.Collection(mongo.MessageRateStatsHourCollectionName)}
}

func (a *adapter) DeleteBefore(ctx context.Context, before types.CpsTime) (int64, error) {
	return a.dbCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$lt": before}})
}
