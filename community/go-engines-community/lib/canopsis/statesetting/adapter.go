package statesetting

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type mongoAdapter struct {
	collection mongo.DbCollection
}

func NewMongoAdapter(client mongo.DbClient) Adapter {
	return &mongoAdapter{
		collection: client.Collection(mongo.StateSettingsMongoCollection),
	}
}

func (a *mongoAdapter) Get(ctx context.Context, settingID string) (StateSetting, error) {
	var stateSettings StateSetting
	res := a.collection.FindOne(ctx, bson.M{"_id": settingID})
	if res.Err() != nil {
		return stateSettings, res.Err()
	}

	err := res.Decode(&stateSettings)
	return stateSettings, err
}
