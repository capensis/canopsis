package entityservice

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
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
			"enabled":  true,
			"services": serviceID,
		}},
	})
}

func (a *mongoAdapter) GetDependenciesCount(
	ctx context.Context,
	serviceID string,
) (int64, error) {
	return a.collection.CountDocuments(ctx, bson.M{
		"services": serviceID,
	})
}

func (a *mongoAdapter) AddToService(ctx context.Context, serviceId string, ids []string) error {
	_, err := a.collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": ids}},
		bson.M{"$push": bson.M{"services": serviceId}},
	)
	return err
}

func (a *mongoAdapter) RemoveFromService(ctx context.Context, serviceId string, ids []string) error {
	_, err := a.collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": ids}},
		bson.M{"$pull": bson.M{"services": serviceId}},
	)
	return err
}

func (a *mongoAdapter) AddToServiceByQuery(ctx context.Context, serviceId string, query bson.M) (int64, error) {
	var count int64
	and := []bson.M{
		{"enabled": true},
		query,
		{"services": bson.M{"$nin": bson.A{serviceId}}},
	}
	match := bson.M{"$and": and}
	for {
		res, err := a.collection.Find(
			ctx,
			match,
			options.Find().
				SetLimit(canopsis.DefaultBulkSize).
				SetProjection(bson.M{"_id": 1}))
		if err != nil {
			return 0, err
		}

		entities := make([]types.Entity, 0)
		err = res.All(ctx, &entities)
		if err != nil {
			return 0, err
		}
		if len(entities) == 0 {
			break
		}

		ids := make([]string, len(entities))
		for i := range entities {
			ids[i] = entities[i].ID
		}

		_, err = a.collection.UpdateMany(
			ctx,
			bson.M{
				"_id":  bson.M{"$in": ids},
				"$and": and,
			},
			bson.M{"$push": bson.M{"services": serviceId}},
		)
		if err != nil {
			return 0, err
		}

		count += int64(len(ids))
	}

	return count, nil
}

func (a *mongoAdapter) RemoveFromServiceByQuery(ctx context.Context, serviceId string, query bson.M) (int64, error) {
	and := []bson.M{{"enabled": true}}
	if query != nil {
		and = append(and, query)
	}
	and = append(and, bson.M{"services": bson.M{"$in": bson.A{serviceId}}})
	match := bson.M{"$and": and}
	var count int64
	for {
		res, err := a.collection.Find(
			ctx,
			match,
			options.Find().
				SetLimit(canopsis.DefaultBulkSize).
				SetProjection(bson.M{"_id": 1}))
		if err != nil {
			return 0, err
		}

		entities := make([]types.Entity, 0)
		err = res.All(ctx, &entities)
		if err != nil {
			return 0, err
		}

		if len(entities) == 0 {
			break
		}

		ids := make([]string, len(entities))
		for i := range entities {
			ids[i] = entities[i].ID
		}

		_, err = a.collection.UpdateMany(
			ctx,
			bson.M{
				"_id":  bson.M{"$in": ids},
				"$and": and,
			},
			bson.M{"$pull": bson.M{"services": serviceId}},
		)
		if err != nil {
			return 0, err
		}

		count += int64(len(ids))
	}

	return count, nil
}

func (a *mongoAdapter) find(ctx context.Context, filter bson.M) ([]EntityService, error) {
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

func (a *mongoAdapter) findOne(ctx context.Context, filter bson.M) (*EntityService, error) {
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
