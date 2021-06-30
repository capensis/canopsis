package action

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoAdapter struct {
	dbCollection mongo.DbCollection
}

func NewAdapter(dbClient mongo.DbClient) Adapter {
	return &mongoAdapter{
		dbCollection: dbClient.Collection(mongo.ScenarioMongoCollection),
	}
}

func (a *mongoAdapter) GetEnabled(ctx context.Context) ([]Scenario, error) {
	cursor, err := a.dbCollection.Find(ctx, bson.M{"$or": []bson.M{
		{"enabled": true},
		{"enabled": bson.M{"$exists": false}},
	}}, options.Find().SetSort(bson.D{
		{Key: PriorityField, Value: 1},
		{Key: IdField, Value: 1},
	}))
	if err != nil {
		return nil, err
	}
	scenarios := make([]Scenario, 0)

	for cursor.Next(ctx) {
		scenario := Scenario{}
		err := cursor.Decode(&scenario)
		if err != nil {
			return nil, err
		}

		scenarios = append(scenarios, scenario)
	}

	return scenarios, nil
}

func (a *mongoAdapter) GetEnabledById(ctx context.Context, id string) (Scenario, error) {
	scenario := Scenario{}
	res := a.dbCollection.FindOne(ctx, bson.M{"$and": []bson.M{
		{"_id": id},
		{"$or": []bson.M{
			{"enabled": true},
			{"enabled": bson.M{"$exists": false}},
		}},
	}})
	if err := res.Err(); err != nil {
		return scenario, err
	}

	err := res.Decode(&scenario)
	if err != nil {
		return scenario, err
	}

	return scenario, nil
}

func (a *mongoAdapter) GetEnabledByIDs(ctx context.Context, ids []string) ([]Scenario, error) {
	cursor, err := a.dbCollection.Find(ctx, bson.M{"$and": []bson.M{
		{"_id": bson.M{"$in": ids}},
		{"$or": []bson.M{
			{"enabled": true},
			{"enabled": bson.M{"$exists": false}},
		}},
	}}, options.Find().SetSort(bson.D{
		{Key: PriorityField, Value: 1},
		{Key: IdField, Value: 1},
	}))
	if err != nil {
		return nil, err
	}
	scenarios := make([]Scenario, 0)

	for cursor.Next(ctx) {
		scenario := Scenario{}
		err := cursor.Decode(&scenario)
		if err != nil {
			return nil, err
		}

		scenarios = append(scenarios, scenario)
	}

	return scenarios, nil
}
