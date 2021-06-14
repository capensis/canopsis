package heartbeat

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type mongoAdapter struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

func (a *mongoAdapter) Get() ([]Heartbeat, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hl := make([]Heartbeat, 0)

	cursor, err := a.dbCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &hl)
	if err != nil {
		return nil, err
	}

	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	return hl, err
}

func NewAdapter(dbClient mongo.DbClient) Adapter {
	return &mongoAdapter{
		dbClient: dbClient,
		dbCollection: dbClient.Collection(mongo.HeartbeatMongoCollection),
	}
}
