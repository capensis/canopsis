package statesetting

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type mongoAdapter struct {
	collection mongo.DbCollection
}

func NewMongoAdapter(client mongo.DbClient) Adapter {
	return &mongoAdapter{
		collection: client.Collection(mongo.StateSettingsMongoCollection),
	}
}

func (a *mongoAdapter) Get(ctx context.Context, settingID string) (*StateSetting, error) {
	var stateSettings StateSetting

	err := a.collection.FindOne(ctx, bson.M{"_id": settingID}).Decode(&stateSettings)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &stateSettings, err
}
