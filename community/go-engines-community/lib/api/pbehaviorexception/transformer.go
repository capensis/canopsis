package pbehaviorexception

import (
	"context"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type ModelTransformer interface {
	TransformCreateRequestToModel(request CreateRequest) (*Exception, error)
	TransformUpdateRequestToModel(request UpdateRequest) (*Exception, error)
	TransformExdatesRequestToModel(request []ExdateRequest) ([]Exdate, error)
}

func NewModelTransformer(dbClient mongo.DbClient) ModelTransformer {
	return &modelTransformer{
		dbClient:       dbClient,
		typeCollection: dbClient.Collection(pbehavior.TypeCollectionName),
	}
}

type modelTransformer struct {
	dbClient       mongo.DbClient
	typeCollection mongo.DbCollection
}

func (t *modelTransformer) TransformCreateRequestToModel(request CreateRequest) (*Exception, error) {
	exception := Exception{}

	exception.ID = request.ID
	exception.Name = request.Name
	exception.Description = request.Description
	exdates, err := t.TransformExdatesRequestToModel(request.Exdates)
	if err != nil {
		return nil, err
	}

	exception.Exdates = exdates

	return &exception, nil
}

func (t *modelTransformer) TransformUpdateRequestToModel(request UpdateRequest) (*Exception, error) {
	return t.TransformCreateRequestToModel(CreateRequest(request))
}

func (t *modelTransformer) TransformExdatesRequestToModel(request []ExdateRequest) ([]Exdate, error) {
	if len(request) == 0 {
		return []Exdate{}, nil
	}

	types := make([]string, len(request))
	for i := range request {
		types[i] = request[i].Type
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
