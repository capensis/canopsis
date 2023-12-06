package contextgraph

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type storage struct {
	client     mongo.DbClient
	collection mongo.DbCollection
}

//todo: decide about cache policy if needed

func NewEntityServiceStorage(client mongo.DbClient) EntityServiceStorage {
	return &storage{client: client, collection: client.Collection(mongo.EntityMongoCollection)}
}

func (s *storage) GetAll(ctx context.Context) ([]entityservice.EntityService, error) {
	var services []entityservice.EntityService
	cursor, err := s.collection.Find(ctx, bson.M{"type": types.EntityTypeService, "enabled": true})
	if err != nil {
		return nil, err
	}

	return services, cursor.All(ctx, &services)
}

func (s *storage) Get(ctx context.Context, serviceID string) (entityservice.EntityService, error) {
	var service entityservice.EntityService

	err := s.collection.FindOne(ctx, bson.M{"_id": serviceID}).Decode(&service)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return entityservice.EntityService{}, nil
		}
	}

	return service, err
}
