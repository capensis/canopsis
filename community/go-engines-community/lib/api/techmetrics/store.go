package techmetrics

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	GetSettings(ctx context.Context) (Settings, error)
	UpdateSettings(ctx context.Context, settings Settings) error
}

type store struct {
	collection mongo.DbCollection
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		collection: dbClient.Collection(mongo.ConfigurationMongoCollection),
	}
}

func (s *store) GetSettings(ctx context.Context) (Settings, error) {
	res := Settings{}
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": config.ConfigKeyName}},
		{"$replaceRoot": bson.M{"newRoot": "$tech_metrics"}},
	})
	if err != nil {
		return res, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
	}

	return res, err
}

func (s *store) UpdateSettings(ctx context.Context, settings Settings) error {
	_, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": config.ConfigKeyName},
		bson.M{"$set": bson.M{
			"tech_metrics.enabled": settings.Enabled,
		}},
	)

	return err
}
