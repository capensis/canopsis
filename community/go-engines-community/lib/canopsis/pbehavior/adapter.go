package pbehavior

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Adapter interface {
	UpdateLastAlarmDate(ctx context.Context, id string, time datetime.CpsTime) error
}

func NewAdapter(client mongo.DbClient) Adapter {
	return &adapter{
		collection: client.Collection(mongo.PbehaviorMongoCollection),
	}
}

type adapter struct {
	collection mongo.DbCollection
}

func (a *adapter) UpdateLastAlarmDate(ctx context.Context, id string, time datetime.CpsTime) error {
	_, err := a.collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"last_alarm_date": time}},
	)
	return err
}
