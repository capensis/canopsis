package statesetting

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
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

func (a *mongoAdapter) Get(settingType string) (StateSetting, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var stateSettings StateSetting

	res := a.collection.FindOne(ctx, bson.M{"type": settingType})
	if res.Err() != nil {
		return stateSettings, res.Err()
	}

	err := res.Decode(&stateSettings)
	return stateSettings, err
}


