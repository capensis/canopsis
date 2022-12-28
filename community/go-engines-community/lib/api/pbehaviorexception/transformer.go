package pbehaviorexception

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type ModelTransformer interface {
	TransformCreateRequestToModel(ctx context.Context, request CreateRequest) (*Exception, error)
	TransformUpdateRequestToModel(ctx context.Context, request UpdateRequest) (*Exception, error)
	TransformExdatesRequestToModel(ctx context.Context, request []ExdateRequest) ([]Exdate, error)
}

func NewModelTransformer(dbClient mongo.DbClient) ModelTransformer {
	return &modelTransformer{
		typeCollection: dbClient.Collection(mongo.PbehaviorTypeMongoCollection),
	}
}

type modelTransformer struct {
	typeCollection mongo.DbCollection
}

func (t *modelTransformer) TransformCreateRequestToModel(ctx context.Context, request CreateRequest) (*Exception, error) {
	exception := Exception{}

	exception.ID = request.ID
	exception.Name = request.Name
	exception.Description = request.Description
	exdates, err := t.TransformExdatesRequestToModel(ctx, request.Exdates)
	if err != nil {
		return nil, err
	}

	exception.Exdates = exdates

	return &exception, nil
}

func (t *modelTransformer) TransformUpdateRequestToModel(ctx context.Context, request UpdateRequest) (*Exception, error) {
	return t.TransformCreateRequestToModel(ctx, CreateRequest(request))
}

func (t *modelTransformer) TransformExdatesRequestToModel(ctx context.Context, request []ExdateRequest) ([]Exdate, error) {
	if len(request) == 0 {
		return []Exdate{}, nil
	}

	types := make([]string, len(request))
	for i := range request {
		types[i] = request[i].Type
	}

	res, err := t.typeCollection.Find(ctx, bson.M{"_id": bson.M{"$in": types}})
	if err != nil {
		return nil, err
	}

	defer res.Close(ctx)
	typesByID := make(map[string]*pbehavior.Type)
	for res.Next(ctx) {
		var t pbehavior.Type
		err = res.Decode(&t)
		if err != nil {
			return nil, err
		}
		if _, ok := typesByID[t.ID]; !ok {
			typesByID[t.ID] = &t
		}
	}

	exdates := make([]Exdate, len(request))
	for i := range request {
		if t, ok := typesByID[request[i].Type]; ok {
			exdates[i] = Exdate{
				Begin: request[i].Begin,
				End:   request[i].End,
				Type:  *t,
			}
		} else {
			return nil, ErrTypeNotExists
		}
	}

	return exdates, nil
}
