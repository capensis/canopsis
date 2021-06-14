package datastorage

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Get(ctx context.Context) (datastorage.DataStorage, error)
	Update(context.Context, datastorage.Config) (datastorage.DataStorage, error)
}

func NewStore(client mongo.DbClient) Store {
	return &store{collection: client.Collection(mongo.ConfigurationMongoCollection)}
}

type store struct {
	collection mongo.DbCollection
}

func (s *store) Get(ctx context.Context) (datastorage.DataStorage, error) {
	data := datastorage.DataStorage{}
	err := s.collection.FindOne(ctx, bson.M{"_id": datastorage.ID}).Decode(&data)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return data, nil
		}

		return data, err
	}

	return data, nil
}

func (s *store) Update(ctx context.Context, conf datastorage.Config) (datastorage.DataStorage, error) {
	data := datastorage.DataStorage{}
	err := s.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": datastorage.ID},
		bson.M{"$set": bson.M{
			"config": conf,
		}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&data)

	return data, err
}
