package entityservice

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (a *mongoAdapter) GetAll(ctx context.Context) ([]EntityService, error) {
	return a.find(ctx, bson.M{"type": types.EntityTypeService})
}

func (a *mongoAdapter) GetEnabled(ctx context.Context) ([]EntityService, error) {
	return a.find(ctx, bson.M{"type": types.EntityTypeService, "enabled": true})
}

func (a *mongoAdapter) GetByID(ctx context.Context, id string) (*EntityService, error) {
	return a.findOne(ctx, bson.M{"type": types.EntityTypeService, "_id": id})
}

func (a *mongoAdapter) AddDepends(ctx context.Context, id string, depends []string) (bool, error) {
	res, err := a.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$addToSet": bson.M{
			"depends": bson.M{"$each": depends},
		},
	})

	if err != nil {
		return false, err
	}

	return res.MatchedCount > 0, nil
}

func (a *mongoAdapter) RemoveDepends(ctx context.Context, id string, depends []string) (bool, error) {
	res, err := a.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$pull": bson.M{
			"depends": bson.M{"$in": depends},
		},
	})
	if err != nil {
		return false, err
	}

	return res.MatchedCount > 0, nil
}

func (a *mongoAdapter) RemoveDependByQuery(ctx context.Context, query interface{}, depend string) ([]string, error) {
	res, err := a.collection.Find(
		ctx,
		bson.M{"$and": []interface{}{
			query,
			bson.M{"depends": depend},
		}},
		options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}

	removedEntities := make([]types.Entity, 0)
	err = res.All(ctx, &removedEntities)
	if err != nil {
		return nil, err
	}

	removedIDs := make([]string, len(removedEntities))
	for i := range removedEntities {
		removedIDs[i] = removedEntities[i].ID
	}

	if len(removedIDs) > 0 {
		_, err = a.collection.UpdateMany(
			ctx,
			bson.M{"_id": bson.M{"$in": removedIDs}},
			bson.M{"$pull": bson.M{"depends": depend}},
		)
		if err != nil {
			return nil, err
		}
	}

	return removedIDs, nil
}

func (a *mongoAdapter) UpdateCounters(ctx context.Context, id string, counters AlarmCounters) error {
	_, err := a.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"counters": counters,
	}})

	return err
}

func (a *mongoAdapter) UpdateBulk(ctx context.Context, models []mongodriver.WriteModel) error {
	_, err := a.collection.BulkWrite(ctx, models)

	return err
}

func (a *mongoAdapter) GetServiceDependencies(
	ctx context.Context,
	serviceID string,
) (mongo.Cursor, error) {
	return a.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"enabled": true,
			"impact":  serviceID,
		}},
	})
}

func (a *mongoAdapter) GetDependenciesCount(
	ctx context.Context,
	serviceID string,
) (int64, error) {
	cursor, err := a.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id": serviceID,
		}},
		{"$project": bson.M{
			"depends_count": bson.M{"$size": "$depends"},
		}},
	})
	if err != nil {
		return 0, err
	}
	if cursor.Next(ctx) {
		res := struct {
			DependsCount int64 `bson:"depends_count"`
		}{}
		err = cursor.Decode(&res)
		return res.DependsCount, err
	}

	return 0, nil
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

func (a *mongoAdapter) findOne(ctx context.Context, filter interface{}) (*EntityService, error) {
	mongoRes := a.collection.FindOne(ctx, filter)
	if err := mongoRes.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	var res EntityService
	err := mongoRes.Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
