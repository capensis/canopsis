package datastorage

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrConfNotFound = errors.New("cannot find configuration _id=" + ID)

type adapter struct {
	collection mongo.DbCollection
}

func NewAdapter(client mongo.DbClient) Adapter {
	return &adapter{
		collection: client.Collection(mongo.ConfigurationMongoCollection),
	}
}

func (a *adapter) Get(ctx context.Context) (DataStorage, error) {
	data := DataStorage{}
	err := a.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&data)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return data, nil
		}

		return data, err
	}

	return data, nil
}

func (a *adapter) UpdateHistoryEntityDisabled(ctx context.Context, history HistoryWithCount) error {
	return a.updateHistoryWithCount(ctx, "history.entity_disabled", history)
}

func (a *adapter) UpdateHistoryEntityUnlinked(ctx context.Context, history HistoryWithCount) error {
	return a.updateHistoryWithCount(ctx, "history.entity_unlinked", history)
}

func (a *adapter) UpdateHistoryEntityCleaned(ctx context.Context, history HistoryWithCount) error {
	return a.updateHistoryWithCount(ctx, "history.entity_cleaned", history)
}

func (a *adapter) updateHistoryWithCount(ctx context.Context, key string, history HistoryWithCount) error {
	return a.updateOne(ctx, bson.M{key: history})
}

func (a *adapter) updateOne(ctx context.Context, upd bson.M) error {
	res, err := a.collection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$set": upd})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrConfNotFound
	}

	return nil
}
