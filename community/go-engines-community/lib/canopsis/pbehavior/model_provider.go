package pbehavior

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// ModelProvider is used to implement fetching models from storage.
type ModelProvider interface {
	// GetTypes returns types by id.
	GetTypes(ctx context.Context) (map[string]*Type, error)
	// GetPbehaviors returns pbehaviors by id.
	GetEnabledPbehaviors(ctx context.Context) (map[string]*PBehavior, error)
	// GetPbehavior returns pbehavior.
	GetEnabledPbehavior(ctx context.Context, id string) (*PBehavior, error)
	// GetExceptions returns exceptions by id.
	GetExceptions(ctx context.Context) (map[string]*Exception, error)
	// GetReasons returns reasons by id.
	GetReasons(ctx context.Context) (map[string]*Reason, error)
}

// modelProvider implements fetching models from mongo db.
type modelProvider struct {
	dbClient mongo.DbClient
}

// NewModelProvider creates new model provider.
func NewModelProvider(dbClient mongo.DbClient) ModelProvider {
	return &modelProvider{dbClient: dbClient}
}

func (p *modelProvider) GetTypes(parentCtx context.Context) (map[string]*Type, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	cursor, err := p.dbClient.Collection(TypeCollectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	typesByID := map[string]*Type{}

	for cursor.Next(ctx) {
		var pbehaviorType Type

		err = cursor.Decode(&pbehaviorType)
		if err != nil {
			return nil, err
		}

		typesByID[pbehaviorType.ID] = &pbehaviorType
	}

	return typesByID, nil
}

func (p *modelProvider) GetEnabledPbehaviors(parentCtx context.Context) (map[string]*PBehavior, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	coll := p.dbClient.Collection(PBehaviorCollectionName)
	cursor, err := coll.Find(ctx, bson.M{"enabled": true})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	pbehaviorsByID := make(map[string]*PBehavior)
	for cursor.Next(ctx) {
		var pbehavior PBehavior

		err = cursor.Decode(&pbehavior)
		if err != nil {
			return nil, err
		}

		pbehaviorsByID[pbehavior.ID] = &pbehavior
	}

	return pbehaviorsByID, nil
}

func (p *modelProvider) GetEnabledPbehavior(parentCtx context.Context, id string) (*PBehavior, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	coll := p.dbClient.Collection(PBehaviorCollectionName)
	res := coll.FindOne(ctx, bson.M{"_id": id, "enabled": true})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	var pbehavior PBehavior
	err := res.Decode(&pbehavior)
	if err != nil {
		return nil, err
	}

	return &pbehavior, nil
}

func (p *modelProvider) GetExceptions(parentCtx context.Context) (map[string]*Exception, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	cursor, err := p.dbClient.Collection(ExceptionCollectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	exceptionsByID := make(map[string]*Exception)
	for cursor.Next(ctx) {
		var exception Exception

		err = cursor.Decode(&exception)
		if err != nil {
			return nil, err
		}

		exceptionsByID[exception.ID] = &exception
	}

	return exceptionsByID, nil
}

func (p *modelProvider) GetReasons(parentCtx context.Context) (map[string]*Reason, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	cursor, err := p.dbClient.Collection(ReasonCollectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	reasonsByID := make(map[string]*Reason)
	for cursor.Next(ctx) {
		var reason Reason

		err = cursor.Decode(&reason)
		if err != nil {
			return nil, err
		}

		reasonsByID[reason.ID] = &reason
	}

	return reasonsByID, nil
}
