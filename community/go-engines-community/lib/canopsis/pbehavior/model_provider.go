package pbehavior

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"go.mongodb.org/mongo-driver/bson"
)

// ModelProvider is used to implement fetching models from storage.
type ModelProvider interface {
	// GetTypes returns types by id.
	GetTypes(ctx context.Context) (map[string]Type, error)
	// GetEnabledPbehaviors returns pbehaviors.
	GetEnabledPbehaviors(ctx context.Context, span timespan.Span) (map[string]PBehavior, error)
	// GetEnabledPbehaviorsByIds returns pbehaviors.
	GetEnabledPbehaviorsByIds(ctx context.Context, ids []string, span timespan.Span) (map[string]PBehavior, error)
	// GetExceptions returns exceptions by id.
	GetExceptions(ctx context.Context) (map[string]Exception, error)
	// GetReasons returns reasons by id.
	GetReasons(ctx context.Context) (map[string]Reason, error)
}

// modelProvider implements fetching models from mongo db.
type modelProvider struct {
	dbClient mongo.DbClient
}

// NewModelProvider creates new model provider.
func NewModelProvider(dbClient mongo.DbClient) ModelProvider {
	return &modelProvider{dbClient: dbClient}
}

func (p *modelProvider) GetTypes(ctx context.Context) (map[string]Type, error) {
	cursor, err := p.dbClient.Collection(mongo.PbehaviorTypeMongoCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	typesByID := make(map[string]Type)

	for cursor.Next(ctx) {
		var pbehaviorType Type
		err = cursor.Decode(&pbehaviorType)
		if err != nil {
			return nil, err
		}

		typesByID[pbehaviorType.ID] = pbehaviorType
	}

	return typesByID, nil
}

func (p *modelProvider) GetEnabledPbehaviors(ctx context.Context, span timespan.Span) (map[string]PBehavior, error) {
	coll := p.dbClient.Collection(mongo.PbehaviorMongoCollection)
	cursor, err := coll.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"enabled": true,
			"$or": []bson.M{
				{
					"rrule": bson.M{"$nin": bson.A{nil, ""}},
					"$or": []bson.M{
						{"rrule_end": nil},
						{"rrule_end": bson.M{"$gte": span.From().Unix()}},
					},
				},
				{
					"tstart": bson.M{"$lte": span.To().Unix()},
					"$or": []bson.M{
						{"tstop": nil},
						{"tstop": bson.M{"$gte": span.From().Unix()}},
					},
				},
			},
		}},
		{"$project": bson.M{
			"comments": 0,
		}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	pbehaviorsByID := make(map[string]PBehavior)
	for cursor.Next(ctx) {
		var pbehavior PBehavior
		err = cursor.Decode(&pbehavior)
		if err != nil {
			return nil, err
		}

		pbehaviorsByID[pbehavior.ID] = pbehavior
	}

	return pbehaviorsByID, nil
}

func (p *modelProvider) GetEnabledPbehaviorsByIds(ctx context.Context, ids []string, span timespan.Span) (map[string]PBehavior, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	coll := p.dbClient.Collection(mongo.PbehaviorMongoCollection)
	cursor, err := coll.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id":     bson.M{"$in": ids},
			"enabled": true,
			"$or": []bson.M{
				{"rrule": bson.M{"$nin": bson.A{nil, ""}}},
				{
					"tstart": bson.M{"$lte": span.To().Unix()},
					"$or": []bson.M{
						{"tstop": nil},
						{"tstop": bson.M{"$gte": span.From().Unix()}},
					},
				},
			},
		}},
		{"$project": bson.M{
			"comments": 0,
		}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	pbehaviorsByID := make(map[string]PBehavior)
	for cursor.Next(ctx) {
		var pbehavior PBehavior
		err = cursor.Decode(&pbehavior)
		if err != nil {
			return nil, err
		}

		pbehaviorsByID[pbehavior.ID] = pbehavior
	}

	return pbehaviorsByID, nil
}

func (p *modelProvider) GetExceptions(ctx context.Context) (map[string]Exception, error) {
	cursor, err := p.dbClient.Collection(mongo.PbehaviorExceptionMongoCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	exceptionsByID := make(map[string]Exception)
	for cursor.Next(ctx) {
		var exception Exception
		err = cursor.Decode(&exception)
		if err != nil {
			return nil, err
		}

		exceptionsByID[exception.ID] = exception
	}

	return exceptionsByID, nil
}

func (p *modelProvider) GetReasons(ctx context.Context) (map[string]Reason, error) {
	cursor, err := p.dbClient.Collection(mongo.PbehaviorReasonMongoCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	reasonsByID := make(map[string]Reason)
	for cursor.Next(ctx) {
		var reason Reason
		err = cursor.Decode(&reason)
		if err != nil {
			return nil, err
		}

		reasonsByID[reason.ID] = reason
	}

	return reasonsByID, nil
}
