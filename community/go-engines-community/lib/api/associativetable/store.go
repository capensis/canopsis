package associativetable

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		dbCollection: dbClient.Collection(mongo.AssociativeTableCollection),
	}
}

type Store interface {
	Update(ctx context.Context, model *AssociativeTable) (bool, error)
	GetByName(ctx context.Context, name string) (*AssociativeTable, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbCollection mongo.DbCollection
}

func (s store) GetByName(ctx context.Context, name string) (*AssociativeTable, error) {
	at := &AssociativeTable{}
	err := s.dbCollection.
		FindOne(ctx, bson.M{"name": name}).
		Decode(at)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return at, nil
}

func (s store) Update(ctx context.Context, model *AssociativeTable) (bool, error) {
	res, err := s.dbCollection.UpdateOne(
		ctx,
		bson.M{"name": model.Name},
		bson.M{"$set": model},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		return false, err
	}

	return res.UpsertedCount > 0, nil
}

func (s store) Delete(ctx context.Context, name string) (bool, error) {
	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}
