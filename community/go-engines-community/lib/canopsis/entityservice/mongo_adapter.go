package entityservice

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func NewAdapter(client mongo.DbClient) Adapter {
	return &mongoAdapter{
		collection:      client.Collection(mongo.EntityMongoCollection),
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
	}
}

type mongoAdapter struct {
	collection      mongo.DbCollection
	alarmCollection mongo.DbCollection
}

func (a *mongoAdapter) GetValid(ctx context.Context) ([]EntityService, error) {
	return a.find(ctx, bson.M{"type": types.EntityTypeService, "enabled": true})
}

func (a *mongoAdapter) UpdateBulk(ctx context.Context, models []mongodriver.WriteModel) error {
	_, err := a.collection.BulkWrite(ctx, models)

	return err
}

func (a *mongoAdapter) find(ctx context.Context, filter interface{}) ([]EntityService, error) {
	cursor, err := a.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	res := make([]EntityService, 0)
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
