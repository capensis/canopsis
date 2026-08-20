package contextgraph

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type EntityService struct {
	ID               string         `bson:"_id"`
	EntityPattern    pattern.Entity `bson:"entity_pattern"`
	InheritedPattern pattern.Entity `bson:"inherited_pattern"`
	Enabled          bool           `bson:"enabled"`
}

type storage struct {
	client     mongo.DbClient
	collection mongo.DbCollection
}

//todo: decide about cache policy if needed

func NewEntityServiceStorage(client mongo.DbClient) EntityServiceStorage {
	return &storage{client: client, collection: client.Collection(mongo.EntityMongoCollection)}
}

func (s *storage) GetAll(ctx context.Context) ([]EntityService, error) {
	var services []EntityService
	cursor, err := s.collection.Find(
		ctx,
		bson.M{"type": types.EntityTypeService, "enabled": true},
		options.Find().SetProjection(bson.M{
			"_id":               1,
			"enabled":           1,
			"entity_pattern":    1,
			"inherited_pattern": "$state_info.inherited_pattern",
		}),
	)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (s *storage) Get(ctx context.Context, serviceID string) (entityservice.EntityService, error) {
	var service entityservice.EntityService

	err := s.collection.FindOne(ctx, bson.M{"_id": serviceID}).Decode(&service)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return entityservice.EntityService{}, nil
		}

		return entityservice.EntityService{}, err
	}

	return service, nil
}
